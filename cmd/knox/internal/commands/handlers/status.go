package handlers

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/common"
	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/common/output"
	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/common/output/formatters/plain"
	"github.com/tomdoesdev/knox/internal/workspace"
)

func StatusHandler() error {
	return common.WithLocalWorkspace(func(ws *workspace.Workspace) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		// Create the output tree using v2 system
		root := output.NewRoot()

		// Workspace Directory section
		workspaceDirHeading := output.NewHeading("Workspace Directory")
		workspaceDirHeading.Add(output.NewText(cwd))
		root.Add(workspaceDirHeading)

		// Projects section
		projects, err := ws.ListProjects()
		if err != nil {
			return err
		}

		projectsHeading := output.NewHeading("Projects")
		if len(projects) == 0 {
			projectsHeading.Add(output.NewText("none"))
		} else {
			projectsList := output.NewList()

			// Get current active project
			currentProject, _ := ws.CurrentProject()

			for _, project := range projects {
				if project == currentProject {
					projectsList.Add(output.NewListItem(project + "   [active]"))
				} else {
					projectsList.Add(output.NewListItem(project))
				}
			}
			projectsHeading.Add(projectsList)
		}
		root.Add(projectsHeading)

		// Linked Vaults section
		vaults, err := ws.GetLinkedVaults()
		if err != nil {
			return err
		}

		vaultsHeading := output.NewHeading("Linked Vaults")
		if len(vaults) == 0 {
			vaultsHeading.Add(output.NewText("None"))
		} else {
			vaultsList := output.NewList()
			for _, vault := range vaults {
				vaultsList.Add(output.NewListItem(vault.Alias + " " + vault.Path))
			}
			vaultsHeading.Add(vaultsList)
		}
		root.Add(vaultsHeading)

		// Create formatter and render
		formatter := output.NewOutputFormatter().
			Using("root", plain.RootFormatter).
			Using("heading", plain.HeadingFormatter).
			Using("list", plain.ListFormatter).
			Using("listitem", plain.ListItemFormatter).
			Using("text", plain.TextWithNewlineFormatter)

		result, err := formatter.Render(root)
		if err != nil {
			return fmt.Errorf("failed to render status output: %w", err)
		}

		fmt.Print(result)
		slog.Debug("knox workspace dir", "dir", ws.DataDir())

		return nil
	})
}
