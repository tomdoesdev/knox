# Knox Application Specification

## Overview

Knox is a secure local development secrets manager designed to solve the problem of accidentally committing sensitive configuration to version control systems. It provides a project-based approach to managing secrets during development by storing them locally and injecting them into applications at runtime.

## Purpose & User Problem

### Primary Problem
Developers frequently commit sensitive information (API keys, database credentials, tokens) to git repositories through `.env` files, creating security vulnerabilities and compliance issues.

### Target Users
- Individual developers working on local development environments
- Development teams who need to share project structure without sharing actual secrets
- Projects that require environment-specific configuration during development

### User Goals
- Remove secrets from version-controlled `.env` files
- Maintain seamless development workflow with environment variables
- Prevent accidental exposure of sensitive data in git history
- Enable secure local development without complex setup

## Success Criteria

### Primary Success Metrics
- **Security**: Zero secrets committed to git repositories using Knox-managed projects
- **Usability**: Developers can replace traditional `.env` workflows with minimal learning curve
- **Reliability**: Consistent secret injection across different development processes
- **Performance**: Fast secret retrieval and process execution (< 100ms overhead)

### User Experience Goals
- Simple CLI commands for daily operations
- Seamless integration with existing development workflows
- Clear error messages and status feedback
- No impact on application performance

## Core Features

### 1. Project Management
- **Initialize Projects**: `knox init` creates project configuration and vault
- **Project Status**: `knox status` displays project health and secret counts
- **Project Isolation**: Each project has unique `project_id` for secret namespacing within shared vault

### 2. Secret Management
- **Store Secrets**: `knox set/add <key> <value>` stores secrets in local vault
- **Retrieve Secrets**: `knox get <key>` retrieves individual secrets
- **Remove Secrets**: `knox remove <key>` deletes secrets from vault
- **List Secrets**: View available secret keys (values never displayed)

### 3. Environment Variable Injection
- **Template Processing**: Parse `.env` files as Go text templates
- **Secret Injection**: Replace template placeholders with actual secrets in-memory
- **Process Execution**: `knox run <command>` executes applications with injected environment
- **Environment Inheritance**: Preserve existing environment variables

### 4. Development Workflow Integration
- **Command Execution**: Run any development command with injected secrets
- **Signal Handling**: Proper process management and signal forwarding
- **Exit Code Preservation**: Maintain application exit codes for CI/CD compatibility

## Technical Architecture

### Storage Layer
- **Backend**: SQLite database for fast, reliable local storage
- **Location**: User's local data directory (`~/.local/share/knox/`) or custom `vault_path`
- **Schema**: Single `vault` table with `(project_id, key)` unique constraint
- **Isolation**: Projects share vault file but are separated by `project_id`
- **Encryption**: No encryption (development-focused, not production)

### Template Engine
- **Engine**: Go's `text/template` package for `.env` file processing
- **Syntax**: Standard Go template syntax for secret references
- **Processing**: In-memory template execution (never writes secrets to disk)
- **Parsing**: `github.com/hashicorp/go-envparse` for environment variable parsing

### Process Management
- **Execution**: `os/exec` for spawning child processes
- **Signals**: Proper signal forwarding (SIGINT, SIGTERM, SIGQUIT)
- **Streams**: Direct stdin/stdout/stderr forwarding
- **Timeouts**: Configurable process timeouts

### CLI Framework
- **Library**: `github.com/urfave/cli/v3` for command structure
- **Logging**: Structured logging with `log/slog`
- **Error Handling**: Wrapped errors with context

## Configuration

### Project Configuration (`knox.json`)
```json
{
  "project_id": "unique-project-identifier",
  "vault_path": "/path/to/vault/file"
}
```

### Environment Templates (`.env.template`)
```bash
DATABASE_URL={{.Secret "DATABASE_URL"}}
API_KEY={{.Secret "API_KEY"}}
DEBUG={{.Secret "DEBUG_MODE"}}
```

### Runtime Configuration
- **Environment Variables**: `LOG_LEVEL=debug` for verbose logging
- **CLI Flags**: Command-specific options and parameters

## Security Model

