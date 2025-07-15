package internal

import (
	"errors"
	"fmt"

	"github.com/tomdoesdev/knox/internal/project"
	"github.com/urfave/cli/v3"
)

func Initialize() error {
	_, err := project.CreateFile()
	if err != nil {
		return fmt.Errorf("knox.init.createProjectFile: %w", err)
	}

	return nil
}

func SetSecret(cmd *cli.Command, p *project.Project) error {
	if cmd.Args().Len() != 2 {
		return fmt.Errorf("expected 2 arguments, got %d", cmd.Args().Len())
	}

	key := cmd.Args().Get(0)
	value := cmd.Args().Get(1)

	err := p.WriteSecret(key, value)
	if errors.Is(err, project.ErrSecretExists) {
		fmt.Printf("secret already exists: %s\n", key)
		return nil
	}

	return err

}
