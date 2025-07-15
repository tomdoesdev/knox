package vault

import (
	"database/sql"
	"fmt"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tomdoesdev/knox/kit/errkit"
	"github.com/tomdoesdev/knox/kit/fskit"
)

const schema = `
CREATE TABLE vault (
id INTEGER PRIMARY KEY,
project_id TEXT NOT NULL,
key TEXT NOT NULL,
value TEXT NOT NULL,

UNIQUE (project_id,key)
)`

func EnsureVaultExists(vPath string) error {
	dirExists, err := fskit.Exists(path.Dir(vPath))
	if err != nil {
		return fmt.Errorf("vault:create:dirExists: %w", err)
	}

	vaultExists, err := fskit.Exists(vPath)

	if err != nil {
		return fmt.Errorf("vault:create:vaultExists: %w", err)
	}

	if !dirExists {
		err = os.MkdirAll(vPath, 0600)
		if err != nil {
			return fmt.Errorf("vault:create:mkDir: %w", err)
		}
	}

	if !vaultExists {
		err = createSqliteStore(vPath)
		if err != nil {
			return fmt.Errorf("vault:create:createSqliteStore: %w", err)
		}
	}

	return nil
}

func createSqliteStore(vaultPath string) error {
	exists, err := fskit.Exists(vaultPath)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("vault already exists %w", errkit.ErrAlreadyExists)
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
