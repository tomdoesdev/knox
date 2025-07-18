package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands"
	"github.com/tomdoesdev/knox/kit/log"
)

func main() {
	log.NewSlog("text")

	cmdRoot := commands.NewKnoxCommand()

	if err := cmdRoot.Run(context.Background(), os.Args); err != nil {
		slog.Error("knox.init.run", "error", err)
		os.Exit(1)
	}
}
