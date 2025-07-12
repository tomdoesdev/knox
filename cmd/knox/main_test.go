package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tomdoesdev/knox/cmd/knox/internal/commands"
	"github.com/tomdoesdev/knox/internal/errors"
)

func TestCliCommands(t *testing.T) {
	t.Run("all sub commands exist", func(t *testing.T) {
		a := assert.New(t)

		tests := []string{"init", "status", "set", "remove"}

		for _, cmd := range tests {
			c := commands.NewRootCommand()
			err := c.Run(context.Background(), []string{"knox", cmd})
			a.NotNil(err)
			a.EqualError(errors.ErrNotImplemented, err.Error())
		}

	})
}
