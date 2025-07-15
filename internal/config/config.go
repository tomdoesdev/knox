package config

import (
	"fmt"
	"path"

	"github.com/adrg/xdg"
)

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
