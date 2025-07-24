package commands

import (
	"context"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/common"
	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/handlers"
	"github.com/urfave/cli/v3"
)

func NewNewCommand() *cli.Command {
	return &cli.Command{
		Name:  "new",
		Usage: "create a new project or vault",
		Commands: []*cli.Command{
			newNewProjectCommand(),
			newVaultCommand(),
			// TODO: newNewSecretCommand(), etc.
		},
	}
}

func newNewProjectCommand() *cli.Command {
	return &cli.Command{
		Name:      "project",
		Usage:     "create a new project",
		ArgsUsage: "<project-name>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "description",
				Usage: "project description",
				Value: "",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if err := common.ExpectExactArgCount(1, "project name is required", cmd.Args()); err != nil {
				return err
			}

			name := cmd.Args().First()
			description := cmd.String("description")

			return handlers.NewProjectHandler(name, description)
		},
	}
}

func newVaultCommand() *cli.Command {
	return &cli.Command{
		Name:      "vault",
		Usage:     "create a new vault",
		ArgsUsage: "<vault-alias>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Usage: "path where vault will be created (defaults to $KNOX_ROOT/vaults or ~/.knox/vaults)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if err := common.ExpectExactArgCount(1, "vault alias is required", cmd.Args()); err != nil {
				return err
			}

			alias := cmd.Args().First()
			path := cmd.String("path")

			return handlers.NewVaultHandler(alias, path)
		},
	}
}
