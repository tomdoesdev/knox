package commands

import (
	"context"
	"fmt"

	"github.com/tomdoesdev/knox/internal/v1"
	"github.com/tomdoesdev/knox/pkg/errs"
	"github.com/urfave/cli/v3"
)

func NewAddCommand() *cli.Command {
	var key string
	var value string

	return &cli.Command{
		Name:  "add",
		Usage: "add a secret to the current knox vault",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "key",
				Destination: &key,
			},
			&cli.StringArg{
				Name:        "value",
				Destination: &value,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			force := cmd.Bool("force")
			k, err := LoadKnoxContextWithOptions(force)
			if err != nil {
				return err
			}
			return addActionHandler(key, value, k)
		},
	}
}

func addActionHandler(key, value string, k *v1.Knox) error {
	if key == "" {
		return errs.Wrap(
			ErrInvalidArguments,
			InvalidArguments,
			"key argument is required").WithContext("missing", "key")
	}

	if value == "" {
		return errs.Wrap(
			ErrInvalidArguments,
			InvalidArguments,
			"value argument is required").WithContext("missing", "value")
	}

	err := k.Store.WriteSecret(key, value)
	if err != nil {
		return err
	}

	fmt.Printf("secret added: %s\n", key)
	return nil
}
