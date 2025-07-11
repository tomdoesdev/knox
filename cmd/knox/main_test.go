package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/tomdoesdev/knox/internal/cli"
	"github.com/tomdoesdev/knox/internal/errors"
	"testing"
)

func TestCliCommands(t *testing.T) {

	t.Run("should run init", func(t *testing.T) {
		cmd := cli.NewRootCommand()
		a := assert.New(t)
		err := cmd.Run(context.Background(), []string{"knox", "init"})

		a.EqualError(errors.ErrNotImplemented, err.Error())
	})
}
