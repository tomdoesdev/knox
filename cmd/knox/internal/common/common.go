package common

import (
	"github.com/tomdoesdev/knox/internal/error_codes"
	"github.com/tomdoesdev/knox/internal/workspace"
	"github.com/tomdoesdev/knox/kit/errs"
)

// WithLocalWorkspace wraps workspace.WithLocalWorkspace with command-specific error handling
func WithLocalWorkspace(handler func(*workspace.Workspace) error) error {
	err := workspace.WithLocalWorkspace(handler)
	if err != nil {
		return errs.Wrap(err, error_codes.SearchFailureErrCode, "workspace operation failed")
	}
	return nil
}
