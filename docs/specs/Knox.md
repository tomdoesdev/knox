# Knox v2 Specification

## Purpose & User Problem

Define the complete redesign of Knox to address lessons learned from v1 implementation and incorporate new requirements for improved architecture, usability, and functionality.

## Success Criteria

- [ ] Clear architectural improvements over v1
- [ ] Address identified design flaws from v1 
- [ ] Meet new functional requirements
- [ ] Maintain backward compatibility where feasible
- [ ] Improved developer experience and code maintainability

## Current State - v1 Analysis

Knox v1 has been implemented with vault file location system, SQLite backend, and basic CLI structure. 

**Major v1 Design Flaws Identified:**
- **Project-Vault Disconnect**: `knox init` creates `knox.json` with project ID in CWD, but no record in main vault database
- **Orphaned Data**: Deleting `knox.json` leaves secrets in vault with no project reference
- **Single Project Per Directory**: One projectId per folder prevents project switching
- **Limited Workspace Flexibility**: Cannot share projects across multiple application directories

**V1 Architecture Problems:**
- Tight coupling between filesystem location and project identity
- No project lifecycle management in vault
- Missing workspace concept
- No project switching capability

## Questions for Requirements Gathering

### Architecture & Design
1. What specific design choices from v1 do you regret? What didn't work well?
2. What architectural patterns would you prefer for v2? (e.g., plugin system, modular design, etc.)
3. How should Knox v2 handle configuration vs. v1?

### Core Functionality
1. What new features/capabilities should Knox v2 have that v1 lacks?
2. Should v2 maintain the same core concept (secret management) or expand scope?
3. How should Knox v2 handle different types of secrets/data?

### Storage & Backend
1. Should v2 stick with SQLite or consider other storage options?
2. How should v2 handle encryption and security differently than v1?
3. What about backup/sync/sharing capabilities?

### CLI/UX Design
1. How should the CLI interface change in v2?
2. What workflow improvements are needed over v1?
3. Should v2 have additional interfaces (web UI, API, etc.)?

### Technical Stack
1. Any changes to Go dependencies or frameworks?
2. How should v2 handle testing and development workflow?
3. What about deployment and distribution improvements?

## Knox v2 Core Design

**Workspace-Project Architecture:**
- **Workspace**: Directory marked with `.knox-workspace/` folder (like git repos)
- **Project**: Managed entity in main vault database with proper lifecycle
- **Project Switching**: Git-branch-like switching between projects within workspace
- **Shared Projects**: Multiple workspaces can share the same project

**Git-Style Commands:**
```bash
# Workspace & Project Management
knox init                    # Creates .knox-workspace/ workspace in CWD
knox new project <name>      # Creates new project in vault DB
knox new vault <name>        # Creates new vault with optional workspace linking
knox link <project>@<vault>  # Link project to workspace using @ notation
knox link <vault-path> --as <alias>  # Link vault to workspace
knox unlink <name>           # Remove vault or project from workspace
knox switch <name>           # Switch active project in workspace
knox status                  # Show current workspace status (active project, linked vaults/projects)
knox ls                      # List available projects
knox ls --linked             # Show linked vaults and projects in workspace
knox delete project <name>@<vault>  # Remove project from specific vault
knox delete vault <name>     # Remove vault (with confirmation)

# Workspace Configuration (Enhanced TEXT with Registry)
knox config <key> <value>    # Set workspace configuration (validated by registry)
knox config <key>            # Get configuration value
knox config --list           # List all workspace configuration
knox config --unset <key>    # Remove configuration setting
# Examples: knox config current_project myapi@staging
#          knox config default_output_format json

# Secret Management (within current project context)
knox set <key> <value>       # Set secret in current project
knox get <key>               # Get secret from current project
knox ls secrets              # List secrets in current project
knox rm <key>                # Remove secret from current project

```

**Workspace Structure:**
```
/path/to/my-app/           # Any directory can become a workspace
├── .knox-workspace/       # Workspace marker directory
│   └── workspace.db       # SQLite database (vaults/projects/settings - everything)
├── src/
└── README.md
```

**Benefits:**
- Clean separation of concerns (workspace ≠ project)
- Project lifecycle management in vault
- Multiple apps can share projects
- No orphaned secrets from deleted knox.json files
- Git-like workflow familiarity
- Enhanced TEXT design: debuggable storage with type-safe API
- Registry-based validation prevents invalid configurations
- Single source of truth for all workspace state

## Technical Considerations

**Project Metadata Schema:**
```go
type Project struct {
    Name        string   // No spaces, CLI-friendly (e.g., "my-api")
    Description string   // Human-readable description
}
```

**Security & Permissions:**
- Knox remains local development tool with no security guarantees
- Files created with sensible permissions (0700 for directories)
- User responsible for project access management
- No cross-project permission enforcement

**Workspace Discovery:**
- Like git: traverse up directory tree to find nearest `.knox-workspace/` folder
- Commands work from any subdirectory within workspace

**Project Context for Secrets:**
- `--project <name>` for specific project
- `-C` or `--current-project` for workspace's active project
- Example: `knox secrets add --project myapi DB_URL postgres://...`
- Example: `knox secrets add -C API_KEY abc123`

