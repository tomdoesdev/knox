package config

import "github.com/tomdoesdev/knox/kit"

type KnoxConfig struct {
	ProjectID string `json:"project_id"`
}

func NewKnoxConfig() (*KnoxConfig, error) {
	pid, err := kit.NewNanoid()
	if err != nil {
		return nil, err
	}

	conf := KnoxConfig{
		ProjectID: pid,
	}

	return &conf, nil
}
