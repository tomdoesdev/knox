package commands

import (
	"log/slog"

	"github.com/tomdoesdev/knox/internal/v1"
	"github.com/tomdoesdev/knox/internal/v1/project"
	secrets2 "github.com/tomdoesdev/knox/internal/v1/secrets"
)

func LoadKnoxContextWithOptions(force bool) (*v1.Knox, error) {
	p, err := project.Load()
	if err != nil {
		slog.Error("commands.LoadKnoxContext", "error", err)
		return nil, err
	}

	workspace := p.Workspace()

	slog.Debug("project loaded",
		slog.String("project.vaultpath", workspace.VaultFilePath),
	)

	e := secrets2.NewNoOpEncryptionHandler()

	s, err := secrets2.NewFileSecretStoreWithOptions(workspace.VaultFilePath, workspace.ProjectID, e, force)
	if err != nil {
		return nil, err
	}

	k, err := v1.NewKnox(&v1.KnoxOptions{
		SecretStore: s,
		Workspace:   p.Workspace(),
	})
	if err != nil {
		return nil, err
	}

	return k, nil
}
