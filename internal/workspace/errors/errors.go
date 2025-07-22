package errors

import (
	knoxErrors "github.com/tomdoesdev/knox/internal/errors"
	"github.com/tomdoesdev/knox/kit/errs"
)

var (
	ErrDirectoryReadFailed = errs.New(knoxErrors.SearchFailureCode, "failed to read directory")
	ErrNoWorkspace         = errs.New(knoxErrors.SearchFailureCode, "no workspace")
	ErrWorkspaceExists     = errs.New(knoxErrors.CreateFailureCode, "workspace already exists")
	ErrNotADirectory       = errs.New(knoxErrors.CreateFailureCode, "not a directory")
	ErrInvalidDatabase     = errs.New(knoxErrors.DatabaseFailureCode, "invalid database")
	ErrDatabasePathInvalid = errs.New(knoxErrors.DatabasePathInvalidCode, "invalid database path")
)
