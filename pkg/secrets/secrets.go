package secrets

import (
	"errors"
	"io"
)

var (
	ErrSecretExists = errors.New("secret already exists")
)

type SecretReader interface {
	ReadSecret(key string) (value string, err error)
}

type SecretWriter interface {
	WriteSecret(key, value string) error
}

type SecretDeleter interface {
	DeleteSecret(key string) error
}

type SecretReaderWriter interface {
	SecretReader
	SecretWriter
}

type SecretStore interface {
	SecretReaderWriter
	SecretDeleter
	io.Closer
}
