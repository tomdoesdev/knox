package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tomdoesdev/knox/cmd/v2/knox/internal/commands"
	"github.com/tomdoesdev/knox/internal/v2/vault"
	"github.com/tomdoesdev/knox/kit/log"
	"github.com/tomdoesdev/knox/pkg/errs"
	"github.com/urfave/cli/v3"
)

func main() {
	log.NewSlog("text")

	v, err := vault.OpenFileSystem()
	if err != nil {
		if errs.Is(err, vault.VaultConnectionCode) {
			fmt.Printf("sqlite connection failed - check permissions: %s\n", err.Error())
		} else if errs.Is(err, vault.VaultCreationCode) {
			fmt.Printf("failed to create vault - check direcroty permissions: %s\n", err.Error())
		} else {
			fmt.Printf("failed to open vault: %s\n", err.Error())
		}

		os.Exit(1)
	}

	defer func(v *vault.Vault) {
		err := v.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(v)

	app := &cli.Command{
		Name:  "knox",
		Usage: "local development secrets manager",
		Commands: []*cli.Command{
			commands.NewInitCommand(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
