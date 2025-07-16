package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands"
	"github.com/tomdoesdev/knox/internal/vault"
	"github.com/tomdoesdev/knox/kit/log"
	"github.com/tomdoesdev/knox/pkg/project"
)

func main() {
	log.NewSlog("text")

	p, err := project.Load()
	if err != nil {
		slog.Error("knox.init.loadProject", "error", err)
		os.Exit(1)
	}

	slog.Debug("project loaded",
		slog.String("project.vaultpath", p.VaultPath),
	)

	vaultpath := p.VaultPath

	if !vault.Exists(vaultpath) {
		if !confirmCreateDefaultVault(vaultpath) {
			fmt.Println("Aborted: Knox vault creation declined.")
			os.Exit(0)
		}

		vaultpath = vault.DefaultPath()
	}

	err = vault.EnsureVaultExists(vaultpath)
	if err != nil {
		slog.Error("knox.init.ensureVaultExists", "error", err)
		os.Exit(1)
	}

	cmdRoot := commands.NewKnoxCommand(p)

	if err := cmdRoot.Run(context.Background(), os.Args); err != nil {
		if errors.Is(err, project.ErrProjectExists) {
			slog.Info("project already exists")
			os.Exit(0)
		}
		os.Exit(1)
	}
}

func confirmCreateDefaultVault(p string) bool {
	fmt.Printf("Knox will create a default vault at: %s\n", p)
	fmt.Print("Do you want to proceed? [y/N]: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false
	}

	response := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return response == "y" || response == "yes"
}
