# Knox Command Layer Refactoring

## Purpose & Problem
The knox command layer has grown organically and now contains significant code duplication, inconsistent patterns, and maintenance issues. Every command handler repeats the same workspace resolution boilerplate, error handling patterns vary, and there's dead code scattered throughout.

## Success Criteria
- [ ] Eliminate duplicate workspace resolution code across all commands
- [ ] Standardize error codes and handling patterns
- [ ] Remove all dead code from command handlers
- [ ] Establish consistent command structure and naming conventions
- [ ] Create reusable utilities for common operations
- [ ] Maintain backward compatibility for all existing commands

## Scope & Implementation

### Phase 1: Common Utilities Package
Create `pkg/commands/common.go` with shared utilities:

```go
// WithWorkspace handles the common pattern of getting workspace
func WithWorkspace(handler func(*workspace.Workspace) error) error

// ValidateArgs validates command arguments with consistent error messages  
func ValidateArgs(cmd *cli.Command, expected int, operation string) error

// FormatOutput provides consistent output formatting
func FormatList(title string, items []string, activeItem string) string
func FormatSuccess(operation, target string) string
func FormatError(err error) string
```

### Phase 2: Error Code Standardization
Create `internal/codes.go` with standardized error codes:

```go
const (
    SearchFailureCode     = "SEARCH_FAILURE"
    ValidationCode        = "VALIDATION_ERROR" 
    CreateFailureCode     = "CREATE_FAILURE"
    VaultCreationCode     = "VAULT_CREATION_ERROR"
    SecretExistsCode      = "SECRET_EXISTS"
    SecretNotFoundCode    = "SECRET_NOT_FOUND"
)
```

### Phase 3: Clean Dead Code
Remove from `cmd_init.go`:
- `projectHandler()` function (lines 94-101)
- `initProjectArgValidator()` function (lines 46-54)
- `InitCmdArgs` struct (lines 17-19)
- `newInitWorkspaceCommand()` function (lines 34-44)

### Phase 4: Extract Common Patterns

#### Workspace Output Pattern
Create shared function for workspace status display (used in `cmd_init.go` lines 78-89):

```go
func PrintWorkspaceStatus(result workspace.Result, path, currentProject string)
```

#### Command Builder Pattern
Standardize command creation with consistent naming:
- `NewInitCommand()` ✓ (already consistent)
- `NewNewCommand()` → Keep (represents "new" subcommands)
- `NewProjectCommand()` ✓ (already consistent) 
- `NewStatusCommand()` ✓ (already consistent)

### Phase 5: Handler Refactoring
Refactor all handlers to use new utilities:

**Before:**
```go
func projectListHandler() error {
    cwd, err := os.Getwd()
    if err != nil {
        return errs.Wrap(err, internal.SearchFailureCode, "failed to get current working directory")
    }
    ws, err := workspace.FindWorkspace(cwd)
    if err != nil {
        return errs.Wrap(err, internal.SearchFailureCode, "failed to find workspace")
    }
    // actual logic...
}
```

**After:**
```go
func projectListHandler() error {
    return WithWorkspace(func(ws *workspace.Workspace) error {
        // actual logic only...
    })
}
```

## Technical Considerations

### Backward Compatibility
- All existing CLI commands must continue to work exactly as before
- Error messages should remain consistent with current behavior
- Output formats should not change

### Dependencies
- No new external dependencies
- Uses existing `github.com/urfave/cli/v3` patterns
- Builds on current `internal/workspace` and `kit/errs` packages

### File Structure
```
cmd/knox/internal/commands/
├── common.go           # New shared utilities
├── cmd_init.go         # Refactored, dead code removed
├── cmd_new.go          # Refactored to use common utilities
├── cmd_project.go      # Refactored to use common utilities  
├── cmd_status.go       # Minimal changes needed
└── internal/
    ├── codes.go        # New error code constants
    └── status.go       # Minor updates for error codes
```

## Out of Scope
- Changes to CLI argument parsing or command structure
- Modifications to workspace or vault core logic
- Performance optimizations beyond code deduplication
- New features or command additions
- Changes to configuration or environment handling

## Constraints
- Must maintain 100% backward compatibility
- Cannot change existing command interfaces
- Must preserve all current error messages and codes
- Refactoring only - no behavior changes
- All existing tests must continue to pass

## Risk Assessment
**Low Risk**: This is pure refactoring with no functional changes. The common utilities will be thoroughly tested and all existing commands will be validated to ensure identical behavior.

## Testing Strategy
- Unit tests for all new common utilities
- Integration tests to verify command behavior unchanged
- Regression testing on all existing command combinations
- Manual testing of error scenarios to ensure consistent messaging