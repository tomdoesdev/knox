package errs

import (
	"fmt"
	"testing"
)

const (
	TestCode1 Code = "TEST_CODE_1"
	TestCode2 Code = "TEST_CODE_2"
	TestCode3 Code = "TEST_CODE_3"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected string
	}{
		{
			name: "error without cause",
			err: &Error{
				Code:    TestCode1,
				Message: "test error message",
			},
			expected: "[TEST_CODE_1] test error message",
		},
		{
			name: "error with cause",
			err: &Error{
				Code:    TestCode2,
				Message: "wrapped error",
				Cause:   fmt.Errorf("original error"),
			},
			expected: "[TEST_CODE_2] wrapped error: original error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestError_Is(t *testing.T) {
	err1 := &Error{Code: TestCode1, Message: "test"}
	err2 := &Error{Code: TestCode1, Message: "different message"}
	err3 := &Error{Code: TestCode2, Message: "test"}

	if !err1.Is(err2) {
		t.Error("Errors with same code should be equal")
	}

	if err1.Is(err3) {
		t.Error("Errors with different codes should not be equal")
	}
}

func TestError_WithContext(t *testing.T) {
	err := New(TestCode1, "test").
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
	structuredErr := New(TestCode1, "test")
	systemErr := fmt.Errorf("system error")

	if !Is(structuredErr, TestCode1) {
		t.Error("Is() should return true for matching structured error")
	}

	if Is(systemErr, TestCode1) {
		t.Error("Is() should return false for non-structured error")
	}

	if Is(structuredErr, TestCode2) {
		t.Error("Is() should return false for different error code")
	}
}

func TestGetCode(t *testing.T) {
	structuredErr := New(TestCode1, "test")
	systemErr := fmt.Errorf("system error")

	if GetCode(structuredErr) != TestCode1 {
		t.Error("GetCode() should return correct code for structured error")
	}

	if GetCode(systemErr) != "" {
		t.Error("GetCode() should return empty string for non-structured error")
	}
}

func TestWrapWithContext(t *testing.T) {
	originalErr := fmt.Errorf("original error")
	err := WrapWithContext(originalErr, TestCode2, "wrapped error",
		"path", "/test/path", "attempt", 3)

	if err.Code != TestCode2 {
		t.Errorf("Expected code %s, got %s", TestCode2, err.Code)
	}

	if err.Message != "wrapped error" {
		t.Errorf("Expected message 'wrapped error', got '%s'", err.Message)
	}

	if err.Cause != originalErr {
		t.Error("Cause should be preserved")
	}

	if err.Context["path"] != "/test/path" {
		t.Errorf("Expected path context '/test/path', got %v", err.Context["path"])
	}

	if err.Context["attempt"] != 3 {
		t.Errorf("Expected attempt context 3, got %v", err.Context["attempt"])
	}
}

func TestErrorBuilder(t *testing.T) {
	originalErr := fmt.Errorf("database error")

	err := Build(TestCode3).
		Message("failed to process at %s", "/test/location").
		Cause(originalErr).
		Path("/test/location").
		Attempt(2).
		Context("checksum", "abc123").
		Error()

	if err.Code != TestCode3 {
		t.Errorf("Expected code %s, got %s", TestCode3, err.Code)
	}

	if err.Message != "failed to process at /test/location" {
		t.Errorf("Expected formatted message, got '%s'", err.Message)
	}

	if err.Cause != originalErr {
		t.Error("Cause should be preserved")
	}

	if err.Context["path"] != "/test/location" {
		t.Errorf("Expected path context, got %v", err.Context["path"])
	}

	if err.Context["attempt"] != 2 {
		t.Errorf("Expected attempt context 2, got %v", err.Context["attempt"])
	}

	if err.Context["checksum"] != "abc123" {
		t.Errorf("Expected checksum context, got %v", err.Context["checksum"])
	}
}

func TestConvenienceMethods(t *testing.T) {
	err := New(TestCode1, "test error").
		WithPath("/test/path").
		WithFile("test.txt").
		WithAttempt(5).
		WithDatabase("test.db").
		WithOperation("read").
		WithID("test123")

	expected := map[string]interface{}{
		"path":      "/test/path",
		"file":      "test.txt",
		"attempt":   5,
		"database":  "test.db",
		"operation": "read",
		"id":        "test123",
	}

	for key, expectedValue := range expected {
		if err.Context[key] != expectedValue {
			t.Errorf("Expected %s context %v, got %v", key, expectedValue, err.Context[key])
		}
	}
}
