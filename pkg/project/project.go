package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/adrg/xdg"
	gonanoid "github.com/matoous/go-nanoid/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tomdoesdev/knox/internal/vault"
	"github.com/tomdoesdev/knox/kit/fs"
)

var (
	ErrProjectExists   = errors.New("project already exists")
	ErrProjectNotFound = errors.New("project not found")
)

const ConfigFilename = "knox.json"

type Project struct {
	ProjectID string `json:"project_id"`
	VaultPath string `json:"vault_path,omitempty"`
}

func New() (*Project, error) {
	pid, err := gonanoid.New(12)
	if err != nil {
		return nil, fmt.Errorf("project: generating project id: %w", err)
	}

	vaultPath := path.Join(path.Join(xdg.DataHome, "knox"), "knox.vault")

	conf := Project{
		ProjectID: pid,
		VaultPath: vaultPath,
	}

	return &conf, nil
}

func CreateFile() (*Project, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	confpath := path.Join(cwd, ConfigFilename)

	exists, err := fs.Exists(confpath)
	if err != nil {
		return nil, fmt.Errorf("project.createProject: %w", err)
	}

	if exists {
		return nil, ErrProjectExists
	}

	p, err := New()
	if err != nil {
		return nil, err
	}

	jsonbytes, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(confpath, jsonbytes, 0660)
	if err != nil {
		return nil, err
	}

	return p, nil
}
func Load() (*Project, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return load(cwd)
}
func load(confDir string) (*Project, error) {
	confPath := path.Join(confDir, ConfigFilename)

	slog.Debug("loading project configuration", "confPath", confPath)
	exists, err := fs.Exists(confPath)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrProjectNotFound
	}

	jsonbytes, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	p := new(Project)

	err = json.Unmarshal(jsonbytes, &p)
	if err != nil {
		return nil, err
	}

	if p.VaultPath == "" {
		dp := vault.DefaultPath()
		slog.Debug("project.load: no vault_path found in configuration, using default",
			slog.String("path", dp),
		)
		p.VaultPath = dp
	}

	return p, nil
}
