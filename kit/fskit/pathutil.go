package fskit

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func ExpandHome(path string) string {
	home := UserHomeDir()
	if path == "" || home == "" {
		return path
	}
	if path[0] == '~' {
		return filepath.Join(home, path[1:])
	}
	if strings.HasPrefix(path, "$HOME") {
		return filepath.Join(home, path[5:])
	}

	return path
}
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || errors.Is(err, os.ErrNotExist)
}

func UserHomeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}

	return "/"
}
