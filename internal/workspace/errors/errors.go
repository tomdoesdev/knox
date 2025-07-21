package errors

import "github.com/tomdoesdev/knox/pkg/errs"

const (
	SearchFailureCode   errs.ErrorCode = "search_failure"
	CreateFailureCode   errs.ErrorCode = "create_failure"
	DatabaseFailureCode errs.ErrorCode = "database_failure"

	DatabasePathInvalidCode errs.ErrorCode = "database_path_invalid"
)

var (
	ErrDirectoryReadFailed = errs.New(SearchFailureCode, "failed to read directory")
	ErrNoWorkspace         = errs.New(SearchFailureCode, "no workspace")
	ErrWorkspaceExists     = errs.New(CreateFailureCode, "workspace already exists")
	ErrNotADirectory       = errs.New(CreateFailureCode, "not a directory")
	ErrInvalidDatabase     = errs.New(DatabaseFailureCode, "invalid database")
	ErrDatabasePathInvalid = errs.New(DatabasePathInvalidCode, "invalid database path")
)
