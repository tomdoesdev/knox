package vault

import (
	"github.com/tomdoesdev/knox/kit/errs"
	"github.com/tomdoesdev/knox/internal/errors"
)

var (
	ErrVaultConnectionFailed   = errs.New(errors.VaultConnectionCode, "failed to connect to vault")
	ErrVaultIntegrityCheck     = errs.New(errors.VaultIntegrityCode, "vault integrity check failed")
	ErrVaultCreationFailed     = errs.New(errors.VaultCreationCode, "failed to create vault")
	ErrDatasourcePathInvalid   = errs.New(errors.DatasourceCode, "invalid datasource path")
	ErrDirectoryCreationFailed = errs.New(errors.VaultCreationCode, "failed to create vault directory")
	ErrDatasourceUnreachable   = errs.New(errors.DatasourceCode, "datasource database unreachable")
)
