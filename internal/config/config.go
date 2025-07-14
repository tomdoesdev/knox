package config

import (
	"fmt"
	"log/slog"
	"path"

	"github.com/adrg/xdg"
	"github.com/tomdoesdev/knox/kit"
)

type ProjectConfig struct {
	ProjectID string `json:"project_id"`
}

func NewProjectConfig() (*ProjectConfig, error) {
	pid, err := kit.NewNanoid()
	if err != nil {
		return nil, fmt.Errorf("config: generating project id: %w", err)
	}

	conf := ProjectConfig{
		ProjectID: pid,
	}

	return &conf, nil
}

type ApplicationConfig struct {
	VaultDir string `json:"vault_dir"`
}

func (conf *ApplicationConfig) String() string {
	return fmt.Sprintf("%+v", *conf)
}

func NewApplicationConfig() ApplicationConfig {

	vaultDir := path.Join(xdg.DataHome, "knox")

	conf := ApplicationConfig{
		VaultDir: vaultDir,
	}

	slog.Info("config: creating new application config", "config", conf)

	return conf
}
