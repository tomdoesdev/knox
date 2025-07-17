package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands"
	"github.com/tomdoesdev/knox/internal"
	"github.com/tomdoesdev/knox/internal/project"
	"github.com/tomdoesdev/knox/internal/secrets"
	"github.com/tomdoesdev/knox/internal/secrets/store"
	"github.com/tomdoesdev/knox/kit/log"
)

func main() {
	log.NewSlog("text")

	p, err := project.Load()
	if err != nil {
		slog.Error("knox.init.loadProject", "error", err)
		os.Exit(1)
	}

	workspace := p.Workspace()

	slog.Debug("project loaded",
		slog.String("project.vaultpath", workspace.VaultFilePath),
	)

	e := secrets.NewNoOpEncryptionHandler()

	s, err := store.NewFileSecretStore(workspace.VaultFilePath, workspace.ProjectID, e)
	if err != nil {
		slog.Error("failed to create vault provider", "error", err)
		os.Exit(1)
	}

	k, err := internal.NewKnox(&internal.KnoxOptions{
		SecretStore: s,
		Workspace:   p.Workspace(),
	})
	if err != nil {
		slog.Error("failed to create Knox execution context", "error", err)
		os.Exit(1)
	}

	defer (func() {
		err := k.Close()
		if err != nil {
			slog.Error("error closing vault: %v", err)
		}
	})()

	cmdRoot := commands.NewKnoxCommand(k)

	if err := cmdRoot.Run(context.Background(), os.Args); err != nil {
		if errors.Is(err, project.ErrProjectExists) {
			slog.Info("project already exists")
			os.Exit(0)
		}

		slog.Error("knox.init.run", "error", err)
		os.Exit(1)
	}
}
