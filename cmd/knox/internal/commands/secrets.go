package commands

import (
	"github.com/urfave/cli/v3"
)

func NewSecretsCommand() *cli.Command {
	return &cli.Command{
		Name:    "secrets",
		Usage:   "manage project secrets",
		Aliases: []string{"sec", "s"},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "force",
				Usage: "force creation of vault file even if it would be placed in a version control directory",
				Value: false,
			},
		},
		Commands: []*cli.Command{
			NewAddCommand(),
			NewGetCommand(),
			NewRemoveCommand(),
			NewListCommand(),
		},
	}
}
