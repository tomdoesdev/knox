package handlers

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/common"
	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/renderers/cli"
	"github.com/tomdoesdev/knox/internal/workspace"
	"github.com/tomdoesdev/knox/kit/ast"
)

func StatusHandler() error {
	return common.WithLocalWorkspace(func(ws *workspace.Workspace) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		// Create the AST tree
		root := ast.NewRoot()

		// Workspace Directory section
		workspaceDirHeading := ast.NewHeading("Workspace Directory")
		workspaceDirHeading.AddChild(ast.Text(cwd))
		root.AddChild(workspaceDirHeading)

		// Projects section
		projects, err := ws.ListProjects()
		if err != nil {
			return err
		}

		projectsHeading := ast.NewHeading("Projects")
		if len(projects) == 0 {
			projectsHeading.AddChild(ast.Text("none"))
		} else {
			projectsList := ast.NewList()

			// Get current active project
			currentProject, _ := ws.CurrentProject()

			for _, project := range projects {
				listItem := ast.NewListItem(project)
				if project == currentProject {
					// Use attributes to mark active project
					listItem.SetAttribute("active", true)
					listItem.SetContent(ast.StringValue(project + "   [active]"))
				}
				projectsList.AddChild(listItem)
			}
			projectsHeading.AddChild(projectsList)
		}
		root.AddChild(projectsHeading)

		// Linked Vaults section
		vaults, err := ws.GetLinkedVaults()
		if err != nil {
			return err
		}

		vaultsHeading := ast.NewHeading("Linked Vaults")
		if len(vaults) == 0 {
			vaultsHeading.AddChild(ast.Text("None"))
		} else {
			vaultsList := ast.NewList()
			for _, vault := range vaults {
				vaultItem := ast.NewListItem(vault.Alias + " " + vault.Path)
				// Add metadata attributes for potential future use
				vaultItem.SetAttribute("alias", vault.Alias)
				vaultItem.SetAttribute("path", vault.Path)
				vaultsList.AddChild(vaultItem)
			}
			vaultsHeading.AddChild(vaultsList)
		}
		root.AddChild(vaultsHeading)

		// Create CLI renderer and render
		renderer := cli.NewRenderer()
		result, err := renderer.Render(root)
		if err != nil {
			return fmt.Errorf("failed to render status output: %w", err)
		}

		fmt.Print(result)
		slog.Debug("knox workspace dir", "dir", ws.DataDir())

		return nil
	})
}
