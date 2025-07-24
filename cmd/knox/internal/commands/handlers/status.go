package handlers

import (
	"fmt"
	"log/slog"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/common"
	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/common/output"
	"github.com/tomdoesdev/knox/internal/workspace"
)

// StatusItem represents a single item in the status output
type StatusItem interface {
	Name() string
	Value(ws *workspace.Workspace, workspaceDir string) (string, error)
}

// StatusRegistry manages all status items
type StatusRegistry struct {
	items []StatusItem
}

// NewStatusRegistry creates a new status registry with default items
func NewStatusRegistry() *StatusRegistry {
	registry := &StatusRegistry{}

	// Add default status items
	registry.AddItem(&WorkspaceDirectoryItem{})
	registry.AddItem(&ProjectsItem{})
	registry.AddItem(&LinkedVaultsItem{})

	return registry
}

// AddItem adds a new status item to the registry
func (sr *StatusRegistry) AddItem(item StatusItem) {
	sr.items = append(sr.items, item)
}

// PrintStatus prints all status items
func (sr *StatusRegistry) PrintStatus(ws *workspace.Workspace, workspaceDir string) error {

	for _, item := range sr.items {
		value, err := item.Value(ws, workspaceDir)
		if err != nil {
			value = fmt.Sprintf("error: %s", err.Error())
		}

		fmt.Printf("%-20s %s\n\n", item.Name()+":", value)
	}

	return nil
}

// WorkspaceDirectoryItem shows the current workspace directory
type WorkspaceDirectoryItem struct{}

func (w *WorkspaceDirectoryItem) Name() string {
	return "Workspace Directory"
}

func (w *WorkspaceDirectoryItem) Value(ws *workspace.Workspace, workspaceDir string) (string, error) {
	return workspaceDir, nil
}

// CurrentProjectItem shows the current active project
type CurrentProjectItem struct{}

func (c *CurrentProjectItem) Name() string {
	return "Current Project"
}

func (c *CurrentProjectItem) Value(ws *workspace.Workspace, workspaceDir string) (string, error) {
	currentProject, err := ws.CurrentProject()
	if err != nil {
		return "none", nil
	}
	return currentProject, nil
}

// ProjectsItem shows detailed information about projects in the workspace
type ProjectsItem struct{}

func (p *ProjectsItem) Name() string {
	return "Projects"
}

func (p *ProjectsItem) Value(ws *workspace.Workspace, workspaceDir string) (string, error) {
	projects, err := ws.ListProjects()
	if err != nil {
		return "unknown", err
	}

	if len(projects) == 0 {
		return "none", nil
	}

	// Get current active project
	currentProject, err := ws.CurrentProject()
	if err != nil {
		currentProject = ""
	}

	// Build multi-line output
	result := ""
	for _, project := range projects {
		if project == currentProject {
			result += fmt.Sprintf("\n  - %s   [active]", project)
		} else {
			result += fmt.Sprintf("\n  - %s", project)
		}
	}

	return result, nil
}

// LinkedVaultsItem shows all linked vaults and their paths
type LinkedVaultsItem struct{}

func (l *LinkedVaultsItem) Name() string {
	return "Linked Vaults"
}

func (l *LinkedVaultsItem) Value(ws *workspace.Workspace, workspaceDir string) (string, error) {
	vaults, err := ws.GetLinkedVaults()
	if err != nil {
		return "unknown", err
	}

	if len(vaults) == 0 {
		return "\n    None", nil // Aligned with bullet points
	}

	result := ""
	for _, vault := range vaults {
		result += fmt.Sprintf("\n  - %s %s", vault.Alias, vault.Path)
	}
	return result, nil
}

// StatusHandler handles the status command
func StatusHandler() error {
	return common.WithLocalWorkspace(func(ws *workspace.Workspace) error {
		root := output.NewRoot()
		p, err := ws.ListProjects()

		if err != nil {
			return err
		}

		l := output.NewList()

		for _, project := range p {
			l.Add(output.Text(project))
		}

		root.Add(l)

		fmt.Print(root.String())
		slog.Debug("knox workspace dir", "dir", ws.DataDir())

		return nil
	})

}
