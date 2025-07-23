package vault

import (
	"database/sql"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tomdoesdev/knox/kit/errs"
)

const vaultTableSchema = `
CREATE TABLE IF NOT EXISTS secrets (
    id INTEGER PRIMARY KEY,
	collection_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,

UNIQUE (collection_id, key),
FOREIGN KEY (collection_id) REFERENCES collections(id)
);

CREATE TABLE IF NOT EXISTS collections(
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	description TEXT
);
`

func createSqliteFile(dsp string) error {
	slog.Debug("creating vault", "path", dsp)
	if dsp == "" {
		return ErrDatasourceUnreachable
	}

	db, err := sql.Open("sqlite3", dsp)
	if err != nil {
		return errs.Wrap(err, ErrVaultCreationFailed.Code, "failed to open database for creation").
			WithContext("path", dsp)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			slog.Warn("closing database failed", "err", err)
		}
	}(db)

	if err := configureSQLiteSecurity(db); err != nil {
		return errs.Wrap(err, ErrVaultCreationFailed.Code, "failed to configure SQLite security").
			WithContext("path", dsp)
	}

	if _, err := db.Exec(vaultTableSchema); err != nil {
		return errs.Wrap(err, ErrVaultCreationFailed.Code, "failed to create database schema").
			WithContext("path", dsp)
	}

	// Create default 'global' collection
	_, err = db.Exec("INSERT INTO collections (name, description) VALUES ('global', 'Default collection for vault')")
	if err != nil {
		return errs.Wrap(err, ErrVaultCreationFailed.Code, "failed to create default global collection").
			WithContext("path", dsp)
	}

	return nil
}

func configureSQLiteSecurity(db *sql.DB) error {
	// Enable secure delete - overwrites deleted data with zeros
	// Options: OFF (default), ON (secure), FAST (compromise)
	_, err := db.Exec("PRAGMA secure_delete = FAST")
	if err != nil {
		return errs.Wrap(err, ErrVaultCreationFailed.Code, "failed to configure secure delete")
	}

	// Enable foreign key constraints for data integrity
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return errs.Wrap(err, ErrVaultCreationFailed.Code, "failed to enable foreign key constraints")
	}

	// Set journal mode to WAL for better concurrency and crash recovery
	_, err = db.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		return errs.Wrap(err, ErrVaultCreationFailed.Code, "failed to set WAL journal mode")
	}

	// Enable synchronous mode for better durability
	_, err = db.Exec("PRAGMA synchronous = NORMAL")
	if err != nil {
		return errs.Wrap(err, ErrVaultCreationFailed.Code, "failed to set synchronous mode")
	}

	return nil
}
