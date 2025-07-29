package internal

import (
	"os"
	"path/filepath"
	"strings"
)

func UserHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}
	if home == "" {
		panic("$HOME not set")
	}
	return home
}

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
