<!--suppress HtmlDeprecatedAttribute -->
<div align="center">
  <img src="knox.svg" alt="Knox Logo" width="150" height="150">
</div>

# Knox

A secure local development secrets manager for modern development workflows.

## Overview

Knox provides a secure, project-based approach to managing sensitive configuration during development. Each project maintains its own encrypted vault of secrets, ensuring sensitive data never leaves your local environment. 

Knox also provides a means to execute other applications with environment variables from a .env file. Knox will parse the .env file and replace any
template calls to Secret with the appropriate, unencrypted secret.


## Features

- **Project-Based Isolation** - Each project has its own secrets vault
- **SQLite Backend** - Fast, reliable local storage with ACID compliance
- **Unique Key Constraints** - Prevents accidental duplicate secrets
- **Clean CLI Interface** - Simple commands for everyday operations
- **JSON Configuration** - Human-readable project metadata

## Installation

### From Source

```bash
git clone https://github.com/tomdoesdev/knox.git
cd knox
just build
```

### Prerequisites

- Go 1.24.5 or later
- SQLite3 (automatically included)

## Quick Start

### Initialize a New Project

```bash
# In your project directory
knox init
```

This creates a `knox.json` configuration file and initializes your project's secret vault.

### Set Secrets

```bash
knox set DATABASE_URL "postgresql://localhost:5432/myapp"
knox set API_KEY "sk-1234567890abcdef"
```

### Run An Applicatiom

```bash
knox get DATABASE_URL
knox get API_KEY
```

### Check Project Status

```bash
knox status
```

## Architecture

Knox uses a clean, layered architecture:

- **CLI Layer** (`cmd/knox/`) - User interface and command routing
- **Business Logic** (`internal/`) - Core functionality and domain models
- **Storage** - SQLite-based encrypted vault storage

### Project Structure

```
knox/
├── cmd/knox/           # CLI application entry point
├── internal/
│   ├── config/         # Configuration management
│   ├── project/        # Project operations
│   ├── vault/          # Secret storage backend
│   └── secrets/        # Secret management logic
└── kit/                # Shared utilities
```

## Configuration

Knox stores project configuration in `knox.json`:

```json
{
  "project_id": "abc123def456",
  "vault_path": "/Users/you/.local/share/knox/knox.vault"
}
```

### Environment Variables

- `LOG_LEVEL=debug` - Enable debug logging

## Security

- **Local Storage Only** - Secrets never leave your machine
- **Project Isolation** - Each project maintains separate secrets
- **Unique Constraints** - Prevents accidental secret overwrites
- **No Network Access** - Pure local operation

## Development

Knox follows Go best practices and maintains a clean architecture.

### Build Commands

```bash
# Build the application
just build

# Run tests
just test

# Run with arguments
just run [command] [args]

# Lint code
just lint

# Format code
just format
```

### Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/vault

# With verbose output
go test -v ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes following the existing patterns
4. Add tests for new functionality
5. Run the test suite (`just test`)
6. Commit your changes (`git commit -am 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Code Style

- Follow Go conventions and existing patterns
- Use structured logging with `slog`
- Error handling with wrapped errors
- Tests for all new functionality

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [urfave/cli](https://github.com/urfave/cli) for CLI framework
- Uses [SQLite](https://sqlite.org/) for secure local storage
- Inspired by modern secret management best practices
