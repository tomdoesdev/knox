package vault

import (
	"database/sql"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tomdoesdev/knox/kit/errs"
)

const vaultTableSchema = `
CREATE TABLE IF NOT EXISTS projects (
id integer primary key,
name text not null unique,
description text
);

CREATE TABLE IF NOT EXISTS secrets (
id INTEGER PRIMARY KEY,
project_id integer not null,
key TEXT NOT NULL,
value TEXT NOT NULL,

UNIQUE (project_id,key),
FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS secret_tags (
secret_id integer not null,
tag text not null,

PRIMARY KEY (secret_id, tag),
FOREIGN KEY (secret_id) REFERENCES secrets(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS project_tags (
project_id integer not null, 
tag text not null,
PRIMARY KEY (project_id, tag),
FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
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
