package commands

import (
	"context"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/internal"
	"github.com/urfave/cli/v3"
)

func NewStatusCommand() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "show workspace status information",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return internal.StatusHandler()
		},
	}
}
