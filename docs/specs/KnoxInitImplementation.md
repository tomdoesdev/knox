# Knox Init Implementation Spec

## Purpose & User Problem

Complete the implementation of `knox init` command for Knox v2, enabling users to initialize workspaces with git-like directory context that tracks current active project and links to vaults.

## Success Criteria

- [ ] `knox init` creates workspace in current directory
- [ ] Workspace auto-links default vault during creation  
- [ ] Database-only storage (no config files)
- [ ] Enhanced TEXT settings with registry-based validation
- [ ] Proper error handling for existing workspaces and permission issues
- [ ] Clear user feedback and guidance on next steps
- [ ] Workspace starts in detached state (no current project)

## Current Implementation Status

**What exists:**
- `CreateWorkspace()` function with basic database creation
- SQLite schema for workspace database (linked_vaults, linked_projects, workspace_meta)
- Directory structure validation (`ContainsDataDirectory()`)
- Error handling infrastructure
- CLI command structure (`knox init`, `knox init workspace`, `knox init project`)

**What needs completion:**
1. Database schema update (`workspace_meta` → `workspace_settings`)
2. Settings registry system for validation
3. Workspace settings methods (GetSetting, SetSetting, etc.)
4. Default vault auto-linking
5. CLI integration (`workspaceHandler()` implementation)
6. Enhanced workspace creation with metadata

## Workspace Structure

```
/path/to/my-app/           # Any directory can become a workspace
├── .knox-workspace/       # Workspace marker directory  
│   └── workspace.db       # Single SQLite database (everything stored here)
├── src/
└── README.md
```

**Database Design - Enhanced TEXT Settings:**
- **Single source of truth**: All workspace state in `workspace_settings` table
- **Registry validation**: Code-based registry defines valid settings and types  
- **Semantic prefixes**: `config.*` for user settings, `meta.*` for system metadata
- **Type safety**: Helper methods with automatic validation
- **Debuggable**: TEXT storage allows direct SQL inspection

## Implementation Tasks

### Task 1: Database Schema Update

**File:** `internal/workspace/database/schema.go`

Update schema to replace `workspace_meta` with `workspace_settings`:

```sql
-- Replace existing workspace_meta table
CREATE TABLE workspace_settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,        -- Always TEXT for debuggability
    category TEXT NOT NULL,     -- 'meta' or 'config'  
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Keep existing tables
CREATE TABLE linked_vaults (...);
CREATE TABLE linked_projects (...);
```

### Task 2: Settings Registry System

**File:** `internal/workspace/workspace.go`

Add settings registry and validation:

```go
// Settings registry for validation and documentation
type SettingDefinition struct {
    DataType     string // "string", "bool", "int"
    Category     string // "config", "meta"
    Description  string
    DefaultValue string
    Validator    func(string) error
}

// Built-in workspace settings registry
var WorkspaceSettings = map[string]SettingDefinition{
    "config.current_project": {
        DataType:    "string",
        Category:    "config", 
        Description: "Active project in project@vault format",
        Validator:   validateProjectName,
    },
    "config.default_output_format": {
        DataType:     "string",
        Category:     "config",
        Description:  "Preferred output format (text, json, yaml)",
        DefaultValue: "text",
        Validator:    validateOutputFormat,
    },
    "meta.created_at": {
        DataType:    "string",
        Category:    "meta",
        Description: "Workspace creation timestamp (RFC3339)",
    },
    "meta.knox_version": {
        DataType:    "string", 
        Category:    "meta",
        Description: "Knox version used to create workspace",
    },
}
```

### Task 3: Workspace Settings Methods

**File:** `internal/workspace/workspace.go`

Add database settings methods to `Workspace` struct:

```go
// Core setting methods (low-level)
func (w *Workspace) GetSetting(key string) (string, error)
func (w *Workspace) SetSetting(key, value, category string) error  
func (w *Workspace) UnsetSetting(key string) error

// Convenience methods for configuration (user-changeable)
func (w *Workspace) GetConfig(key string) (string, error)
func (w *Workspace) SetConfig(key, value string) error
func (w *Workspace) UnsetConfig(key string) error

// Convenience methods for metadata (system-managed)
func (w *Workspace) GetMeta(key string) (string, error)
func (w *Workspace) SetMeta(key, value string) error

// Type-safe helper methods
func (w *Workspace) GetConfigBool(key string) (bool, error)
func (w *Workspace) SetConfigBool(key string, value bool) error

// Project-specific convenience methods
func (w *Workspace) GetCurrentProject() (string, error)
func (w *Workspace) SetCurrentProject(projectName string) error
func (w *Workspace) ClearCurrentProject() error
func (w *Workspace) IsDetached() bool
```

**Implementation Requirements:**
- Use database transactions for atomic operations
- Handle missing keys gracefully (return default from registry)
- Registry-based validation for all setting keys
- INSERT OR REPLACE pattern with automatic updated_at timestamp

### Task 4: Default Vault Auto-Linking

