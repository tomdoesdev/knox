package vault

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"path"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tomdoesdev/knox/internal/config"
	"github.com/tomdoesdev/knox/kit/fskit"
)

var (
	ErrVaultExists = errors.New("vault already exists")
)

var schema = `
CREATE TABLE vault_info (
version INTEGER
);

CREATE TABLE vault (
id INTEGER PRIMARY KEY,
project_id TEXT NOT NULL,
namespace TEXT NOT NULL DEFAULT 'default',
key TEXT NOT NULL,
value TEXT NOT NULL,
description TEXT NOT NULL DEFAULT '',

UNIQUE (project_id,namespace, key)
)`

func createSqliteStore(conf *config.ApplicationConfig) error {
	dbPath := path.Join(conf.VaultDir, conf.VaultFileName)

	exists, err := fskit.Exists(dbPath)
	if err != nil {
		return err
	}
	if exists {
		return ErrVaultExists
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("createSqliteStore: sql.Open:  %w", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			slog.Error("createSqliteStore: defer db.Close", err)
		}
	}(db)

	_, err = db.Exec(schema)
	if err != nil {
		return fmt.Errorf("createSqliteStore: execSchema: %w", err)
	}
	return nil
}
