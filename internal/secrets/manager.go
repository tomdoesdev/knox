package secrets

import "github.com/tomdoesdev/knox/internal/errors"

type VaultManager struct {
}
type VaultOptions struct {
	Path string
	Name string
}

func NewVaultManager(path string) VaultManager {
	return VaultManager{}
}

func (v *VaultManager) NewVault() error {
	return errors.ErrNotImplemented
}

type VaultInitializer interface {
	NewVault(opts *VaultOptions) error
}

func NewVaultInitializer() VaultInitializer {
	return NewSqliteStore()
}
