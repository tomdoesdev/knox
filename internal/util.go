package internal

import (
	"os"
	"path"
)

func GetKnoxDataDir() string {
	dir := "knox"
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return path.Join(home, dir)
}
