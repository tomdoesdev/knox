package vault

import "github.com/tomdoesdev/knox/pkg/errs"

const (
	VaultConnectionCode errs.ErrorCode = "vault_connection_error"
	VaultIntegrityCode  errs.ErrorCode = "vault_integrity_error"
	VaultCreationCode   errs.ErrorCode = "vault_creation_error"
	DatasourceCode      errs.ErrorCode = "datasource_error"
)

var (
	ErrVaultConnectionFailed   = errs.New(VaultConnectionCode, "failed to connect to vault")
	ErrVaultIntegrityCheck     = errs.New(VaultIntegrityCode, "vault integrity check failed")
	ErrVaultCreationFailed     = errs.New(VaultCreationCode, "failed to create vault")
	ErrDatasourcePathInvalid   = errs.New(DatasourceCode, "invalid datasource path")
	ErrDirectoryCreationFailed = errs.New(VaultCreationCode, "failed to create vault directory")
	ErrDatasourceUnreachable   = errs.New(DatasourceCode, "datasource database unreachable")
)
