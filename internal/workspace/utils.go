package workspace

import (
	"os"
	"path/filepath"

	"github.com/tomdoesdev/knox/internal/error_codes"
	"github.com/tomdoesdev/knox/kit/errs"
	"github.com/tomdoesdev/knox/kit/fs"
)

const (
	MissingWorkspaceDataDir errs.Code = "NO_WORKSPACE_DATA_DIR"
)

var (
	ErrWorkspaceMissing = errs.New(MissingWorkspaceDataDir, "missing workspace directory")
)

// WithLocalWorkspace handles getting the current working directory and finding
// the workspace, then calls the provided handler with the workspace.
// Returns raw errors without wrapping for maximum flexibility.
func WithLocalWorkspace(handler func(*Workspace) error) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	ws, err := FindWorkspace(cwd)
	if err != nil {
		return err
	}

	return handler(ws)
}

func FindWorkspace(path string) (*Workspace, error) {
	currentDir := path

	for {
		if ContainsDataDirectory(currentDir) {
			w, err := OpenWorkspace(currentDir)
			if err != nil {
				return nil, err
			}
			return w, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return nil, errs.Wrap(ErrWorkspaceMissing,
				error_codes.SearchFailureErrCode,
				"workspace directory search returned no data directory").
				WithPath(path).WithContext("root", parentDir)
		}
		currentDir = parentDir
	}
}

func IsDataDirectory(path string) bool {
	if fs.IsDir(path) && filepath.Base(path) == DataDirectoryName {
		return true
	}
	return false
}

func ContainsDataDirectory(path string) bool {
	if !fs.IsDir(path) {
		return false
	}

	path = filepath.Join(path, DataDirectoryName)

	return IsDataDirectory(path)
}

func OpenWorkspace(path string) (*Workspace, error) {
	panic("not implemented")
}
