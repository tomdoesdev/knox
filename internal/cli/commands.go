package cli

import (
	"context"
	"github.com/tomdoesdev/knox/internal/errors"
	"github.com/urfave/cli/v3"
)

func NewRootCommand() *cli.Command {
	cmd := &cli.Command{
		Name:        "knox",
		Usage:       "manage local development secrets",
		Description: "local development secrets manager",
		Commands: []*cli.Command{
			initCommand(),
		},
	}

	return cmd
}

func initCommand() *cli.Command {
	return &cli.Command{
		Name:    "init",
		Aliases: nil,
		Usage:   "initialise new knox project vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return errors.ErrNotImplemented
		},
	}
}
