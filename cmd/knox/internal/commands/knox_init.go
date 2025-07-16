package commands

import (
	"context"
	"fmt"

	"github.com/tomdoesdev/knox/pkg/project"
	"github.com/urfave/cli/v3"
)

func NewInitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "initialise new knox project vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return initialise()
		},
	}
}

func initialise() error {
	_, err := project.CreateFile()
	if err != nil {
		return fmt.Errorf("knox.init.createProjectFile: %w", err)
	}

	return nil
}
