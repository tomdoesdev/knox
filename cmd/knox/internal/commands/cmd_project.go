package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/tomdoesdev/knox/internal"
	"github.com/tomdoesdev/knox/internal/workspace"
	"github.com/tomdoesdev/knox/kit/errs"
	"github.com/urfave/cli/v3"
)

func NewProjectCommand() *cli.Command {
	return &cli.Command{
		Name:  "project",
		Usage: "manage workspace projects",
		Commands: []*cli.Command{
			newProjectListCommand(),
			newProjectDeleteCommand(),
			newProjectListSecretsCommand(),
			newProjectAddSecretCommand(),
			newProjectRemoveSecretCommand(),
		},
	}
}

func newProjectListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "list all projects",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return projectListHandler()
		},
	}
}

func newProjectDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "delete a project",
		ArgsUsage: "<project-name>",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return errs.New(internal.ValidationCode, "project name is required")
			}

			name := cmd.Args().First()
			return projectDeleteHandler(name)
		},
	}
}

func newProjectListSecretsCommand() *cli.Command {
	return &cli.Command{
		Name:      "list-secrets",
		Usage:     "list secrets in a project",
		ArgsUsage: "[project-name]",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			var name string
			if cmd.Args().Len() == 1 {
				name = cmd.Args().First()
			}
			return projectListSecretsHandler(name)
		},
	}
}

func newProjectAddSecretCommand() *cli.Command {
	return &cli.Command{
		Name:      "add-secret",
		Usage:     "add a secret to a project",
		ArgsUsage: "<project-name> <logical-name> <vault-path@vault-alias>",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 3 {
				return errs.New(internal.ValidationCode, "project name, logical name, and secret reference are required")
			}

			projectName := cmd.Args().Get(0)
			logicalName := cmd.Args().Get(1)
			secretRef := cmd.Args().Get(2)

			return projectAddSecretHandler(projectName, logicalName, secretRef)
		},
	}
}

func newProjectRemoveSecretCommand() *cli.Command {
	return &cli.Command{
		Name:      "remove-secret",
		Usage:     "remove a secret from a project",
		ArgsUsage: "<project-name> <logical-name>",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 2 {
				return errs.New(internal.ValidationCode, "project name and logical name are required")
			}

			projectName := cmd.Args().Get(0)
			logicalName := cmd.Args().Get(1)

			return projectRemoveSecretHandler(projectName, logicalName)
		},
	}
}

// Handler functions

func projectListHandler() error {
	cwd, err := os.Getwd()
	if err != nil {
		return errs.Wrap(err, internal.SearchFailureCode, "failed to get current working directory")
	}

	ws, err := workspace.FindWorkspace(cwd)
	if err != nil {
		return errs.Wrap(err, internal.SearchFailureCode, "failed to find workspace")
	}

	projects, err := ws.ListProjects()
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		fmt.Println("No projects found")
		return nil
	}

	fmt.Printf("Projects (%d):\n", len(projects))
	for _, project := range projects {
		fmt.Printf("  %s\n", project)
	}

	return nil
}

func projectDeleteHandler(name string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return errs.Wrap(err, internal.SearchFailureCode, "failed to get current working directory")
	}

	ws, err := workspace.FindWorkspace(cwd)
	if err != nil {
		return errs.Wrap(err, internal.SearchFailureCode, "failed to find workspace")
	}

	err = ws.DeleteProject(name)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted project '%s'\n", name)
	return nil
}

func projectListSecretsHandler(projectName string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return errs.Wrap(err, internal.SearchFailureCode, "failed to get current working directory")
	}

	ws, err := workspace.FindWorkspace(cwd)
	if err != nil {
		return errs.Wrap(err, internal.SearchFailureCode, "failed to find workspace")
	}

	// If no project name provided, try to get current project
	if projectName == "" {
		// TODO: Get current project from workspace settings
		return errs.New(internal.ValidationCode, "project name is required (current project tracking not yet implemented)")
	}

	project, err := ws.LoadProject(projectName)
	if err != nil {
		return err
	}

	secrets := project.ListSecrets()
	if len(secrets) == 0 {
		fmt.Printf("Project '%s' has no secrets\n", projectName)
		return nil
	}

	fmt.Printf("Secrets in project '%s' (%d):\n", projectName, len(secrets))
	for _, logicalName := range secrets {
		secretRef, _ := project.GetSecret(logicalName)
		fmt.Printf("  %s -> %s\n", logicalName, secretRef)
	}

	return nil
}

func projectAddSecretHandler(projectName, logicalName, secretRef string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return errs.Wrap(err, internal.SearchFailureCode, "failed to get current working directory")
	}

	ws, err := workspace.FindWorkspace(cwd)
	if err != nil {
		return errs.Wrap(err, internal.SearchFailureCode, "failed to find workspace")
	}

	// Validate secret reference format
	_, err = workspace.ParseSecretReference(secretRef)
	if err != nil {
		return err
	}

	project, err := ws.LoadProject(projectName)
	if err != nil {
		return err
	}

	// Check if secret already exists
	if _, exists := project.GetSecret(logicalName); exists {
		return errs.New(internal.SecretExistsCode, "secret already exists in project").WithContext("logical_name", logicalName)
	}

	project.AddSecret(logicalName, secretRef)

	err = ws.UpdateProject(project)
	if err != nil {
		return err
	}

	fmt.Printf("Added secret '%s' -> '%s' to project '%s'\n", logicalName, secretRef, projectName)
	return nil
}

func projectRemoveSecretHandler(projectName, logicalName string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return errs.Wrap(err, internal.SearchFailureCode, "failed to get current working directory")
	}

	ws, err := workspace.FindWorkspace(cwd)
	if err != nil {
		return errs.Wrap(err, internal.SearchFailureCode, "failed to find workspace")
	}

	project, err := ws.LoadProject(projectName)
	if err != nil {
		return err
	}

	// Check if secret exists
	if _, exists := project.GetSecret(logicalName); !exists {
		return errs.New(internal.SecretNotFoundCode, "secret not found in project").WithContext("logical_name", logicalName)
	}

	project.RemoveSecret(logicalName)

	err = ws.UpdateProject(project)
	if err != nil {
		return err
	}

	fmt.Printf("Removed secret '%s' from project '%s'\n", logicalName, projectName)
	return nil
}
