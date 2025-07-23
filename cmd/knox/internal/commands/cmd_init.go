package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	internal2 "github.com/tomdoesdev/knox/cmd/knox/internal/commands/internal"
	"github.com/tomdoesdev/knox/internal"
	"github.com/tomdoesdev/knox/internal/project"
	"github.com/tomdoesdev/knox/internal/workspace"
	"github.com/tomdoesdev/knox/kit/errs"
	"github.com/urfave/cli/v3"
)

type InitCmdArgs struct {
	Name string
}

func NewInitCommand() *cli.Command {
	c := &cli.Command{
		Name:  "init",
		Usage: "initialize a new workspace or project",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// initialise a workspace if no sub command is specified
			return workspaceHandler()
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

func initProjectArgValidator(cmd *cli.Command, args *InitCmdArgs) error {
	name := strings.Join(cmd.Args().Slice(), " ")
	if name == "" {
		return internal2.ErrInvalidInitArgs
	}

	args.Name = name
	return nil
}

func workspaceHandler() error {
	cwd, err := os.Getwd()
	if err != nil {
		return errs.Wrap(err, internal.CreateFailureCode, "failed to get current working directory")
	}

	result, err := workspace.EnsureWorkspace(cwd)
	if err != nil {
		return errs.Wrap(err, internal.CreateFailureCode, "failed to ensure workspace")
	}

	// Get workspace to show current project
	ws, err := workspace.FindWorkspace(cwd)
	if err != nil {
		return errs.Wrap(err, internal.SearchFailureCode, "failed to find workspace")
	}

	currentProject, err := ws.GetCurrentProject()
	if err != nil {
		currentProject = "none"
	}

	switch result {
	case workspace.Created:
		fmt.Printf("initialized empty workspace in %s\n", cwd)
		fmt.Printf("current project: %s\n", currentProject)
		break
	case workspace.Existed:
		fmt.Printf("workspace already exists in %s\n", cwd)
		fmt.Printf("current project: %s\n", currentProject)
		break
	default:
		panic(fmt.Sprintf("unexpected result: %s", result))
	}

	return nil
}

func projectHandler(args *InitCmdArgs) error {
	name, err := project.NewName(args.Name)
	if err != nil {
		return err
	}
	_, _ = fmt.Printf("init project called with name=%s\n", name)
	return nil
}
