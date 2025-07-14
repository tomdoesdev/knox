package internal

import (
	"fmt"

	"github.com/tomdoesdev/knox/internal/config"
	"github.com/tomdoesdev/knox/internal/project"
	"github.com/tomdoesdev/knox/internal/vault"
)

func Initialize(conf *config.ApplicationConfig) error {

	err := vault.EnsureVaultExists(conf)
	if err != nil {
		return fmt.Errorf("initialize: vault.EnsureVaultExists: %w", err)
	}

	_, err = project.Create()
	if err != nil {
		return fmt.Errorf("knox.init.createProjectFile: %w", err)
	}

	return nil
}
