package commands

import (
	"context"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v3"
)

func NewKnoxCommand() *cli.Command {
	return &cli.Command{
		Name:        "knox",
		Usage:       "manage local development vault",
		Description: "local development vault manager",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("knox run")
			return nil
		},
		Commands: []*cli.Command{
			NewInitCommand(),
			NewStatusCommand(),
			NewAddCommand(),
			NewRemoveCommand(),
			NewRunCommand(),
			NewGetCommand(),
		},
	}
}
