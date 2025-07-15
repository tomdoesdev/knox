package errkit

import (
	"errors"
)

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrAlreadyExists  = errors.New("already exists")
)
