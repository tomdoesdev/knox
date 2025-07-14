package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/tomdoesdev/knox/cmd/knox/internal"
	"github.com/tomdoesdev/knox/internal/config"
	"github.com/tomdoesdev/knox/internal/project"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var appConfig = config.NewApplicationConfig()

	cmdRoot := internal.NewKnoxCommand(&appConfig)

	if err := cmdRoot.Run(context.Background(), os.Args); err != nil {
		if errors.Is(err, project.ErrProjectExists) {
			logger.Info("project already exists")
			os.Exit(0)
		}
		logger.Error(err.Error())
		os.Exit(1)
	}
}
