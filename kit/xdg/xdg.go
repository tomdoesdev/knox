package xdg

import (
	"path/filepath"

	"github.com/tomdoesdev/knox/kit/errs"
	"github.com/tomdoesdev/knox/kit/fs"
	"github.com/tomdoesdev/knox/kit/xdg/internal"
)

const (
	envDataHome   = "XDG_DATA_HOME"
	envDataDirs   = "XDG_DATA_DIRS"
	envConfigHome = "XDG_CONFIG_HOME"
)

var (
	Home       = internal.UserHomeDir()
	ConfigHome string
	DataHome   string
)

var (
	BadPathErrCode        errs.Code = "BAD_PATH"
	MakeDirFailureErrCode errs.Code = "CANT_MK_DIR"
	TouchFailureErrCode   errs.Code = "TOUCH_FAILURE"
)

func init() {
	initBaseDirs(Home)
}

func ConfigFile(appName, filename string) (string, error) {
	p := filepath.Join(ConfigHome, appName)
	if p == "" {
		return "", errs.New(BadPathErrCode, "path for config file is empty").WithPath(p)
	}

	err := fs.MkdirAll(p, 0600)
	if err != nil {
		return "", errs.Wrap(err, MakeDirFailureErrCode, "xdg: failed to make config dirs").WithPath(p)
	}

	configPath := filepath.Join(p, filename)
	if fs.IsFile(configPath) {
		return configPath, nil
	}

	err = fs.Touch(configPath)
	if err != nil {
		return "", errs.Wrap(err, TouchFailureErrCode, "xdg: failed to touch config file").WithPath(configPath)
	}

	return configPath, nil

}
