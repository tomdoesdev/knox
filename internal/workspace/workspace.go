package workspace

import (
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tomdoesdev/knox/internal/error_codes"
	"github.com/tomdoesdev/knox/internal/workspace/internal"
	"github.com/tomdoesdev/knox/internal/workspace/internal/database"

	"github.com/tomdoesdev/knox/kit/errs"
	"github.com/tomdoesdev/knox/kit/fs"
)

type InitResult string

const (
	Created InitResult = "workspace_created"
	Existed InitResult = "workspace_existed"
)

type Workspace struct {
	db   *database.Database
	path string
}

// LinkedVault represents a vault linked to the workspace
type LinkedVault struct {
	Alias     string `json:"alias"`
	Path      string `json:"path"`
	CreatedAt string `json:"created_at"`
}

// DataDir returns the full path to the workspace data directory
func (w *Workspace) DataDir() string {
	if filepath.Base(w.path) == internal.DataDirectoryName {
		return w.path
	}
	return filepath.Clean(filepath.Join(w.path, internal.DataDirectoryName))
}

// Dir returns the path to the parent directory containing the workspace data directory
func (w *Workspace) Dir() string {
	return filepath.Clean(w.path)
}

func NewWorkspace(db *database.Database, path string) *Workspace {
	return &Workspace{db: db, path: path}
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
		if internal.ContainsDataDirectory(currentDir) {
			w, err := OpenWorkspace(currentDir)
			if err != nil {
				return nil, err
			}
			return w, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return nil, internal.ErrNoWorkspace
		}
		currentDir = parentDir
	}
}

// EnsureWorkspace ensures that a workspace directory exists at the specified path.
func EnsureWorkspace(path string) (InitResult, error) {
	if internal.ContainsDataDirectory(path) {
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
	if internal.ContainsDataDirectory(path) {
		slog.Debug("creating workspace", slog.String("path", path), slog.Bool("exists", true))
		return nil, internal.ErrWorkspaceExists
	}

	dir := filepath.Join(path, internal.DataDirectoryName)

	if !fs.IsDir(dir) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			return nil, errs.Wrap(err, error_codes.CreateFailureErrCode, "failed to create workspace directory").WithContext("path", path)
		}
	}

	// Create projects directory
	projectsDir := filepath.Join(dir, internal.ProjectsDirectoryName)
	if !fs.IsDir(projectsDir) {
		err := os.MkdirAll(projectsDir, 0700)
		if err != nil {
			return nil, errs.Wrap(err, error_codes.CreateFailureErrCode, "failed to create projects directory").WithContext("path", projectsDir)
		}
	}

	dbPath := database.NewPath(path)

	db, err := database.EnsureWorkspaceDatabase(dbPath)
	if err != nil {
		return nil, err
	}

	workspace := NewWorkspace(db, path)

	// Create default project
	defaultProject := NewProject("default", "Default project for workspace")
	err = workspace.CreateProject(defaultProject)
	if err != nil {
		return nil, errs.Wrap(err, error_codes.CreateFailureErrCode, "failed to create default project")
	}

	// Set default project as current
	err = workspace.SetCurrentProject("default")
	if err != nil {
		return nil, errs.Wrap(err, error_codes.CreateFailureErrCode, "failed to set default project as current")
	}

	return workspace, nil

}

func OpenWorkspace(path string) (*Workspace, error) {
	if !internal.ContainsDataDirectory(path) {
		return nil, internal.ErrNoWorkspace.WithContext("path", path)
	}

	// Extract workspace root path (parent of .knox-workspace)
	workspaceRoot := filepath.Dir(filepath.Join(path, internal.DataDirectoryName))

	dbPath := database.NewPath(path)

	db, err := database.OpenWorkspaceDatabase(dbPath)
	if err != nil {
		return nil, err
	}
	return NewWorkspace(db, workspaceRoot), nil
}

// ProjectsPath returns the path to the projects directory
func (w *Workspace) ProjectsPath() string {
	return filepath.Join(w.path, internal.DataDirectoryName, internal.ProjectsDirectoryName)
}

