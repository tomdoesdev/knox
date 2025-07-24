package common

import (
	"github.com/tomdoesdev/knox/internal/error_codes"
	"github.com/tomdoesdev/knox/kit/errs"
	"github.com/urfave/cli/v3"
)

func ExpectExactArgCount(count int, message string, args cli.Args) error {
	argsLen := args.Len()

	if argsLen != count {
		return errs.New(error_codes.ValidationErrCode, message).
			WithContext("got", argsLen).
			WithContext("expected", count)

	}

	return nil
}
