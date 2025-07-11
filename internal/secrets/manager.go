package secrets

import "github.com/tomdoesdev/knox/internal/errors"

type VaultManager struct {
}

func NewVaultManager(path string) VaultManager {
	return VaultManager{}
}

func (v *VaultManager) Init() error {
	return errors.ErrNotImplemented
}
