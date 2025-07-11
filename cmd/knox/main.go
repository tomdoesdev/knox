package main

import (
	"context"
	"github.com/tomdoesdev/knox/cmd/knox/internal/commands"
	"log/slog"

	"log"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cmd := commands.NewRootCommand()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