### Threat Model
- **In Scope**: Prevention of accidental git commits of secrets
- **Out of Scope**: Protection against local system compromise
- **Assumption**: Developer workstation is trusted environment

### Security Measures
- **Local Storage**: Secrets never leave developer machine
- **Project Isolation**: Each project has separate secret namespace
- **Memory Processing**: Templates processed in-memory, never written to disk
- **Unique Constraints**: Prevent accidental secret overwrites

### Security Boundaries
- **Git Repository**: No secrets stored in version control
- **Process Isolation**: Each application runs with its own environment
- **Project Isolation**: Secrets separated by `project_id` within shared vault
- **User Isolation**: Secrets accessible only to user who created them

## API Design

### Command Interface
```bash
# Project lifecycle
knox init                    # Initialize project
knox status                  # Show project status

# Secret management
knox set KEY VALUE          # Store secret
knox get KEY                # Retrieve secret
knox remove KEY             # Delete secret

# Application execution
knox run [--env FILE] COMMAND [ARGS...]  # Run with injected secrets
```

### Template API
```go
// Available template functions
{{.Secret "KEY"}}           // Retrieve secret by key
{{.Env "KEY"}}             // Get environment variable
{{.Default "KEY" "VALUE"}} // Get with default value
```

## Implementation Phases

### Phase 1: Core Secret Management ✅
- [x] Project initialization (`knox init`)
- [x] Secret storage (`knox set/add`)
- [x] Secret retrieval (`knox get`)
- [x] Secret removal (`knox remove`)
- [x] Project status (`knox status`)

### Phase 2: Template Processing
- [ ] Go template engine integration
- [ ] Template function library
- [ ] Error handling for template parsing
- [ ] Template validation

### Phase 3: Process Execution
- [ ] Command execution with environment injection
- [ ] Signal handling and forwarding
- [ ] Exit code preservation
- [ ] Timeout management

### Phase 4: Developer Experience
- [ ] Enhanced error messages
- [ ] Configuration validation
- [ ] Performance optimization
- [ ] Documentation and examples

## Constraints & Limitations

### Technical Constraints
- **Go Version**: Requires Go 1.24.5 or later
- **Platform**: Cross-platform (Linux, macOS, Windows)
- **Dependencies**: Minimal external dependencies
- **Performance**: Single-threaded operation acceptable for development use

### Functional Limitations
- **No Encryption**: Secrets stored in plaintext locally
- **No Sharing**: No built-in team secret sharing
- **No Versioning**: No secret version history
- **No Audit**: No access logging or audit trails

### Development Constraints
- **Local Only**: No network operations
- **Single User**: No multi-user support
- **Development Focus**: Not suitable for production use

## Out of Scope

### Explicitly Excluded Features
- **Production Deployment**: Not designed for production secret management
- **Network Synchronization**: No cloud or team secret sharing
- **Advanced Encryption**: No encryption at rest or in transit
- **Audit Logging**: No access tracking or compliance features
- **Secret Rotation**: No automatic secret rotation
- **Integration**: No CI/CD or third-party service integration
- **Backup/Recovery**: No secret backup or disaster recovery

### Future Considerations
- Team secret sharing (separate product)
- Integration with production secret managers
- Enhanced security features for sensitive environments
- Cloud-based secret synchronization

## Testing Strategy

### Unit Testing
- Secret storage and retrieval operations
- Template processing and error handling
- Process execution and signal handling
- CLI command parsing and validation

### Integration Testing
- End-to-end workflow testing
- Cross-platform compatibility
- Performance benchmarking
- Error scenario testing

### Manual Testing
- Developer workflow validation
- CLI usability testing
- Error message clarity
- Documentation accuracy

## Success Metrics

### Quantitative Metrics
- **Performance**: < 100ms for secret operations
- **Reliability**: 99.9% uptime for local operations
- **Compatibility**: Support for major development platforms
- **Coverage**: > 90% test coverage for core functionality

### Qualitative Metrics
- **Developer Satisfaction**: Positive feedback on workflow integration
- **Security**: Zero reported secret leaks via git commits
- **Usability**: Developers can onboard without documentation
- **Maintenance**: Low bug report rate and quick resolution

---

*This specification serves as the authoritative reference for Knox development and should be updated as requirements evolve.*