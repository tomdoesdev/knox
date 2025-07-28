package fs

import (
	"fmt"
	"os"
)

// TempDir creates a temporary directory and returns its path along with a cleanup function.
func TempDir(pattern string) (string, func(), error) {
	dir, err := os.MkdirTemp("", pattern)
	if err != nil {
		return "", nil, fmt.Errorf("fs: create temp directory: %w", err)
	}

	cleanup := func() {
		os.RemoveAll(dir)
	}

	return dir, cleanup, nil
}

// TempFile creates a temporary file and returns the file along with a cleanup function.
func TempFile(pattern string) (*os.File, func(), error) {
	file, err := os.CreateTemp("", pattern)
	if err != nil {
		return nil, nil, fmt.Errorf("fs: create temp file: %w", err)
	}

	cleanup := func() {
		file.Close()
		os.Remove(file.Name())
	}

	return file, cleanup, nil
}
