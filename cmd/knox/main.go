package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands"
	"github.com/tomdoesdev/knox/internal"
	vault2 "github.com/tomdoesdev/knox/internal/vault"
	"github.com/tomdoesdev/knox/kit/errs"
	"github.com/tomdoesdev/knox/kit/log"
	"github.com/urfave/cli/v3"
)

func main() {
	log.NewSlog("text")

	v, err := vault2.OpenFileSystem()
	if err != nil {
		if errs.Is(err, internal.VaultConnectionCode) {
			fmt.Printf("sqlite connection failed - check permissions: %s\n", err.Error())
		} else if errs.Is(err, internal.VaultCreationCode) {
			fmt.Printf("failed to create vault - check direcroty permissions: %s\n", err.Error())
		} else {
			fmt.Printf("failed to open vault: %s\n", err.Error())
		}

		os.Exit(1)
	}

	defer func(v *vault2.Vault) {
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
			commands.NewNewCommand(),
			commands.NewProjectCommand(),
			commands.NewLinkCommand(),
			commands.NewStatusCommand(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
