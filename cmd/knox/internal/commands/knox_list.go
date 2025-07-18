package commands

import (
	"context"
	"fmt"

	secrets2 "github.com/tomdoesdev/knox/internal/secrets"
	"github.com/urfave/cli/v3"
)

func NewListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Flags: []cli.Flag{&cli.BoolFlag{Name: "print-secrets", Value: false, Usage: "print both keys and secrets to stdout"}},
		Usage: "list the secrets for the current project vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return listActionHandler(cmd.Bool("print-secrets"))
		},
	}
}

func listActionHandler(printSecrets bool) error {
	knox, err := LoadKnoxContext()
	if err != nil {
		return err
	}

	var secrets []secrets2.Secret

	if printSecrets {
		s, err := knox.Store.ListSecrets()
		if err != nil {
			return err
		}
		for _, ss := range s {
			secrets = append(secrets, ss)
		}
	} else {
		s, err := knox.Store.ListKeys()
		if err != nil {
			return err
		}

		for _, ss := range s {
			secrets = append(secrets, secrets2.Secret{Key: ss, Value: ""})
		}

	}

	if len(secrets) == 0 {
		fmt.Println("No secrets for current project")
	}

	for _, kv := range secrets {
		if printSecrets {
			fmt.Printf("%s=%s", kv.Key, kv.Value)
		} else {
			fmt.Println(kv.Key)
		}
	}

	return nil
}
