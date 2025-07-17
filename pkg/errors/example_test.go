package errors_test

import (
	"fmt"
	"testing"

	errors2 "github.com/tomdoesdev/knox/pkg/errors"
)

// Example of how to use Knox errors in production code
func Example_usage() {
	// Creating a new error with a specific code
	err1 := errors2.New(errors2.ProjectExistsCode, "project 'myapp' already exists")
	fmt.Println(err1.Error())

	// Creating an error with the default message
	err2 := errors2.NewWithDefaultMessage(errors2.SecretNotFoundCode)
	fmt.Println(err2.Error())

	// Wrapping an existing error
	originalErr := fmt.Errorf("file not found")
	err3 := errors2.Wrap(originalErr, errors2.FileNotFoundCode, "config file missing")
	fmt.Println(err3.Error())

	// Checking error codes
	if errors2.Is(err1, errors2.ProjectExistsCode) {
		fmt.Println("Project already exists!")
	}

	// Output:
	// [PROJECT_EXISTS] project 'myapp' already exists
	// [SECRET_NOT_FOUND] secret not found
	// [FILE_NOT_FOUND] config file missing: file not found
	// Project already exists!
}

// Example of how to use Knox errors in tests
func TestErrorUsageExample(t *testing.T) {
	// Simulate a function that returns a Knox error
	mockFunction := func() error {
		return errors2.NewWithDefaultMessage(errors2.ProjectExistsCode)
	}

	err := mockFunction()

	// Old way (brittle - string comparison)
	// if err.Error() != "[PROJECT_EXISTS] project already exists" {
	//     t.Error("wrong error message")
	// }

	// New way (robust - error code checking)
	errors2.AssertErrorCode(t, err, errors2.ProjectExistsCode)

	// You can also check multiple properties
	errors2.AssertKnoxError(t, err, errors2.ProjectExistsCode, "project already exists")
}

// Example of error with context
func TestErrorContext(t *testing.T) {
	err := errors2.New(errors2.VaultCorruptedCode, "vault checksum mismatch").
		WithContext("vault_path", "/path/to/vault").
		WithContext("expected_checksum", "abc123").
		WithContext("actual_checksum", "def456")

	var knoxErr *errors2.KnoxError
	if errors2.AsKnoxError(err, &knoxErr) {
		if knoxErr.Context["vault_path"] != "/path/to/vault" {
			t.Error("Context not preserved")
		}
	}
}

// Example showing error chaining
func TestErrorChaining(t *testing.T) {
	// Original system error
	systemErr := fmt.Errorf("permission denied")

	// Knox error wrapping system error
	knoxErr := errors2.Wrap(systemErr, errors2.VaultPermissionCode, "cannot access vault")

	// Check that we can still access the original error
	if knoxErr.Unwrap().Error() != "permission denied" {
		t.Error("Original error not preserved")
	}

	// Check that the Knox error code is correct
	errors2.AssertErrorCode(t, knoxErr, errors2.VaultPermissionCode)
}
