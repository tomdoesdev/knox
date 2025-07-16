package vault

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tomdoesdev/knox/kit/fs"
)

const schema = `
CREATE TABLE vault (
id INTEGER PRIMARY KEY,
project_id TEXT NOT NULL,
key TEXT NOT NULL,
value TEXT NOT NULL,

UNIQUE (project_id,key)
)`

func EnsureVaultExists(vaultpath string) error {
	dir := path.Dir(vaultpath)

	dirExists, err := fs.Exists(dir)
	if err != nil {
		return fmt.Errorf("vault:create:dirExists: %w", err)
	}

	vaultExists, err := fs.Exists(vaultpath)

	if err != nil {
		return fmt.Errorf("vault:create:vaultExists: %w", err)
	}

	if !dirExists {
		err = os.MkdirAll(dir, 0600)
		if err != nil {
			return fmt.Errorf("vault:create:mkDir: %w", err)
		}
	}

	if !vaultExists {
		err = createSqliteStore(vaultpath)
		if err != nil {
			return fmt.Errorf("vault:create:createSqliteStore: %w", err)
		}
	}

	return nil
}

func createSqliteStore(vaultPath string) error {
	exists, err := fs.Exists(vaultPath)
	if err != nil {
		return err
	}
	if exists {
		slog.Info("vault already exists", slog.String("path", vaultPath))
		return nil
	}

	db, err := sql.Open("sqlite3", vaultPath)

	if err != nil {
		return fmt.Errorf("createSqliteStore: sql.Open:  %w", err)
	}
	_, err = db.Exec(schema)
	if err != nil {
		return fmt.Errorf("createSqliteStore: execSchema: %w", err)
	}
	return nil
}