**File:** `internal/workspace/database/database.go`

Add method to auto-link default vault during workspace creation:

```go
// LinkDefaultVault automatically links default vault during workspace creation  
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
    // Use KNOX_ROOT env var if set, otherwise ~/.knox/vaults/default.db
}
```

**Requirements:**
- Respect `KNOX_ROOT` environment variable
- Create default vault if it doesn't exist
- Use "default" as alias
- Handle case where default vault already linked

### Task 5: Enhanced Workspace Creation

**File:** `internal/workspace/workspace.go`

Update `CreateWorkspace()` function to include settings and vault linking:

```go
func CreateWorkspace(path string) (*Workspace, error) {
    // Check if workspace already exists
    if ContainsDataDirectory(path) {
        return nil, errors.ErrWorkspaceExists.WithContext("path", path)
    }

    workspacePath := filepath.Join(path, constants.DataDirectoryName)

    // Create workspace directory
    if err := os.MkdirAll(workspacePath, 0700); err != nil {
        return nil, errs.Wrap(err, errors.CreateFailureCode, "failed to create workspace directory")
    }

    // Create database 
    db, err := database.CreateWorkspaceDatabase(workspacePath)
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

    // Initialize in detached state (no config.current_project setting)
    
    return workspace, nil
}
```

### Task 6: CLI Integration

**File:** `cmd/knox/internal/commands/init.go`

Implement `workspaceHandler()` function:

```go
func workspaceHandler() error {
    currentDir, err := os.Getwd()
    if err != nil {
        return fmt.Errorf("failed to get current directory: %w", err)
    }

    workspace, err := workspace.CreateWorkspace(currentDir)
    if err != nil {
        if errs.Is(err, errors.ErrWorkspaceExists) {
            return fmt.Errorf("workspace already exists in current directory")
        }
        return fmt.Errorf("failed to create workspace: %w", err)
    }
    defer workspace.Close()

    workspacePath := filepath.Join(currentDir, constants.DataDirectoryName)
    fmt.Printf("Initialized empty Knox workspace in %s\n", workspacePath)
    fmt.Println("Linked default vault 'default' (path: ~/.knox/vaults/default.db)")
    fmt.Println("Workspace is in detached state. Use 'knox switch <project>' to attach to a project.")
    
    return nil
}
```

**Requirements:**
- Handle errors gracefully with user-friendly messages
- Provide clear success feedback  
- Include guidance on next steps
- Import workspace package and constants

## CLI Workflow

```bash
# Initialize workspace
$ knox init
Initialized empty Knox workspace in /path/to/my-app/.knox-workspace
Linked default vault 'default' (path: ~/.knox/vaults/default.db)
Workspace is in detached state. Use 'knox switch <project>' to attach to a project.

# Create a project in default vault
$ knox new project backend-api

# Switch to the project  
$ knox switch backend-api

# Now can use secrets
$ knox set DB_URL postgres://localhost/mydb
$ knox get DB_URL
```

## Error Scenarios & Handling

1. **Workspace already exists**: Clear message, suggest using existing workspace
2. **Permission denied**: Clear message about directory permissions  
3. **Default vault creation fails**: Specific error about vault setup
4. **Settings initialization fails**: Error about workspace metadata setup
5. **Database creation fails**: Database-specific error message

## File Structure After Creation

```
/path/to/my-app/
├── .knox-workspace/       # Workspace marker directory
│   └── workspace.db       # Single SQLite database
├── src/
└── README.md
```

**Database contents after creation:**
- `linked_vaults` table: one entry for "default" vault
- `workspace_settings` table: `meta.knox_version`, `meta.created_at`
- No `config.current_project` setting (detached state)

## Implementation Order

1. **Database schema update** (rename table, add columns)
2. **Settings registry system** (WorkspaceSettings map and validation)
3. **Database settings methods** (GetSetting, SetSetting, convenience methods)
4. **Default vault auto-linking** (LinkDefaultVault method)
5. **Enhanced workspace creation** (update CreateWorkspace function)
6. **CLI integration** (implement workspaceHandler)
7. **Testing and validation**

## Design Benefits

1. **Unified settings model**: Current project is just another workspace setting
2. **Single source of truth**: All workspace state in one database table
3. **Registry-based validation**: Type safety with debuggable TEXT storage
4. **Atomic operations**: Database transactions prevent inconsistent state
5. **Extensible**: Easy to add more settings using same pattern
6. **Self-contained**: Entire workspace state in single database file

## Testing Strategy

1. **Unit tests** for each new function
2. **Integration test** for full workspace creation flow  
3. **Error scenario tests** (permissions, existing workspace)
4. **CLI test** with mock filesystem

## Out of Scope

- Workspace migration/import functionality
- Multiple workspace support in single directory
- Project creation during workspace init
- Workspace linking commands (`knox link`, `knox unlink`)
- Workspace status command (`knox status`)

**Note:** This spec focuses specifically on `knox init` implementation. Other workspace features like linking, status, and discovery are covered in separate specs.