// CreateProject creates a new project file
func (w *Workspace) CreateProject(project *Project) error {
	// Validate project structure
	if err := project.Validate(); err != nil {
		return err
	}

	// Validate against available vaults
	vaults, err := w.GetLinkedVaultAliases()
	if err != nil {
		return errs.Wrap(err, error_codes.DatabaseFailureErrCode, "failed to get linked vaults")
	}

	if err := project.ValidateWithVaults(vaults); err != nil {
		return err
	}

	projectPath := filepath.Join(w.ProjectsPath(), project.Name+".json")

	// Check if project already exists
	if fs.IsFile(projectPath) {
		return errs.New(error_codes.ProjectExistsErrCode, "project already exists").WithContext("name", project.Name)
	}

	data, err := project.ToJSON()
	if err != nil {
		return errs.Wrap(err, error_codes.ProjectInvalidErrCode, "failed to serialize project")
	}

	err = os.WriteFile(projectPath, data, 0600)
	if err != nil {
		return errs.Wrap(err, error_codes.FilePermissionErrCode, "failed to write project file").WithContext("path", projectPath)
	}

	return nil
}

// LoadProject loads a project by name
func (w *Workspace) LoadProject(name string) (*Project, error) {
	projectPath := filepath.Join(w.ProjectsPath(), name+".json")

	if !fs.IsFile(projectPath) {
		return nil, errs.New(error_codes.ProjectNotFoundErrCode, "project not found").WithContext("name", name)
	}

	data, err := os.ReadFile(projectPath)
	if err != nil {
		return nil, errs.Wrap(err, error_codes.FileNotFoundErrCode, "failed to read project file").WithContext("path", projectPath)
	}

	project, err := FromJSON(data)
	if err != nil {
		return nil, errs.Wrap(err, error_codes.ProjectInvalidErrCode, "failed to parse project file").WithContext("path", projectPath)
	}

	return project, nil
}

// UpdateProject updates an existing project file
func (w *Workspace) UpdateProject(project *Project) error {
	// Validate project structure
	if err := project.Validate(); err != nil {
		return err
	}

	// Validate against available vaults
	vaults, err := w.GetLinkedVaultAliases()
	if err != nil {
		return errs.Wrap(err, error_codes.DatabaseFailureErrCode, "failed to get linked vaults")
	}

	if err := project.ValidateWithVaults(vaults); err != nil {
		return err
	}

	projectPath := filepath.Join(w.ProjectsPath(), project.Name+".json")

	if !fs.IsFile(projectPath) {
		return errs.New(error_codes.ProjectNotFoundErrCode, "project not found").WithContext("name", project.Name)
	}

	data, err := project.ToJSON()
	if err != nil {
		return errs.Wrap(err, error_codes.ProjectInvalidErrCode, "failed to serialize project")
	}

	err = os.WriteFile(projectPath, data, 0600)
	if err != nil {
		return errs.Wrap(err, error_codes.FilePermissionErrCode, "failed to write project file").WithContext("path", projectPath)
	}

	return nil
}

// DeleteProject removes a project file
func (w *Workspace) DeleteProject(name string) error {
	projectPath := filepath.Join(w.ProjectsPath(), name+".json")

	if !fs.IsFile(projectPath) {
		return errs.New(error_codes.ProjectNotFoundErrCode, "project not found").WithContext("name", name)
	}

	err := os.Remove(projectPath)
	if err != nil {
		return errs.Wrap(err, error_codes.FilePermissionErrCode, "failed to delete project file").WithContext("path", projectPath)
	}

	return nil
}

// ListProjects returns a list of all project names
func (w *Workspace) ListProjects() ([]string, error) {
	projectsDir := w.ProjectsPath()

	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return nil, errs.Wrap(err, error_codes.DirectoryInvalidErrCode, "failed to read projects directory").WithContext("path", projectsDir)
	}

	var projects []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			name := strings.TrimSuffix(entry.Name(), ".json")
			projects = append(projects, name)
		}
	}

	return projects, nil
}

// GetLinkedVaultAliases returns a list of all linked vault aliases
func (w *Workspace) GetLinkedVaultAliases() ([]string, error) {
	// TODO: Implement database query to get vault aliases
	// For now, return empty slice - this will be implemented when we add vault linking
	return []string{}, nil
}

// CurrentProject returns the currently active project name
func (w *Workspace) CurrentProject() (string, error) {
	return w.GetSetting("current_project")
}

// SetCurrentProject sets the currently active project
func (w *Workspace) SetCurrentProject(projectName string) error {
	return w.SetSetting("current_project", projectName)
}

// GetSetting retrieves a setting value from the workspace_settings table
func (w *Workspace) GetSetting(key string) (string, error) {
	query := "SELECT value FROM workspace_settings WHERE key = ? AND category = 'config'"

	var value string
	err := w.db.DB().QueryRow(query, key).Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.New(error_codes.SearchFailureErrCode, "setting not found").WithContext("key", key)
		}
		return "", errs.Wrap(err, error_codes.DatabaseFailureErrCode, "failed to get setting").WithContext("key", key)
	}

	return value, nil
}

