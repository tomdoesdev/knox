package workspace

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tomdoesdev/knox/internal/workspace/constants"
	"github.com/tomdoesdev/knox/internal/workspace/database"
	"github.com/tomdoesdev/knox/internal/workspace/errors"

	"github.com/tomdoesdev/knox/kit/fs"
	"github.com/tomdoesdev/knox/pkg/errs"
)

type InitResult string

const (
	Created InitResult = "workspace_created"
	Existed InitResult = "workspace_existed"
)

type Workspace struct {
	db *database.Database
}

func newWorkspace(db *database.Database) *Workspace {
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

// EnsureWorkspace ensures that a workspace directory exists at the specified path.
func EnsureWorkspace(path string) (InitResult, error) {
	if ContainsDataDirectory(path) {
		slog.Debug("ensuring workspace", slog.String("path", path), slog.Bool("exists", true))
		_, err := OpenWorkspace(path)
		if err != nil {
			return "", err
		}

		return Existed, nil
	}

	slog.Debug("ensuring workspace", slog.String("path", path), slog.Bool("exists", false))
	_, err := CreateWorkspace(path)
	if err != nil {
		return "", err
	}

	return Created, nil
}

// CreateWorkspace creates a workspace if one doesn't already exist
// Returns ErrWorkspaceExists if the path is already a workspace
func CreateWorkspace(path string) (*Workspace, error) {
	if ContainsDataDirectory(path) {
		slog.Debug("creating workspace", slog.String("path", path), slog.Bool("exists", true))
		return nil, errors.ErrWorkspaceExists
	}

	dir := filepath.Join(path, constants.DataDirectoryName)

	if !fs.IsDir(dir) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			return nil, errs.Wrap(err, errors.CreateFailureCode, "failed to create workspace directory").WithContext("path", path)
		}
	}

	dbPath := database.NewPath(path)

	db, err := database.EnsureWorkspaceDatabase(dbPath)
	if err != nil {
		return nil, err
	}

	return newWorkspace(db), nil

}

func OpenWorkspace(path string) (*Workspace, error) {
	if !ContainsDataDirectory(path) {
		return nil, errors.ErrNoWorkspace.WithContext("path", path)
	}

	dbPath := database.NewPath(path)

	db, err := database.OpenWorkspaceDatabase(dbPath)
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
