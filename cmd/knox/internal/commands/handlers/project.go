package handlers

import (
	"fmt"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/common"
	"github.com/tomdoesdev/knox/internal/error_codes"
	"github.com/tomdoesdev/knox/internal/workspace"
	"github.com/tomdoesdev/knox/kit/errs"
)

func ProjectListHandler() error {
	return common.WithLocalWorkspace(func(ws *workspace.Workspace) error {
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
	})

}

func ProjectDeleteHandler(name string) error {
	return common.WithLocalWorkspace(func(ws *workspace.Workspace) error {
		err := ws.DeleteProject(name)
		if err != nil {
			return err
		}

		fmt.Printf("Deleted project '%s'\n", name)
		return nil
	})
}

func ProjectListSecretsHandler(projectName string) error {

	return common.WithLocalWorkspace(func(ws *workspace.Workspace) error {

		// If no project name provided, try to get current project
		if projectName == "" {
			// TODO: Get current project from workspace settings
			return errs.New(error_codes.ValidationErrCode, "project name is required (current project tracking not yet implemented)")
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
	})
}
func ProjectAddSecretHandler(projectName, logicalName, secretRef string) error {
	return common.WithLocalWorkspace(func(ws *workspace.Workspace) error {
		// Validate secret reference format
		_, err := workspace.ParseSecretReference(secretRef)
		if err != nil {
			return err
		}

		project, err := ws.LoadProject(projectName)
		if err != nil {
			return err
		}

		// Check if secret already exists
		if _, exists := project.GetSecret(logicalName); exists {
			return errs.New(error_codes.SecretExistsErrCode, "secret already exists in project").WithContext("logical_name", logicalName)
		}

		project.AddSecret(logicalName, secretRef)

		err = ws.UpdateProject(project)
		if err != nil {
			return err
		}

		fmt.Printf("Added secret '%s' -> '%s' to project '%s'\n", logicalName, secretRef, projectName)
		return nil
	})
}

func ProjectRemoveSecretHandler(projectName, logicalName string) error {

	return common.WithLocalWorkspace(func(ws *workspace.Workspace) error {
		project, err := ws.LoadProject(projectName)
		if err != nil {
			return err
		}

		// Check if secret exists
		if _, exists := project.GetSecret(logicalName); !exists {
			return errs.New(error_codes.SecretNotFoundErrCode, "secret not found in project").WithContext("logical_name", logicalName)
		}

		project.RemoveSecret(logicalName)

		err = ws.UpdateProject(project)
		if err != nil {
			return err
		}

		fmt.Printf("Removed secret '%s' from project '%s'\n", logicalName, projectName)
		return nil
	})
}
