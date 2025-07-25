package handlers

import (
	"fmt"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/common"
	"github.com/tomdoesdev/knox/internal/workspace"
	"github.com/tomdoesdev/knox/kit/ast"
)

func InitHandler() (ast.Node, error) {
	err := common.WithEnsuredLocalWorkspace(func(ws *workspace.Workspace, result workspace.InitResult) error {
		currentProject, err := ws.CurrentProject()
		if err != nil {
			currentProject = "none"
		}

		switch result {
		case workspace.Created:
			fmt.Printf("initialized empty workspace in %s\n", ws.Dir())
			fmt.Printf("current project: %s\n", currentProject)
			break
		case workspace.Existed:
			fmt.Printf("workspace already exists in %s\n", ws.Dir())
			fmt.Printf("current project: %s\n", currentProject)
			break
		default:
			panic(fmt.Sprintf("unexpected result: %s", result))
		}

		return nil
	})

	return nil, err
}
