# Knox Workspace Mini-Spec

## Purpose & User Problem

Define how Knox v2 workspaces function - the local directory context that tracks which project is currently active, similar to how git repositories work.

## Success Criteria

- [ ] `knox init` creates workspace in current directory
- [ ] Knox commands work from any subdirectory within workspace (like git)
- [ ] Workspace tracks current active project
- [ ] Clean error messages when not in workspace
- [ ] Workspace discovery traverses up directory tree

## Workspace Structure

```
/path/to/my-app/           # Any directory can become a workspace
├── .knox-workspace/       # Workspace marker directory
│   ├── workspace.db       # SQLite database for linked vaults/projects
│   ├── config             # Minimal JSON config (version, current project)
│   └── current-project    # Active project name (plain text)
├── src/
└── README.md
```

## Core Functions Needed

### 1. Workspace Discovery
```go
func FindWorkspace() (*Workspace, error)
```
- Start from current directory
- Traverse up directory tree looking for `.knox-workspace/` folder
- Return workspace details or error if not found
- Similar to `git rev-parse --show-toplevel`

### 2. Workspace Creation
```go
func CreateWorkspace(path string) (*Workspace, error)
```
- Create `.knox-workspace/` directory in specified path
- Create `workspace.db` SQLite database with schema
- Auto-link 'default' vault: `~/.knox/vault.db` or `$KNOX_ROOT/.knox/vault.db`
- Create minimal JSON config with workspace version
- Initialize with no active project (detached state)
- Return workspace instance

### 3. Workspace Validation
```go
func IsWorkspace(path string) bool
```
- Check if directory contains `.knox-workspace/` folder
- Validate workspace structure

### 4. Project Context
```go
func (w *Workspace) GetCurrentProject() (project.Name, error)
func (w *Workspace) SetCurrentProject(name project.Name) error
```
- Read/write current project from workspace config
- Handle case where no project is set

## Workspace Configuration

**Database: `.knox-workspace/workspace.db` (SQLite schema)**
```sql
CREATE TABLE linked_vaults (
    id INTEGER PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,
    path TEXT NOT NULL UNIQUE,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE linked_projects (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    vault_id INTEGER NOT NULL,
    project_name TEXT, -- actual name in vault if different
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vault_id) REFERENCES linked_vaults(id) ON DELETE CASCADE
);

CREATE TABLE workspace_meta (
    key TEXT PRIMARY KEY,
    value TEXT
);
-- Example: INSERT INTO workspace_meta VALUES ('version', 'v2'), ('default_project', 'backend-api');
```

**File: `.knox-workspace/config` (minimal JSON, extensible design)**
```json
{
  "workspace_version": "v2"
}
```

**File: `.knox-workspace/current-project` (plain text, git-style)**
```
my-api
```

**Implementation notes:**
- **SQLite database**: Stores linked vaults/projects with relational integrity
- **Stable integer keys**: Enables easy vault alias renaming without cascade updates
- **JSON config**: Minimal metadata, extensible for future workspace settings
- **Current project**: Plain text file for fast atomic updates (git HEAD pattern)
- **Foreign key constraints**: Automatic cleanup of orphaned project links

## Workspace Validation

**Validation Strategy (middle ground):**
1. Check `.knox-workspace/` directory exists
2. Validate `workspace.db` contains valid schema
3. If `current-project` file exists:
   - Query workspace database for linked project
   - Find corresponding vault via foreign key relationship
   - Validate project exists in that specific vault
4. If project not found or vault unreachable: warn user and enter detached state
5. Handle missing/corrupted files and database gracefully

## Error Scenarios

- **Not in workspace**: Clear message "Not in a Knox workspace. Run 'knox init' to create one."
- **Invalid current project**: Warn and unset: "Project 'old-project' no longer exists. Current project unset."
- **Permission errors**: Cannot create/read `.knox-workspace/` directory
- **Detached workspace state**: When no current-project file exists:
  - Warn: "Workspace is in detached state. Use 'knox switch <project>' to attach to a project."
  - Commands requiring project context return error until user switches

## CLI Integration

**Auto-discovery behavior:**
- All commands auto-discover workspace from current directory (git-style)
- Global `--workspace <path>` flag to override workspace location
- Project context flags where applicable: `--project <name>` to override current project

### Workspace Linking Commands:
- `knox workspace link <vault-path> --as <alias>` - Link vault to workspace (auto-detected by filepath)
- `knox workspace link <project-name> --from <vault-alias>` - Link project from vault (auto-detected)
- `knox workspace link <project-name> --from <vault-alias> --as <alias>` - Link project with custom name
- `knox workspace unlink <name>` - Remove vault or project from workspace
- `knox ls --linked` - Show linked vaults and projects
- `knox tidy` - Clean up workspace (remove deleted/unreachable vaults and projects)

### Commands that require workspace:
- `knox switch <project>` (only works with linked projects)
- `knox set <key> <value>` (uses current project, or `--project <name>`)
- `knox get <key>` (uses current project, or `--project <name>`)
- `knox tag add <entity> <tag>` (current project context for secrets)

### Commands that work without workspace:
- `knox init` (creates workspace)
- `knox ls` (lists all projects from vault)
- `knox new project <name>` (creates project in vault)

### CLI Workflow Examples:
```bash
# Setup workspace (auto-links default vault)
$ knox init
Initialized empty Knox workspace in /path/to/my-app/.knox-workspace
Linked default vault 'default' (path: ~/.knox/vault.db)

# Link additional vaults
$ knox workspace link /team-staging/.knox/vault.db --as staging
Linked vault 'staging' to workspace (path: /team-staging/.knox/vault.db)

# Link projects from vaults
$ knox workspace link backend-api --from default
Linked project 'backend-api' from vault 'default'

$ knox workspace link backend-api --from staging --as staging-backend
Linked project 'staging-backend' (backend-api) from vault 'staging'

# Switch between linked projects only
$ knox switch backend-api
Switched to project 'backend-api' (vault: default)

$ knox switch staging-backend
Switched to project 'staging-backend' (vault: staging)

$ knox switch frontend-web
Error: Project 'frontend-web' is not linked to this workspace.

# Use current project
$ knox set DB_URL postgres://...

# Override current project (must be linked)
$ knox set --project staging-backend DB_URL postgres://...

# Workspace maintenance  
$ knox tidy
Warning: Vault '/old/path/vault.db' unreachable for alias 'old-staging'
Removed 1 unreachable vault link
Removed 2 orphaned project links
Workspace cleaned up successfully

# Detached state example
$ knox set DB_URL postgres://...
Error: Workspace is in detached state. Use 'knox switch <project>' to attach to a project.

$ knox switch backend-api
Switched to project 'backend-api' (vault: default)
```

## Implementation Priority

1. **Core workspace detection and creation**
2. **Project context management** 
3. **Integration with existing CLI commands**
4. **Error handling and user guidance**

## Out of Scope

- Workspace migration/import
- Multiple workspace support in single directory