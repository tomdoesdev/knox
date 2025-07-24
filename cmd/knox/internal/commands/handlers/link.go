package handlers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tomdoesdev/knox/internal/error_codes"
	"github.com/tomdoesdev/knox/internal/vault"
	"github.com/tomdoesdev/knox/internal/workspace"
	"github.com/tomdoesdev/knox/kit/errs"
)

func LinkVaultHandler(vaultPath, alias string) error {
	// Get current working directory to find workspace
	cwd, err := os.Getwd()
	if err != nil {
		return errs.Wrap(err, error_codes.SearchFailureErrCode, "failed to get current working directory")
	}

	// Find workspace
	ws, err := workspace.FindWorkspace(cwd)
	if err != nil {
		return errs.Wrap(err, error_codes.SearchFailureErrCode, "failed to find workspace")
	}

	// Convert relative path to absolute path
	absPath, err := filepath.Abs(vaultPath)
	if err != nil {
		return errs.Wrap(err, error_codes.ValidationErrCode, "failed to resolve vault path").WithContext("path", vaultPath)
	}

	// Verify the vault exists and is valid
	exists, err := vault.IsVault(absPath)
	if err != nil {
		return errs.Wrap(err, error_codes.VaultConnectionErrCode, "failed to check vault").WithContext("path", absPath)
	}
	if !exists {
		return errs.New(error_codes.VaultConnectionErrCode, "no valid vault found at path").WithContext("path", absPath)
	}

	// Link the vault to the workspace
	err = ws.LinkVault(alias, absPath)
	if err != nil {
		return errs.Wrap(err, error_codes.VaultCreationErrCode, "failed to link vault to workspace")
	}

	fmt.Printf("Linked vault at %s with alias '%s'\n", absPath, alias)
	return nil
}
