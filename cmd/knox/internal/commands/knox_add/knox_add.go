package knox_add

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/tomdoesdev/knox/internal/store"
	"github.com/tomdoesdev/knox/pkg/project"
	"github.com/tomdoesdev/knox/pkg/secrets"
	"github.com/urfave/cli/v3"
)

func NewAddCommand(p *project.Project) *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: fmt.Sprintf("add a secret to the current knox vault (%s)", p.VaultPath),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return addActionfunc(cmd, p)
		},
	}
}

func addActionfunc(cmd *cli.Command, p *project.Project) error {
	if cmd.Args().Len() != 2 {
		return fmt.Errorf("expected 2 arguments, got %d", cmd.Args().Len())
	}

	key := cmd.Args().Get(0)
	value := cmd.Args().Get(1)
	s := store.NewSqlite(p.VaultPath, p.ProjectID)

	defer func() {
		err := s.Close()
		if err != nil {
			slog.Error("error closing secret store", "error", err)
		}
	}()

	err := add(key, value, s)
	if err != nil {
		return err
	}

	fmt.Printf("secret added: %s\n", key)
	return nil
}
func add(key, value string, s secrets.SecretReaderWriter) error {
	err := s.WriteSecret(key, value)
	if errors.Is(err, secrets.ErrSecretExists) {
		fmt.Printf("secret already exists: %s\n", key)
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}