// SetSetting stores a setting value in the workspace_settings table
func (w *Workspace) SetSetting(key, value string) error {
	query := `
		INSERT INTO workspace_settings (key, value, category, updated_at) 
		VALUES (?, ?, 'config', CURRENT_TIMESTAMP)
		ON CONFLICT(key) DO UPDATE SET 
			value = excluded.value,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := w.db.DB().Exec(query, key, value)
	if err != nil {
		return errs.Wrap(err, error_codes.DatabaseFailureErrCode, "failed to set setting").WithContext("key", key).WithContext("value", value)
	}

	return nil
}

// GetMeta retrieves a metadata value from the workspace_settings table
func (w *Workspace) GetMeta(key string) (string, error) {
	query := "SELECT value FROM workspace_settings WHERE key = ? AND category = 'meta'"

	var value string
	err := w.db.DB().QueryRow(query, key).Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.New(error_codes.SearchFailureErrCode, "metadata not found").WithContext("key", key)
		}
		return "", errs.Wrap(err, error_codes.DatabaseFailureErrCode, "failed to get metadata").WithContext("key", key)
	}

	return value, nil
}

// SetMeta stores a metadata value in the workspace_settings table
func (w *Workspace) SetMeta(key, value string) error {
	query := `
		INSERT INTO workspace_settings (key, value, category, updated_at) 
		VALUES (?, ?, 'meta', CURRENT_TIMESTAMP)
		ON CONFLICT(key) DO UPDATE SET 
			value = excluded.value,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := w.db.DB().Exec(query, key, value)
	if err != nil {
		return errs.Wrap(err, error_codes.DatabaseFailureErrCode, "failed to set metadata").WithContext("key", key).WithContext("value", value)
	}

	return nil
}

// LinkVault links a vault to the workspace with the given alias
func (w *Workspace) LinkVault(alias, vaultPath string) error {
	// Validate alias
	if alias == "" {
		return errs.New(error_codes.ValidationErrCode, "vault alias cannot be empty")
	}

	// Validate vault path
	if vaultPath == "" {
		return errs.New(error_codes.ValidationErrCode, "vault path cannot be empty")
	}

	// Check if alias already exists
	query := "SELECT COUNT(*) FROM linked_vaults WHERE alias = ?"
	var count int
	err := w.db.DB().QueryRow(query, alias).Scan(&count)
	if err != nil {
		return errs.Wrap(err, error_codes.DatabaseFailureErrCode, "failed to check vault alias").WithContext("alias", alias)
	}
	if count > 0 {
		return errs.New(error_codes.ValidationErrCode, "vault alias already exists").WithContext("alias", alias)
	}

	// Check if path already exists
	query = "SELECT COUNT(*) FROM linked_vaults WHERE path = ?"
	err = w.db.DB().QueryRow(query, vaultPath).Scan(&count)
	if err != nil {
		return errs.Wrap(err, error_codes.DatabaseFailureErrCode, "failed to check vault path").WithContext("path", vaultPath)
	}
	if count > 0 {
		return errs.New(error_codes.ValidationErrCode, "vault path already linked").WithContext("path", vaultPath)
	}

	// Insert the vault link
	insertQuery := "INSERT INTO linked_vaults (alias, path) VALUES (?, ?)"
	_, err = w.db.DB().Exec(insertQuery, alias, vaultPath)
	if err != nil {
		return errs.Wrap(err, error_codes.DatabaseFailureErrCode, "failed to link vault").WithContext("alias", alias).WithContext("path", vaultPath)
	}

	return nil
}

// GetLinkedVaults returns all linked vaults
func (w *Workspace) GetLinkedVaults() ([]LinkedVault, error) {
	query := "SELECT alias, path, created_at FROM linked_vaults ORDER BY alias"

	rows, err := w.db.DB().Query(query)
	if err != nil {
		return nil, errs.Wrap(err, error_codes.DatabaseFailureErrCode, "failed to query linked vaults")
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var vaults []LinkedVault
	for rows.Next() {
		var vault LinkedVault
		err := rows.Scan(&vault.Alias, &vault.Path, &vault.CreatedAt)
		if err != nil {
			return nil, errs.Wrap(err, error_codes.DatabaseFailureErrCode, "failed to scan vault row")
		}
		vaults = append(vaults, vault)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.Wrap(err, error_codes.DatabaseFailureErrCode, "error iterating vault rows")
	}

	return vaults, nil
}
