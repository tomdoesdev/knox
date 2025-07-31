package fs

import (
	"io"
	"os"
	"time"

	"github.com/tomdoesdev/knox/kit/errs"
)

// Copy copies a file from src to dst, preserving permissions.
func Copy(src, dst string) error {
	if src == "" {
		return errs.New(ECodeInvalidPath, "cannot copy from empty source path")
	}
	if dst == "" {
		return errs.New(ECodeInvalidPath, "cannot copy to empty destination path")
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return errs.Wrap(err, ECodeFileReadFailure, "open source file failed").
			WithPath(src).
			WithOperation("open")
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return errs.Wrap(err, ECodeFileReadFailure, "stat source file failed").
			WithPath(src).
			WithOperation("stat")
	}

	if !srcInfo.Mode().IsRegular() {
		return errs.New(ECodeInvalidPath, "source is not a regular file").
			WithPath(src)
	}

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return errs.Wrap(err, ECodeFileWriteFailure, "create destination file failed").
			WithPath(dst).
			WithOperation("create")
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return errs.Wrap(err, ECodeFileMoveFailure, "copy data failed").
			WithContext("src", src).
			WithContext("dst", dst).
			WithOperation("copy")
	}

	return nil
}

// Move moves/renames a file from src to dst.
func Move(src, dst string) error {
	if src == "" {
		return errs.New(ECodeInvalidPath, "cannot move from empty source path")
	}
	if dst == "" {
		return errs.New(ECodeInvalidPath, "cannot move to empty destination path")
	}

	err := os.Rename(src, dst)
	if err != nil {
		return errs.Wrap(err, ECodeFileMoveFailure, "move file failed").
			WithContext("src", src).
			WithContext("dst", dst).
			WithOperation("move")
	}
	return nil
}

// Touch creates an empty file or updates the timestamp of an existing file.
func Touch(path string, perm os.FileMode) error {
	if path == "" {
		return errs.New(ECodeInvalidPath, "cannot touch file with empty path")
	}

	now := time.Now()
	err := os.Chtimes(path, now, now)
	if err != nil {
		// File doesn't exist, create it
		file, createErr := os.OpenFile(path, os.O_CREATE, perm)
		if createErr != nil {
			return errs.Wrap(createErr, ECodeFileWriteFailure, "touch file failed").
				WithPath(path).
				WithOperation("touch")
		}
		file.Close()
	}
	return nil
}

// ReadFile reads the contents of a file.
func ReadFile(path string) ([]byte, error) {
	if path == "" {
		return nil, errs.New(ECodeInvalidPath, "cannot read file with empty path")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errs.Wrap(err, ECodeFileReadFailure, "read file failed").
			WithPath(path).
			WithOperation("read")
	}
	return data, nil
}

/*
WriteFile writes data to a file with the given permissions.

If the file exists it results in an ECodeEntityExists error.

Returns:
  - ECodeEntityExists
  - ECodeInvalidPath
  - ECodeFileWriteFailure
*/
func WriteFile(path string, data []byte, perm os.FileMode) error {
	if path == "" {
		return errs.New(ECodeInvalidPath, "cannot write file with empty path")
	}

	if IsFile(path) {
		return errs.New(ECodeEntityExists, "cannot write to existing file").
			WithPath(path).WithContext("hint", "use 'OverwriteFile' function")
	}

	err := os.WriteFile(path, data, perm)
	if err != nil {
		return errs.Wrap(err, ECodeFileWriteFailure, "write file failed").
			WithPath(path).
			WithOperation("write")
	}
	return nil
}

func OverwriteFile(path string, data []byte, perm os.FileMode) error {
	if path == "" {
		return errs.New(ECodeInvalidPath, "cannot write file with empty path")
	}

	err := os.WriteFile(path, data, perm)
	if err != nil {
		return errs.Wrap(err, ECodeFileWriteFailure, "write file failed").
			WithPath(path).
			WithOperation("write")
	}
	return nil
}
