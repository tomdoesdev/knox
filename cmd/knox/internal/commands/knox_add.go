package commands

import (
	"context"
	"fmt"

	"github.com/tomdoesdev/knox/internal"
	"github.com/tomdoesdev/knox/pkg/errs"
	"github.com/urfave/cli/v3"
)

func NewAddCommand() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "add a secret to the current knox vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			k, err := LoadKnoxContext()
			if err != nil {
				return err
			}
			return addActionHandler(cmd, k)
		},
	}
}

func addActionHandler(cmd *cli.Command, k *internal.Knox) error {
	if cmd.Args().Len() != 2 {
		return errs.Wrap(ErrInvalidArguments, InvalidArguments, fmt.Sprintf("expected 2 arguments, got %d", cmd.Args().Len()))
	}

	key := cmd.Args().Get(0)
	value := cmd.Args().Get(1)

	err := k.Store.WriteSecret(key, value)
	if err != nil {
		return err
	}

	fmt.Printf("secret added: %s\n", key)
	return nil
}
