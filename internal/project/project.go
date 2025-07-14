package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/tomdoesdev/knox/kit/fskit"
)

var (
	ErrProjectExists = errors.New("project already exists")
)

const ConfigFilename = "knox.json"

type Project struct {
	ProjectID string `json:"project_id"`
	VaultPath string `json:"vault_path,omitempty"`
}

func Create() (*Project, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	confpath := path.Join(cwd, ConfigFilename)

	exists, err := fskit.Exists(confpath)
	if err != nil {
		return nil, fmt.Errorf("project.createProject: %w", err)
	}

	if exists {
		return nil, ErrProjectExists
	}

	slog.Info("creating project file", "path", confpath)
	p, err := newProject()
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

func Load(confPath string) (*Project, error) {
	jsonbytes, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	p := new(Project)

	err = json.Unmarshal(jsonbytes, &p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func newProject() (*Project, error) {

	pid, err := gonanoid.New(12)
	if err != nil {
		return nil, fmt.Errorf("project: generating project id: %w", err)
	}

	conf := Project{
		ProjectID: pid,
	}

	return &conf, nil
}
