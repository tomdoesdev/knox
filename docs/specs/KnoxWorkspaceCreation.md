# Knox Workspace Creation Implementation Spec

## Purpose & User Problem

Complete the implementation of workspace creation functionality for Knox v2, building on the existing workspace infrastructure to enable users to initialize workspaces with `knox init`.

## Success Criteria

- [ ] `knox init` successfully creates workspace in current directory
- [ ] Workspace includes all required files: database and config
- [ ] Auto-links default vault during workspace creation
- [ ] Proper error handling for existing workspaces and permission issues
- [ ] Clear user feedback on success/failure
- [ ] Clean code with consistent imports and no redundant functionality

## Current Implementation Status

**What exists:**
- `CreateWorkspace()` function with database creation
- SQLite schema for workspace database
- Directory structure validation (`ContainsDataDirectory()`)
- Error handling infrastructure
- Config interface definition (`ConfigManager`)

**What's missing:**
1. Config file implementation (JSON I/O)
2. Workspace configuration methods (database-stored settings)
3. Default vault auto-linking
4. CLI integration (`workspaceHandler()` is placeholder)
5. User feedback and proper error messages

## Implementation Tasks

### Task 1: Implement Config File Management

**File:** `internal/workspace/config.go`

Add concrete implementation of `ConfigManager` interface:

```go
type FileConfigManager struct {
    configPath string
}

func NewFileConfigManager(workspacePath string) *FileConfigManager {
    configPath := filepath.Join(workspacePath, "config")
    return &FileConfigManager{configPath: configPath}
}

func (f *FileConfigManager) Read() (*Config, error) {
    // Read JSON config file
    // Return default config if file doesn't exist
}

func (f *FileConfigManager) Write(config *Config) error {
    // Write config as JSON to file
    // Create file with 0600 permissions
}

// Helper function to create default config
func DefaultConfig() *Config {
    return &Config{
        WorkspaceVersion: "v2",
    }
}
```

**Requirements:**
- Handle missing config file gracefully (return default config)
- Create config file with appropriate permissions (0600)
- Use JSON format for extensibility
- Error handling for file I/O operations

### Task 2: Workspace Configuration Methods

**File:** `internal/workspace/workspace.go`

Add methods to `Workspace` struct for database-stored configuration:

```go
// GetSetting reads a setting value from workspace_settings table
func (w *Workspace) GetSetting(key string) (string, error) {
    // Query workspace_settings table for the key
    // Return empty string if key doesn't exist
}

// SetSetting writes a setting value to workspace_settings table
func (w *Workspace) SetSetting(key, value, category string) error {
    // INSERT OR REPLACE into workspace_settings table
    // Handle database transaction properly
    // Auto-set updated_at timestamp
}

// UnsetSetting removes a setting key from workspace_settings table
func (w *Workspace) UnsetSetting(key string) error {
    // DELETE from workspace_settings table
}

// Convenience methods for configuration (user-changeable settings)
func (w *Workspace) GetConfig(key string) (string, error) {
    return w.GetSetting("config." + key)
}

func (w *Workspace) SetConfig(key, value string) error {
    return w.SetSetting("config." + key, value, "config")
}

func (w *Workspace) UnsetConfig(key string) error {
    return w.UnsetSetting("config." + key)
}

// Convenience methods for metadata (system-managed settings)
func (w *Workspace) GetMeta(key string) (string, error) {
    return w.GetSetting("meta." + key)
}

func (w *Workspace) SetMeta(key, value string) error {
    return w.SetSetting("meta." + key, value, "meta")
}

// Project-specific convenience methods
func (w *Workspace) GetCurrentProject() (string, error) {
    return w.GetConfig("current_project")
}

func (w *Workspace) SetCurrentProject(projectName string) error {
    return w.SetConfig("current_project", projectName)
}

func (w *Workspace) ClearCurrentProject() error {
    return w.UnsetConfig("current_project")
}

func (w *Workspace) IsDetached() bool {
    project, err := w.GetCurrentProject()
    return err != nil || project == ""
}
```

**Requirements:**
- Use database transactions for atomic operations
- Handle missing keys gracefully (return empty string)
- Support INSERT OR REPLACE pattern with automatic updated_at timestamp
- Semantic key prefixes: `config.*` for user settings, `meta.*` for system metadata
- Category column for easy filtering and future extensibility
- Validation of setting keys and values

### Task 3: Database Schema Update

**File:** `internal/workspace/database/schema.go`

Update the schema to use the new `workspace_settings` table:

```sql
CREATE TABLE workspace_settings (
    key TEXT PRIMARY KEY,
    value TEXT,
    category TEXT NOT NULL, -- 'meta' or 'config'
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### Task 4: Default Vault Auto-Linking

**File:** `internal/workspace/database/database.go`

Add method to auto-link default vault:

```go
// LinkDefaultVault automatically links the default vault during workspace creation
func (db *Database) LinkDefaultVault() error {
    defaultVaultPath := getDefaultVaultPath()
    
    // Check if default vault exists, create if needed
    if !vaultExists(defaultVaultPath) {
        if err := createDefaultVault(defaultVaultPath); err != nil {
            return err
        }
    }
    
    // Insert into linked_vaults table
    query := `INSERT INTO linked_vaults (alias, path) VALUES (?, ?)`
    _, err := db.db.Exec(query, "default", defaultVaultPath)
    return err
}

