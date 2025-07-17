package commands

import (
	"context"
	"fmt"

	"github.com/tomdoesdev/knox/internal"
	"github.com/urfave/cli/v3"
)

func NewRemoveCommand(k *internal.Knox) *cli.Command {
	return &cli.Command{
		Name:    "remove",
		Aliases: []string{"rm"},
		Usage:   "remove a secret from the vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return removeActionfunc(cmd, k)
		},
	}
}

func removeActionfunc(cmd *cli.Command, k *internal.Knox) error {
	if cmd.Args().Len() != 1 {
		return fmt.Errorf("expected 1 argument, got %d", cmd.Args().Len())
	}

	key := cmd.Args().Get(0)

	err := k.Store.DeleteSecret(key)
	if err != nil {
		return fmt.Errorf("removing secret: %w", err)
	}

	fmt.Println("removed secret:", key)
	return nil

}
