package workspace

import (
	"context"
	"os"
	"path/filepath"

	"github.com/tomdoesdev/knox/kit/errs"
	"github.com/tomdoesdev/knox/kit/fs"
)
import "github.com/google/uuid"

const (
	DataDirectoryName = ".knox-workspace"
)

const (
	ECodeWorkspaceExists     errs.Code = "WORKSPACE_EXISTS"
	ECodeWorkspaceInitFailed errs.Code = "WORKSPACE_INIT_FAILED"
)

const (
	dataDirectoryPermissions os.FileMode = 0700
	workspaceFilePermissions os.FileMode = 0600
)

type Workspace struct {
	Id      uuid.UUID
	dataDir string
	rootDir string
}

func (w *Workspace) DataDir() string {
	return w.dataDir
}

func NewWorkspace(id uuid.UUID, path string) (*Workspace, error) {
	rootDir := path

	if filepath.Base(path) == DataDirectoryName {
		rootDir = filepath.Dir(path)
	}

	return &Workspace{dataDir: path, rootDir: rootDir, Id: id}, nil
}

// CreateWorkspace creates a new knox workspace directory and the required files.
// If a workspace already exists at the given path the function will abort and return an error.
func CreateWorkspace(path string) (*Workspace, error) {
	if IsDataDirectory(path) {
		return nil, errs.New(ECodeWorkspaceExists, "workspace already exists").WithPath(path)
	}

	dataDir := filepath.Join(path, DataDirectoryName)
	projectsDir := filepath.Join(dataDir, "projects")

	// Create knox workspace data directory and 'projects' subdirectory
	if err := fs.MkdirAll(projectsDir, dataDirectoryPermissions); err != nil {
		return nil, errs.Wrap(err, ECodeWorkspaceInitFailed, "failed to create workspace directory").
			WithPath(filepath.Dir(dataDir)).
			WithContext("permissions", dataDirectoryPermissions)
	}

	files := []string{"workspace.toml", "state.json"}
	// Create empty workspace and state files.
	for _, file := range files {
		file = filepath.Join(dataDir, file)

		err := fs.Touch(file, workspaceFilePermissions)
		if err != nil {
			return nil, errs.Wrap(err, ECodeWorkspaceInitFailed, "failed to create workspace file").
				WithPath(filepath.Dir(file)).
				WithFile(file)
		}
	}
}
