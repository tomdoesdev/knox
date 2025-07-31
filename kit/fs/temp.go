package fs

import (
	"os"

	"github.com/tomdoesdev/knox/kit/errs"
)

// TempDir creates a temporary directory and returns its path along with a cleanup function.
func TempDir(pattern string) (string, func(), error) {
	dir, err := os.MkdirTemp("", pattern)
	if err != nil {
		return "", nil, errs.Wrap(err, ECodeTempFailure, "create temp directory failed").
			WithOperation("mkdtemp")
	}

	cleanup := func() {
		os.RemoveAll(dir)
	}

	return dir, cleanup, nil
}

// TempFile creates a temporary file and returns the file along with a cleanup function.
func TempFile(pattern string) (*os.File, func(), error) {
	file, err := os.CreateTemp("", pattern)
	if err != nil {
		return nil, nil, errs.Wrap(err, ECodeTempFailure, "create temp file failed").
			WithOperation("mktemp")
	}

	cleanup := func() {
		file.Close()
		os.Remove(file.Name())
	}

	return file, cleanup, nil
}
