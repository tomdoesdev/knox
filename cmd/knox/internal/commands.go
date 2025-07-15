package internal

import (
	"context"

	"github.com/tomdoesdev/knox/internal/config"
	"github.com/tomdoesdev/knox/kit/errkit"
	"github.com/urfave/cli/v3"
)

func NewKnoxCommand(appConfig *config.ApplicationConfig) *cli.Command {
	return &cli.Command{
		Name:        "knox",
		Usage:       "manage local development vault",
		Description: "local development vault manager",
		Commands: []*cli.Command{
			newInitCommand(appConfig),
			newStatusCommand(appConfig),
			newSetCommand(appConfig),
			newRemoveCommand(appConfig),
		},
	}
}

func newInitCommand(appConfig *config.ApplicationConfig) *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "initialise new knox project vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return Initialize(appConfig)
		},
	}
}

func newStatusCommand(appConfig *config.ApplicationConfig) *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "show status of the current knox project vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return errkit.ErrNotImplemented
		},
	}
}

func newSetCommand(appConfig *config.ApplicationConfig) *cli.Command {
	return &cli.Command{
		Name:  "set",
		Usage: "set the value of a secret in the vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return errkit.ErrNotImplemented
		},
	}
}

func newRemoveCommand(appConfig *config.ApplicationConfig) *cli.Command {
	return &cli.Command{
		Name:  "remove",
		Usage: "remove a secret from the vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return errkit.ErrNotImplemented
		},
	}
}
