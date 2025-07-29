package internal

import (
	"os"
	"path/filepath"
)

func EnvPath(name string, fallbackPaths ...string) string {
	dir := ExpandHome(os.Getenv(name))
	if dir != "" && filepath.IsAbs(dir) {
		return dir
	}

	return FirstAbsPath(fallbackPaths)
}

func FirstAbsPath(paths []string) string {
	for _, p := range paths {
		if p = ExpandHome(p); p != "" && filepath.IsAbs(p) {
			return p
		}
	}

	return ""
}
