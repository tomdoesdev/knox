package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// MkdirAll creates a directory and all necessary parents with the given permissions.
func MkdirAll(path string, perm os.FileMode) error {
	if path == "" {
		return fmt.Errorf("fs: cannot create directory with empty path")
	}

	err := os.MkdirAll(path, perm)
	if err != nil {
		return fmt.Errorf("fs: create directory %q: %w", path, err)
	}
	return nil
}

// RemoveAll removes a directory and all its contents.
func RemoveAll(path string) error {
	if path == "" {
		return fmt.Errorf("fs: cannot remove directory with empty path")
	}

	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("fs: remove directory %q: %w", path, err)
	}
	return nil
}

// ListDir returns a list of directory entry names in the given directory.
func ListDir(path string) ([]string, error) {
	if path == "" {
		return nil, fmt.Errorf("fs: cannot list directory with empty path")
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("fs: read directory %q: %w", path, err)
	}

	names := make([]string, len(entries))
	for i, entry := range entries {
		names[i] = entry.Name()
	}

	return names, nil
}

// WalkDir walks the file tree rooted at root, calling fn for each file or directory.
func WalkDir(root string, fn func(path string, d os.DirEntry) error) error {
	if root == "" {
		return fmt.Errorf("fs: cannot walk directory with empty root")
	}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return fn(path, d)
	})

	if err != nil {
		return fmt.Errorf("fs: walk directory %q: %w", root, err)
	}
	return nil
}
