package commands

import "github.com/tomdoesdev/knox/pkg/errs"

const (
	InvalidArguments errs.ErrorCode = "INVALID_ARGUMENTS"
)

var (
	ErrInvalidArguments = errs.New(InvalidArguments, "invalid arguments")
)
