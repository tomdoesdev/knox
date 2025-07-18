# Template Processing Specification

## Overview

This specification covers the implementation of .env.template file parsing and processing for Knox's `run` command. This is Phase 2 of the Knox Application roadmap, enabling developers to use Go template syntax in environment files for secure secret injection.

## Purpose & User Problem

### Primary Problem
Developers need a way to reference secrets stored in Knox vault within their environment configuration files without exposing actual secret values in version control.

### Target Use Case
- Developer has `.env.template` file in their project with template placeholders
- Knox processes the template, injects secrets from vault, and provides clean environment to application
- Template file can be safely committed to git (contains no actual secrets)

### User Goals
- Write environment templates using familiar Go template syntax
- Reference Knox secrets without hardcoding values
- Maintain readable, documented environment configuration
- Enable team sharing of environment structure without secret values

## Success Criteria

### Primary Success Metrics
- **Template Processing**: Successfully parse and execute .env.template files
- **Secret Injection**: Accurately replace template placeholders with vault secrets
- **Error Handling**: Clear error messages for template syntax errors and missing secrets
- **Performance**: Template processing completes in < 50ms for typical files

### User Experience Goals
- Intuitive template syntax matching Go conventions
- Helpful error messages with line numbers for template errors
- Seamless integration with existing Knox workflow
- Support for common environment variable patterns

## Core Features

### 1. Template File Discovery
- **Default Location**: Look for `.env.template` in current directory
- **Custom Path**: Support `--env FILE` flag to specify template file
- **Fallback Behavior**: If no template file found, proceed without environment injection
- **Multiple Files**: Support for multiple template files (future enhancement)

### 2. Template Syntax Support
- **Secret Function**: `{{.Secret "KEY"}}` - Retrieve secret from Knox vault
- **Environment Function**: `{{.Env "KEY"}}` - Get existing environment variable
- **Default Function**: `{{.Default "KEY" "fallback"}}` - Get secret with fallback value
- **Comments**: Standard Go template comments `{{/* comment */}}`

### 3. Template Processing
- **Engine**: Go's `text/template` package for reliable, fast processing
- **Context**: Provide vault access and environment functions to template
- **Output**: Generate complete environment variable list in memory
- **Validation**: Validate template syntax before processing

### 4. Environment Variable Handling
- **Parsing**: Use `github.com/hashicorp/go-envparse` for env var parsing
- **Inheritance**: Preserve existing environment variables
- **Override**: Template variables override existing environment
- **Format**: Support standard KEY=VALUE format

## Technical Architecture

### Template Engine
```go
// Template context structure
type TemplateContext struct {
    secretStore SecretStore
    environment map[string]string
}

// Template functions
func (ctx *TemplateContext) Secret(key string) (string, error)
func (ctx *TemplateContext) Env(key string) string
func (ctx *TemplateContext) Default(key, fallback string) string
```

### Processing Pipeline
1. **File Discovery**: Locate .env.template file
2. **Template Parsing**: Parse template syntax and validate
3. **Context Preparation**: Setup vault access and environment
4. **Template Execution**: Execute template with context
5. **Environment Parsing**: Parse resulting environment variables
6. **Validation**: Ensure all secrets are resolved

### Error Handling
- **Template Syntax Errors**: Line numbers and syntax descriptions
- **Missing Secrets**: Clear indication of undefined secret keys
- **File Errors**: Readable file access error messages
- **Validation Errors**: Template validation failures with context

## Implementation Details

### Template Functions

#### Secret Function
```go
// {{.Secret "API_KEY"}}
func (ctx *TemplateContext) Secret(key string) (string, error) {
    value, err := ctx.secretStore.ReadSecret(key)
    if err != nil {
        return "", fmt.Errorf("secret %q not found: %w", key, err)
    }
    return value, nil
}
```

#### Environment Function  
```go
// {{.Env "PATH"}}
func (ctx *TemplateContext) Env(key string) string {
    return os.Getenv(key)
}
```

#### Default Function
```go
// {{.Default "DEBUG_MODE" "false"}}
func (ctx *TemplateContext) Default(key, fallback string) string {
    if value, err := ctx.secretStore.ReadSecret(key); err == nil {
        return value
    }
    return fallback
}
```

### File Structure
```
internal/
  template/
    processor.go      # Main template processing logic
    context.go        # Template context and functions
    errors.go         # Template-specific error types
    processor_test.go # Unit tests
```

### Example Template File
```bash
# .env.template - Safe to commit to git
DATABASE_URL={{.Secret "DATABASE_URL"}}
API_KEY={{.Secret "API_KEY"}}
DEBUG={{.Default "DEBUG_MODE" "false"}}
PATH={{.Env "PATH"}}
APP_NAME=myapp
```

