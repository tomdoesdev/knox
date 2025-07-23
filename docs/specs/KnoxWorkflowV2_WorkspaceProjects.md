# Knox Workflow V2: Workspace Projects

## Purpose
Redesign Knox's secret management workflow from vault-centric projects to workspace-centric contexts that can compose secrets from multiple vaults.

## Success Criteria
- [ ] Users can create workspace projects that reference secrets from any linked vault
- [ ] Projects are defined as JSON files in the workspace, not tied to vault structure
- [ ] Users can switch between workspace projects seamlessly
- [ ] Secret resolution works across multiple vaults within a single project
- [ ] Migration path exists from current vault-project model

## Current State Problems
1. **Vault Lock-in**: Users can only access secrets from one vault project at a time
2. **Rigid Organization**: Secret organization forced by vault structure, not team needs
3. **Complex Data Model**: `linked_projects` table creates unnecessary coupling
4. **Limited Composition**: Cannot combine secrets from multiple vaults/projects

## Proposed Solution: Workspace Projects

### Core Concept
Replace vault projects with workspace projects - JSON configuration files that map logical secret names to vault paths.

### Architecture Overview
```
Workspace
├── .knox-workspace/
│   ├── workspace.db          # Only linked_vaults + workspace_settings
│   └── projects/
│       ├── staging.json      # Project: staging environment
│       ├── production.json   # Project: production environment
│       └── development.json  # Project: development environment
└── (user files)

Project File Format:
{
  "name": "staging",
  "description": "Staging environment secrets",
  "secret_map": {
    "API_KEY": "api-keys/staging@vault-prod",
    "DB_PASSWORD": "database/staging-db@vault-shared",
    "REDIS_URL": "cache/redis-staging@vault-infra"
  }
}
```

## Technical Considerations

### Database Schema Changes
**Remove:**
- `linked_projects` table (no longer needed)

**Keep:**
- `linked_vaults` table (still need vault registry)
- `workspace_settings` table (add current_project setting)

### Project File Structure
```json
{
  "name": "project-name",
  "description": "Human readable description",
  "created_at": "2025-07-23T15:08:28Z",
  "secret_map": {
    "LOGICAL_NAME": "secret-path@vault-alias",
    "API_KEY": "services/api/key@production-vault",
    "DATABASE_URL": "databases/main@shared-infra"
  }
}
```

### Secret Path Format
`<secret-path>@<vault-alias>`
- `secret-path`: Path to secret within the vault
- `vault-alias`: Alias of linked vault (from `linked_vaults` table)

### CLI Commands
```bash
# Project management
knox project create <name>              # Create new project file
knox project list                       # List available projects
knox project delete <name>              # Remove project file
knox project edit <name>                # Open project file in editor

# Project usage
knox use <project-name>                 # Switch to project
knox status                            # Show current project + secret status

# Secret discovery
knox list secrets --all                 # List all secret keys grouped by vault
knox list secrets <vault-alias>         # List secret keys for specific vault

# Secret management within project
knox project add-secret <logical-name> <vault-path@vault-alias>
knox project remove-secret <logical-name>
knox project list-secrets              # Show current project's secret map
```

## Implementation Plan

### Phase 1: Project File System
1. Create `.knox-workspace/projects/` directory structure
2. Implement project file I/O operations (JSON read/write)
3. Add project validation (vault exists, secret paths format)
4. Implement basic project CRUD operations

### Phase 2: Database Migration
1. Add migration to remove `linked_projects` table
2. Add `current_project` to workspace_settings
3. Update workspace initialization

### Phase 3: CLI Integration
1. Implement `knox project` commands
2. Update `knox use` to work with projects
3. Update `knox status` to show project info

### Phase 4: Secret Resolution
1. Implement multi-vault secret fetching
2. Add secret validation/existence checking
3. Environment variable export functionality

## Out of Scope
- Automatic migration of existing vault projects to workspace projects
- Secret caching/performance optimization
- Secret encryption at rest in project files
- Project templates or scaffolding
- Project sharing between workspaces

## Open Questions
1. **Validation**: Should we validate secret paths exist when switching projects?
2. **Project Creation UX**: Should `knox project create` open an editor, use interactive prompts, or create empty template?
3. **Secret Resolution Error Handling**: How should we handle missing vaults or secret paths when using a project?

## Migration Strategy
- Keep current vault project functionality during transition
- Add deprecation warnings to vault project commands
- Provide migration tool to convert vault projects to workspace projects
- Remove vault project support in future major version