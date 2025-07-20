package errs

import (
	"errors"
	"fmt"
	"strings"
)

// ErrorCode represents different types of Knox errs
type ErrorCode string

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
	var msg strings.Builder

	msg.WriteString(fmt.Sprintf("[%s] ", e.Code))
	msg.WriteString(e.Message)
	if e.Cause != nil {
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
func Wrap(err error, code ErrorCode, message string, params ...any) *KnoxError {
	return &KnoxError{
		Code:    code,
		Message: fmt.Sprintf(message, params...),
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