**Database Schema Updates:**

*Vault Database Schema:*
```sql
-- Core tables (stored in each vault)
CREATE TABLE projects (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,       -- No spaces, CLI-friendly
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE secrets (
    id INTEGER PRIMARY KEY,
    project_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (project_id, key),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

```

*Workspace Database Schema:*
```sql
-- Workspace database schema (stored in .knox/workspace.db)
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

CREATE TABLE workspace_settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    category TEXT NOT NULL, -- 'meta' or 'config'
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- Examples:
-- INSERT INTO workspace_settings VALUES ('meta.knox_version', 'v2', 'meta', '2024-01-15T10:30:00Z');
-- INSERT INTO workspace_settings VALUES ('config.current_project', 'myproject@myvault', 'config', '2024-01-15T10:30:00Z');
-- INSERT INTO workspace_settings VALUES ('config.auto_vault_discovery', 'true', 'config', '2024-01-15T10:30:00Z');
```

**Multi-Vault Support:**
- **Default vault**: `~/.knox/vaults/default.db` (standardized location)
- **Named vaults**: `~/.knox/vaults/<name>.db` for organization
- **Workspace linking**: Link multiple vaults to workspaces with aliases
- **@ notation**: Explicit vault specification (`knox delete project api@staging`)
- **Vault creation**: `knox new vault <name>` with optional `--path` and `--link-as`
- **Project targeting**: Each project exists in exactly one vault

**File Organization:**
```
~/.knox/
└── vaults/               # All vaults directory
    ├── default.db        # Default vault
    ├── staging.db
    ├── production.db
    └── team_backend.db

# Workspace directories (created per project)
/path/to/my-app/.knox-workspace/
└── workspace.db         # Links vaults + settings (Enhanced TEXT design - everything)
```

**Enhanced TEXT Configuration Design:**
- **Registry-based validation**: All workspace settings defined in code-based registry
- **Type-safe API**: Dedicated methods for bool, int, string types with automatic validation
- **Debuggable storage**: Settings stored as human-readable TEXT in database
- **Semantic organization**: `config.*` for user settings, `meta.*` for system metadata
- **Single table design**: All settings in `workspace_settings` table with category column
- **SQL-friendly**: Can query and filter settings directly with standard SQL
- **Self-documenting**: Settings registry serves as live documentation

**Settings Registry Examples:**
```go
// Built-in workspace settings registry
var WorkspaceSettings = map[string]SettingDefinition{
    "config.current_project": {
        DataType: "string", Category: "config",
        Description: "Active project name in project@vault format",
    },
    "config.default_output_format": {
        DataType: "string", Category: "config",
        DefaultValue: "text", Validator: validateOutputFormat,
    },
    "meta.created_at": {
        DataType: "string", Category: "meta",
        Description: "Workspace creation timestamp",
    },
}
```


**Output Formats:**
- Modular output system supporting multiple formats
- Built-in: human-readable text, JSON, YAML
- Extensible architecture for future formats (XML, CSV, etc.)
- Global `--output` or `-o` flag: `knox ls -o json`

**Tool Separation (Unix Philosophy):**
- **knox**: Core secret/project/workspace management only
- **knox-run-env**: Separate tool for template processing and execution
- Clean separation of concerns, independent evolution

**Error Handling:**
- Start with simple 0 (success) / 1 (error) exit codes
- Design error handling architecture to easily support specific exit codes in future
- Use structured error types internally for potential exit code mapping

**Command Equivalence:**
- `knox switch <project@vault>` ≡ `knox config current_project <project@vault>`
- Unified configuration model where project switching is just setting configuration
- Type-safe helpers: `workspace.SetConfigBool("auto_discovery", true)`
- String-based API: `workspace.SetConfig("current_project", "api@staging")`

**Migration Strategy:**
- No automatic migration needed (single user)
- User will manually re-add secrets to new v2 projects

## Implementation Strategy

*To be defined after requirements gathering*

## Scope

**Knox v2.0 Core Tool:**
- Workspace management (`knox init`, `knox config`, `knox status`)
- Project lifecycle (`knox new project`, `knox switch`, `knox delete project`)
- Vault lifecycle (`knox new vault`, `knox delete vault`)
- Secret management (`knox set`, `knox get`, `knox ls secrets`, `knox rm`)
- Workspace linking (`knox link`, `knox unlink`, `knox tidy`)
- Multi-format output (text, JSON, YAML)
- Cross-platform workspace discovery

**Future Features:**
- Tagging system for projects and secrets (`knox tag add/rm/ls`, filtering with `--tag`)
- Entity lookup and name collision prevention for tags
- Advanced search and filtering capabilities

**Future Tools:**
- `knox-run-env`: Template processing and execution tool

## Out of Scope (for v2.0)

- Template processing and execution (delegated to future `knox-run-env` tool)
- Encryption/security beyond file permissions
- Remote secret storage or sync
- Web UI or API interfaces
- Backup/restore functionality
- Migration from v1 (manual re-entry)