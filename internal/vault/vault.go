package vault

import (
	"github.com/tomdoesdev/knox/kit/errs"
)

const (
	ECodeFailedToPersist errs.Code = "FAILED_TO_PERSIST"
)

type (
	Vault interface {
		Save() error
		SecretGetter
		SecretSetter
		SecretDeleter
	}
	SecretGetter interface {
		Get(string) string
	}
	SecretSetter interface {
		Set(string, string) error
	}

	SecretDeleter interface {
		Delete(string) error
	}
)
