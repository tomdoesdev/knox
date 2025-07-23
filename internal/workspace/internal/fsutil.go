package internal

import (
	"path/filepath"

	"github.com/tomdoesdev/knox/kit/fs"
)

func IsDataDirectory(path string) bool {
	if fs.IsDir(path) && filepath.Base(path) == DataDirectoryName {
		return true
	}
	return false
}

// ContainsDataDirectory returns true if the
func ContainsDataDirectory(path string) bool {
	if !fs.IsDir(path) {
		return false
	}

	path = filepath.Join(path, DataDirectoryName)

	return IsDataDirectory(path)
}
