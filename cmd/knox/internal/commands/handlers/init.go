package handlers

import (
	"fmt"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/common"
	"github.com/tomdoesdev/knox/internal/workspace"
	"github.com/tomdoesdev/knox/kit/ast"
)

func InitHandler() (ast.Node, error) {
	b := ast.NewBuilder("result")

	err := common.WithEnsuredLocalWorkspace(func(ws *workspace.Workspace, result workspace.InitResult) error {
		currentProject, err := ws.CurrentProject()
		if err != nil {
			currentProject = "none"
		}

		switch result {
		case workspace.Created:
			{
				b.Attr("created", true).
					Attr("project", currentProject).
					Node("message").
					Content(fmt.Sprintf("initialized empty workspace in %s\n", ws.Dir())).
					Up().
					Node("message").
					Content(fmt.Sprintf("current project: %s\n", currentProject))
			}
			break
		case workspace.Existed:
			b.Attr("created", true).
				Attr("project", currentProject).
				Node("message").
				Content(fmt.Sprintf("workspace already exists in %s\n", ws.Dir())).
				Up().
				Node("message").
				Content(fmt.Sprintf("current project: %s\n", currentProject))
			break
		default:
			panic(fmt.Sprintf("unexpected result: %s", result))
		}

		return nil
	})

	return b.Build(), err
}
