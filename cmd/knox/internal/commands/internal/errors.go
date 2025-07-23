package internal

import "github.com/tomdoesdev/knox/kit/errs"

var (
	ErrInvalidInitArgs = errs.New(InvalidCommandArgsCode, "invalid arguments for init command")
)

const (
	InvalidCommandArgsCode errs.Code = "invalid_command_args"
)
