package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands"
	"github.com/tomdoesdev/knox/internal/project"
	"github.com/tomdoesdev/knox/kit/log"
)

func main() {
	log.NewSlog("text")

	cmdRoot := commands.NewKnoxCommand()

	if err := cmdRoot.Run(context.Background(), os.Args); err != nil {
		if errors.Is(err, project.ErrProjectExists) {
			slog.Info("project already exists")
			os.Exit(0)
		}

		slog.Error("knox.init.run", "error", err)
		os.Exit(1)
	}
}
