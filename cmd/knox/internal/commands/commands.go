package commands

import (
	"log/slog"

	"github.com/tomdoesdev/knox/internal"
	"github.com/tomdoesdev/knox/internal/project"
	"github.com/tomdoesdev/knox/internal/secrets"
)

func LoadKnoxContextWithOptions(force bool) (*internal.Knox, error) {
	p, err := project.Load()
	if err != nil {
		slog.Error("commands.LoadKnoxContext", "error", err)
		return nil, err
	}

	workspace := p.Workspace()

	slog.Debug("project loaded",
		slog.String("project.vaultpath", workspace.VaultFilePath),
	)

	e := secrets.NewNoOpEncryptionHandler()

	s, err := secrets.NewFileSecretStoreWithOptions(workspace.VaultFilePath, workspace.ProjectID, e, force)
	if err != nil {
		return nil, err
	}

	k, err := internal.NewKnox(&internal.KnoxOptions{
		SecretStore: s,
		Workspace:   p.Workspace(),
	})
	if err != nil {
		return nil, err
	}

	return k, nil
}
