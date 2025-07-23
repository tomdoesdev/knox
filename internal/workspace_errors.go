package internal

import "github.com/tomdoesdev/knox/kit/errs"

const (
	ProjectExistsCode   errs.Code = "PROJECT_EXISTS"
	ProjectNotFoundCode errs.Code = "PROJECT_NOT_FOUND"
	ProjectInvalidCode  errs.Code = "PROJECT_INVALID"

	SecretNotFoundCode errs.Code = "SECRET_NOT_FOUND"
	SecretExistsCode   errs.Code = "SECRET_EXISTS"
	SecretInvalidCode  errs.Code = "SECRET_INVALID"

	VaultCorruptedCode  errs.Code = "VAULT_CORRUPTED"
	VaultLockedCode     errs.Code = "VAULT_LOCKED"
	VaultPermissionCode errs.Code = "VAULT_PERMISSION"
	VaultCreationCode   errs.Code = "VAULT_CREATION"
	VaultConnectionCode errs.Code = "VAULT_CONNECTION"
	VaultIntegrityCode  errs.Code = "VAULT_INTEGRITY"

	TemplateParseCode errs.Code = "TEMPLATE_PARSE"
	TemplateExecCode  errs.Code = "TEMPLATE_EXEC"

	ProcessFailedCode  errs.Code = "PROCESS_FAILED"
	ProcessTimeoutCode errs.Code = "PROCESS_TIMEOUT"

	FileNotFoundCode     errs.Code = "FILE_NOT_FOUND"
	FilePermissionCode   errs.Code = "FILE_PERMISSION"
	DirectoryInvalidCode errs.Code = "DIRECTORY_INVALID"

	InternalCode   errs.Code = "INTERNAL"
	ValidationCode errs.Code = "VALIDATION"

	// Workspace specific codes
	SearchFailureCode       errs.Code = "SEARCH_FAILURE"
	CreateFailureCode       errs.Code = "CREATE_FAILURE"
	DatabaseFailureCode     errs.Code = "DATABASE_FAILURE"
	DatabasePathInvalidCode errs.Code = "DATABASE_PATH_INVALID"

	DatasourceCode errs.Code = "DATASOURCE_ERROR"

	InvalidIdentifierCode errs.Code = "INVALID_IDENTIFIER"
)

// NewWithDefaultMessage creates a new error with a default message for Knox error codes
func NewWithDefaultMessage(code errs.Code) *errs.Error {
	message := defaultMessages[code]
	if message == "" {
		message = string(code)
	}
	return errs.New(code, message)
}

// defaultMessages maps Knox error codes to their default messages
var defaultMessages = map[errs.Code]string{
	ProjectExistsCode:       "project already exists",
	ProjectNotFoundCode:     "project not found",
	ProjectInvalidCode:      "project invalid",
	SecretNotFoundCode:      "secret not found",
	SecretExistsCode:        "secret already exists",
	SecretInvalidCode:       "secret invalid",
	VaultCorruptedCode:      "vault corrupted",
	VaultLockedCode:         "vault locked",
	VaultPermissionCode:     "vault permission denied",
	VaultCreationCode:       "vault creation failed",
	VaultConnectionCode:     "vault connection failed",
	VaultIntegrityCode:      "vault integrity check failed",
	TemplateParseCode:       "template parse error",
	TemplateExecCode:        "template execution error",
	ProcessFailedCode:       "process failed",
	ProcessTimeoutCode:      "process timeout",
	FileNotFoundCode:        "file not found",
	FilePermissionCode:      "file permission denied",
	DirectoryInvalidCode:    "directory invalid",
	InternalCode:            "internal error",
	ValidationCode:          "validation error",
	SearchFailureCode:       "search failure",
	CreateFailureCode:       "create failure",
	DatabaseFailureCode:     "database failure",
	DatabasePathInvalidCode: "database path invalid",
	DatasourceCode:          "datasource error",
	InvalidIdentifierCode:   "invalid identifier",
}
