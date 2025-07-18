package secrets

import (
	"github.com/tomdoesdev/knox/kit/fs"
	"github.com/tomdoesdev/knox/pkg/errs"
)

const (
	StoreOpenFailureCode errs.ErrorCode = "STORE_FAILURE"
	StoreExecFailureCode errs.ErrorCode = "STORE_EXEC_FAILURE"
)

func NewFileSecretStore(source fs.FilePath, projectId string, handler EncryptionHandler) (*SqliteSecretStore, error) {
	s, err := newSqliteStore(source, projectId, handler)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func NewMemorySecretStore(projectId string, handler EncryptionHandler) (*SqliteSecretStore, error) {
	s, err := newSqliteStore(":memory:", projectId, handler)
	if err != nil {
		return nil, err
	}

	return s, nil
}
