package envar

import (
	"fmt"
	"io"

	"github.com/hashicorp/go-envparse"
)

type EnvVars = map[string]string

func Parse(r io.Reader) (EnvVars, error) {
	env, err := envparse.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse environment: %w", err)
	}

	return env, nil
}
