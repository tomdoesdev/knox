package project

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/adrg/xdg"
	gonanoid "github.com/matoous/go-nanoid/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tomdoesdev/knox/kit/fskit"
)

var (
	ErrProjectExists   = errors.New("project already exists")
	ErrProjectNotFound = errors.New("project not found")
	ErrSecretExists    = errors.New("secret already exists")
)

const ConfigFilename = "knox.json"

const setSecretStmt = `
	INSERT INTO vault (project_id, key, value) VALUES ($1, $2, $3);
`

const getSecretStmt = `
	SELECT value FROM vault WHERE project_id = $1 AND key = $2 LIMIT 1;
`

type Project struct {
	ProjectID string `json:"project_id"`
	VaultPath string `json:"vault_path,omitempty"`
	db        *sql.DB
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

func (p *Project) ReadSecret(key string) (value string, err error) {
	db, err := p.ensureDb()
	if err != nil {
		return "", err
	}

	var secret string
	err = db.QueryRow(getSecretStmt, p.ProjectID, key).Scan(&secret)
	if err != nil {
		return "", err
	}

	return secret, nil
}

func (p *Project) WriteSecret(key string, value string) error {
	db, err := p.ensureDb()
	if err != nil {
		return err
	}

	_, err = db.Exec(setSecretStmt, p.ProjectID, key, value)
	if err != nil {
		if isUniqueConstraintError(err) {
			return fmt.Errorf("secret key '%s': %w", key, ErrSecretExists)
		}
		return err
	}

	return nil
}

func isUniqueConstraintError(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "UNIQUE constraint failed") ||
		strings.Contains(errStr, "constraint failed")
}

func (p *Project) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func CreateFile() (*Project, error) {
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

func Load(confDir string) (*Project, error) {
	confPath := path.Join(confDir, ConfigFilename)

	slog.Debug("loading project configuration", "confPath", confPath)
	exists, err := fskit.Exists(confPath)
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

	return p, nil
}

func LoadCwd() (*Project, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return Load(cwd)
}

func (p *Project) ensureDb() (*sql.DB, error) {
	if p.db != nil {
		return p.db, nil
	}

	db, err := sql.Open("sqlite3", p.VaultPath)
	if err != nil {
		fmt.Printf("project.ensureDb: %v\n", err)
		return nil, err
	}
	p.db = db
	return p.db, nil
}
