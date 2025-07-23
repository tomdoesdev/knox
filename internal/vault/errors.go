package vault

import (
	"github.com/tomdoesdev/knox/internal"
	"github.com/tomdoesdev/knox/kit/errs"
)

var (
	ErrVaultConnectionFailed   = errs.New(internal.VaultConnectionCode, "failed to connect to vault")
	ErrVaultIntegrityCheck     = errs.New(internal.VaultIntegrityCode, "vault integrity check failed")
	ErrVaultCreationFailed     = errs.New(internal.VaultCreationCode, "failed to create vault")
	ErrDatasourcePathInvalid   = errs.New(internal.DatasourceCode, "invalid datasource path")
	ErrDirectoryCreationFailed = errs.New(internal.VaultCreationCode, "failed to create vault directory")
	ErrDatasourceUnreachable   = errs.New(internal.DatasourceCode, "datasource database unreachable")
)
