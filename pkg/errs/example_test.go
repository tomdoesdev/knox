package errs_test

import (
	"fmt"
	"testing"

	"github.com/tomdoesdev/knox/pkg/errs"
)

// Example of how to use Knox errs in production code
func Example_usage() {
	// Creating a new error with a specific code
	err1 := errs.New(errs.ProjectExistsCode, "project 'myapp' already exists")
	fmt.Println(err1.Error())

	// Creating an error with the default message
	err2 := errs.NewWithDefaultMessage(errs.SecretNotFoundCode)
	fmt.Println(err2.Error())

	// Wrapping an existing error
	originalErr := fmt.Errorf("file not found")
	err3 := errs.Wrap(originalErr, errs.FileNotFoundCode, "config file missing")
	fmt.Println(err3.Error())

	// Checking error codes
	if errs.Is(err1, errs.ProjectExistsCode) {
		fmt.Println("Project already exists!")
	}

	// Output:
	// [PROJECT_EXISTS] project 'myapp' already exists
	// [SECRET_NOT_FOUND] secret not found
	// [FILE_NOT_FOUND] config file missing: file not found
	// Project already exists!
}

// Example of how to use Knox errs in tests
func TestErrorUsageExample(t *testing.T) {
	// Simulate a function that returns a Knox error
	mockFunction := func() error {
		return errs.NewWithDefaultMessage(errs.ProjectExistsCode)
	}

	e := mockFunction()

	// New way (robust - error code checking)
	errs.AssertErrorCode(t, e, errs.ProjectExistsCode)

	// You can also check multiple properties
	errs.AssertKnoxError(t, e, errs.ProjectExistsCode, "project already exists")
}

// Example of error with context
func TestErrorContext(t *testing.T) {
	err := errs.New(errs.VaultCorruptedCode, "vault checksum mismatch").
		WithContext("vault_path", "/path/to/vault").
		WithContext("expected_checksum", "abc123").
		WithContext("actual_checksum", "def456")

	var knoxErr *errs.KnoxError
	if errs.AsKnoxError(err, &knoxErr) {
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
	knoxErr := errs.Wrap(systemErr, errs.VaultPermissionCode, "cannot access vault")

	// Check that we can still access the original error
	if knoxErr.Unwrap().Error() != "permission denied" {
		t.Error("Original error not preserved")
	}

	// Check that the Knox error code is correct
	errs.AssertErrorCode(t, knoxErr, errs.VaultPermissionCode)
}
