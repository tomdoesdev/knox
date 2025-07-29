package xdg

import (
	"path/filepath"

	"github.com/tomdoesdev/knox/kit/xdg/internal"
)

func initBaseDirs(home string) {
	ConfigHome = internal.EnvPath(envConfigHome, filepath.Join(home, ".config"))
}
