package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/tomdoesdev/knox/internal/project"
	"github.com/tomdoesdev/knox/internal/workspace"
	"github.com/tomdoesdev/knox/internal/workspace/errors"
	"github.com/tomdoesdev/knox/pkg/errs"
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
	cwd, err := os.Getwd()
	if err != nil {
		return errs.Wrap(err, errors.CreateFailureCode, "failed to get current working directory")
	}

	result, err := workspace.EnsureWorkspace(cwd)
	if err != nil {
		return errs.Wrap(err, errors.CreateFailureCode, "failed to ensure workspace")
	}

	switch result {
	case workspace.Created:
		fmt.Printf("initialized empty workspace in %s\n", cwd)
		fmt.Println("workspace is in a detached state.")
		break
	case workspace.Existed:
		//TODO: When we implement status we need to show the current state of the workspace
		fmt.Printf("workspace already exists in %s\n", cwd)
		break
	default:
		panic(fmt.Sprintf("unexpected result: %s", result))
	}

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
