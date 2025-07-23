package fs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type FilePath = string
type DirPath = string

func IsExist(path string) (bool, error) {
	if path == "" {
		return false, nil
	}
	_, err := os.Lstat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, fmt.Errorf("pathutil: stat %q: %w", path, err)
}

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
