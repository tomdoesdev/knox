# Structured Errors Package

A generic, reusable Go package for structured error handling with codes, messages, context, and error chaining.

## Features

- **Structured Errors**: Errors with codes, messages, underlying causes, and contextual information
- **Error Codes**: Type-safe error codes for programmatic error handling
- **Context Support**: Add arbitrary key-value context to errors for debugging
- **Error Chaining**: Wrap and unwrap errors while preserving the original cause
- **Builder Pattern**: Fluent interface for building complex errors
- **Testing Utilities**: Helper functions for testing error conditions
- **Generic Design**: No application-specific dependencies, suitable for any Go project

## Basic Usage

```go
package main

import (
    "fmt"
    "github.com/tomdoesdev/knox/kit/errs"
)

const (
    DatabaseError errs.Code = "DATABASE_ERROR"
    FileNotFound  errs.Code = "FILE_NOT_FOUND"
)

func main() {
    // Create a simple error
    err1 := errs.New(DatabaseError, "connection failed")
    fmt.Println(err1) // [DATABASE_ERROR] connection failed

    // Wrap an existing error
    originalErr := fmt.Errorf("file not found")
    err2 := errs.Wrap(originalErr, FileNotFound, "config file missing")
    fmt.Println(err2) // [FILE_NOT_FOUND] config file missing: file not found

    // Add context to an error
    err3 := errs.New(DatabaseError, "query failed").
        WithContext("table", "users").
        WithContext("query", "SELECT * FROM users")
    fmt.Println(err3) // [DATABASE_ERROR] query failed | context: table=users, query=SELECT * FROM users

    // Check error codes
    if errs.Is(err1, DatabaseError) {
        fmt.Println("Database error detected")
    }
}
```

## Advanced Usage

### Builder Pattern

```go
err := errs.Build(NetworkError).
    Message("failed to connect to %s", "api.example.com").
    Cause(originalErr).
    Path("/api/users").
    Attempt(3).
    Context("timeout", "30s").
    Error()
```

### Convenience Methods

```go
err := errs.New(FileNotFound, "config missing").
    WithPath("/etc/config").
    WithFile("app.yaml").
    WithAttempt(2).
    WithDatabase("main.db").
    WithOperation("read").
    WithID("user123")
```

### WrapWithContext

```go
err := errs.WrapWithContext(originalErr, ProcessFailed, "command execution failed",
    "command", "git clone",
    "directory", "/tmp/repo",
    "exit_code", 1)
```

## Testing

```go
func TestMyFunction(t *testing.T) {
    err := myFunction()
    
    // Assert specific error code
    errs.AssertErrorCode(t, err, DatabaseError)
    
    // Assert no error occurred
    errs.AssertNoError(t, err)
    
    // Assert error contains text
    errs.AssertErrorContains(t, err, "connection failed")
    
    // Assert structured error properties
    errs.AssertError(t, err, DatabaseError, "connection failed")
}
```

## Error Handling Patterns

### Checking Error Types

```go
// Check if error has specific code
if errs.Is(err, DatabaseError) {
    // Handle database error
}

// Extract error code
code := errs.GetCode(err)

// Check if it's a structured error
var structuredErr *errs.Error
if errs.AsError(err, &structuredErr) {
    fmt.Printf("Context: %+v\n", structuredErr.Context)
}
```

### Error Propagation

```go
func processFile(path string) error {
    data, err := readFile(path)
    if err != nil {
        return errs.Wrap(err, FileNotFound, "failed to read config").
            WithPath(path).
            WithOperation("read")
    }
    
    if err := validateData(data); err != nil {
        return errs.Wrap(err, ValidationError, "invalid config format").
            WithPath(path).
            WithFile(filepath.Base(path))
    }
    
    return nil
}
```

## Types

- `errs.Code`: String-based error code type
- `errs.Error`: Main structured error type
- `errs.ErrorBuilder`: Builder for complex error construction

## Functions

- `New(code, message)`: Create new structured error
- `Wrap(err, code, message)`: Wrap existing error
- `WrapWithContext(err, code, message, pairs...)`: Wrap with inline context
- `Build(code)`: Create error builder
- `Is(err, code)`: Check if error has specific code
- `GetCode(err)`: Extract error code
- `AsError(err, target)`: Check if error is structured error

## Integration

This package is designed to be generic and reusable. To integrate it into your project:

1. Import the package
2. Define your application-specific error codes as constants
3. Use the provided functions and methods for error creation and handling
4. Optionally create wrapper functions for common error patterns in your domain

The package has no external dependencies and follows Go error handling conventions, making it compatible with existing error handling patterns and libraries.