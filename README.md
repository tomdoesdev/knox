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
- **Template-Based Injection** - Use Go templates in `.env.template` files with in-memory secret replacement
- **Clean Environment Control** - Only template variables are passed to processes by default
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
knox secrets add DATABASE_URL "postgresql://localhost:5432/myapp"
knox secrets add API_KEY "sk-1234567890abcdef"

# Retrieve secrets
knox secrets get DATABASE_URL
knox secrets get API_KEY

# Remove secrets
knox secrets remove OLD_API_KEY

# List all secret keys (values never displayed)
knox secrets list

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
# Run with clean environment (only template variables)
knox run npm start                     # Uses .env.template by default
knox run --env production.env.template npm start  # Custom template file

# Run with inherited environment (adds template variables to current env)
knox run --inherit-env npm start       # Include current environment variables
knox run --inherit-env --allow-override npm start  # Allow template to override existing vars

# Run with timeout
knox run --timeout 30s npm start      # Kill process after 30 seconds
```

Knox processes the template in-memory, injects secrets, and runs your application with a controlled environment containing only the variables you specify.

## Environment Control

Knox provides precise control over what environment variables your applications receive:

### Clean Environment (Default)
By default, Knox starts with a clean environment and only includes variables from your template:

```bash
# With .env.template containing: API_KEY={{.Secret "API_KEY"}}
knox run env  # Only shows: API_KEY=your-secret-value
```

### Inherited Environment  
When you need access to system environment variables (like `PATH`), use `--inherit-env`:

```bash
# Include current environment + template variables
knox run --inherit-env npm start

# Or explicitly include PATH in your template:
# PATH={{.Env "PATH"}}
knox run npm start
```

### Template Functions
- `{{.Secret "KEY"}}` - Retrieve secret from Knox vault (fails if missing)
- `{{.Env "KEY"}}` - Get current environment variable
- `{{.Default "KEY" "fallback"}}` - Get secret with fallback if not found

## How It Works

1. **Store Secrets Locally** - Secrets are stored in a local SQLite database
2. **Template Processing** - `.env` files are processed as Go templates
3. **In-Memory Injection** - Secrets are injected into templates in-memory (never written to disk)
4. **Process Execution** - Applications run with the constructed environment variables

## Architecture

### Storage
- **Backend**: SQLite database for local storage
- **Location**: `~/.local/share/knox/` (or custom `vault_path`)
- **Schema**: Single `secrets` table with `(project_id, key)` unique constraint
- **Isolation**: Projects share vault file but are separated by `project_id`
- **Security**: Secure delete (FAST), WAL journal mode, foreign key constraints

### Template Engine
- **Engine**: Go's `text/template` package
- **Functions**: `{{.Secret "KEY"}}`, `{{.Env "KEY"}}`, `{{.Default "KEY" "VALUE"}}`
- **Processing**: In-memory execution (secrets never written to disk)
- **Parsing**: `github.com/hashicorp/go-envparse` for environment variable parsing
- **Environment**: Clean environment by default (only template variables)

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
knox init                                    # Initialize project

# Secret management
knox secrets add KEY VALUE                   # Store secret
knox secrets get KEY                         # Retrieve secret
knox secrets remove KEY                      # Delete secret
knox secrets list                            # List all secret keys

# Application execution
knox run [OPTIONS] COMMAND [ARGS...]         # Run with injected secrets

# Run command options:
--env FILE                                   # Template file path (default: .env.template)
--inherit-env                                # Inherit current environment variables
--allow-override                             # Allow template to override existing vars
--timeout DURATION                           # Process timeout

# Secrets command options:
--force                                      # Force vault creation in version-controlled directories
```

## Development Status

### ✅ Phase 1: Core Secret Management (Complete)
- [x] Project initialization (`knox init`)
- [x] Secret storage (`knox secrets add`)
- [x] Secret retrieval (`knox secrets get`)
- [x] Secret removal (`knox secrets remove`)
- [x] Secret listing (`knox secrets list`)

### ✅ Phase 2: Template Processing (Complete)
- [x] Go template engine integration
- [x] Template function library (`Secret`, `Env`, `Default`)
- [x] Error handling for template parsing
- [x] Template validation
- [x] Environment variable parsing with `go-envparse`

### ✅ Phase 3: Process Execution (Complete)
- [x] Command execution with environment injection
- [x] Signal handling and forwarding
- [x] Exit code preservation
- [x] Timeout management
- [x] Clean environment control
- [x] Environment inheritance options

## Best Practices

### Template Files
- ✅ **Commit** `.env.template` files to version control
- ❌ **Never commit** actual `.env` files with secrets
- ✅ Use descriptive secret names: `{{.Secret "DATABASE_URL"}}` not `{{.Secret "DB"}}`
- ✅ Include `PATH={{.Env "PATH"}}` if your application needs system commands

### Environment Control
- ✅ Use clean environment by default for better security
- ✅ Only add `--inherit-env` when you need system environment variables
- ✅ Be explicit about what environment variables your app receives

### Secret Management
- ✅ Use namespaced keys: `aws:api_key`, `db:password`
- ✅ Rotate secrets regularly in your vault
- ✅ Use `knox secrets list` to audit what secrets exist
- ✅ Knox prevents vault creation in version-controlled directories (use `--force` to override)

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