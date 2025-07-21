package workspace

type ConfigManager interface {
	Read() (*Config, error)
	Write(config *Config) error
}

// Config file types
type Config struct {
	WorkspaceVersion string `json:"workspace_version"`
}
