package config

type (
	WorkspaceSettings struct {
		DefaultProject string `toml:"default_project"`
	}

	WorkspaceConfig struct {
		Version  string            `toml:"version"`
		Settings WorkspaceSettings `toml:"settings"`
		Vaults   []VaultConfig     `toml:"vaults"`
	}

	VaultConfig struct {
		Alias       string `toml:"alias"`
		Resource    string `toml:"resource"`
		Description string `toml:"description"`
	}
)
