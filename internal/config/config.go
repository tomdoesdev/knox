package config

import (
	"fmt"
	"path"

	"github.com/adrg/xdg"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type ProjectConfig struct {
	ProjectID string `json:"project_id"`
}

func NewProjectConfig() (*ProjectConfig, error) {

	pid, err := gonanoid.New(12)
	if err != nil {
		return nil, fmt.Errorf("config: generating project id: %w", err)
	}

	conf := ProjectConfig{
		ProjectID: pid,
	}

	return &conf, nil
}

type ApplicationConfig struct {
	VaultDir      string `json:"vault_dir"`
	VaultFileName string `json:"vault_file_name"`
	VaultPath     string `json:"vault_path"`
}

func (conf *ApplicationConfig) String() string {
	return fmt.Sprintf("%+v", *conf)
}

func NewApplicationConfig() ApplicationConfig {

	fileName := "knox.vault"
	vaultDir := path.Join(xdg.DataHome, "knox")
	vaultPath := path.Join(vaultDir, fileName)

	conf := ApplicationConfig{
		VaultDir:      vaultDir,
		VaultFileName: fileName,
		VaultPath:     vaultPath,
	}

	return conf
}
