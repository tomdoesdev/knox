package commands

import (
	"context"

	"github.com/tomdoesdev/knox/internal/v1/project"
	"github.com/tomdoesdev/knox/internal/v1/runner"
	secrets2 "github.com/tomdoesdev/knox/internal/v1/secrets"
	"github.com/tomdoesdev/knox/internal/v1/template"
	"github.com/tomdoesdev/knox/pkg/errs"
	"github.com/urfave/cli/v3"
)

func NewRunCommand() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "run application with environment variables from template file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "env",
				Aliases: []string{"e"},
				Usage:   "path to environment template file",
				Value:   ".env.template",
			},
			&cli.DurationFlag{
				Name:  "timeout",
				Usage: "timeout for command execution",
				Value: 0,
			},
			&cli.BoolFlag{
				Name:  "allow-override",
				Usage: "allow template variables to override existing environment variables (only used with --inherit-env)",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "inherit-env",
				Usage: "inherit current environment variables (default: clean environment with only template variables)",
				Value: false,
			},
		},
		Action: runCommand,
	}
}

func runCommand(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() == 0 {
		return errs.Wrap(ErrInvalidArguments, InvalidArguments, "no command specified")
	}

	proj, err := project.Load()
	if err != nil {
		return err
	}

	workspace := proj.Workspace()
	secretStore, err := secrets2.NewFileSecretStore(workspace.VaultFilePath, workspace.ProjectID, secrets2.NewNoOpEncryptionHandler())
	if err != nil {
		return err
	}
	defer secretStore.Close()

	processor := template.NewProcessor(secretStore)
	envVars, err := processor.ProcessFileOptional(cmd.String("env"))
	if err != nil {
		return err
	}

	runnerConfig := runner.Config{
		Command:       cmd.Args().Slice(),
		Timeout:       cmd.Duration("timeout"),
		AllowOverride: cmd.Bool("allow-override"),
		InheritEnv:    cmd.Bool("inherit-env"),
	}

	runner := runner.NewEnvRunner(runnerConfig)
	return runner.Run(envVars)
}
