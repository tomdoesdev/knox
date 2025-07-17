package secrets

import (
	"database/sql"

	"github.com/tomdoesdev/knox/pkg/errs"
)

const Schema = `
CREATE TABLE vault (
id INTEGER PRIMARY KEY,
project_id TEXT NOT NULL,
key TEXT NOT NULL,
value TEXT NOT NULL,

UNIQUE (project_id,key)
)`

const (
	VaultRead       errs.ErrorCode = "VAULT_READ"
	VaultWrite      errs.ErrorCode = "VAULT_WRITE"
	VaultDelete     errs.ErrorCode = "VAULT_DELETE"
	VaultOpen       errs.ErrorCode = "VAULT_OPEN"
	VaultInitSchema errs.ErrorCode = "VAULT_INIT_SCHEMA"
)

type SqliteSecretStore struct {
	db        *sql.DB
	projectId string
	encryptor EncryptionHandler
}

func (s *SqliteSecretStore) ReadSecret(key string) (string, error) {
	query := `SELECT value FROM vault WHERE project_id = $1 AND key = $2 LIMIT 1;`
	var secret string

	err := s.db.QueryRow(query, s.projectId, key).Scan(&secret)
	if err != nil {
		return "", errs.Wrap(err, VaultRead, "can not read secret from database")
	}

	value, err := s.encryptor.Decrypt(secret)
	if err != nil {
		return "", errs.Wrap(err, VaultRead, "can not decrypt secret")
	}

	return value, nil
}

func (s *SqliteSecretStore) WriteSecret(key, value string) error {
	query := `INSERT INTO vault (project_id, key, value) VALUES ($1, $2, $3);`
	value, err := s.encryptor.Encrypt(value)
	if err != nil {
		return errs.Wrap(err, VaultWrite, "can not encrypt secret")
	}
	_, err = s.db.Exec(query, s.projectId, key, value)
	if err != nil {
		return errs.Wrap(err, VaultWrite, "can not write secret to database")
	}
	return nil
}

func (s *SqliteSecretStore) DeleteSecret(key string) error {
	query := `DELETE FROM vault WHERE project_id = $1 AND key = $2;`

	_, err := s.db.Exec(query, s.projectId, key)
	if err != nil {
		return errs.Wrap(err, VaultDelete, "can not delete secret from database")
	}

	return nil
}

func (s *SqliteSecretStore) Close() error {
	return s.db.Close()
}

func newSqliteStore(dataSourceName, projectId string, handler EncryptionHandler) (*SqliteSecretStore, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, errs.Wrap(err, VaultOpen, "failed to open database")
	}

	// Apply schema and return early if we are using an in memory sqlite instance
	if dataSourceName == ":memory:" {
		if _, err := db.Exec(Schema); err != nil {
			return nil, errs.Wrap(err, VaultInitSchema, "failed to initialize database")
		}
		return &SqliteSecretStore{db: db, projectId: projectId, encryptor: handler}, nil
	}

	if isInitialized(db) {
		return &SqliteSecretStore{db: db, projectId: projectId, encryptor: handler}, nil
	}

	if _, err := db.Exec(Schema); err != nil {
		return nil, errs.Wrap(err, VaultInitSchema, "failed to initialize database")
	}

	return &SqliteSecretStore{db: db, projectId: projectId, encryptor: handler}, nil

}

func isInitialized(db *sql.DB) bool {
	var result string

	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='vault'").Scan(&result)
	if err != nil {
		return false
	}

	return result == "vault"
}
