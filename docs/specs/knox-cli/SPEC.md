# Knox CLI Specification

## Overview

Knox is a local secrets management system for development environments, similar to .NET's user-secrets tool but written in Go and configuration-format agnostic. Knox allows developers to securely store and manage secrets outside of their project repositories while maintaining easy access during development.

## Purpose & User Problem

Developers need a way to:
- Store sensitive configuration (API keys, database credentials, etc.) outside of source control
- Access secrets easily during local development
- Manage different sets of secrets for different environments (dev, staging, testing)
- Share secret management patterns across different programming languages and frameworks

## Architecture

Knox operates on a three-tier architecture:

### 1. Workspace
- Represents the root directory of a development environment
- Contains a `.knox-workspace/` directory with configuration and state
- Has a unique GUID identifier
- Can contain multiple projects

### 2. Projects  
- Logical groupings of secrets (e.g., "dev", "staging", "testing")
- Each project has its own vault for storing secrets
- Projects are defined in JSON files within the workspace
- Each project (except default) has a unique GUID that identifies its vault

### 3. Vaults
- Storage containers for secrets
- Currently implemented as JSON files (extensible for future SQLite, HTTP server, etc.)
- Located at `~/knox/vaults/<workspace-guid>/`
- One vault per project

## File Structure

```
Workspace Root/
├── .knox-workspace/
│   ├── workspace.toml          # Workspace configuration + GUID
│   ├── state.json             # Application state (active project, etc.)
│   └── projects/
│       ├── default.json       # Default project configuration
│       ├── dev.json          # Custom project configurations
│       └── staging.json
│
~/knox/vaults/<workspace-guid>/
├── default.json              # Default project vault (literal name)
├── <dev-project-guid>.json   # Project-specific vaults (GUID names)
└── <staging-project-guid>.json
```

## User Stories

### US1: Initialize Workspace
**As a** developer  
**I want** to initialize a Knox workspace in my project directory  
**So that** I can start managing secrets for this project

**Acceptance Criteria:**
- Given I'm in any directory
- When I run `knox init`
- Then a `.knox-workspace/` directory is created
- And `workspace.toml` is created with a unique workspace GUID
- And `state.json` is created with default project as active
- And `projects/default.json` is created
- And `~/knox/vaults/<workspace-guid>/default.json` vault is created

### US2: Create Projects
**As a** developer  
**I want** to create different projects for different environments  
**So that** I can separate secrets by context (dev, staging, etc.)

**Acceptance Criteria:**
- Given I'm in a Knox workspace
- When I run `knox new project <name>`
- Then a `projects/<name>.json` file is created with a unique project GUID
- And a vault at `~/knox/vaults/<workspace-guid>/<project-guid>.json` is created
- And the project configuration includes the vault identifier

### US3: Switch Between Projects
**As a** developer  
**I want** to switch between different projects  
**So that** I can work with different sets of secrets

**Acceptance Criteria:**
- Given I have multiple projects in my workspace
- When I run `knox use <project-name>`
- Then the active project in `state.json` is updated
- And subsequent vault operations target the active project's vault

### US4: Manage Secrets
**As a** developer  
**I want** to set, get, and delete secrets  
**So that** I can manage my configuration values

**Acceptance Criteria:**
- Given I'm in a Knox workspace with an active project
- When I run `knox vault set KEY=VALUE KEY2=VALUE2`
- Then both secrets are stored in the active project's vault
- When I run `knox vault set KEY=VALUE` and KEY already exists
- Then an error is displayed and the secret is not overwritten
- When I run `knox vault set --replace KEY=NEWVALUE` and KEY already exists
- Then the existing secret is replaced with the new value
- When I run `knox vault get KEY` 
- Then the secret value is displayed
- When I run `knox vault del KEY`
- Then the secret is removed from the vault

### US5: Export Secrets via Templates
**As a** developer  
**I want** to process template files with my secrets  
**So that** I can generate configuration files with actual secret values

**Acceptance Criteria:**
- Given I have a template file `.env` containing `API_KEY={{Secret "API_KEY"}}`
- And I have a secret `API_KEY=secret123` in my active project
- When I run `knox export ./.env`
- Then the processed template is output to stdout with `API_KEY=secret123`
- When I run `knox export ./.env ./output/` and `./output/` is outside the Knox workspace
- Then a file `.env` is created in `./output/` with the processed content
- When I run `knox export ./.env ./workspace-subdir/` and `./workspace-subdir/` is within the Knox workspace
- Then an error is displayed warning about exporting to a potentially version-controlled location
- When I run `knox export --force ./.env ./workspace-subdir/`
- Then the export proceeds despite being within the workspace
- When I run `knox export` without a template file argument
- Then an error is displayed requiring a template file path
- Template files use Go text/template syntax with `{{Secret "KEY_NAME"}}` function

