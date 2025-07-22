package errs

import (
	"fmt"
	"testing"
)

func TestKnoxErrorWrapperWorks(t *testing.T) {
	// Test that our wrapper re-exports work correctly
	err := New(ProjectExistsCode, "test project exists")

	if err.Code != ProjectExistsCode {
		t.Errorf("Expected code %s, got %s", ProjectExistsCode, err.Code)
	}

	if err.Message != "test project exists" {
		t.Errorf("Expected message 'test project exists', got '%s'", err.Message)
	}
}

func TestNewWithDefaultMessageWrapper(t *testing.T) {
	err := NewWithDefaultMessage(SecretNotFoundCode)

	if err.Code != SecretNotFoundCode {
		t.Errorf("Expected code %s, got %s", SecretNotFoundCode, err.Code)
	}

	if err.Message != "secret not found" {
		t.Errorf("Expected default message 'secret not found', got '%s'", err.Message)
	}
}

func TestWrapperFunctions(t *testing.T) {
	originalErr := fmt.Errorf("original error")

	// Test Wrap function
	wrappedErr := Wrap(originalErr, VaultCreationCode, "vault creation failed")

	if !Is(wrappedErr, VaultCreationCode) {
		t.Error("Is() function not working through wrapper")
	}

	if Code(wrappedErr) != VaultCreationCode {
		t.Error("Code() function not working through wrapper")
	}
}
