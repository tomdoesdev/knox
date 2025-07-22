package errs_test

import (
	"fmt"
	"testing"

	"github.com/tomdoesdev/knox/kit/errs"
)

const (
	DatabaseError errs.Code = "DATABASE_ERROR"
	FileNotFound  errs.Code = "FILE_NOT_FOUND"
	NetworkError  errs.Code = "NETWORK_ERROR"
)

// Example of how to use structured errors in production code
func Example_usage() {
	// Creating a new error with a specific code
	err1 := errs.New(DatabaseError, "connection failed to database 'users'")
	fmt.Println(err1.Error())

	// Wrapping an existing error
	originalErr := fmt.Errorf("file not found")
	err2 := errs.Wrap(originalErr, FileNotFound, "config file missing")
	fmt.Println(err2.Error())

	// Using the builder pattern
	err3 := errs.Build(NetworkError).
		Message("failed to connect to %s", "api.example.com").
		Path("/api/users").
		Attempt(3).
		Error()
	fmt.Println(err3.Error())

	// Checking error codes
	if errs.Is(err1, DatabaseError) {
		fmt.Println("Database error detected!")
	}

	// Output:
	// [DATABASE_ERROR] connection failed to database 'users'
	// [FILE_NOT_FOUND] config file missing: file not found
	// [NETWORK_ERROR] failed to connect to api.example.com | context: path=/api/users, attempt=3
	// Database error detected!
}

// Example of how to use structured errors in tests
func TestErrorUsageExample(t *testing.T) {
	// Simulate a function that returns a structured error
	mockFunction := func() error {
		return errs.New(DatabaseError, "connection timeout")
	}

	err := mockFunction()

	// Check error code
	if !errs.Is(err, DatabaseError) {
		t.Error("Expected DatabaseError")
	}

	// Check error code extraction
	if errs.GetCode(err) != DatabaseError {
		t.Errorf("Expected code %s, got %s", DatabaseError, errs.GetCode(err))
	}
}

// Example of error with context
func TestErrorContext(t *testing.T) {
	err := errs.New(NetworkError, "connection timeout").
		WithContext("host", "api.example.com").
		WithContext("port", 443).
		WithContext("timeout", "30s")

	var structuredErr *errs.Error
	if errs.AsError(err, &structuredErr) {
		if structuredErr.Context["host"] != "api.example.com" {
			t.Error("Context not preserved")
		}
	}
}

// Example showing error chaining
func TestErrorChaining(t *testing.T) {
	// Original system error
	systemErr := fmt.Errorf("permission denied")

	// Structured error wrapping system error
	structuredErr := errs.Wrap(systemErr, FileNotFound, "cannot access config")

	// Check that we can still access the original error
	if structuredErr.Unwrap().Error() != "permission denied" {
		t.Error("Original error not preserved")
	}

	// Check that the structured error code is correct
	if errs.GetCode(structuredErr) != FileNotFound {
		t.Errorf("Expected code %s, got %s", FileNotFound, errs.GetCode(structuredErr))
	}
}
