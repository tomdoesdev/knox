package secrets

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
	knoxerr "github.com/tomdoesdev/knox/internal/errors"
)

type SecretStorer interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Del(key string) error
}

type SqliteStore struct {
}

func NewSqliteStore() *SqliteStore {
	return &SqliteStore{}
}

func (p *SqliteStore) NewVault(opts *VaultOptions) error {
	fmt.Println(opts.Path)
	info, err := os.Stat(path.Join(opts.Path, "knox.db"))

	fmt.Printf("%#v\n", info)
	if err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
			slog.Info("creating a new sqlite store")
			err := createVaultFile(opts)
			if err != nil {
				return err
			}
		default:
			return err
		}
	}

	return nil
}

func createVaultFile(opts *VaultOptions) error {
	fmt.Printf("Creating vault file at %s\n", path.Join(opts.Path, "knox.db"))
	db, err := sql.Open("sqlite3", path.Join(opts.Path, "knox.db"))
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

	CREATE TABLE secrets (
		id INTEGER PRIMARY KEY,
		namespace TEXT NOT NULL DEFAULT 'global',
		key TEXT NOT NULL,
		description TEXT NOT NULL DEFAULT ''
	)`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	slog.Info(sqlStmt)
	return nil
}

func (p *SqliteStore) Set(key, value string) error {
	return knoxerr.ErrNotImplemented
}

func (p *SqliteStore) Get(key string) (string, error) {
	return "", knoxerr.ErrNotImplemented
}

func (p *SqliteStore) Del(key string) error {
	return knoxerr.ErrNotImplemented
}
