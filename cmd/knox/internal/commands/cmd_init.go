package commands

import (
	"context"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/handlers"
	"github.com/urfave/cli/v3"
)

func NewInitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "initialize a new workspace",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// initialise a workspace if no sub command is specified
			return handlers.InitHandler()
		},
	}
}
