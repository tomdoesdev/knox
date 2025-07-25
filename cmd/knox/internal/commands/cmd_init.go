package commands

import (
	"context"
	"fmt"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/handlers"
	"github.com/tomdoesdev/knox/kit/ast"
	"github.com/urfave/cli/v3"
)

func initVisitor(n ast.Node) error {
	var created bool

	if boolVal, ok := n.GetAttribute("created"); ok {
		if flag, isBool := boolVal.(bool); isBool {
			created = flag
		}
	}

	fmt.Printf("%s - created=%b path=%s\n", n.Content().String(), created, n.Content().String())
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
