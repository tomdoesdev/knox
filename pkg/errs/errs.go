package errs

import (
	"errors"
	"fmt"
)

// ErrorCode represents different types of Knox errs
type ErrorCode string

const (
	ProjectExistsCode   ErrorCode = "PROJECT_EXISTS"
	ProjectNotFoundCode ErrorCode = "PROJECT_NOT_FOUND"
	ProjectInvalidCode  ErrorCode = "PROJECT_INVALID"

	SecretNotFoundCode ErrorCode = "SECRET_NOT_FOUND"
	SecretExistsCode   ErrorCode = "SECRET_EXISTS"
	SecretInvalidCode  ErrorCode = "SECRET_INVALID"

	VaultCorruptedCode  ErrorCode = "VAULT_CORRUPTED"
	VaultLockedCode     ErrorCode = "VAULT_LOCKED"
	VaultPermissionCode ErrorCode = "VAULT_PERMISSION"

	TemplateParseCode ErrorCode = "TEMPLATE_PARSE"
	TemplateExecCode  ErrorCode = "TEMPLATE_EXEC"

	ProcessFailedCode  ErrorCode = "PROCESS_FAILED"
	ProcessTimeoutCode ErrorCode = "PROCESS_TIMEOUT"

	FileNotFoundCode     ErrorCode = "FILE_NOT_FOUND"
	FilePermissionCode   ErrorCode = "FILE_PERMISSION"
	DirectoryInvalidCode ErrorCode = "DIRECTORY_INVALID"

	InternalCode   ErrorCode = "INTERNAL"
	ValidationCode ErrorCode = "VALIDATION"
)

var (
	ErrProjectExists   = New(ProjectExistsCode, "project already exists")
	ErrProjectNotFound = New(ProjectNotFoundCode, "project not found")
	ErrSecretNotFound  = New(SecretNotFoundCode, "secret not found")
	ErrVaultCorrupted  = New(VaultCorruptedCode, "vault is corrupted")
)

// KnoxError represents a structured error in Knox
type KnoxError struct {
	Code    ErrorCode
	Message string
	Cause   error
	Context map[string]interface{}
}

func New(code ErrorCode, message string) *KnoxError {
	return &KnoxError{
		Code:    code,
		Message: message,
	}
}

func (e *KnoxError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *KnoxError) Unwrap() error {
	return e.Cause
}

func (e *KnoxError) Is(target error) bool {
	var t *KnoxError
	if errors.As(target, &t) {
		return e.Code == t.Code
	}
	return false
}

func (e *KnoxError) WithContext(key string, value interface{}) *KnoxError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// Wrap wraps an existing error with Knox error context
func Wrap(err error, code ErrorCode, message string) *KnoxError {
	return &KnoxError{
		Code:    code,
		Message: message,
		Cause:   err,
	}
}

// Is checks if an error has a specific Knox error code
func Is(err error, code ErrorCode) bool {
	var knoxErr *KnoxError
	if errors.As(err, &knoxErr) {
		return knoxErr.Code == code
	}
	return false
}

// Code extracts the error code from an error
func Code(err error) ErrorCode {
	var knoxErr *KnoxError
	if errors.As(err, &knoxErr) {
		return knoxErr.Code
	}
	return ""
}
