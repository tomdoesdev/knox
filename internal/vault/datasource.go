package vault

import (
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/tomdoesdev/knox/kit/errs"
	"github.com/tomdoesdev/knox/kit/fs"
)

type Datasource string

func (d Datasource) String() string {
	return string(d)
}

func IsVault(datasourcePath string) (bool, error) {

	if datasourcePath == ":memory:" {
		return true, nil
	}

	if datasourcePath == "" {
		return false, nil
	}

	if !fs.IsFile(datasourcePath) {
		return false, nil
	}

	db, err := sql.Open("sqlite3", datasourcePath)
	if err != nil {
		return false, errs.Wrap(err, ErrVaultConnectionFailed.Code,
			"failed to open database").
			WithContext("datasourcePath", datasourcePath)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	// Check for both secrets and collections tables
	var tableCount int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name IN ('secrets', 'collections')").Scan(&tableCount)
	if err != nil {
		return false, errs.Wrap(err, ErrVaultIntegrityCheck.Code, "failed to check vault tables").
			WithContext("datasourcePath", datasourcePath)
	}

	return tableCount == 2, nil

}

type DatasourceProvider interface {
	Datasource() (Datasource, error)
}

type filesystemDatasource struct {
	datasource Datasource
	dsPath     string
}

func NewFileSystemDatasource() (DatasourceProvider, error) {
	dsp := &filesystemDatasource{}

	err := updateFilesystemDatasourcePath(dsp)
	if err != nil {
		return nil, errs.Wrap(err, ErrDatasourcePathInvalid.Code, "failed to get datasource path")
	}

	if !fs.IsDir(filepath.Dir(dsp.dsPath)) {
		slog.Debug("creating directories")
		err := os.MkdirAll(filepath.Dir(dsp.dsPath), 0700)
		if err != nil {
			return nil, errs.Wrap(err, ErrVaultCreationFailed.Code, "failed to create datasource directory")
		}
	}

	if !fs.IsFile(dsp.dsPath) {
		slog.Debug("creating file")
		err := createSqliteFile(dsp.dsPath)
		if err != nil {
			return nil, err
		}
	}

	return dsp, nil
}

func (f *filesystemDatasource) Datasource() (Datasource, error) {
	if f.datasource != "" {
		return f.datasource, nil
	}

	if isVault, err := IsVault(f.dsPath); err != nil || !isVault {
		return "", errs.Wrap(err, ErrDatasourceUnreachable.Code,
			"datasource is unreachable or not a valid sqlite file",
		).WithContext("path", f.dsPath)
	}

	f.datasource = Datasource(f.dsPath)

	return f.datasource, nil
}

func updateFilesystemDatasourcePath(f *filesystemDatasource) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return errs.Wrap(err, ErrDatasourcePathInvalid.Code, "failed to get user home directory")
	}

	if override, exists := os.LookupEnv(KnoxRootEnvVar); exists {
		home = override
	}

	f.dsPath = filepath.Join(home, DefaultDirName, DefaultFileName)
	return nil
}
