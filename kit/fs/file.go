package fs

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Copy copies a file from src to dst, preserving permissions.
func Copy(src, dst string) error {
	if src == "" {
		return fmt.Errorf("fs: cannot copy from empty source path")
	}
	if dst == "" {
		return fmt.Errorf("fs: cannot copy to empty destination path")
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("fs: open source file %q: %w", src, err)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("fs: stat source file %q: %w", src, err)
	}

	if !srcInfo.Mode().IsRegular() {
		return fmt.Errorf("fs: source %q is not a regular file", src)
	}

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("fs: create destination file %q: %w", dst, err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("fs: copy data from %q to %q: %w", src, dst, err)
	}

	return nil
}

// Move moves/renames a file from src to dst.
func Move(src, dst string) error {
	if src == "" {
		return fmt.Errorf("fs: cannot move from empty source path")
	}
	if dst == "" {
		return fmt.Errorf("fs: cannot move to empty destination path")
	}

	err := os.Rename(src, dst)
	if err != nil {
		return fmt.Errorf("fs: move %q to %q: %w", src, dst, err)
	}
	return nil
}

// Touch creates an empty file or updates the timestamp of an existing file.
func Touch(path string, perm os.FileMode) error {
	if path == "" {
		return fmt.Errorf("fs: cannot touch file with empty path")
	}

	now := time.Now()
	err := os.Chtimes(path, now, now)
	if err != nil {
		// File doesn't exist, create it
		file, createErr := os.OpenFile(path, os.O_CREATE, perm)
		if createErr != nil {
			return fmt.Errorf("fs: touch file %q: %w", path, createErr)
		}
		file.Close()
	}
	return nil
}

// ReadFile reads the contents of a file.
func ReadFile(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("fs: cannot read file with empty path")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("fs: read file %q: %w", path, err)
	}
	return data, nil
}

// WriteFile writes data to a file with the given permissions.
func WriteFile(path string, data []byte, perm os.FileMode) error {
	if path == "" {
		return fmt.Errorf("fs: cannot write file with empty path")
	}

	err := os.WriteFile(path, data, perm)
	if err != nil {
		return fmt.Errorf("fs: write file %q: %w", path, err)
	}
	return nil
}
