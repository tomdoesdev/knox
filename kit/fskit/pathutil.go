package fskit

import (
	"errors"
	"fmt"
	"os"
)

func Exists(path string) (bool, error) {
	_, err := os.Lstat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, fmt.Errorf("pathutil: stat %q: %w", path, err)
}
