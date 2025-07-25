package commands

import (
	"context"
	"fmt"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/handlers"
	"github.com/tomdoesdev/knox/kit/ast"
	"github.com/urfave/cli/v3"
)

func initVisitor(n ast.Node) error {
	if n.Type() == "result" {
		//Ignore the root as it has no content
		return nil
	}
	var path string
	var created bool

	if attr, ok := n.GetAttribute("created"); ok {
		v, err := attr.AsBool()
		if err != nil {
			return err
		}
		created = v
	}

	if attr, ok := n.GetAttribute("path"); ok {
		v, err := attr.AsString()
		if err != nil {
			return err
		}
		path = v
	}

	fmt.Printf("%s - created=%t path=%s\n", n.Content().String(), created, path)
	return nil
}

func NewInitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "initialize a new workspace",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			n, err := handlers.InitHandler()
			if err != nil {
				return err
			}

			return ast.WithPreOrderTraversal(n, initVisitor)
		},
	}
}
