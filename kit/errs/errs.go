package errs

import (
	"errors"
	"fmt"
	"strings"
)

// Code represents different types of application errors
type Code string

// Error represents a structured error with code, message, cause, and context
type Error struct {
	Code    Code
	Message string
	Cause   error
	Context map[string]interface{}
}

// New creates a new Error with the specified code and message
func New(code Code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Error implements the error interface
func (e *Error) Error() string {
	var msg strings.Builder

	msg.WriteString(fmt.Sprintf("[%s] ", e.Code))
	msg.WriteString(e.Message)
	if e.Cause != nil {
		msg.WriteString(": ")
		msg.WriteString(e.Cause.Error())
	}

	if e.Context != nil {
		msg.WriteString(" | context: ")
		first := true
		for key, value := range e.Context {
			if !first {
				msg.WriteString(", ")
			}
			msg.WriteString(fmt.Sprintf("%s=%v", key, value))
			first = false
		}
	}
	return msg.String()
}

// Unwrap returns the underlying cause error
func (e *Error) Unwrap() error {
	return e.Cause
}

// Is checks if this error matches the target error by comparing codes
func (e *Error) Is(target error) bool {
	var t *Error
	if errors.As(target, &t) {
		return e.Code == t.Code
	}
	return false
}

// WithContext adds a context key-value pair to the error
func (e *Error) WithContext(key string, value interface{}) *Error {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// Wrap wraps an existing error with structured error context
func Wrap(err error, code Code, message string, params ...any) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(message, params...),
		Cause:   err,
	}
}

// Is checks if an error has a specific error code
func Is(err error, code Code) bool {
	var structuredErr *Error
	if errors.As(err, &structuredErr) {
		return structuredErr.Code == code
	}
	return false
}

// GetCode extracts the error code from an error
func GetCode(err error) Code {
	var structuredErr *Error
	if errors.As(err, &structuredErr) {
		return structuredErr.Code
	}
	return ""
}

// WrapWithContext wraps an existing error with structured error context and additional context pairs
// should be provided as key-value pairs: "key1", value1, "key2", value2, ...
func WrapWithContext(err error, code Code, message string, contextPairs ...interface{}) *Error {
	structuredErr := &Error{
		Code:    code,
		Message: message,
		Cause:   err,
	}

	if len(contextPairs) > 0 {
		structuredErr.Context = make(map[string]interface{})
		for i := 0; i < len(contextPairs)-1; i += 2 {
			if key, ok := contextPairs[i].(string); ok {
				structuredErr.Context[key] = contextPairs[i+1]
			}
		}
	}

	return structuredErr
}

// ErrorBuilder provides a fluent interface for building complex errors
type ErrorBuilder struct {
	err *Error
}

// Build creates a new ErrorBuilder with the specified error code
func Build(code Code) *ErrorBuilder {
	return &ErrorBuilder{
		err: &Error{
			Code: code,
		},
	}
}

// Message sets the error message with optional formatting
func (b *ErrorBuilder) Message(msg string, args ...interface{}) *ErrorBuilder {
	if len(args) > 0 {
		b.err.Message = fmt.Sprintf(msg, args...)
	} else {
		b.err.Message = msg
	}
	return b
}

// Cause sets the underlying cause error
func (b *ErrorBuilder) Cause(err error) *ErrorBuilder {
	b.err.Cause = err
	return b
}

// Context adds a single context key-value pair
func (b *ErrorBuilder) Context(key string, value interface{}) *ErrorBuilder {
	if b.err.Context == nil {
		b.err.Context = make(map[string]interface{})
	}
	b.err.Context[key] = value
	return b
}

// Contexts adds multiple context pairs: "key1", value1, "key2", value2, ...
func (b *ErrorBuilder) Contexts(pairs ...interface{}) *ErrorBuilder {
	if len(pairs) > 0 {
		if b.err.Context == nil {
			b.err.Context = make(map[string]interface{})
		}
		for i := 0; i < len(pairs)-1; i += 2 {
			if key, ok := pairs[i].(string); ok {
				b.err.Context[key] = pairs[i+1]
			}
		}
	}
	return b
}

// Error returns the built Error
func (b *ErrorBuilder) Error() *Error {
	return b.err
}

// Convenience methods for common context types

// WithPath adds a file/directory path to the error context
func (e *Error) WithPath(path string) *Error {
	return e.WithContext("path", path)
}

// WithFile adds a filename to the error context
func (e *Error) WithFile(filename string) *Error {
	return e.WithContext("file", filename)
}

// WithAttempt adds an attempt count to the error context
func (e *Error) WithAttempt(count int) *Error {
	return e.WithContext("attempt", count)
}

// WithDatabase adds a database name/path to the error context
func (e *Error) WithDatabase(database string) *Error {
	return e.WithContext("database", database)
}

// WithOperation adds an operation name to the error context
func (e *Error) WithOperation(operation string) *Error {
	return e.WithContext("operation", operation)
}

// WithID adds an identifier to the error context
func (e *Error) WithID(id string) *Error {
	return e.WithContext("id", id)
}

// Convenience methods for ErrorBuilder

// Path adds a file/directory path to the error context
func (b *ErrorBuilder) Path(path string) *ErrorBuilder {
	return b.Context("path", path)
}

// File adds a filename to the error context
func (b *ErrorBuilder) File(filename string) *ErrorBuilder {
	return b.Context("file", filename)
}

// Attempt adds an attempt count to the error context
func (b *ErrorBuilder) Attempt(count int) *ErrorBuilder {
	return b.Context("attempt", count)
}

// Database adds a database name/path to the error context
func (b *ErrorBuilder) Database(database string) *ErrorBuilder {
	return b.Context("database", database)
}

// Operation adds an operation name to the error context
func (b *ErrorBuilder) Operation(operation string) *ErrorBuilder {
	return b.Context("operation", operation)
}

// ID adds an identifier to the error context
func (b *ErrorBuilder) ID(id string) *ErrorBuilder {
	return b.Context("id", id)
}
