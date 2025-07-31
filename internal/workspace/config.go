package workspace

import (
	"github.com/google/uuid"
	"github.com/tomdoesdev/knox/internal"
)

const (
	defaultProjectName string = "default"
)

type (
	Settings struct {
		DefaultProject string `toml:"default_project"`
	}

	Config struct {
		Version  string    `toml:"version"`
		Settings *Settings `toml:"settings"`
		Id       string    `toml:"id"`
	}

	ProjectConfig struct {
		Name    string `json:"name"`
		VaultId string `json:"vault_id"`
	}
)

func NewConfigDefault() *Config {
	id := uuid.New()

	return &Config{
		Id:      id.String(),
		Version: internal.KnoxVersion,
		Settings: &Settings{
			DefaultProject: defaultProjectName,
		},
	}
}
