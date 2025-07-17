package secrets

import (
	"io"
	"log/slog"
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

// EncryptionHandler defines the interface for encrypting/decrypting secrets
type EncryptionHandler interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type noOpEncryption struct{}

func (p *noOpEncryption) Encrypt(plaintext string) (string, error) {
	return plaintext, nil
}

func (p *noOpEncryption) Decrypt(ciphertext string) (string, error) {
	return ciphertext, nil
}

func NewNoOpEncryptionHandler() EncryptionHandler {
	slog.Warn("Using noop encryption handler")
	return &noOpEncryption{}
}
