package errs

import (
	"fmt"
	"testing"
)

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

func TestKnoxErr_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *KnoxError
		expected string
	}{
		{
			name: "error without cause",
			err: &KnoxError{
				Code:    ProjectExistsCode,
				Message: "project already exists",
			},
			expected: "[PROJECT_EXISTS] project already exists",
		},
		{
			name: "error with cause",
			err: &KnoxError{
				Code:    FileNotFoundCode,
				Message: "config file missing",
				Cause:   fmt.Errorf("no such file or directory"),
			},
			expected: "[FILE_NOT_FOUND] config file missing: no such file or directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("KnoxError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestKnoxErr_Is(t *testing.T) {
	err1 := &KnoxError{Code: ProjectExistsCode, Message: "test"}
	err2 := &KnoxError{Code: ProjectExistsCode, Message: "different message"}
	err3 := &KnoxError{Code: SecretNotFoundCode, Message: "test"}

	if !err1.Is(err2) {
		t.Error("Errors with same code should be equal")
	}

	if err1.Is(err3) {
		t.Error("Errors with different codes should not be equal")
	}
}

func TestKnoxErr_WithContext(t *testing.T) {
	err := New(VaultCorruptedCode, "test").
		WithContext("key1", "value1").
		WithContext("key2", 42)

	if err.Context["key1"] != "value1" {
		t.Error("String context not preserved")
	}

	if err.Context["key2"] != 42 {
		t.Error("Integer context not preserved")
	}
}

func TestIs(t *testing.T) {
	knoxErr := New(ProjectExistsCode, "test")
	systemErr := fmt.Errorf("system error")

	if !Is(knoxErr, ProjectExistsCode) {
		t.Error("Is() should return true for matching Knox error")
	}

	if Is(systemErr, ProjectExistsCode) {
		t.Error("Is() should return false for non-Knox error")
	}

	if Is(knoxErr, SecretNotFoundCode) {
		t.Error("Is() should return false for different Knox error code")
	}
}

func TestCode(t *testing.T) {
	knoxErr := New(ProjectExistsCode, "test")
	systemErr := fmt.Errorf("system error")

	if Code(knoxErr) != ProjectExistsCode {
		t.Error("Code() should return correct code for Knox error")
	}

	if Code(systemErr) != "" {
		t.Error("Code() should return empty string for non-Knox error")
	}
}
