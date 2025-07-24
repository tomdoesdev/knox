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
		// Check if this is a workspace-related error and wrap appropriately
		return errs.Wrap(err, error_codes.SearchFailureErrCode, "workspace operation failed")
	}
	return nil
}

func WithEnsuredLocalWorkspace(handler func(*workspace.Workspace, workspace.InitResult) error) error {
	err := workspace.WithEnsuredLocalWorkspace(handler)
	if err != nil {
		return errs.Wrap(err, error_codes.SearchFailureErrCode, "workspace operation failed")
	}
	return nil
}
