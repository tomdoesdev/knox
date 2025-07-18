package secrets

import (
	"log/slog"
)

// EncryptionHandler defines the interface for encrypting/decrypting secrets
type EncryptionHandler interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type noOpEncryption struct{}

func (p *noOpEncryption) Encrypt(plaintext string) (string, error) {
	slog.Debug("Using no-op encryption handler")
	return plaintext, nil
}

func (p *noOpEncryption) Decrypt(ciphertext string) (string, error) {
	slog.Debug("Using no-op encryption handler")
	return ciphertext, nil
}

func NewNoOpEncryptionHandler() EncryptionHandler {
	slog.Debug("Initializing knox with no-op encryption handler")
	return &noOpEncryption{}
}
