package runner

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	ErrNoCommand = errors.New("no command")
)

type EnvVars = map[string]string

type Config struct {
	EnvFile       string
	Command       []string
	Timeout       time.Duration
	LogLevel      string
	MaskSecrets   bool
	AllowOverride bool
	InheritEnv    bool
}

// EnvRunner manages process execution with environment variables
type EnvRunner struct {
	config Config
	logger *log.Logger
}

func NewEnvRunner(config Config) *EnvRunner {
	logger := log.New(os.Stderr, "[env-runner] ", log.LstdFlags)

	return &EnvRunner{
		config: config,
		logger: logger,
	}
}

func (er *EnvRunner) Run(envVars EnvVars) error {
	if len(er.config.Command) == 0 {
		return ErrNoCommand
	}

	env := er.buildEnvironment(envVars)

	// Create context with timeout if specified
	ctx := context.Background()
	if er.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, er.config.Timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, er.config.Command[0], er.config.Command[1:]...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	er.logger.Printf("Starting command: %s", strings.Join(er.config.Command, " "))

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	go func() {
		sig := <-sigChan
		er.logger.Printf("Received signal %v, forwarding to child process", sig)

		if cmd.Process != nil {
			// Forward signal to child process
			err := cmd.Process.Signal(sig)
			if err != nil {
				er.logger.Printf("Failed to send signal to child: %v", err)
			}
		}
	}()

	err := cmd.Wait()

	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitCode := exitError.ExitCode()
			er.logger.Printf("Command exited with code %d", exitCode)
			os.Exit(exitCode)
		}
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}

func (er *EnvRunner) buildEnvironment(envVars EnvVars) []string {
	var env []string

	if er.config.InheritEnv {
		// Start with current environment
		env = os.Environ()

		if !er.config.AllowOverride {
			// Add only new variables, don't override existing ones
			existing := make(map[string]bool)
			for _, e := range env {
				if idx := strings.Index(e, "="); idx != -1 {
					existing[e[:idx]] = true
				}
			}

			for key, value := range envVars {
				if !existing[key] {
					env = append(env, fmt.Sprintf("%s=%s", key, value))
				} else {
					er.logger.Printf("Skipping override of existing variable: %s", key)
				}
			}
		} else {
			// Allow overrides - add all variables
			for key, value := range envVars {
				env = append(env, fmt.Sprintf("%s=%s", key, value))
			}
		}
	} else {
		// Start with clean environment - only use template variables
		env = make([]string, 0, len(envVars))
		for key, value := range envVars {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	return env
}
