package sqlite

import "strings"

func IsUniqueConstraintError(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "UNIQUE constraint failed") ||
		strings.Contains(errStr, "constraint failed")
}
