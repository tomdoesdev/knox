package errors

import (
	"errors"
	"testing"
)

// AssertErrorCode checks if an error has the expected Knox error code
func AssertErrorCode(t *testing.T, err error, expectedCode ErrorCode) {
	t.Helper()

	if err == nil {
		t.Fatalf("Expected error with code %s, got nil", expectedCode)
	}

	if !Is(err, expectedCode) {
		actualCode := Code(err)
		t.Errorf("Expected error code %s, got %s (error: %v)", expectedCode, actualCode, err)
	}
}

// AssertNoError checks that no error occurred
func AssertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

// AssertErrorContains checks if an error contains specific text
func AssertErrorContains(t *testing.T, err error, expectedText string) {
	t.Helper()

	if err == nil {
		t.Fatalf("Expected error containing '%s', got nil", expectedText)
	}

	if !contains(err.Error(), expectedText) {
		t.Errorf("Expected error to contain '%s', got: %v", expectedText, err)
	}
}

// AssertKnoxError checks if an error is a Knox error with specific properties
func AssertKnoxError(t *testing.T, err error, expectedCode ErrorCode, expectedMessage string) {
	t.Helper()

	var knoxErr *KnoxError
	if !AsKnoxError(err, &knoxErr) {
		t.Fatalf("Expected Knox error, got: %T (%v)", err, err)
	}

	if knoxErr.Code != expectedCode {
		t.Errorf("Expected error code %s, got %s", expectedCode, knoxErr.Code)
	}

	if expectedMessage != "" && knoxErr.Message != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, knoxErr.Message)
	}
}

// AsKnoxError is a helper that checks if an error is a Knox error
func AsKnoxError(err error, target **KnoxError) bool {
	return errors.As(err, target)
}

// contains is a simple string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (len(substr) == 0 || indexOf(s, substr) >= 0)
}

// indexOf returns the index of substr in s, or -1 if not found
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