### US6: List Resources
**As a** developer  
**I want** to list available projects and vaults  
**So that** I can understand my workspace structure

**Acceptance Criteria:**
- Given I'm in a Knox workspace
- When I run `knox list projects`
- Then all available projects are displayed
- When I run `knox list vaults`
- Then all vaults are displayed with their type (json, sqlite, etc.) and associated projects

### US7: Workspace Discovery
**As a** developer  
**I want** Knox to find my workspace automatically  
**So that** I can run commands from subdirectories

**Acceptance Criteria:**
- Given I'm in a subdirectory of a Knox workspace
- When I run any Knox command
- Then Knox searches upward for `.knox-workspace/` directory
- And if I'm in a subdirectory of `~`, search only up to `~`
- And if I'm outside `~`, don't search upward
- And if no workspace is found, display appropriate error

### US8: Workspace Status
**As a** developer  
**I want** to see my current workspace status  
**So that** I can understand which workspace and project I'm working with

**Acceptance Criteria:**
- Given I'm in a Knox workspace
- When I run `knox status`
- Then I see the workspace GUID, workspace path, and current active project
- And I see the number of available projects
- And I see the vault location for the active project

## Commands

### Core Commands
- `knox init` - Initialize a new workspace
- `knox new project <name>` - Create a new project
- `knox use <project>` - Switch active project
- `knox status` - Display workspace status and current active project

### Vault Management
- `knox vault set KEY=VALUE [KEY2=VALUE2 ...]` - Set one or more secrets (fails if secret exists unless --replace/-r flag is used)
- `knox vault set --replace/-r KEY=VALUE [KEY2=VALUE2 ...]` - Set/replace one or more secrets
- `knox vault get KEY` - Get a secret value  
- `knox vault del KEY` - Delete a secret

### Export/Import
- `knox export <template-file> [output-directory]` - Process template file with secrets and output to directory or stdout
- `knox export --force/-f <template-file> [output-directory]` - Force export even to potentially unsafe locations (within workspace)

### Listing
- `knox list projects` - List all projects
- `knox list vaults` - List all vaults with metadata

## Technical Considerations

### Storage
- Vaults are stored at `~/knox/vaults/<workspace-guid>/`
- Default vault uses literal name "default.json"
- Other vaults use project GUID as filename
- Secrets are stored unencrypted (development use only)
- Storage system is extensible for future vault types

### Configuration Files
- `workspace.toml` - TOML format for workspace configuration
- `state.json` - JSON format for runtime application state
- `projects/*.json` - JSON format for project configurations

### Template Processing
- Uses Go `text/template` package for processing template files
- Provides `Secret` function: `{{Secret "KEY_NAME"}}` to access vault secrets
- Template files can be any text format (.env, .yaml, .json, .config, etc.)
- Output preserves original filename when writing to directory
- Missing secrets in templates should result in clear error messages

### Export Safety
- By default, prevents exporting processed templates to locations within Knox workspace
- Assumes workspace directories are version-controlled and should not contain secrets
- `--force/-f` flag allows overriding safety check for intentional exports within workspace
- Detects workspace boundaries by searching for `.knox-workspace` directory

### Workspace Discovery
- Search upward from current directory
- Limit search to user home directory if within it
- Stop at filesystem boundaries for security

### Error Handling
- Graceful handling when outside workspace
- Clear error messages for missing projects/secrets
- Validation of project names and secret keys

## Security Considerations

- Secrets stored unencrypted (development environment assumption)
- Vaults stored outside project directories to prevent accidental commits
- No network communication (local-only tool)
- File permissions should restrict access to user only

## Out of Scope

### V1 Release
- Encryption of stored secrets
- Remote vault synchronization
- Team/shared secret management
- Integration with external secret management systems
- Automatic secret rotation
- Audit logging

### Future Considerations
- SQLite vault backend
- HTTP server vault backend
- Secret encryption options
- Team collaboration features
- Integration with CI/CD systems
- Git/VCS detection for export safety (beyond just Knox workspace detection)

## Success Criteria

1. Users can initialize workspaces in any directory
2. Projects can be created and managed independently
3. Secrets can be stored, retrieved, and deleted reliably
4. Multiple output formats are supported for integration
5. Workspace discovery works intuitively
6. No secrets are accidentally committed to version control
7. Tool works consistently across macOS, Linux, and Windows

## Constraints

- Go programming language
- Local filesystem storage only (V1)
- Single-user focused
- Development environment targeted
- No external dependencies for core functionality