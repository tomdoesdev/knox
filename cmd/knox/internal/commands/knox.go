package commands

import (
	"github.com/tomdoesdev/knox/pkg/project"
	"github.com/urfave/cli/v3"
)

func NewKnoxCommand(p *project.Project) *cli.Command {
	return &cli.Command{
		Name:        "knox",
		Usage:       "manage local development vault",
		Description: "local development vault manager",

		Commands: []*cli.Command{
			NewInitCommand(),
			NewStatusCommand(p),
			NewAddCommand(p),
			NewRemoveCommand(p),
			NewRunCommand(),
		},
	}
}
