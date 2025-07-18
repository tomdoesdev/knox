package secrets

import (
	"database/sql"
	"log/slog"

	"github.com/tomdoesdev/knox/pkg/errs"
)

const Schema = `
create table if not exists vault_config (
	key text not null primary Key,
	value text not null
);

CREATE TABLE secrets (
id INTEGER PRIMARY KEY,
project_id TEXT NOT NULL,
key TEXT NOT NULL,
value TEXT NOT NULL,

UNIQUE (project_id,key)
)`

type SqliteSecretStore struct {
	db        *sql.DB
	projectId string
	encryptor EncryptionHandler
}

func (s *SqliteSecretStore) ListKeys() ([]string, error) {
	query := `SELECT key FROM secrets WHERE project_id = ?`

	rows, err := s.db.Query(query, s.projectId)
	if err != nil {
		return nil, errs.Wrap(err, SecretListFailureCode, "failed to query database")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			slog.Error("failed to close database rows", "err", err)
		}
	}(rows)

	keys := make([]string, 0)
	for rows.Next() {
		var key string

		if err := rows.Scan(&key); err != nil {
			return nil, errs.Wrap(err, SecretReadFailureCode, "failed to scan Key")
		}

		keys = append(keys, key)
	}
	return keys, nil
}

func (s *SqliteSecretStore) ListSecrets() ([]Secret, error) {
	query := `SELECT key, value FROM secrets WHERE project_id = ?`

	rows, err := s.db.Query(query, s.projectId)
	if err != nil {
		return nil, errs.Wrap(err, SecretListFailureCode, "failed to query database")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			slog.Error("failed to close database rows", "err", err)
		}
	}(rows)

	secrets := make([]Secret, 0)
	for rows.Next() {
		secret := Secret{}

		if err := rows.Scan(&secret.Key, &secret.Value); err != nil {
			return nil, errs.Wrap(err, SecretReadFailureCode, "failed to scan secret")
		}

		secrets = append(secrets, secret)
	}
	return secrets, nil
}

func (s *SqliteSecretStore) ReadSecret(key string) (string, error) {
	query := `SELECT value FROM secrets WHERE project_id = $1 AND key = $2 LIMIT 1;`
	var secret string

	err := s.db.QueryRow(query, s.projectId, key).Scan(&secret)
	if err != nil {
		return "", errs.Wrap(err, SecretReadFailureCode, "can not read secret from database")
	}

	value, err := s.encryptor.Decrypt(secret)
	if err != nil {
		return "", errs.Wrap(err, SecretReadFailureCode, "can not decrypt secret")
	}

	return value, nil
}

func (s *SqliteSecretStore) WriteSecret(key, value string) error {
	query := `INSERT INTO secrets (project_id, key, value) VALUES ($1, $2, $3);`
	value, err := s.encryptor.Encrypt(value)
	if err != nil {
		return errs.Wrap(err, SecretWriteFailureCode, "can not encrypt secret")
	}
	_, err = s.db.Exec(query, s.projectId, key, value)
	if err != nil {
		return errs.Wrap(err, SecretWriteFailureCode, "can not write secret to database")
	}
	return nil
}

func (s *SqliteSecretStore) DeleteSecret(key string) error {
	query := `DELETE FROM secrets WHERE project_id = $1 AND key = $2;`

	_, err := s.db.Exec(query, s.projectId, key)
	if err != nil {
		return errs.Wrap(err, SecreteDeleteFailureCode, "can not delete secret from database")
	}

	return nil
}

func (s *SqliteSecretStore) Close() error {
	return s.db.Close()
}

func newSqliteStore(dataSourceName, projectId string, handler EncryptionHandler) (*SqliteSecretStore, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, errs.Wrap(err, StoreOpenFailureCode, "failed to open database")
	}

	// Configure SQLite security settings
	if err := configureSQLiteSecurity(db); err != nil {
		return nil, errs.Wrap(err, StoreExecFailureCode, "failed to configure SQLite security")
	}

	// Apply schema and return early if we are using an in memory sqlite instance
	if dataSourceName == ":memory:" {
		if _, err := db.Exec(Schema); err != nil {
			return nil, errs.Wrap(err, StoreExecFailureCode, "failed to initialize database")
		}
		return &SqliteSecretStore{db: db, projectId: projectId, encryptor: handler}, nil
	}

	if isInitialized(db) {
		return &SqliteSecretStore{db: db, projectId: projectId, encryptor: handler}, nil
	}

	if _, err := db.Exec(Schema); err != nil {
		return nil, errs.Wrap(err, StoreExecFailureCode, "failed to initialize database")
	}

	return &SqliteSecretStore{db: db, projectId: projectId, encryptor: handler}, nil

}

func configureSQLiteSecurity(db *sql.DB) error {
	// Enable secure delete - overwrites deleted data with zeros
	// Options: OFF (default), ON (secure), FAST (compromise)
	_, err := db.Exec("PRAGMA secure_delete = FAST")
	if err != nil {
		return err
	}

	// Enable foreign key constraints for data integrity
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return err
	}

	// Set journal mode to WAL for better concurrency and crash recovery
	_, err = db.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		return err
	}

	// Enable synchronous mode for better durability
	_, err = db.Exec("PRAGMA synchronous = NORMAL")
	if err != nil {
		return err
	}

	return nil
}

func isInitialized(db *sql.DB) bool {
	var result string

	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='secrets'").Scan(&result)
	if err != nil {
		return false
	}

	return result == "secrets"
}
