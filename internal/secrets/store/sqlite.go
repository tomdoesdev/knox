package store

import (
	"database/sql"
	"fmt"

	"github.com/tomdoesdev/knox/internal/secrets"
)

const Schema = `
CREATE TABLE vault (
id INTEGER PRIMARY KEY,
project_id TEXT NOT NULL,
key TEXT NOT NULL,
value TEXT NOT NULL,

UNIQUE (project_id,key)
)`

type SqliteSecretStore struct {
	db        *sql.DB
	projectId string
	encryptor secrets.EncryptionHandler
}

func (s *SqliteSecretStore) ReadSecret(key string) (*secrets.SecureString, error) {
	query := `SELECT value FROM vault WHERE project_id = $1 AND key = $2 LIMIT 1;`
	var secret string

	err := s.db.QueryRow(query, s.projectId, key).Scan(&secret)
	if err != nil {
		return nil, err
	}

	return s.encryptor.Decrypt(secret)
}

func (s *SqliteSecretStore) WriteSecret(key, value string) error {
	query := `INSERT INTO vault (project_id, key, value) VALUES ($1, $2, $3);`
	value, err := s.encryptor.Encrypt(value)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(query, s.projectId, key, value)
	return err
}

func (s *SqliteSecretStore) DeleteSecret(key string) error {
	query := `DELETE FROM vault WHERE project_id = $1 AND key = $2;`
	_, err := s.db.Exec(query, s.projectId, key)
	return err
}

func (s *SqliteSecretStore) Close() error {
	return s.db.Close()
}

func newSqliteStore(dataSourceName, projectId string, handler secrets.EncryptionHandler) (*SqliteSecretStore, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open in-memory vault: %w", err)
	}

	// Apply schema and return early if we are using an in memory sqlite instance
	if dataSourceName == ":memory:" {
		if _, err := db.Exec(Schema); err != nil {
			return nil, fmt.Errorf("failed to initialize sqlite schema: %w", err)
		}
		return &SqliteSecretStore{db: db, projectId: projectId, encryptor: handler}, nil
	}

	if isInitialized(db) {
		return &SqliteSecretStore{db: db, projectId: projectId, encryptor: handler}, nil
	}

	//Apply schema to unintialized db
	if _, err := db.Exec(Schema); err != nil {
		return nil, fmt.Errorf("failed to initialize sqlite schema: %w", err)
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
