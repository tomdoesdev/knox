package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cmd := commands.NewRootCommand()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logger.Error(err.Error())
	}
}
