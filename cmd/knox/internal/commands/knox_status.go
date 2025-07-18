package commands

import (
	"context"

	"github.com/urfave/cli/v3"
)

func NewStatusCommand() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "show status of the current knox project vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return statusAction()
		},
	}
}

func statusAction() error {
	return nil
}
