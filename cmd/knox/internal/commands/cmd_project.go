package commands

import (
	"context"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/common"
	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/handlers"
	"github.com/urfave/cli/v3"
)

func NewProjectCommand() *cli.Command {
	return &cli.Command{
		Name:  "project",
		Usage: "manage workspace projects",
		Commands: []*cli.Command{
			newProjectListCommand(),
			newProjectDeleteCommand(),
			newProjectListSecretsCommand(),
			newProjectAddSecretCommand(),
			newProjectRemoveSecretCommand(),
		},
	}
}

func newProjectListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "list all projects",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return handlers.ProjectListHandler()
		},
	}
}

func newProjectDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "delete a project",
		ArgsUsage: "<project-name>",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if err := common.ExpectExactArgCount(1, "project name is required", cmd.Args()); err != nil {
				return err
			}

			name := cmd.Args().First()
			return handlers.ProjectDeleteHandler(name)
		},
	}
}

func newProjectListSecretsCommand() *cli.Command {
	return &cli.Command{
		Name:      "list-secrets",
		Usage:     "list secrets in a project",
		ArgsUsage: "[project-name]",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			var name string
			if cmd.Args().Len() == 1 {
				name = cmd.Args().First()
			}
			return handlers.ProjectListSecretsHandler(name)
		},
	}
}

func newProjectAddSecretCommand() *cli.Command {
	return &cli.Command{
		Name:      "add-secret",
		Usage:     "add a secret to a project",
		ArgsUsage: "<project-name> <logical-name> <vault-path@vault-alias>",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if err := common.ExpectExactArgCount(3,
				"project name, logical name, and secret reference are required", cmd.Args()); err != nil {
				return err
			}

			projectName := cmd.Args().Get(0)
			logicalName := cmd.Args().Get(1)
			secretRef := cmd.Args().Get(2)

			return handlers.ProjectAddSecretHandler(projectName, logicalName, secretRef)
		},
	}
}

func newProjectRemoveSecretCommand() *cli.Command {
	return &cli.Command{
		Name:      "remove-secret",
		Usage:     "remove a secret from a project",
		ArgsUsage: "<project-name> <logical-name>",
		Action: func(ctx context.Context, cmd *cli.Command) error {

			if err := common.ExpectExactArgCount(2,
				"project name and logical name are required", cmd.Args()); err != nil {
				return err
			}

			projectName := cmd.Args().Get(0)
			logicalName := cmd.Args().Get(1)

			return handlers.ProjectRemoveSecretHandler(projectName, logicalName)
		},
	}
}
