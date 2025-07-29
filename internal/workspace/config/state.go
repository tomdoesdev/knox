package config

type (
	State struct {
		Version       string `json:"version"`
		ActiveProject string `json:"active_project"`
	}
)
