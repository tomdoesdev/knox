package vault

import (
	"github.com/tomdoesdev/knox/internal/error_codes"
	"github.com/tomdoesdev/knox/kit/errs"
)

var (
	ErrVaultConnectionFailed = errs.New(error_codes.VaultConnectionErrCode, "failed to connect to vault")
	ErrVaultIntegrityCheck   = errs.New(error_codes.VaultIntegrityErrCode, "vault integrity check failed")
	ErrVaultCreationFailed   = errs.New(error_codes.VaultCreationErrCode, "failed to create vault")
	ErrDatasourcePathInvalid = errs.New(error_codes.DatasourceErrCode, "invalid datasource path")
	ErrDatasourceUnreachable = errs.New(error_codes.DatasourceErrCode, "datasource database unreachable")
)