### Generated Environment
```bash
DATABASE_URL=postgres://user:pass@localhost/db
API_KEY=sk-1234567890abcdef
DEBUG=true
PATH=/usr/local/bin:/usr/bin:/bin
APP_NAME=myapp
```

## Error Scenarios & Handling

### Template Syntax Errors
```
Error: template syntax error in .env.template:3
  {{.Secret "API_KEY}  # Missing closing quote
                   ^
Expected: closing quote (")
```

### Missing Secret Errors
```
Error: template execution failed in .env.template:1
  DATABASE_URL={{.Secret "DATABASE_URL"}}
Secret "DATABASE_URL" not found in vault
Hint: Use 'knox set DATABASE_URL <value>' to add this secret
```

### File Access Errors
```
Error: cannot read template file
  File: .env.template
  Cause: permission denied
Hint: Check file permissions or specify different file with --env
```

## Configuration & Integration

### Command Line Integration
```bash
# Use default .env.template file
knox run npm start

# Use custom template file
knox run --env production.env.template npm start

# No template processing (existing behavior)
knox run npm start  # when no .env.template exists
```

### Template Context Configuration
- **Vault Access**: Use project's configured vault and secret store
- **Environment**: Inherit from current process environment
- **Error Mode**: Fail fast on missing secrets (no partial execution)

## Testing Strategy

### Unit Testing
- Template parsing with valid/invalid syntax
- Template function execution (Secret, Env, Default)
- Error handling for missing secrets and syntax errors
- Environment variable parsing and merging

### Integration Testing
- End-to-end template processing with real vault
- File discovery and custom path handling
- Environment inheritance and override behavior
- Error message accuracy and helpfulness

### Test Cases
```go
// Valid template processing
func TestTemplateProcessing_ValidTemplate(t *testing.T)

// Missing secret handling
func TestTemplateProcessing_MissingSecret(t *testing.T)

// Syntax error handling
func TestTemplateProcessing_SyntaxError(t *testing.T)

// Environment inheritance
func TestTemplateProcessing_EnvironmentInheritance(t *testing.T)
```

## Security Considerations

### Template Security
- **No Code Execution**: Template functions are read-only, no arbitrary code execution
- **Secret Isolation**: Only project secrets accessible via template context
- **Memory Processing**: Template output never written to disk
- **Input Validation**: Validate template syntax before processing

### Error Information
- **No Secret Leakage**: Error messages never include secret values
- **Limited Context**: Error messages show template structure, not secret content
- **Safe Logging**: Template processing errors safe to log

## Performance Requirements

### Processing Speed
- **Template Parsing**: < 10ms for typical .env.template files
- **Secret Retrieval**: < 1ms per secret lookup from vault
- **Total Processing**: < 50ms for files with 10-20 variables
- **Memory Usage**: Minimal memory footprint, no persistent template cache

### Scalability
- **File Size**: Support templates up to 1MB (thousands of variables)
- **Secret Count**: Handle 100+ secret references per template
- **Concurrent Access**: Thread-safe template processing

## Implementation Phases

### Phase 2.1: Core Template Engine
- [ ] Template parsing and validation
- [ ] Secret, Env, Default function implementation
- [ ] Basic error handling with line numbers
- [ ] Unit tests for template functions

### Phase 2.2: File Processing
- [ ] .env.template file discovery
- [ ] Custom file path support (--env flag)
- [ ] Environment variable parsing integration
- [ ] Integration tests

### Phase 2.3: Error Experience
- [ ] Enhanced error messages with context
- [ ] Template validation before execution
- [ ] Helpful hints for common errors
- [ ] Error message testing

### Phase 2.4: Integration & Polish
- [ ] CLI command integration
- [ ] Performance optimization
- [ ] Documentation and examples
- [ ] Cross-platform testing

## Out of Scope

### Explicitly Excluded
- **Complex Template Logic**: No conditional statements, loops, or complex control flow
- **Template Includes**: No support for including other template files
- **Custom Functions**: No user-defined template functions
- **Template Caching**: No persistent template compilation caching
- **Template Inheritance**: No template extension or inheritance
- **Dynamic Templates**: No runtime template modification

### Future Considerations
- Multiple environment file support
- Template validation CLI command
- IDE integration for template syntax highlighting
- Template debugging tools

## Success Metrics

### Quantitative Metrics
- **Processing Time**: < 50ms for typical templates
- **Error Rate**: < 1% template processing failures
- **Test Coverage**: > 95% coverage for template processing code
- **Memory Usage**: < 10MB peak memory during template processing

### Qualitative Metrics
- **Developer Experience**: Intuitive template syntax and clear error messages
- **Reliability**: Consistent template processing across platforms
- **Integration**: Seamless fit with existing Knox workflow
- **Documentation**: Clear examples and error resolution guides

---

*This specification covers Phase 2 of Knox development, enabling secure environment template processing as outlined in the main KnoxApplication.md specification.*