# Filesystem Utilities Enhancement Spec

## Overview
Enhance the `kit/fs` package to provide a comprehensive set of filesystem utilities that can be shared across multiple Go projects, focusing on reducing boilerplate for common filesystem operations while maintaining idiomatic Go patterns.

## Current State
The `kit/fs` package currently contains basic path utilities:
- `IsExist(path string) (bool, error)` - Check if path exists (inconsistent - returns error)
- `IsDir(path string) bool` - Check if path is directory  
- `IsFile(path string) bool` - Check if path is regular file

## Purpose & User Problem
Go developers frequently write repetitive filesystem code across projects. Common pain points include:
- Verbose error handling for simple existence checks
- Boilerplate for creating directories with parents
- Manual cleanup of temporary files/directories
- Repetitive file copying and moving operations
- Complex directory walking with filtering

## Success Criteria
- Reduce common filesystem boilerplate by 50-70%
- Maintain 100% compatibility with Go standard library patterns
- Zero external dependencies
- All functions follow consistent naming and return patterns
- Safe for concurrent use where applicable

## Proposed Features

### 1. Consistent Path Checking (Fix existing + add new)
- `IsExist(path string) bool` - Fix to return boolean only
- `IsEmpty(path string) bool` - Check if file/directory is empty
- `IsSymlink(path string) bool` - Check if path is symbolic link
- `IsExecutable(path string) bool` - Check if file is executable

### 2. Directory Operations
- `MkdirAll(path string, perm os.FileMode) error` - Create directory with parents (wraps os.MkdirAll)
- `RemoveAll(path string) error` - Remove directory and contents (wraps os.RemoveAll)
- `ListDir(path string) ([]string, error)` - List directory contents (names only)
- `WalkDir(root string, fn func(path string, d os.DirEntry) error) error` - Walk directory tree

### 3. File Operations
- `Copy(src, dst string) error` - Copy file with proper permissions
- `Move(src, dst string) error` - Move/rename file
- `Touch(path string) error` - Create empty file or update timestamp
- `ReadFile(path string) ([]byte, error)` - Alias for os.ReadFile (consistency)
- `WriteFile(path string, data []byte, perm os.FileMode) error` - Alias for os.WriteFile

### 4. Temporary File/Directory Management
- `TempDir(pattern string) (string, func(), error)` - Create temp dir with cleanup function
- `TempFile(pattern string) (*os.File, func(), error)` - Create temp file with cleanup function

### 5. Path Utilities (enhance existing)
- Keep existing `FilePath` and `DirPath` type aliases
- `CleanPath(path string) string` - Alias for filepath.Clean
- `JoinPath(parts ...string) string` - Alias for filepath.Join
- `BaseName(path string) string` - Alias for filepath.Base
- `DirName(path string) string` - Alias for filepath.Dir

## Technical Considerations
- **Consistency**: All `Is*` functions return `bool` only, no errors
- **Error Handling**: Non-boolean functions return descriptive errors with path context
- **Concurrency**: All functions are safe for concurrent use (stateless)
- **Performance**: Leverage standard library implementations, avoid unnecessary allocations
- **Cleanup Pattern**: Temp functions return cleanup functions to avoid resource leaks

## API Design Principles
1. **Idiomatic**: Follow Go naming conventions and patterns
2. **Consistent**: Similar operations have similar signatures
3. **Convenient**: Reduce boilerplate without hiding complexity
4. **Safe**: Proper error handling and resource cleanup

## Out of Scope
- Advanced file operations (compression, encryption)
- Network filesystem support
- File watching/monitoring
- Complex permission management beyond basic modes
- File locking mechanisms

## Implementation Notes
- Update `IsExist` to return `bool` only for consistency
- Use `filepath` package for all path operations
- Wrap standard library functions where they provide convenience
- Include path in all error messages for debugging
- Cleanup functions should be safe to call multiple times