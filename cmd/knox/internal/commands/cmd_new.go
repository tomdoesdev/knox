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
	"github.com/tomdoesdev/knox/kit/fs"
	"github.com/urfave/cli/v3"
)

func NewNewCommand() *cli.Command {
	return &cli.Command{
		Name:  "new",
		Usage: "create new resources",
		Commands: []*cli.Command{
			newNewProjectCommand(),
			newNewVaultCommand(),
			// TODO: newNewSecretCommand(), etc.
		},
	}
}

func newNewProjectCommand() *cli.Command {
	return &cli.Command{
		Name:      "project",
		Usage:     "create a new project",
		ArgsUsage: "<project-name>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "description",
				Usage: "project description",
				Value: "",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return errs.New(internal.ValidationCode, "project name is required")
			}

			name := cmd.Args().First()
			description := cmd.String("description")

			return newProjectHandler(name, description)
		},
	}
}

func newProjectHandler(name, description string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return errs.Wrap(err, internal.SearchFailureCode, "failed to get current working directory")
	}

	ws, err := workspace.FindWorkspace(cwd)
	if err != nil {
		return errs.Wrap(err, internal.SearchFailureCode, "failed to find workspace")
	}

	project := workspace.NewProject(name, description)

	err = ws.CreateProject(project)
	if err != nil {
		return err
	}

	fmt.Printf("Created project '%s'\n", name)
	if description != "" {
		fmt.Printf("Description: %s\n", description)
	}
	fmt.Printf("Project file: .knox-workspace/projects/%s.json\n", name)

	return nil
}

func newNewVaultCommand() *cli.Command {
	return &cli.Command{
		Name:      "vault",
		Usage:     "create a new vault",
		ArgsUsage: "<vault-alias>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Usage:    "path where vault will be created",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return errs.New(internal.ValidationCode, "vault alias is required")
			}

			alias := cmd.Args().First()
			path := cmd.String("path")

			return newVaultHandler(alias, path)
		},
	}
}

func newVaultHandler(alias, vaultPath string) error {
	// Get current working directory to find workspace
	cwd, err := os.Getwd()
	if err != nil {
		return errs.Wrap(err, internal.SearchFailureCode, "failed to get current working directory")
	}

	// Find workspace for vault linking
	ws, err := workspace.FindWorkspace(cwd)
	if err != nil {
		return errs.Wrap(err, internal.SearchFailureCode, "failed to find workspace")
	}

	// Convert relative path to absolute path
	absPath, err := filepath.Abs(vaultPath)
	if err != nil {
		return errs.Wrap(err, internal.ValidationCode, "failed to resolve vault path").WithContext("path", vaultPath)
	}

	// Check if vault already exists at this path
	if exists, err := vault.IsVault(absPath); err == nil && exists {
		return errs.New(internal.VaultCreationCode, "vault already exists at this path").WithContext("path", absPath)
	}

	// Create the vault directory if it doesn't exist
	vaultDir := filepath.Dir(absPath)
	if err := os.MkdirAll(vaultDir, 0755); err != nil {
		return errs.Wrap(err, internal.VaultCreationCode, "failed to create vault directory").WithContext("path", vaultDir)
	}

	// Check if vault file already exists
	if fs.IsFile(absPath) {
		if exists, _ := vault.IsVault(absPath); exists {
			return errs.New(internal.VaultCreationCode, "vault already exists at this path").WithContext("path", absPath)
		}
	}

	// Create a simple vault file by manually creating the SQLite database
	// This is a simplified approach until we have proper vault creation API
	err = createVaultAtPath(absPath)
	if err != nil {
		return errs.Wrap(err, internal.VaultCreationCode, "failed to create vault")
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
