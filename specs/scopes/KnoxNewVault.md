# Knox New Vault Mini-Spec

## Purpose & User Problem

Enable users to create new named vaults for organizing secrets across different environments, teams, or projects. With Knox v2's multi-vault support, users need a way to create additional vaults beyond the default `~/.knox/vault.db`.

## Success Criteria

- [ ] `knox new vault <name>` creates a new vault at a specified location
- [ ] Vault names are validated and normalized
- [ ] Support for custom vault locations via `--path` flag
- [ ] Integration with workspace linking workflow
- [ ] Clear error messages for naming conflicts and permission issues

## Command Syntax

```bash
# Create vault in default location with name
knox new vault <vault-name>

# Create vault at custom path
knox new vault <vault-name> --path /custom/path/

# Create vault and immediately link to current workspace
knox new vault <vault-name> --link-as <alias>
```

## Vault Naming & Location Strategy

**Default location pattern:**
- `~/.knox/vaults/<vault-name>.db` or `$KNOX_ROOT/.knox/vaults/<vault-name>.db`
- Creates vaults directory if it doesn't exist

**Naming rules:**
- Same validation as project names: `^[a-z0-9_-]+$`
- Automatically normalize: lowercase, spaces→underscores
- No path separators allowed in vault names

**Examples:**
```bash
# Creates ~/.knox/vaults/staging.db
knox new vault staging

# Creates ~/.knox/vaults/team_backend.db  
knox new vault "Team Backend"

# Creates /custom/path/production.db
knox new vault production --path /custom/path/
```

## Implementation Requirements

### 1. Vault Creation Function
```go
func CreateVault(name string, path string) (*Vault, error)
```
- Validate vault name using existing project validation logic
- Create vault directory if needed (0700 permissions)
- Initialize SQLite database with Knox v2 schema
- Return opened vault instance

### 2. Default Path Resolution
```go
func GetDefaultVaultPath(name string) (string, error)
```
- Respect `KNOX_ROOT` environment variable
- Create `~/.knox/vaults/` directory structure
- Return absolute path: `{root}/.knox/vaults/{name}.db`

### 3. Conflict Detection
- Check if vault file already exists at target path
- Validate vault name doesn't conflict with existing vaults in default location
- Option to `--force` overwrite existing vault (with confirmation)

## CLI Integration

### Command Structure
```go
// cmd/v2/knox/internal/commands/new.go
func NewVaultCommand() *cli.Command {
    return &cli.Command{
        Name: "vault",
        Usage: "Create a new vault",
        ArgsUsage: "<vault-name>",
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name: "path",
                Usage: "Custom path for vault (defaults to ~/.knox/vaults/)",
            },
            &cli.StringFlag{
                Name: "link-as",
                Usage: "Automatically link vault to current workspace with alias",
            },
        },
        Action: createVaultAction,
    }
}
```

### Workflow Examples
```bash
# Basic vault creation
$ knox new vault staging
Created vault 'staging' at ~/.knox/vaults/staging.db

# Custom location
$ knox new vault production --path /team-shared/
Created vault 'production' at /team-shared/production.db

# Create and link to workspace
$ knox new vault staging --link-as staging
Created vault 'staging' at ~/.knox/vaults/staging.db
Linked vault 'staging' to current workspace

# Handle conflicts
$ knox new vault staging
Error: Vault 'staging' already exists at ~/.knox/vaults/staging.db
Use 'knox delete vault staging' first or choose a different name
```

## Error Scenarios

- **Invalid vault name**: "Vault name 'My/Vault' is invalid. Use only letters, numbers, underscores, and hyphens."
- **Permission denied**: "Cannot create vault directory ~/.knox/vaults/ - check permissions"
- **Already exists**: "Vault 'staging' already exists at ~/.knox/vaults/staging.db"
- **Not in workspace** (when using `--link-as`): "Cannot link vault - not in a Knox workspace"
- **Invalid path**: "Cannot create vault at '/invalid/path/' - directory doesn't exist"

## Delete Commands Integration

### Knox Delete Vault
```bash
knox delete vault <vault-name>
```
- Deletes vault file from default location (`~/.knox/vaults/<name>.db`)
- Requires confirmation prompt for safety
- Removes vault from any workspace links automatically (via foreign key constraints)

### Knox Delete Project  
```bash
knox delete project <project-name>@<vault-name>
```
- Uses @ notation to explicitly specify vault
- Handles project names with `/` correctly
- Removes project and all associated secrets/tags

**Examples:**
```bash
# Delete vault (with confirmation)
$ knox delete vault staging
This will permanently delete vault 'staging' and all its projects and secrets.
Continue? (y/N): y
Deleted vault 'staging'

# Delete specific project from specific vault
$ knox delete project backend/api@staging
Deleted project 'backend/api' from vault 'staging'

# Error handling
$ knox delete project backend/api@nonexistent
Error: Vault 'nonexistent' not found
```

## Integration with Workspace Linking

When `--link-as` flag is used:
1. Create the vault successfully
2. Verify current directory is in a workspace
3. Add vault to workspace's `linked_vaults` table
4. Show success message with workspace context

## File Organization Impact

**Current structure:**
```
~/.knox/
└── vault.db          # Default vault
```

**New structure:**
```
~/.knox/
├── vault.db          # Default vault (backward compatibility)
└── vaults/           # Named vaults directory
    ├── staging.db
    ├── production.db
    └── team_backend.db
```

## Implementation Priority

1. **Core vault creation logic** - vault name validation, path resolution
2. **CLI command implementation** - flags, validation, error handling  
3. **Workspace integration** - automatic linking with `--link-as`
4. **Delete commands** - vault and project deletion with @ notation

## Out of Scope

- Vault migration/copying between locations
- Vault templates or initialization with default projects
- Vault listing (covered by `knox ls --vaults` in main spec)

## Notes

- Reuse existing project name validation logic for vault names
- Follow same error handling patterns as vault module
- Maintain backward compatibility with default vault location
- Consider vault discovery for future `knox ls --vaults` command