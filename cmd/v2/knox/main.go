package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tomdoesdev/knox/cmd/v2/knox/internal/commands"
	"github.com/urfave/cli/v3"
)

func main() {
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
