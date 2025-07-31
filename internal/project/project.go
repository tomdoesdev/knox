package project

import (
	"github.com/go-playground/validator/v10"
	"github.com/tomdoesdev/knox/kit/errs"
)

const (
	ECodeInvalidConfig errs.Code = "INVALID_CONFIG"
)

type (
	Project struct {
		name      string
		vaultPath string
	}

	Config struct {
		Name   string     `json:"name" validate:"required"`
		Source Datasource `json:"source" validate:"required"`
	}

	Datasource struct {
		Type string `json:"type" validate:"required,oneof=json"`
		URI  string `json:"uri" validate:"required"`
	}
)

var (
	validate = validator.New(validator.WithRequiredStructEnabled())
)

func ValidateConfig(config *Config) error {
	err := validate.Struct(config)
	if err != nil {
		return errs.Wrap(err, ECodeInvalidConfig, "failed to validate project config")
	}
	return nil
}

func CreateDefault(name string, vaultPath string) error {
	panic("implement me")
}

func Open(path string) (*Project, error) {
	panic("implement me")
}
