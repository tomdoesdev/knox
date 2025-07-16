package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/knox"
	"github.com/tomdoesdev/knox/internal/config"
	"github.com/tomdoesdev/knox/internal/vault"
	"github.com/tomdoesdev/knox/kit/log"
	"github.com/tomdoesdev/knox/pkg/project"
)

func main() {
	log.NewSlog("text")

	var appConfig = config.NewApplicationConfig()

	p, err := project.Load()
	if err != nil {
		slog.Error("knox.init.loadProject", "error", err)
		os.Exit(1)
	}

	if p.VaultPath != "" {
		err = vault.EnsureVaultExists(p.VaultPath)

	} else {
		err = vault.EnsureVaultExists(appConfig.VaultPath)
	}

	if err != nil {
		slog.Error("knox.init.ensureVaultExists", "error", err)
		os.Exit(1)
	}

	cmdRoot := knox.NewKnoxCommand(p)

	if err := cmdRoot.Run(context.Background(), os.Args); err != nil {
		if errors.Is(err, project.ErrProjectExists) {
			slog.Info("project already exists")
			os.Exit(0)
		}
		slog.Error(err.Error())
		os.Exit(1)
	}
}
