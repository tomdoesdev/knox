package vault

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/adrg/xdg"
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

func DefaultPath() string {
	fileName := "knox.vault"
	vaultDir := path.Join(xdg.DataHome, "knox")
	return path.Join(vaultDir, fileName)
}

func Exists(vaultPath string) bool {
	p := vaultPath
	if p == "" {
		p = DefaultPath()
	}
	exists, err := fs.Exists(p)
	if err != nil {
		slog.Debug("error checking vault existence", "path", vaultPath, "error", err)
		return false
	}
	return exists
}
func EnsureVaultExists(vaultpath string) error {
	dir := path.Dir(vaultpath)

	slog.Debug("ensure vault dir exists", slog.String("path", dir))

	dirExists, err := fs.Exists(dir)
	if err != nil {
		return fmt.Errorf("vault:create:dirExists: %w", err)
	}

	vaultExists, err := fs.Exists(vaultpath)
	if err != nil {
		return fmt.Errorf("vault:create:vaultExists: %w", err)
	}

	slog.Debug("things exists",
		slog.Bool("dir", dirExists),
		slog.Bool("vault", vaultExists),
	)

	if !dirExists {
		slog.Debug("vault dir does not exist.", slog.String("path", dir))
		err = os.MkdirAll(dir, 0600)
		if err != nil {
			return fmt.Errorf("vault:create:mkDir: %w", err)
		}
	}

	if !vaultExists {
		slog.Debug("vault does not exist at path.",
			slog.String("dirpath", dir),
			slog.String("vaultpath", vaultpath),
		)
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

	slog.Debug("creating sqlite store", slog.String("path", vaultPath))
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
