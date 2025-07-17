<!--suppress HtmlDeprecatedAttribute -->
<div align="center">
  <img src="knox.svg" alt="Knox Logo" width="150" height="150">
</div>

# Knox

A secure local development secrets manager designed to prevent accidental commits of sensitive configuration to version control systems.

## Overview

Knox solves the problem of accidentally committing sensitive information (API keys, database credentials, tokens) to git repositories through `.env` files. It provides a project-based approach to managing secrets during development by storing them locally and injecting them into applications at runtime.

## Key Features

- **Prevent Secret Leaks** - Remove secrets from version-controlled `.env` files
- **Template-Based Injection** - Use Go templates in `.env` files with in-memory secret replacement
- **Project Isolation** - Each project has unique secret namespace within shared vault
- **Seamless Workflow** - Replace traditional `.env` workflows with minimal changes
- **Local Storage** - SQLite backend for fast, reliable local storage
- **Development Focused** - Designed for local development, not production

## Quick Start

### Installation

```bash
# From source (requires Go 1.24.5+)
git clone https://github.com/tomdoesdev/knox.git
cd knox
go build -o knox ./cmd/knox
```

### Initialize a Project

```bash
# In your project directory
knox init
```

This creates a `knox.json` configuration file and initializes your project's secret vault.

### Manage Secrets

```bash
# Store secrets
knox set DATABASE_URL "postgresql://localhost:5432/myapp"
knox set API_KEY "sk-1234567890abcdef"

# Retrieve secrets
knox get DATABASE_URL
knox get API_KEY

# Remove secrets
knox remove OLD_API_KEY

# Check project status
knox status
```

### Template-Based Environment Files

Create `.env.template` files using Go template syntax:

```bash
# .env.template
DATABASE_URL={{.Secret "DATABASE_URL"}}
API_KEY={{.Secret "API_KEY"}}
DEBUG={{.Secret "DEBUG_MODE"}}
PORT={{.Default "PORT" "3000"}}
```

### Run Applications with Injected Secrets

```bash
# Run your application with secrets injected
knox run --env .env.template npm start
knox run --env .env.template go run main.go
knox run --env .env.template python app.py
```

Knox processes the template in-memory, injects secrets, and runs your application with the resulting environment variables.

## How It Works

1. **Store Secrets Locally** - Secrets are stored in a local SQLite database
2. **Template Processing** - `.env` files are processed as Go templates
3. **In-Memory Injection** - Secrets are injected into templates in-memory (never written to disk)
4. **Process Execution** - Applications run with the constructed environment variables

## Architecture

### Storage
- **Backend**: SQLite database for local storage
- **Location**: `~/.local/share/knox/` (or custom `vault_path`)
- **Schema**: Single `vault` table with `(project_id, key)` unique constraint
- **Isolation**: Projects share vault file but are separated by `project_id`

### Template Engine
- **Engine**: Go's `text/template` package
- **Functions**: `{{.Secret "KEY"}}`, `{{.Env "KEY"}}`, `{{.Default "KEY" "VALUE"}}`
- **Processing**: In-memory execution (secrets never written to disk)
- **Parsing**: `github.com/hashicorp/go-envparse` for environment variable parsing

### Security Model
- **Threat Model**: Prevents accidental git commits of secrets
- **Local Only**: Secrets never leave developer machine
- **No Encryption**: Plaintext storage (development-focused)
- **Project Isolation**: Secrets separated by unique `project_id`

## Configuration

### Project Configuration (`knox.json`)
```json
{
  "project_id": "unique-project-identifier",
  "vault_path": "/path/to/vault/file"
}
```

### Environment Variables
- `LOG_LEVEL=debug` - Enable debug logging

## Commands

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

## Development Status

### ✅ Phase 1: Core Secret Management (Complete)
- [ ] Project initialization (`knox init`)
- [ ] Secret storage (`knox set/add`)
- [ ] Secret retrieval (`knox get`)
- [ ] Secret removal (`knox remove`)
- [ ] Project status (`knox status`)

### 🚧 Phase 2: Template Processing (In Progress)
- [ ] Go template engine integration
- [ ] Template function library
- [ ] Error handling for template parsing
- [ ] Template validation

### 📋 Phase 3: Process Execution (Planned)
- [ ] Command execution with environment injection
- [ ] Signal handling and forwarding
- [ ] Exit code preservation
- [ ] Timeout management

## Limitations

- **Development Only**: Not designed for production secret management
- **No Encryption**: Secrets stored in plaintext locally
- **No Sharing**: No built-in team secret sharing
- **Local Only**: No network operations or cloud integration

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes following existing patterns
4. Add tests for new functionality
5. Run tests (`go test ./...`)
6. Commit your changes
7. Push to the branch
8. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.