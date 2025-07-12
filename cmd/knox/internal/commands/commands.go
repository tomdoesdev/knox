package commands

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path"

	"github.com/tomdoesdev/knox/internal/config"
	"github.com/tomdoesdev/knox/kit/fskit"
)

var (
	ErrAlreadyInitialized = errors.New("knox.json already initialized")
)

func Init() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	confpath := path.Join(cwd, "knox.json")
	slog.Info(confpath)
	if exists, _ := fskit.Exists(confpath); exists {
		return ErrAlreadyInitialized
	}

	c, err := config.NewKnoxConfig()
	if err != nil {
		return err
	}

	jsonbytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(confpath, jsonbytes, 0660)
	if err != nil {
		return err
	}

	return nil

}
