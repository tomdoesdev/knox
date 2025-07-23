package project

import (
	"regexp"
	"strings"

	"github.com/tomdoesdev/knox/internal"
	"github.com/tomdoesdev/knox/kit/errs"
)

var (
	ErrInvalidProjectName = errs.New(internal.InvalidIdentifierCode, "invalid project name")
)

type Name string

func NewName(s string) (Name, error) {
	name := Name(s)

	if err := isValidProjectName(name); err != nil {
		return "", err
	}

	return name, nil
}
func (n Name) String() string {
	return strings.ToLower(strings.TrimSpace(strings.ReplaceAll(string(n), " ", "_")))
}

func isValidProjectName(name Name) error {

	if name == "" {
		return ErrInvalidProjectName
	}

	/*
		- ^[a-z0-9_-]+ - Must start with lowercase alphanumeric, underscore, or hyphen (no slash)
		- (?:/[a-z0-9_-]+)* - Zero or more groups of: slash followed by lowercase alphanumeric/underscore/hyphen
		- $ - Must end with alphanumeric, underscore, or hyphen (no slash)
	*/
	pattern := `^[a-z0-9_-]+(?:/[a-z0-9_-]+)*$`

	matched, err := regexp.MatchString(pattern, name.String())
	if err != nil {
		return errs.Wrap(err,
			internal.InvalidIdentifierCode,
			"invalid project name").WithContext("value", name).WithContext("pattern", pattern)
	}

	if !matched {
		return ErrInvalidProjectName.WithContext("value", name).WithContext("pattern", pattern)
	}

	return nil
}
