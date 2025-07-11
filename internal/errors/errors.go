package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

type VaultError struct {
	Op   string
	Path string
	Err  error
}

func (e *VaultError) Error() string {
	return fmt.Sprintf("%s %s: %v", e.Op, e.Path, e.Err)
}
