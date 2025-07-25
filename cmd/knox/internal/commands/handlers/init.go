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

		// Set path and project attributes

		switch result {
		case workspace.Created:
			{
				b.Node("message").Attr("path", ws.Dir()).Attr("project", currentProject).
					Attr("created", true).
					Content(fmt.Sprintf("initialized empty workspace in %s", ws.Dir())).
					Up().
					Node("message").
					Content(fmt.Sprintf("current project: %s\n", currentProject))
			}
			break
		case workspace.Existed:

			b.Node("message").Attr("path", ws.Dir()).Attr("project", currentProject).
				Attr("created", true).
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
