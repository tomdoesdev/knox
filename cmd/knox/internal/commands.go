package internal

import (
	"context"

	"github.com/tomdoesdev/knox/internal/project"
	"github.com/tomdoesdev/knox/kit/errkit"
	"github.com/urfave/cli/v3"
)

func NewKnoxCommand(p *project.Project) *cli.Command {
	return &cli.Command{
		Name:        "knox",
		Usage:       "manage local development vault",
		Description: "local development vault manager",

		Commands: []*cli.Command{
			newInitCommand(p),
			newStatusCommand(p),
			newSetCommand(p),
			newRemoveCommand(p),
		},
	}
}

func newInitCommand(p *project.Project) *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "initialise new knox project vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return Initialize()
		},
	}
}

func newStatusCommand(p *project.Project) *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "show status of the current knox project vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return errkit.ErrNotImplemented
		},
	}
}

func newSetCommand(p *project.Project) *cli.Command {
	return &cli.Command{
		Name:  "set",
		Usage: "set the value of a secret in the vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return SetSecret(cmd, p)
		},
	}
}

func newRemoveCommand(p *project.Project) *cli.Command {
	return &cli.Command{
		Name:  "remove",
		Usage: "remove a secret from the vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return errkit.ErrNotImplemented
		},
	}
}
