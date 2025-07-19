package secrets

import (
	"fmt"
	"io"

	"github.com/tomdoesdev/knox/pkg/errs"
)

const (
	SecretReadFailureCode    errs.ErrorCode = "SECRET_READ_FAILURE"
	SecretWriteFailureCode   errs.ErrorCode = "SECRET_WRITE_FAILURE"
	SecreteDeleteFailureCode errs.ErrorCode = "SECRET_DELETE_FAILURE"

	SecretListFailureCode errs.ErrorCode = "SECRET_LIST_FAILURE"
)

type Secret struct {
	Key   string
	Value string
}

func (s *Secret) String() string {
	return fmt.Sprintf("%s=%s", s.Key, s.Value)
}

type SecretLister interface {
	ListKeys() ([]string, error)
	ListSecrets() ([]Secret, error)
}
type SecretReader interface {
	ReadSecret(key string) (string, error)
}
type SecretWriter interface {
	WriteSecret(key, value string) error
}

type SecretDeleter interface {
	DeleteSecret(key string) error
}

type SecretReadWriter interface {
	SecretReader
	SecretWriter
}

type SecretStore interface {
	SecretReadWriter
	SecretDeleter
	SecretLister
	io.Closer
}
