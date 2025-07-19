package project

import (
	"regexp"

	v2 "github.com/tomdoesdev/knox/internal/v2"
	"github.com/tomdoesdev/knox/pkg/errs"
)

var (
	ErrInvalidProjectName = errs.New(v2.InvalidIdentifierCode, "invalid project name")
)

func IsValidProjectName(name string) error {
	if name == "" {
		return ErrInvalidProjectName
	}

	//  - ^[a-zA-Z0-9_-]+ - Must start with alphanumeric, underscore, or hyphen (no slash)
	// - (?:/[a-zA-Z0-9_-]+)* - Zero or more groups of: slash followed by alphanumeric/underscore/hyphen
	//  - $ - Must end with alphanumeric, underscore, or hyphen (no slash)
	pattern := `^[a-zA-Z0-9_-]+(?:/[a-zA-Z0-9_-]+)*$`

	matched, err := regexp.MatchString(pattern, name)
	if err != nil {
		return errs.Wrap(err,
			v2.InvalidIdentifierCode,
			"invalid project name").WithContext("value", name).WithContext("pattern", pattern)
	}

	if !matched {
		return ErrInvalidProjectName.WithContext("value", name).WithContext("pattern", pattern)
	}

	return nil
}
