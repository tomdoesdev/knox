package vault

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tomdoesdev/knox/internal/v1/constants"
	"github.com/tomdoesdev/knox/kit/fs"
	"github.com/tomdoesdev/knox/pkg/errs"
)

var (
	CreationFailedCode errs.ErrorCode
)

var (
	ErrDirectoryNotCreated = errs.New(CreationFailedCode, "failed to create directory")
	ErrCantOpenVault       = errs.New(CreationFailedCode, "failed to open vault database")
	ErrBadLocation         = errs.New(CreationFailedCode, "invalid location struct")
)

type LocationDetails struct {
	Directory string
	Filepath  string
}

const Schema = `
create table if not exists vault_config (
	key text not null primary Key,
	value text not null
);

CREATE TABLE secrets (
id INTEGER PRIMARY KEY,
project_id TEXT NOT NULL,
key TEXT NOT NULL,
value TEXT NOT NULL,

UNIQUE (project_id,key)
)`

func FindLocation() (*LocationDetails, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	if override, exists := os.LookupEnv(constants.EnvVarKnoxRoot); exists {
		home = override
	}

	fp := filepath.Join(home, constants.DefaultKnoxDirName, constants.DefaultVaultFileName)
	dir := filepath.Join(home, constants.DefaultKnoxDirName)

	return &LocationDetails{
		Directory: dir,
		Filepath:  fp,
	}, nil
}

func Exists(details *LocationDetails) (bool, error) {
	if details == nil {
		return false, nil
	}

	// Check if the directory exists and create it if need be
	if !fs.IsDir(details.Directory) {
		err := os.MkdirAll(details.Directory, 0700)
		if err != nil {
			return false, ErrDirectoryNotCreated
		}
	}

	isVault, err := IsVault(details)
	if err != nil {
		return false, err
	}

	return isVault, nil
}

// ExistsOrCreate checks if a vault exists at the specified location, otherwise it creates it
func ExistsOrCreate(details *LocationDetails) (bool, error) {
	exists, err := Exists(details)
	if err != nil {
		return false, err
	}

	if exists {
		return true, nil
	}

	err = createVault(details)
	if err != nil {
		return false, err
	}

	return true, nil
}

func IsVault(details *LocationDetails) (bool, error) {
	var result string

	if details == nil {
		return false, nil
	}

	if !fs.IsFile(details.Filepath) {
		return false, nil
	}

	db, err := sql.Open("sqlite3", details.Filepath)
	if err != nil {
		return false, errs.Wrap(err, CreationFailedCode, "failed to open database")
	}

	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='secrets'").Scan(&result)
	if err != nil {
		return false, errs.Wrap(err, CreationFailedCode, "failed to check if secrets exists")
	}

	return result == "secrets", nil

}
func configureSQLiteSecurity(db *sql.DB) error {
	// Enable secure delete - overwrites deleted data with zeros
	// Options: OFF (default), ON (secure), FAST (compromise)
	_, err := db.Exec("PRAGMA secure_delete = FAST")
	if err != nil {
		return err
	}

	// Enable foreign key constraints for data integrity
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return err
	}

	// Set journal mode to WAL for better concurrency and crash recovery
	_, err = db.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		return err
	}

	// Enable synchronous mode for better durability
	_, err = db.Exec("PRAGMA synchronous = NORMAL")
	if err != nil {
		return err
	}

	return nil
}

func isInitialized(db *sql.DB) bool {
	var result string

	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='secrets'").Scan(&result)
	if err != nil {
		return false
	}

	return result == "secrets"
}
func createVault(details *LocationDetails) error {
	if details == nil {
		return ErrBadLocation
	}
	db, err := sql.Open("sqlite3", details.Filepath)
	if err != nil {
		return errs.Wrap(err, CreationFailedCode, "failed to open database")
	}

	// Configure SQLite security settings
	if err := configureSQLiteSecurity(db); err != nil {
		return errs.Wrap(err, CreationFailedCode, "failed to configure SQLite security")
	}

	if isInitialized(db) {
		return nil
	}

	if _, err := db.Exec(Schema); err != nil {
		return errs.Wrap(err, CreationFailedCode, "failed to initialize database")
	}

	return nil
}
