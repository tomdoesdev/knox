package store

import (
	"database/sql"
	"fmt"

	"github.com/tomdoesdev/knox/kit/sqlite"
	"github.com/tomdoesdev/knox/pkg/secrets"
)

const (
	setSecretStmt = `
	INSERT INTO vault (project_id, key, value) VALUES ($1, $2, $3);
`

	getSecretStmt = `
	SELECT value FROM vault WHERE project_id = $1 AND key = $2 LIMIT 1;
`
	delSecretStmt = `
		DELETE FROM vault WHERE project_id = $1 AND key = $2;`
)

type sqliteSecretStore struct {
	db     *sql.DB
	dbPath string
	projId string
}

func (s *sqliteSecretStore) DeleteSecret(key string) error {
	db, err := s.ensureDb()
	if err != nil {
		return err
	}

	_, err = db.Exec(delSecretStmt, s.projId, key)
	return err
}

func (s *sqliteSecretStore) ReadSecret(key string) (value string, err error) {
	db, err := s.ensureDb()
	if err != nil {
		return "", err
	}

	var secret string
	err = db.QueryRow(getSecretStmt, s.projId, key).Scan(&secret)
	if err != nil {
		return "", err
	}

	return secret, nil
}

func (s *sqliteSecretStore) WriteSecret(key string, value string) error {
	db, err := s.ensureDb()
	if err != nil {
		return err
	}

	_, err = db.Exec(setSecretStmt, s.projId, key, value)
	if err != nil {
		if sqlite.IsUniqueConstraintError(err) {
			return fmt.Errorf("secret key '%s': %w", key, secrets.ErrSecretExists)
		}
		return err
	}

	return nil
}

func (s *sqliteSecretStore) ensureDb() (*sql.DB, error) {
	if s.db != nil {
		return s.db, nil
	}

	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		fmt.Printf("project.ensureDb: %v\n", err)
		return nil, err
	}
	s.db = db
	return s.db, nil
}

func NewSqlite(vaultPath, projectId string) secrets.SecretStore {
	store := &sqliteSecretStore{
		dbPath: vaultPath,
		projId: projectId,
	}

	_, err := store.ensureDb()
	if err != nil {
		return nil
	}
	return store
}
func (s *sqliteSecretStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
