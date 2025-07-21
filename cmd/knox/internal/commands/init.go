package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/tomdoesdev/knox/internal/project"
	"github.com/urfave/cli/v3"
)

func NewInitCommand() *cli.Command {
	c := &cli.Command{
		Name:  "init",
		Usage: "initialize a new workspace or project",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// initialise a workspace if no sub command is specified
			return workspaceHandler()
		},
		Commands: []*cli.Command{
			newInitProjectCommand(),
			newInitWorkspaceCommand(),
		},
	}

	return c
}

func newInitWorkspaceCommand() *cli.Command {
	c := &cli.Command{
		Name:  "workspace",
		Usage: "initialize a new workspace",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return workspaceHandler()
		},
	}

	return c
}

func newInitProjectCommand() *cli.Command {
	c := &cli.Command{
		Name:  "project",
		Usage: "initialize a new project",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			name := strings.Join(cmd.Args().Slice(), " ")
			if name == "" {
				return fmt.Errorf("project name is required")
			}
			return projectHandler(name)
		},
	}

	return c
}

func workspaceHandler() error {
	fmt.Println("init workspace called")
	return nil
}

func projectHandler(s string) error {
	name, err := project.NewName(s)
	if err != nil {
		return err
	}
	_, _ = fmt.Printf("init project called with name=%s\n", name)
	return nil
}
