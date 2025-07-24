package handlers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/common"
	"github.com/tomdoesdev/knox/internal/error_codes"
	"github.com/tomdoesdev/knox/internal/vault"
	"github.com/tomdoesdev/knox/internal/workspace"
	"github.com/tomdoesdev/knox/kit/errs"
	"github.com/tomdoesdev/knox/kit/fs"
)

func NewProjectHandler(name, description string) error {

	return common.WithLocalWorkspace(func(ws *workspace.Workspace) error {
		project := workspace.NewProject(name, description)

		err := ws.CreateProject(project)
		if err != nil {
			return err
		}

		fmt.Printf("Created project '%s'\n", name)
		if description != "" {
			fmt.Printf("Description: %s\n", description)
		}
		fmt.Printf("Project file: .knox-workspace/projects/%s.json\n", name)

		return nil
	})

}

func NewVaultHandler(alias, vaultPath string) error {
	// If no path provided, generate default path
	if vaultPath == "" {
		defaultPath, err := getDefaultVaultPath(alias)
		if err != nil {
			return errs.Wrap(err, error_codes.ValidationErrCode, "failed to generate default vault path")
		}
		vaultPath = defaultPath
	}

	return common.WithLocalWorkspace(func(ws *workspace.Workspace) error {
		// Convert relative path to absolute path
		absPath, err := filepath.Abs(vaultPath)
		if err != nil {
			return errs.Wrap(err, error_codes.ValidationErrCode, "failed to resolve vault path").WithContext("path", vaultPath)
		}

		// Check if vault already exists at this path
		if exists, err := vault.IsVault(absPath); err == nil && exists {
			return errs.New(error_codes.VaultCreationErrCode, "vault already exists at this path").WithContext("path", absPath)
		}

		// Create the vault directory if it doesn't exist
		vaultDir := filepath.Dir(absPath)
		if err := ensureVaultDirectory(vaultDir); err != nil {
			return errs.Wrap(err, error_codes.VaultCreationErrCode, "failed to create vault directory").WithContext("path", vaultDir)
		}

		// Check if vault file already exists
		if fs.IsFile(absPath) {
			if exists, _ := vault.IsVault(absPath); exists {
				return errs.New(error_codes.VaultCreationErrCode, "vault already exists at this path").WithContext("path", absPath)
			}
		}

		// Create a simple vault file by manually creating the SQLite database
		// This is a simplified approach until we have proper vault creation API
		err = createVaultAtPath(absPath)
		if err != nil {
			return errs.Wrap(err, error_codes.VaultCreationErrCode, "failed to create vault")
		}

		// Link the vault to the workspace
		err = ws.LinkVault(alias, absPath)
		if err != nil {
			// Vault was created but linking failed - inform user
			fmt.Printf("Created vault '%s' at %s\n", alias, absPath)
			fmt.Printf("Default 'global' collection created\n")
			fmt.Printf("Warning: Failed to link vault to workspace: %s\n", err.Error())
			fmt.Printf("You can link it manually with: knox link %s --alias %s\n", absPath, alias)
			return nil
		}

		fmt.Printf("Created vault '%s' at %s\n", alias, absPath)
		fmt.Printf("Default 'global' collection created\n")
		fmt.Printf("Vault linked to workspace with alias '%s'\n", alias)

		return nil
	})

}

// createVaultAtPath creates a vault file at the specified path
func createVaultAtPath(path string) error {
	// Use a temporary KNOX_ROOT to create vault at specific location
	originalRoot := os.Getenv("KNOX_ROOT")

	// Set KNOX_ROOT to parent directory and create .knox/vault.db structure
	dir := filepath.Dir(path)

	// Create a temporary environment to use NewFileSystemDatasource
	tempKnoxDir := filepath.Join(dir, ".temp-knox")
	_ = os.Setenv("KNOX_ROOT", tempKnoxDir)

	defer func() {
		if originalRoot != "" {
			_ = os.Setenv("KNOX_ROOT", originalRoot)
		} else {
			_ = os.Unsetenv("KNOX_ROOT")
		}
		// Clean up temp directory
		_ = os.RemoveAll(tempKnoxDir)
	}()

	// Create the datasource which will create the vault
	datasourceProvider, err := vault.NewFileSystemDatasource()
	if err != nil {
		return err
	}

	// Open to ensure it's created properly
	v, err := vault.Open(datasourceProvider)
	if err != nil {
		return err
	}
	_ = v.Close()

	// Move the created vault.db to our desired location
	tempVaultPath := filepath.Join(tempKnoxDir, ".knox", "vault.db")
	return os.Rename(tempVaultPath, path)
}

// getDefaultVaultPath generates the default vault path based on KNOX_ROOT or user home
func getDefaultVaultPath(alias string) (string, error) {
	var baseDir string

	// Check for KNOX_ROOT environment variable first
	if knoxRoot := os.Getenv("KNOX_ROOT"); knoxRoot != "" {
		baseDir = knoxRoot
	} else {
		// Fall back to user home directory
		home, err := os.UserHomeDir()
		if err != nil {
			return "", errs.Wrap(err, error_codes.ValidationErrCode, "failed to get user home directory")
		}
		baseDir = filepath.Join(home, ".knox")
	}

	// Create vaults directory path
	vaultsDir := filepath.Join(baseDir, "vaults")

	// Generate filename using alias
	filename := alias + ".db"

	return filepath.Join(vaultsDir, filename), nil
}

// ensureVaultDirectory ensures the vault directory exists and handles conflicts
func ensureVaultDirectory(dir string) error {
	// Check if path exists
	info, err := os.Stat(dir)
	if err == nil {
		// Path exists - check if it's a directory
		if !info.IsDir() {
			return errs.New(error_codes.VaultCreationErrCode, "vault directory path exists but is not a directory").WithContext("path", dir)
		}
		// Directory already exists, we're good
		return nil
	}

	// Path doesn't exist - check if parent directories have conflicts
	parent := filepath.Dir(dir)
	if parent != dir { // Avoid infinite recursion at root
		if err := ensureVaultDirectory(parent); err != nil {
			return err
		}
	}

	// Create the directory
	if err := os.Mkdir(dir, 0755); err != nil {
		return errs.Wrap(err, error_codes.VaultCreationErrCode, "failed to create directory")
	}

	return nil
}
