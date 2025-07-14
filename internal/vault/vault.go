package vault

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/tomdoesdev/knox/internal/config"
	"github.com/tomdoesdev/knox/kit/errkit"
	"github.com/tomdoesdev/knox/kit/fskit"
)

type Options struct {
	Path string
	Name string
}
type Vault interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Del(key string) error
	Clear() error
}

func New(opts *Options) Vault {
	//storePath := path.Join(opts.Path, "knox.db")
	//
	//if fskit.Exists(storePath) == false {
	//
	//}
	return nil
}

type vaultImpl struct{}

func (s *vaultImpl) Set(key, value string) error {
	return errkit.ErrNotImplemented
}
func (s *vaultImpl) Get(key string) (string, error) {
	return "", errkit.ErrNotImplemented
}

func (s *vaultImpl) Del(key string) error {
	return errkit.ErrNotImplemented
}
func (s *vaultImpl) Clear() error {
	return errkit.ErrNotImplemented
}

func Create(conf *config.ApplicationConfig) error {
	exists, err := fskit.Exists(conf.VaultDir)
	if err != nil {
		return fmt.Errorf("knox.init.createKnoxVault: %w", err)
	}

	if exists {
		slog.Warn("user vault already exists")
		return nil
	}

	err = os.MkdirAll(conf.VaultDir, 0600)
	if err != nil {
		return fmt.Errorf("knox.init.createKnoxVault.mkDir: %w", err)
	}

	err = initSqliteVault(conf)
	if err != nil {
		return fmt.Errorf("knox.init.createKnoxVault: %w", err)
	}

	return nil
}
