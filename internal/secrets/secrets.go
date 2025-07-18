package secrets

import (
	"io"

	"github.com/tomdoesdev/knox/pkg/errs"
)

const (
	SecretReadFailureCode    errs.ErrorCode = "SECRET_READ_FAILURE"
	SecretWriteFailureCode   errs.ErrorCode = "SECRET_WRITE_FAILURE"
	SecreteDeleteFailureCode errs.ErrorCode = "SECRET_DELETE_FAILURE"
)

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
	io.Closer
}
