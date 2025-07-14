package vault

import (
	"fmt"
	"os"

	"github.com/tomdoesdev/knox/internal/config"
	"github.com/tomdoesdev/knox/kit/fskit"
)

type Vault interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Del(key string) error
	Clear() error
}

func EnsureVaultExists(conf *config.ApplicationConfig) error {
	dirExists, err := fskit.Exists(conf.VaultDir)
	if err != nil {
		return fmt.Errorf("vault:create:dirExists: %w", err)
	}

	vaultExists, err := fskit.Exists(conf.VaultPath)

	if err != nil {
		return fmt.Errorf("vault:create:vaultExists: %w", err)
	}

	if !dirExists {
		err = os.MkdirAll(conf.VaultDir, 0600)
		if err != nil {
			return fmt.Errorf("vault:create:mkDir: %w", err)
		}
	}

	if !vaultExists {
		err = createSqliteStore(conf)
		if err != nil {
			return fmt.Errorf("vault:create:createSqliteStore: %w", err)
		}
	}

	return nil
}
