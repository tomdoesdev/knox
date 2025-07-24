package commands

import (
	"context"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/common"
	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/handlers"
	"github.com/urfave/cli/v3"
)

func NewLinkCommand() *cli.Command {
	return &cli.Command{
		Name:      "link",
		Usage:     "link a vault to the current workspace",
		ArgsUsage: "<vault-path>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "alias",
				Usage:    "alias for the vault in the workspace",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if err := common.ExpectExactArgCount(1, "vault path is required", cmd.Args()); err != nil {
				return err
			}
			vaultPath := cmd.Args().First()
			alias := cmd.String("alias")

			return handlers.LinkVaultHandler(vaultPath, alias)
		},
	}
}