func getDefaultVaultPath() string {
    // Use KNOX_ROOT env var if set, otherwise ~/.knox/vault.db
}
```

**Requirements:**
- Respect `KNOX_ROOT` environment variable
- Create default vault if it doesn't exist
- Use "default" as alias for the default vault
- Handle case where default vault is already linked

### Task 5: Enhanced Workspace Creation

**File:** `internal/workspace/workspace.go`

Update `CreateWorkspace()` function:

```go
func CreateWorkspace(path string) (*Workspace, error) {
    // Check if workspace already exists
    if ContainsDataDirectory(path) {
        return nil, errors.ErrWorkspaceExists.WithContext("path", path)
    }

    workspacePath := filepath.Join(path, constants.DataDirectoryName)

    // Create workspace directory
    if err := os.MkdirAll(workspacePath, 0700); err != nil {
        return nil, errs.Wrap(err, errors.CreateFailureCode, "failed to create workspace directory").WithContext("path", workspacePath)
    }

    // Create database
    db, err := database2.CreateWorkspaceDatabase(workspacePath)
    if err != nil {
        return nil, err
    }

    // Create workspace instance
    workspace := newWorkspace(db)

    // Set initial metadata
    if err := workspace.SetMeta("created_at", time.Now().Format(time.RFC3339)); err != nil {
        return nil, errs.Wrap(err, errors.CreateFailureCode, "failed to set creation metadata")
    }
    if err := workspace.SetMeta("knox_version", "v2"); err != nil {
        return nil, errs.Wrap(err, errors.CreateFailureCode, "failed to set version metadata")
    }

    // Auto-link default vault
    if err := db.LinkDefaultVault(); err != nil {
        return nil, errs.Wrap(err, errors.CreateFailureCode, "failed to link default vault")
    }

    // Create config file
    configManager := NewFileConfigManager(workspacePath)
    if err := configManager.Write(DefaultConfig()); err != nil {
        return nil, errs.Wrap(err, errors.CreateFailureCode, "failed to create config file")
    }

    // Initialize in detached state (no config.current_project setting)
    
    return workspace, nil
}
```

### Task 6: CLI Integration

**File:** `cmd/knox/internal/commands/init.go`

Implement `workspaceHandler()`:

```go
func workspaceHandler() error {
    currentDir, err := os.Getwd()
    if err != nil {
        return fmt.Errorf("failed to get current directory: %w", err)
    }

    workspace, err := workspace.CreateWorkspace(currentDir)
    if err != nil {
        if errs.Is(err, workspace.ErrWorkspaceExists) {
            return fmt.Errorf("workspace already exists in current directory")
        }
        return fmt.Errorf("failed to create workspace: %w", err)
    }
    defer workspace.Close()

    workspacePath := filepath.Join(currentDir, constants.DataDirectoryName)
    fmt.Printf("Initialized empty Knox workspace in %s\n", workspacePath)
    fmt.Println("Linked default vault 'default' (path: ~/.knox/vault.db)")
    fmt.Println("Workspace is in detached state. Use 'knox switch <project>' to attach to a project.")
    fmt.Println("Note: 'knox switch' is equivalent to 'knox config current_project <project@vault>'")
    
    return nil
}
```

**Requirements:**
- Import workspace package and constants
- Handle errors gracefully with user-friendly messages
- Provide clear success feedback
- Include guidance on next steps

## Error Scenarios & Handling

1. **Workspace already exists**: Clear message, suggest using existing workspace
2. **Permission denied**: Clear message about directory permissions
3. **Default vault creation fails**: Specific error about vault setup
4. **Config file creation fails**: Error about config initialization
5. **Database creation fails**: Database-specific error message

## File Structure After Creation

```
/path/to/my-app/
├── .knox-workspace/       # Workspace marker directory
│   ├── workspace.db       # SQLite database with default vault linked
│   └── config             # JSON config file {"workspace_version": "v2"}
├── src/
└── README.md
```

**Notes:**
- No `workspace.current_project` config initially (detached state)
- Database contains one entry in `linked_vaults` table for default vault
- Database `workspace_settings` table stores both configuration and metadata with semantic prefixes
- Config file uses JSON for future extensibility

## Testing Strategy

1. **Unit tests** for each new function
2. **Integration test** for full workspace creation flow
3. **Error scenario tests** (permissions, existing workspace)
4. **CLI test** with mock filesystem

## Implementation Order

1. Config file management (`FileConfigManager`)
2. Database schema update (rename `workspace_meta` to `workspace_settings`, add `category` and `updated_at`)
3. Workspace settings methods (`GetSetting`, `SetSetting`, config/meta convenience methods, project methods)
4. Default vault auto-linking (`LinkDefaultVault`)
5. Enhanced workspace creation (update `CreateWorkspace`)
6. CLI integration (implement `workspaceHandler`)
7. Testing and validation

## Out of Scope

- Workspace migration/import functionality
- Multiple workspace support in single directory
- Workspace configuration beyond version tracking
- Project creation during workspace init

## Design Benefits

This database-centric approach provides several advantages:

1. **Unified settings model**: Current project is just another workspace setting stored in `workspace_settings`
2. **Command equivalence**: `knox switch <project@vault>` ≡ `knox config current_project <project@vault>`
3. **Semantic organization**: `config.*` keys for user settings, `meta.*` keys for system metadata
4. **Single source of truth**: All workspace state in one database table with clear categorization
5. **Atomic operations**: Database transactions prevent inconsistent state
6. **Extensible**: Easy to add more settings using same table pattern with category filtering
7. **Backup friendly**: One file contains all workspace state
8. **No file synchronization issues**: Database handles concurrent access properly
9. **Rich metadata**: Track creation time, Knox version, usage patterns, and more

## Notes

- Maintain backward compatibility with existing vault structure
- Follow established error handling patterns
- Use consistent file permissions throughout
- Imports have been cleaned up in `workspace.go` (v2 import removed)
- Current project stored as `config.current_project` key in `workspace_settings` table
- Metadata like creation time stored as `meta.created_at` in same table
- Category column enables filtering between user config and system metadata