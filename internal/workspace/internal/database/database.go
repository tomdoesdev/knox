package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tomdoesdev/knox/internal/error_codes"
	workspaceErrors "github.com/tomdoesdev/knox/internal/workspace/internal"
	"github.com/tomdoesdev/knox/kit/errs"
	"github.com/tomdoesdev/knox/kit/fs"
)

type Path struct {
	baseConnString string
	path           string
}

func NewPath(path string) *Path {
	isDir := fs.IsDir(path)

	if isDir {
		if filepath.Base(path) != workspaceErrors.DataDirectoryName {
			path = filepath.Join(path, workspaceErrors.DataDirectoryName, workspaceErrors.DatabaseFileName)
		}
		if filepath.Base(path) != workspaceErrors.DatabaseFileName {
			path = filepath.Join(path, workspaceErrors.DatabaseFileName)
		}
	}

	slog.Debug("created database path", slog.String("given_path", path), slog.String("used_path", path))
	dbp := Path{path: path}

	return &dbp
}

func (p *Path) String() string {
	return p.path
}

func (p *Path) ConnectionString(params ...string) string {

	if p.baseConnString == "" {
		p.baseConnString = fmt.Sprintf("file:%s", p.path)
	}

	var builder strings.Builder
	builder.WriteString(p.baseConnString)

	if len(params) != 0 {
		builder.WriteString("?")

		for i, param := range params {
			if i < len(params)-1 {
				builder.WriteString("&")
			}
			builder.WriteString(param)
		}
	}

	return builder.String()
}

const validationQuery = `
 SELECT COUNT(name) from main.sqlite_master WHERE type='table' AND name IN ('linked_vaults','workspace_settings');	
`

type Database struct {
	db *sql.DB
}

func (w *Database) Close() error {
	if w.db != nil {
		return w.db.Close()
	}

	return nil
}

// DB returns the underlying database connection for queries
func (w *Database) DB() *sql.DB {
	return w.db
}

func EnsureWorkspaceDatabase(path *Path) (*Database, error) {
	if exists := IsDatabaseExists(path); exists {
		slog.Debug("ensure database: database exists", slog.String("path", path.ConnectionString()))
		return OpenWorkspaceDatabase(path)
	}

	slog.Debug("ensure database: database doesnt exist", slog.String("path", path.ConnectionString()))
	return CreateWorkspaceDatabase(path)

}
func CreateWorkspaceDatabase(path *Path) (*Database, error) {
	if exists := IsDatabaseExists(path); exists {
		return nil, errs.Wrap(workspaceErrors.ErrInvalidDatabase, error_codes.DatabaseFailureErrCode, "database already exists")
	}

	db, err := sql.Open("sqlite3", path.ConnectionString())
	if err != nil {
		return nil, errs.Wrap(err, error_codes.CreateFailureErrCode, "failed to open database").WithContext("path", path)
	}

	_, err = db.Exec(tablesSchema)
	if err != nil {
		return nil, errs.Wrap(err, error_codes.CreateFailureErrCode, "failed to create database").WithContext("path", path)
	}

	if !isValidWorkspaceDatabase(db) {
		return nil, workspaceErrors.ErrInvalidDatabase.WithContext("path", path)
	}

	return &Database{db: db}, nil
}

func OpenWorkspaceDatabase(path *Path) (*Database, error) {
	db, err := sql.Open("sqlite3", path.ConnectionString())
	if err != nil {
		return nil, errs.Wrap(err, error_codes.CreateFailureErrCode, "failed to open database").WithContext("path", path)
	}

	if !isValidWorkspaceDatabase(db) {
		return nil, workspaceErrors.ErrInvalidDatabase.WithContext("path", path)
	}

	return &Database{db: db}, nil
}

func isValidWorkspaceDatabase(db *sql.DB) bool {
	err := db.Ping()
	if err != nil {
		return false
	}

	var count int
	err = db.QueryRow(validationQuery).Scan(&count)
	if err != nil {
		return false
	}

	return count == 2

}

func IsDatabaseExists(path *Path) bool {
	exists, err := fs.IsExist(path.String())
	if err != nil {
		return false
	}

	if !exists {
		return false
	}

	if fs.IsDir(path.String()) && filepath.Base(path.String()) != workspaceErrors.DataDirectoryName {
		// No workspace data dir
		return false
	}

	if fs.IsFile(path.String()) && filepath.Base(path.String()) != workspaceErrors.DatabaseFileName {
		// No database file
		return false
	}

	return true
}
