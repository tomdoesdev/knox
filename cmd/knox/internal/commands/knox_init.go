package commands

import (
	"context"

	"github.com/tomdoesdev/knox/internal/project"
	"github.com/tomdoesdev/knox/pkg/errs"
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
	_, err := project.CreateFile("")
	if err != nil {
		// If it's already a Knox error, return it as-is
		if errs.Is(err, project.CreateFailureCode) {
			return err
		}
		// Otherwise wrap it with a more appropriate error
		return errs.Wrap(err, project.CreateFailureCode, "failed to create project file")
	}

	return nil
}
