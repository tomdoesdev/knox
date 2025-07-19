package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/tomdoesdev/knox/cmd/v1/knox/internal/commands"
	"github.com/tomdoesdev/knox/internal/v1/vault"
	"github.com/tomdoesdev/knox/kit/log"
)

func main() {
	log.NewSlog("text")

	vloc, err := vault.FindLocation()
	if err != nil {
		// Add better error handling
		fmt.Println("Error checking for vault")
		os.Exit(1)
	}

	_, err = vault.ExistsOrCreate(vloc)
	if err != nil {
		//Add better error
		fmt.Println("Error creating vault")
	}

	cmdRoot := commands.NewKnoxCommand()

	if err := cmdRoot.Run(context.Background(), os.Args); err != nil {
		slog.Error("knox.init.run", "error", err)
		os.Exit(1)
	}
}
