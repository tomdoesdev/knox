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
4. What new features/capabilities should Knox v2 have that v1 lacks?
5. Should v2 maintain the same core concept (secret management) or expand scope?
6. How should Knox v2 handle different types of secrets/data?

### Storage & Backend
7. Should v2 stick with SQLite or consider other storage options?
8. How should v2 handle encryption and security differently than v1?
9. What about backup/sync/sharing capabilities?

### CLI/UX Design
10. How should the CLI interface change in v2?
11. What workflow improvements are needed over v1?
12. Should v2 have additional interfaces (web UI, API, etc.)?

### Technical Stack
13. Any changes to Go dependencies or frameworks?
14. How should v2 handle testing and development workflow?
15. What about deployment and distribution improvements?

## Knox v2 Core Design

**Workspace-Project Architecture:**
- **Workspace**: Directory marked with `.knox/` folder (like git repos)
- **Project**: Managed entity in main vault database with proper lifecycle
- **Project Switching**: Git-branch-like switching between projects within workspace
- **Shared Projects**: Multiple workspaces can share the same project

**Git-Style Commands:**
```bash
# Workspace & Project Management
knox init                    # Creates .knox/ workspace in CWD
knox new project <name>      # Creates new project in vault DB
knox switch <name>           # Switch active project in workspace
knox ls                      # List available projects
knox rm project <name>       # Remove project from vault (with safeguards)

# Secret Management (within current project context)
knox set <key> <value>       # Set secret in current project
knox get <key>               # Get secret from current project
knox ls secrets              # List secrets in current project
knox rm <key>                # Remove secret from current project

# Tagging (works on both projects and secrets via entity lookup)
knox tag add <entity> <tag>  # Add tag to project or secret
knox tag rm <entity> <tag>   # Remove tag from project or secret
knox tag ls <entity>         # List tags for project or secret
knox ls --tag <tag>          # Filter projects by tag
knox ls secrets --tag <tag>  # Filter secrets by tag
```

**Workspace Structure:**
```
/path/to/app/
├── .knox/
│   ├── config              # Workspace configuration
│   └── current-project     # Tracks active project
├── src/
└── README.md
```

**Benefits:**
- Clean separation of concerns (workspace ≠ project)
- Project lifecycle management in vault
- Multiple apps can share projects
- No orphaned secrets from deleted knox.json files
- Git-like workflow familiarity

## Technical Considerations

**Project Metadata Schema:**
```go
type Project struct {
    Name        string   // No spaces, CLI-friendly (e.g., "my-api")
    Description string   // Human-readable description
    Tags        []string // Optional: for categorization/search
}
```

**Security & Permissions:**
- Knox remains local development tool with no security guarantees
- Files created with sensible permissions (0700 for directories)
- User responsible for project access management
- No cross-project permission enforcement

**Workspace Discovery:**
- Like git: traverse up directory tree to find nearest `.knox/` folder
- Commands work from any subdirectory within workspace

**Project Context for Secrets:**
- `--project <name>` for specific project
- `-C` or `--current-project` for workspace's active project
- Example: `knox secrets add --project myapi DB_URL postgres://...`
- Example: `knox secrets add -C API_KEY abc123`

**Database Schema Updates:**
```sql
-- Core tables
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

-- Direct tagging tables (no entity lookup needed)
CREATE TABLE project_tags (
    project_id INTEGER NOT NULL,
    tag TEXT NOT NULL,
    PRIMARY KEY (project_id, tag),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE TABLE secret_tags (
    secret_id INTEGER NOT NULL,
    tag TEXT NOT NULL,
    PRIMARY KEY (secret_id, tag),
    FOREIGN KEY (secret_id) REFERENCES secrets(id) ON DELETE CASCADE
);
```

**Simplified Tagging:**
- Direct relationships between projects/secrets and tags
- No entity lookup table needed
- No name collision concerns (projects and secrets can share names)
- Cleaner CLI implementation

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

**Migration Strategy:**
- No automatic migration needed (single user)
- User will manually re-add secrets to new v2 projects

## Implementation Strategy

*To be defined after requirements gathering*

## Scope

**Knox v2.0 Core Tool:**
- Workspace management (`knox init`)
- Project lifecycle (`knox new project`, `knox switch`, `knox rm project`)
- Secret management (`knox set`, `knox get`, `knox ls secrets`, `knox rm`)
- Tagging system (`knox tag add/rm/ls`)
- Entity lookup and name collision prevention
- Multi-format output (text, JSON, YAML)
- Cross-platform workspace discovery

**Future Tools:**
- `knox-run-env`: Template processing and execution tool

## Out of Scope (for v2.0)

- Template processing and execution (delegated to future `knox-run-env` tool)
- Encryption/security beyond file permissions
- Remote secret storage or sync
- Web UI or API interfaces
- Backup/restore functionality
- Migration from v1 (manual re-entry)