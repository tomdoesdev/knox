package commands

import (
	"context"
	"github.com/tomdoesdev/knox/internal/errors"
	"github.com/urfave/cli/v3"
)

var (
	cmdInit = &cli.Command{
		Name:    "init",
		Aliases: nil,
		Usage:   "initialise new knox project vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return errors.ErrNotImplemented
		},
	}
	cmdStatus = &cli.Command{
		Name:    "status",
		Aliases: nil,
		Usage:   "show status of the current knox project vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return errors.ErrNotImplemented
		},
	}
	cmdSet = &cli.Command{
		Name:    "set",
		Aliases: nil,
		Usage:   "set the value of a secret in the vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return errors.ErrNotImplemented
		},
	}

	cmdRemove = &cli.Command{
		Name:    "remove",
		Aliases: nil,
		Usage:   "remove a secret from the vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return errors.ErrNotImplemented
		},
	}
)

func NewRootCommand() *cli.Command {
	cmd := &cli.Command{
		Name:        "knox",
		Usage:       "manage local development secrets",
		Description: "local development secrets manager",
		Commands: []*cli.Command{
			cmdInit,
			cmdStatus,
			cmdSet,
			cmdRemove,
		},
	}

	return cmd
}
