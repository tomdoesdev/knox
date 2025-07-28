package fs

import (
	"io"
	"os"
	"path/filepath"
)

// IsExist checks if a path exists.
func IsExist(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Lstat(path)
	return err == nil
}

// IsDir checks if a path exists and is a directory.
func IsDir(path string) bool {
	if path == "" {
		return false
	}

	path = filepath.Clean(path)

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

// IsFile checks if a path exists and is a regular file.
func IsFile(path string) bool {
	if path == "" {
		return false
	}
	path = filepath.Clean(path)

	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

// IsEmpty checks if a file or directory is empty.
func IsEmpty(path string) bool {
	if path == "" {
		return false
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if info.IsDir() {
		// Check if directory is empty
		f, err := os.Open(path)
		if err != nil {
			return false
		}
		defer f.Close()

		_, err = f.Readdirnames(1)
		return err == io.EOF
	}

	// Check if file is empty
	return info.Size() == 0
}

// IsSymlink checks if a path is a symbolic link.
func IsSymlink(path string) bool {
	if path == "" {
		return false
	}

	info, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeSymlink != 0
}

// IsExecutable checks if a file is executable.
func IsExecutable(path string) bool {
	if path == "" {
		return false
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Check if it's a regular file and has execute permission
	return info.Mode().IsRegular() && info.Mode().Perm()&0111 != 0
}
