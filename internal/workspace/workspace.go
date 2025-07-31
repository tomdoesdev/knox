package workspace

import (
	"encoding/json"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/tomdoesdev/knox/internal"
	"github.com/tomdoesdev/knox/kit/errs"
	"github.com/tomdoesdev/knox/kit/fs"
)
import "github.com/google/uuid"

const (
	DataDirectoryName = ".knox-workspace"
	ConfigFilename    = "workspace.toml"
	StateFilename     = "state.json"
)

const (
	ECodeWorkspaceExists     errs.Code = "WORKSPACE_EXISTS"
	ECodeWorkspaceInitFailed errs.Code = "WORKSPACE_INIT_FAILED"
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

/*
CreateWorkspaceDefault creates a new knox workspace directory and the required files.

If a workspace already exists at the given path the function will abort and return ECodeWorkspaceExists.

Returns:
  - ECodeWorkspaceExists
  - ECodeWorkspaceInitFailed
*/
func CreateWorkspaceDefault(path string) error {
	if IsDataDirectory(path) {
		return errs.New(ECodeWorkspaceExists, "workspace already exists").WithPath(path)
	}

	dataDir := filepath.Join(path, DataDirectoryName)
	projectsDir := filepath.Join(dataDir, "projects")

	// Create knox workspace data directory and 'projects' subdirectory
	if err := fs.MkdirAll(projectsDir, internal.DataDirectoryPermissions); err != nil {
		return errs.Wrap(err, ECodeWorkspaceInitFailed, "failed to create workspace directory").
			WithPath(filepath.Dir(dataDir)).
			WithContext("permissions", internal.DataDirectoryPermissions)
	}

	configBytes, err := toml.Marshal(NewConfigDefault())
	if err != nil {
		return errs.Wrap(err, ECodeWorkspaceInitFailed, "failed to marshal default workspace config")
	}

	stateBytes, err := json.Marshal(NewStateDefault())
	if err != nil {
		return errs.Wrap(err, ECodeWorkspaceInitFailed, "failed to marshal default workspace state")
	}

	err = fs.WriteFile(filepath.Join(dataDir, ConfigFilename), configBytes, internal.WorkspaceFilePermissions)
	if err != nil {
		if errs.Is(err, fs.ECodeEntityExists) {
			return errs.New(ECodeWorkspaceExists, "workspace config already exists")
		}
		return errs.Wrap(err, ECodeWorkspaceInitFailed, "failed to write workspace config")
	}

	err = fs.WriteFile(filepath.Join(dataDir, StateFilename), stateBytes, internal.WorkspaceFilePermissions)
	if err != nil {
		if errs.Is(err, fs.ECodeEntityExists) {
			return errs.New(ECodeWorkspaceExists, "workspace state already exists")
		}
		return errs.Wrap(err, ECodeWorkspaceInitFailed, "failed to write workspace state")
	}

	return nil

}
