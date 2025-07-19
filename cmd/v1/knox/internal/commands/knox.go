package commands

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v3"
)

func NewKnoxCommand() *cli.Command {
	return &cli.Command{
		Name:        "knox",
		Usage:       "manage local development vault",
		Description: "local development vault manager",
		Commands: []*cli.Command{
			NewInitCommand(),

			NewRunCommand(),
			NewSecretsCommand(),
		},
	}
}
