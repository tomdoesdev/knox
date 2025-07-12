package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands"
	"github.com/tomdoesdev/knox/kit/errkit"
	"github.com/urfave/cli/v3"
)

var (
	cmdInit = &cli.Command{
		Name:    "init",
		Aliases: nil,
		Usage:   "initialise new knox project vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return commands.Init()
		},
	}
	cmdStatus = &cli.Command{
		Name:    "status",
		Aliases: nil,
		Usage:   "show status of the current knox project vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return errkit.ErrNotImplemented
		},
	}
	cmdSet = &cli.Command{
		Name:    "set",
		Aliases: nil,
		Usage:   "set the value of a secret in the vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return errkit.ErrNotImplemented
		},
	}

	cmdRemove = &cli.Command{
		Name:    "remove",
		Aliases: nil,
		Usage:   "remove a secret from the vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return errkit.ErrNotImplemented
		},
	}

	cmdRoot = &cli.Command{
		Name:        "knox",
		Usage:       "manage local development vault",
		Description: "local development vault manager",
		Commands: []*cli.Command{
			cmdInit,
			cmdStatus,
			cmdSet,
			cmdRemove,
		},
	}
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := cmdRoot.Run(context.Background(), os.Args); err != nil {
		logger.Error(err.Error())
	}
}
