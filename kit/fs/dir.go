package fs

import (
	"os"
	"path/filepath"

	"github.com/tomdoesdev/knox/kit/errs"
)

// MkdirAll creates a directory and all necessary parents with the given permissions.
func MkdirAll(path string, perm os.FileMode) error {
	if path == "" {
		return errs.New(ECodeInvalidPath, "cannot create directory with empty path")
	}

	err := os.MkdirAll(path, perm)
	if err != nil {
		return errs.Wrap(err, ECodeDirectoryFailure, "create directory failed").
			WithPath(path).
			WithOperation("mkdir")
	}
	return nil
}

// RemoveAll removes a directory and all its contents.
func RemoveAll(path string) error {
	if path == "" {
		return errs.New(ECodeInvalidPath, "cannot remove directory with empty path")
	}

	err := os.RemoveAll(path)
	if err != nil {
		return errs.Wrap(err, ECodeDirectoryFailure, "remove directory failed").
			WithPath(path).
			WithOperation("remove")
	}
	return nil
}

// ListDir returns a list of directory entry names in the given directory.
func ListDir(path string) ([]string, error) {
	if path == "" {
		return nil, errs.New(ECodeInvalidPath, "cannot list directory with empty path")
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, errs.Wrap(err, ECodeDirectoryFailure, "read directory failed").
			WithPath(path).
			WithOperation("readdir")
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
		return errs.New(ECodeInvalidPath, "cannot walk directory with empty root")
	}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return fn(path, d)
	})

	if err != nil {
		return errs.Wrap(err, ECodeDirectoryFailure, "walk directory failed").
			WithPath(root).
			WithOperation("walk")
	}
	return nil
}
