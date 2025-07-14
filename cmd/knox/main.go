package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/tomdoesdev/knox/cmd/knox/internal"
	"github.com/tomdoesdev/knox/internal/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var appConfig = config.NewApplicationConfig()

	cmdRoot := internal.NewKnoxCommand(&appConfig)

	if err := cmdRoot.Run(context.Background(), os.Args); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
