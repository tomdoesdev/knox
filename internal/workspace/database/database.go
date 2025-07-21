package database

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tomdoesdev/knox/internal/workspace/constants"
	"github.com/tomdoesdev/knox/internal/workspace/errors"
	"github.com/tomdoesdev/knox/kit/fs"
	"github.com/tomdoesdev/knox/pkg/errs"
)

type databasePath struct {
	baseConnString string
	path           string
}

func (p *databasePath) String() string {
	return p.path
}

func (p *databasePath) ConnectionString(params ...string) string {

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

func newDatabasePath(path string) (*databasePath, error) {
	err := IsDatabaseExists(path)
	if err != nil {
		return nil, err
	}

	dbp := databasePath{path: path}

	db, err := sql.Open("sqlite3", dbp.ConnectionString("mode=ro"))
	if err != nil {
		return nil, errs.Wrap(errors.ErrInvalidDatabase, errors.DatabaseFailureCode, "failed to open database")
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	err = db.Ping()
	if err != nil {
		return nil, errs.Wrap(errors.ErrInvalidDatabase, errors.DatabaseFailureCode, "failed to open database")
	}

	return &dbp, nil
}

const validationQuery = `
 SELECT COUNT(name) from main.sqlite_master WHERE type='table' AND name IN ('linked_vaults', 'linked_projects','workspace_meta');	
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

func CreateWorkspaceDatabase(path string) (*Database, error) {

	if err := IsDatabaseExists(path); err != nil {
		return nil, errs.Wrap(errors.ErrInvalidDatabase, errors.DatabaseFailureCode, "database already exists")
	}

	dbp, err := newDatabasePath(path)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbp.ConnectionString())
	if err != nil {
		return nil, errs.Wrap(err, errors.CreateFailureCode, "failed to open database").WithContext("path", path)
	}

	_, err = db.Exec(tablesSchema)
	if err != nil {
		return nil, errs.Wrap(err, errors.CreateFailureCode, "failed to create database").WithContext("path", path)
	}

	if !isValidWorkspaceDatabase(db) {
		return nil, errors.ErrInvalidDatabase.WithContext("path", path)
	}

	return &Database{db: db}, nil
}

func OpenWorkspaceDatabase(path *databasePath) (*Database, error) {
	db, err := sql.Open("sqlite3", path.ConnectionString())
	if err != nil {
		return nil, errs.Wrap(err, errors.CreateFailureCode, "failed to open database").WithContext("path", path)
	}

	if !isValidWorkspaceDatabase(db) {
		return nil, errors.ErrInvalidDatabase.WithContext("path", path)
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

	return count == 3

}

func IsDatabaseExists(path string) error {
	exists, err := fs.IsExist(path)
	if err != nil {
		return errors.ErrDatabasePathInvalid.WithContext("path", path)
	}

	if !exists {
		return errors.ErrDatabasePathInvalid.WithContext("path", path).WithContext("exists", exists)
	}

	if fs.IsDir(path) && filepath.Base(path) != constants.DataDirectoryName {
		return errors.ErrDatabasePathInvalid.WithContext("msg", "no workspace directory")
	}

	if fs.IsFile(path) && filepath.Base(path) != constants.DatabaseFileName {
		return errs.Wrap(errors.ErrInvalidDatabase, errors.DatabaseFailureCode, "no database file")
	}
	return nil
}
