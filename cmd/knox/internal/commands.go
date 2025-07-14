package internal

import (
	"context"
	"errors"
	"log/slog"

	"github.com/tomdoesdev/knox/internal/config"
	"github.com/tomdoesdev/knox/kit/errkit"
	"github.com/urfave/cli/v3"
)

func NewKnoxCommand(appConfig *config.ApplicationConfig) *cli.Command {
	var (
		cmdInit = &cli.Command{
			Name:    "init",
			Aliases: nil,
			Usage:   "initialise new knox project vault",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				err := Initialize(appConfig)

				if errors.Is(err, ErrAlreadyInitialized) {
					slog.Warn("knox is already initialized")
					return nil
				}

				return err
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
	)

	return &cli.Command{
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
}
