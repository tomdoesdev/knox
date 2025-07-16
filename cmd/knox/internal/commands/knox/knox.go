package knox

import (
	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/knox_add"
	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/knox_init"
	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/knox_remove"
	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/knox_run"
	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/knox_status"
	"github.com/tomdoesdev/knox/pkg/project"
	"github.com/urfave/cli/v3"
)

func NewKnoxCommand(p *project.Project) *cli.Command {
	return &cli.Command{
		Name:        "knox",
		Usage:       "manage local development vault",
		Description: "local development vault manager",

		Commands: []*cli.Command{
			knox_init.NewInitCommand(),
			knox_status.NewStatusCommand(p),
			knox_add.NewAddCommand(p),
			knox_remove.NewRemoveCommand(p),
			knox_run.NewRunCommand(),
		},
	}
}
