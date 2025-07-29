package internal

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func NewInitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "initialize a new workspace",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			r, err := initHandler()
			if err != nil {
				return err
			}

			fmt.Println(r)

			return nil
		},
	}
}

func initHandler() (fmt.Stringer, error) {

	return nil, nil
}
