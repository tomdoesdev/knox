package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tomdoesdev/knox/cmd/v2/knox/internal/commands"
	"github.com/tomdoesdev/knox/internal/v2/vault"
	"github.com/tomdoesdev/knox/kit/log"
	"github.com/urfave/cli/v3"
)

func main() {

	log.NewSlog("text")

	_, err := vault.OpenFileSystem()
	if err != nil {
		fmt.Println(err)
	}

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
