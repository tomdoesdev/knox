package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tomdoesdev/knox/internal"
	"github.com/tomdoesdev/knox/internal/vault"
	"github.com/tomdoesdev/knox/internal/workspace"
	"github.com/tomdoesdev/knox/kit/errs"
	"github.com/urfave/cli/v3"
)

func NewLinkCommand() *cli.Command {
	return &cli.Command{
		Name:      "link",
		Usage:     "link a vault to the current workspace",
		ArgsUsage: "<vault-path>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "alias",
				Usage:    "alias for the vault in the workspace",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return errs.New(internal.ValidationCode, "vault path is required")
			}

			vaultPath := cmd.Args().First()
			alias := cmd.String("alias")

			return linkVaultHandler(vaultPath, alias)
		},
	}
}

func linkVaultHandler(vaultPath, alias string) error {
	// Get current working directory to find workspace
	cwd, err := os.Getwd()
	if err != nil {
		return errs.Wrap(err, internal.SearchFailureCode, "failed to get current working directory")
	}

	// Find workspace
	ws, err := workspace.FindWorkspace(cwd)
	if err != nil {
		return errs.Wrap(err, internal.SearchFailureCode, "failed to find workspace")
	}

	// Convert relative path to absolute path
	absPath, err := filepath.Abs(vaultPath)
	if err != nil {
		return errs.Wrap(err, internal.ValidationCode, "failed to resolve vault path").WithContext("path", vaultPath)
	}

	// Verify the vault exists and is valid
	exists, err := vault.IsVault(absPath)
	if err != nil {
		return errs.Wrap(err, internal.VaultConnectionCode, "failed to check vault").WithContext("path", absPath)
	}
	if !exists {
		return errs.New(internal.VaultConnectionCode, "no valid vault found at path").WithContext("path", absPath)
	}

	// Link the vault to the workspace
	err = ws.LinkVault(alias, absPath)
	if err != nil {
		return errs.Wrap(err, internal.VaultCreationCode, "failed to link vault to workspace")
	}

	fmt.Printf("Linked vault at %s with alias '%s'\n", absPath, alias)
	return nil
}
