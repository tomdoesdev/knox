package commands

import (
	"context"

	"github.com/tomdoesdev/knox/pkg/project"
	"github.com/urfave/cli/v3"
)

func NewStatusCommand(p *project.Project) *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "show status of the current knox project vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return statusAction(p)
		},
	}
}

func statusAction(p *project.Project) error {
	return nil
}
