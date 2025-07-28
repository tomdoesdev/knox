package fs

import "path/filepath"

// CleanPath cleans a file path, removing redundant elements.
func CleanPath(path string) string {
	return filepath.Clean(path)
}

// JoinPath joins path elements into a single path.
func JoinPath(parts ...string) string {
	return filepath.Join(parts...)
}

// BaseName returns the last element of a path.
func BaseName(path string) string {
	return filepath.Base(path)
}

// DirName returns all but the last element of a path.
func DirName(path string) string {
	return filepath.Dir(path)
}
