package store

import (
	"github.com/tomdoesdev/knox/internal/secrets"
	"github.com/tomdoesdev/knox/kit/fs"
)

func NewFileSecretStore(source fs.FilePath, projectId string, handler secrets.EncryptionHandler) (*SqliteSecretStore, error) {
	s, err := newSqliteStore(source, projectId, handler)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func NewMemorySecretStore(projectId string, handler secrets.EncryptionHandler) (*SqliteSecretStore, error) {
	s, err := newSqliteStore(":memory:", projectId, handler)
	if err != nil {
		return nil, err
	}

	return s, nil
}
