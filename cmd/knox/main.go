package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/tomdoesdev/knox/cmd/knox/internal"
	"github.com/tomdoesdev/knox/internal/config"
	"github.com/tomdoesdev/knox/internal/project"
	"github.com/tomdoesdev/knox/internal/vault"
)

func main() {
	var logger *slog.Logger

	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	default:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	slog.SetDefault(logger)

	var appConfig = config.NewApplicationConfig()

	p, err := project.LoadCwd()
	if err != nil {
		slog.Error("knox.init.loadProject", "error", err)
		os.Exit(1)
	}

	defer func() {
		err := p.Close()
		if err != nil {
			slog.Error("error closing project", "err", fmt.Errorf("project.Close: %w", err))
		}
	}()

	if p.VaultPath != "" {
		err = vault.EnsureVaultExists(p.VaultPath)

	} else {
		err = vault.EnsureVaultExists(appConfig.VaultPath)
	}

	if err != nil {
		slog.Error("knox.init.ensureVaultExists", "error", err)
		os.Exit(1)
	}

	cmdRoot := internal.NewKnoxCommand(p)

	if err := cmdRoot.Run(context.Background(), os.Args); err != nil {
		if errors.Is(err, project.ErrProjectExists) {
			logger.Info("project already exists")
			os.Exit(0)
		}
		logger.Error(err.Error())
		os.Exit(1)
	}
}
