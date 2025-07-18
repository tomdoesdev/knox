package commands

import (
	"github.com/urfave/cli/v3"
)

func NewSecretsCommand() *cli.Command {
	return &cli.Command{
		Name:    "secrets",
		Usage:   "manage project secrets",
		Aliases: []string{"sec", "s"},
		Commands: []*cli.Command{
			NewAddCommand(),
			NewGetCommand(),
			NewRemoveCommand(),
		},
	}
}
