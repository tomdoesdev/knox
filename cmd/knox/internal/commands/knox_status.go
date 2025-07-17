package commands

import (
	"context"

	"github.com/tomdoesdev/knox/internal"
	"github.com/urfave/cli/v3"
)

func NewStatusCommand(k *internal.Knox) *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "show status of the current knox project vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return statusAction(k)
		},
	}
}

func statusAction(k *internal.Knox) error {
	return nil
}
