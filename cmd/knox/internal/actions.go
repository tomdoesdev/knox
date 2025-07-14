package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/tomdoesdev/knox/internal/config"
	"github.com/tomdoesdev/knox/internal/vault"
	"github.com/tomdoesdev/knox/kit/fskit"
)

var (
	ErrAlreadyInitialized = errors.New("knox.json already initialized")
)

func Initialize(conf *config.ApplicationConfig) error {

	err := vault.EnsureVaultExists(conf)

	err = createProjectFile()
	if err != nil {
		return fmt.Errorf("knox.init.createProjectFile: %w", err)
	}

	return nil
}

func createProjectFile() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	confpath := path.Join(cwd, "knox.json")

	if exists, _ := fskit.Exists(confpath); exists {
		return ErrAlreadyInitialized
	}
	slog.Info("creating project file", "path", confpath)
	c, err := config.NewProjectConfig()
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
