package main

import (
	"context"
	"github.com/tomdoesdev/knox/internal/cli"
	"log/slog"

	"log"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cmd := cli.NewRootCommand()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
