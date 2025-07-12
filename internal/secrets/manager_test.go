package secrets

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVaultManager(t *testing.T) {
	t.Run("should initialise a new knox project vault", func(t *testing.T) {
		cwd, _ := os.Getwd()

		initializer := NewVaultInitializer()
		err := initializer.NewVault(&VaultOptions{Path: cwd})
		assert.Nil(t, err)

	})
}
