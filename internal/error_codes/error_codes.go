package error_codes

import "github.com/tomdoesdev/knox/kit/errs"

const (
	ProjectExistsErrCode   errs.Code = "PROJECT_EXISTS"
	ProjectNotFoundErrCode errs.Code = "PROJECT_NOT_FOUND"
	ProjectInvalidErrCode  errs.Code = "PROJECT_INVALID"

	SecretNotFoundErrCode errs.Code = "SECRET_NOT_FOUND"
	SecretExistsErrCode   errs.Code = "SECRET_EXISTS"
	SecretInvalidErrCode  errs.Code = "SECRET_INVALID"

	VaultCreationErrCode   errs.Code = "VAULT_CREATION"
	VaultConnectionErrCode errs.Code = "VAULT_CONNECTION"
	VaultIntegrityErrCode  errs.Code = "VAULT_INTEGRITY"

	FileNotFoundErrCode     errs.Code = "FILE_NOT_FOUND"
	FilePermissionErrCode   errs.Code = "FILE_PERMISSION"
	DirectoryInvalidErrCode errs.Code = "DIRECTORY_INVALID"

	ValidationErrCode errs.Code = "VALIDATION"

	SearchFailureErrCode   errs.Code = "SEARCH_FAILURE"
	CreateFailureErrCode   errs.Code = "CREATE_FAILURE"
	DatabaseFailureErrCode errs.Code = "DATABASE_FAILURE"

	DatasourceErrCode errs.Code = "DATASOURCE_ERROR"
)
