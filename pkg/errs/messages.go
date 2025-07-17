package errs

// ErrorMessages contains all user-facing error messages
// This allows for easy localization and consistent messaging
var ErrorMessages = map[ErrorCode]string{
	// Project messages
	ProjectExistsCode:   "project already exists",
	ProjectNotFoundCode: "project not found",
	ProjectInvalidCode:  "project configuration is invalid",

	// Secret messages
	SecretNotFoundCode: "secret not found",
	SecretExistsCode:   "secret already exists",
	SecretInvalidCode:  "secret name or value is invalid",

	// Vault messages
	VaultCorruptedCode:  "vault is corrupted or unreadable",
	VaultLockedCode:     "vault is locked",
	VaultPermissionCode: "insufficient permissions to access vault",

	// Template messages
	TemplateParseCode: "template parsing failed",
	TemplateExecCode:  "template execution failed",

	// Process messages
	ProcessFailedCode:  "process execution failed",
	ProcessTimeoutCode: "process execution timed out",

	// File system messages
	FileNotFoundCode:     "file not found",
	FilePermissionCode:   "insufficient file permissions",
	DirectoryInvalidCode: "directory is invalid or inaccessible",

	// Internal messages
	InternalCode:   "internal error occurred",
	ValidationCode: "validation failed",
}

// GetMessage returns the default message for an error code
func GetMessage(code ErrorCode) string {
	if msg, exists := ErrorMessages[code]; exists {
		return msg
	}
	return "unknown error"
}

// NewWithDefaultMessage creates an error with the default message for the code
func NewWithDefaultMessage(code ErrorCode) *KnoxError {
	return &KnoxError{
		Code:    code,
		Message: GetMessage(code),
	}
}

// WrapWithDefaultMessage wraps an error with the default message for the code
func WrapWithDefaultMessage(err error, code ErrorCode) *KnoxError {
	return &KnoxError{
		Code:    code,
		Message: GetMessage(code),
		Cause:   err,
	}
}
