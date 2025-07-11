package secrets

import (
	"log/slog"
	"os"
	"testing"
)

func TestVaultManager(t *testing.T) {
	t.Run("should initialise a new knox project vault", func(t *testing.T) {
		if cwd, err := os.Getwd(); err == nil {
			slog.Info("dir", "cwd", cwd)
			manager := NewVaultManager(cwd)
			err := manager.Init()
			if err != nil {
				return
			}
		}

	})
}
