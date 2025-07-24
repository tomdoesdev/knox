package internal

import (
	"github.com/tomdoesdev/knox/internal/error_codes"
	"github.com/tomdoesdev/knox/kit/errs"
)

var (
	ErrNoWorkspace     = errs.New(error_codes.SearchFailureErrCode, "no workspace")
	ErrWorkspaceExists = errs.New(error_codes.CreateFailureErrCode, "workspace already exists")
	ErrInvalidDatabase = errs.New(error_codes.DatabaseFailureErrCode, "invalid database")
)
