package project

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/adrg/xdg"
	gonanoid "github.com/matoous/go-nanoid/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tomdoesdev/knox/internal"
	"github.com/tomdoesdev/knox/internal/constants"
	"github.com/tomdoesdev/knox/kit/fs"
	"github.com/tomdoesdev/knox/pkg/errs"
)

const (
	CreateFailureCode errs.ErrorCode = "PROJECT_EXISTS"
	NotFoundCode      errs.ErrorCode = "PROJECT_NOT_FOUND"
)

var (
	ErrProjectExists   = errs.New(CreateFailureCode, "project already exists")
	ErrProjectNotFound = errs.New(NotFoundCode, "project not found")
)

type Config struct {
	ProjectID     string      `json:"project_id"`
	ProjectPath   fs.FilePath `json:"-"`
	VaultFilePath fs.FilePath `json:"vault_file_path,omitempty"`
}

type Project struct {
	config Config
}

func NewProject(config *Config) (*Project, error) {
	//If no vaultpath is provided by the project config fall back to the 'global' vault
	if config.VaultFilePath == "" {
		config.VaultFilePath = path.Join(xdg.DataHome, constants.DefaultKnoxDirName, constants.DefaultVaultFileName)
	}

	return &Project{config: *config}, nil
}

func (p *Project) Workspace() *internal.Workspace {
	return &internal.Workspace{
		ProjectID:       p.config.ProjectID,
		VaultFilePath:   p.config.VaultFilePath,
		ProjectFilePath: p.config.ProjectPath,
	}
}

func NewConfig() (*Config, error) {
	pid, err := gonanoid.New(12)
	if err != nil {
		return nil, fmt.Errorf("project: generating project id: %w", err)
	}

	conf := Config{
		ProjectID: pid,
	}

	return &conf, nil
}

func CreateFile(projectPath fs.FilePath) (*Config, error) {
	if strings.TrimSpace(projectPath) == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, errs.Wrap(err, CreateFailureCode, "could not get current working directory")
		}
		projectPath = cwd
	}

	confpath := path.Join(projectPath, constants.DefaultProjectFileName)

	exists, err := fs.IsExist(confpath)
	if err != nil {
		return nil, errs.Wrap(err, CreateFailureCode,
			"faild to check if %s exists at %s", constants.DefaultProjectFileName, confpath)
	}

	if exists {
		return nil, ErrProjectExists
	}

	p, err := NewConfig()
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
	confPath := path.Join(confDir, constants.DefaultProjectFileName)
	slog.Debug("project.load.fromConfPath", "confPath", confPath)

	exists, err := fs.IsExist(confPath)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrProjectNotFound
	}

	jsonbytes, err := os.ReadFile(confPath)
	if err != nil {
		return nil, errs.Wrap(err, CreateFailureCode, "failed to read %s", constants.DefaultProjectFileName)
	}

	c := new(Config)

	err = json.Unmarshal(jsonbytes, &c)
	if err != nil {
		return nil, err
	}

	p, err := NewProject(c)
	if err != nil {
		return nil, err
	}

	return p, nil
}
