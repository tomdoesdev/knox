package commands

import (
	"context"
	"fmt"

	"github.com/tomdoesdev/knox/internal"
	"github.com/urfave/cli/v3"
)

func NewGetCommand() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "get a secret to the current knox vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			force := cmd.Bool("force")
			k, err := LoadKnoxContextWithOptions(force)
			if err != nil {
				return err
			}
			return getActionHandler(cmd, k)
		},
	}
}

func getActionHandler(cmd *cli.Command, k *internal.Knox) error {
	if cmd.Args().Len() != 1 {
		return fmt.Errorf("expected 1 arguments, got %d", cmd.Args().Len())
	}

	key := cmd.Args().Get(0)

	secret, err := k.Store.ReadSecret(key)
	if err != nil {
		return fmt.Errorf("knox_get.getAction.readStore: %w", err)
	}

	_, err = fmt.Printf("%s\n", secret)

	if err != nil {
		return fmt.Errorf("knox_get.getAction.readSecret: %w", err)
	}

	return nil
}
