package vault

import (
	"database/sql"
	"log/slog"
	"path"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tomdoesdev/knox/internal/config"
	"github.com/tomdoesdev/knox/kit/errkit"
)

type SqliteStore struct {
}

func initSqliteVault(conf *config.ApplicationConfig) error {
	db, err := sql.Open("sqlite3", path.Join(conf.VaultDir, "knox.db"))
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}(db)

	sqlStmt := `
	CREATE TABLE vault_info (
		version INTEGER
	);

	CREATE TABLE vault (
		id INTEGER PRIMARY KEY,
		namespace TEXT NOT NULL DEFAULT 'default',
		key TEXT NOT NULL,
		value TEXT NOT NULL,
		description TEXT NOT NULL DEFAULT ''
	)`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}

func (p *SqliteStore) Set(key, value string) error {
	return errkit.ErrNotImplemented
}

func (p *SqliteStore) Get(key string) (string, error) {
	return "", errkit.ErrNotImplemented
}

func (p *SqliteStore) Del(key string) error {
	return errkit.ErrNotImplemented
}
