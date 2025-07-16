package knox_remove

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tomdoesdev/knox/internal/store"
	"github.com/tomdoesdev/knox/pkg/project"
	"github.com/tomdoesdev/knox/pkg/secrets"
	"github.com/urfave/cli/v3"
)

func NewRemoveCommand(p *project.Project) *cli.Command {
	return &cli.Command{
		Name:    "remove",
		Aliases: []string{"rm"},
		Usage:   "remove a secret from the vault",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return removeActionfunc(cmd, p)
		},
	}
}

func removeActionfunc(cmd *cli.Command, p *project.Project) error {
	if cmd.Args().Len() != 1 {
		return fmt.Errorf("expected 1 argument, got %d", cmd.Args().Len())
	}

	s := store.NewSqlite(p.VaultPath, p.ProjectID)
	defer (func() {
		err := s.Close()
		if err != nil {
			slog.Error("error closing secret s: %v", err)

		}
	})()

	key := cmd.Args().Get(0)

	err := remove(key, s)
	if err != nil {
		return err
	}

	fmt.Println("removed secret:", key)
	return nil

}
func remove(key string, s secrets.SecretDeleter) error {
	return s.DeleteSecret(key)
}
