package workspace

import (
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tomdoesdev/knox/internal/workspace/constants"
	"github.com/tomdoesdev/knox/internal/workspace/database"
	database2 "github.com/tomdoesdev/knox/internal/workspace/database"
	"github.com/tomdoesdev/knox/internal/workspace/errors"

	"github.com/tomdoesdev/knox/kit/fs"
	"github.com/tomdoesdev/knox/pkg/errs"
)

type Workspace struct {
	db *database2.Database
}

func newWorkspace(db *database2.Database) *Workspace {
	return &Workspace{db: db}
}

type LinkedVault struct {
	ID        int       `json:"id"`
	Alias     string    `json:"alias"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
}

type LinkedProject struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	VaultID     int       `json:"vault_id"`
	ProjectName *string   `json:"project_name,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// FindWorkspace finds the nearest .knox directory, traversing up the directory tree until it finds it.
func FindWorkspace(path string) (*Workspace, error) {
	currentDir := path

	for {
		dataDir := filepath.Join(currentDir, constants.DataDirectoryName)
		if ContainsDataDirectory(dataDir) {
			w, err := OpenWorkspace(dataDir)
			if err != nil {
				return nil, err
			}
			return w, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return nil, errors.ErrNoWorkspace
		}
		currentDir = parentDir
	}
}

// CreateWorkspace creates a workspace if one doesn't already exist
// Returns ErrWorkspaceExists if the path is already a workspace
func CreateWorkspace(path string) (*Workspace, error) {
	if ContainsDataDirectory(path) {
		// If a workspace already exists we abort.
		return nil, errors.ErrWorkspaceExists.WithContext("path", path)
	}

	path = filepath.Join(path, constants.DataDirectoryName)

	if !fs.IsDir(path) {
		err := os.MkdirAll(path, 0700)
		if err != nil {
			return nil, errs.Wrap(err, errors.CreateFailureCode, "failed to create workspace directory").WithContext("path", path)
		}
	}

	db, err := database2.CreateWorkspaceDatabase(path)
	if err != nil {
		return nil, err
	}

	return newWorkspace(db), nil

}

func OpenWorkspace(path string) (*Workspace, error) {
	if !ContainsDataDirectory(path) {
		return nil, errors.ErrNoWorkspace.WithContext("path", path)
	}

	dbPath, err := database.newDatabasePath(path)
	if err != nil {
		return nil, err
	}

	db, err := database2.OpenWorkspaceDatabase(dbPath)
	if err != nil {
		return nil, err
	}
	return newWorkspace(db), nil
}

func IsDataDirectory(path string) bool {
	if fs.IsDir(path) && filepath.Base(path) == constants.DataDirectoryName {
		return true
	}
	return false
}

// ContainsDataDirectory returns true if the
func ContainsDataDirectory(path string) bool {
	if !fs.IsDir(path) {
		return false
	}

	return IsDataDirectory(path)
}
