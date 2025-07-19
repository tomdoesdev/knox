package vault

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tomdoesdev/knox/internal/v1/constants"
)

func setEnv() error {
	err := os.Setenv(constants.EnvVarKnoxRoot, ".test")
	if err != nil {
		return err
	}

	err = os.Setenv("XDG_DATA_HOME", ".test/knox.vault")
	return err
}
func TestExists(t *testing.T) {
	err := setEnv()

	if err != nil {
		t.Fatal(err)
	}

	p, err := FindLocation()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Path: %s\n", p)
	assert.NotEmpty(t, p)
}
