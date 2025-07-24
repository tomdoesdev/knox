package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands"
	"github.com/tomdoesdev/knox/kit/log"
	"github.com/urfave/cli/v3"
)

func main() {
	log.NewSlog(log.Text)

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
