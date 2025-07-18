package secrets

import (
	"path/filepath"

	"github.com/tomdoesdev/knox/kit/fs"
	"github.com/tomdoesdev/knox/kit/vcs"
	"github.com/tomdoesdev/knox/pkg/errs"
)

const (
	StoreOpenFailureCode errs.ErrorCode = "STORE_FAILURE"
	StoreExecFailureCode errs.ErrorCode = "STORE_EXEC_FAILURE"
	StoreVCSWarningCode  errs.ErrorCode = "STORE_VCS_WARNING"
)

var (
	ErrVaultInVCS = errs.New(StoreVCSWarningCode, "vault file would be created in a version control directory")
)

func NewFileSecretStore(source fs.FilePath, projectId string, handler EncryptionHandler) (*SqliteSecretStore, error) {
	return NewFileSecretStoreWithOptions(source, projectId, handler, false)
}

func NewFileSecretStoreWithOptions(source fs.FilePath, projectId string, handler EncryptionHandler, force bool) (*SqliteSecretStore, error) {
	// Check if vault file would be created in a VCS directory
	if err := checkVCSProtection(source, force); err != nil {
		return nil, err
	}

	s, err := newSqliteStore(source, projectId, handler)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func checkVCSProtection(vaultPath fs.FilePath, force bool) error {
	// Skip check if force flag is set
	if force {
		return nil
	}

	// Check if file already exists - if so, don't block opening it
	if exists, err := fs.IsExist(vaultPath); err == nil && exists {
		return nil
	}

	// Get the directory that would contain the vault file
	vaultDir := filepath.Dir(vaultPath)

	// Check if this directory is under version control
	if vcs.IsUnderVCS(vaultDir) {
		return ErrVaultInVCS.WithContext("vault_path", vaultPath).WithContext("suggestion", "use --force to override or specify a vault_path outside of version control")
	}

	return nil
}
