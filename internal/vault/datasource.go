package vault

import (
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/tomdoesdev/knox/kit/fs"
	"github.com/tomdoesdev/knox/pkg/errs"
)

type Datasource string

func (d Datasource) String() string {
	return string(d)
}

func IsVault(datasourcePath string) (bool, error) {
	var result string

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
		return false, errs.Wrap(err,
			VaultConnectionCode,
			"failed to open database").
			WithContext("datasourcePath", datasourcePath)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='secrets'").Scan(&result)
	if err != nil {
		return false, errs.Wrap(err, VaultIntegrityCode, "failed to check vault tables").
			WithContext("datasourcePath", datasourcePath)
	}

	return result == "secrets", nil

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
		return nil, errs.Wrap(err, DatasourceCode, "failed to get datasource path")
	}

	if !fs.IsDir(filepath.Dir(dsp.dsPath)) {
		slog.Debug("creating directories")
		err := os.MkdirAll(filepath.Dir(dsp.dsPath), 0700)
		if err != nil {
			return nil, errs.Wrap(err, VaultCreationCode, "failed to create datasource directory")
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
		return "", errs.Wrap(err,
			DatasourceCode,
			"datasource is unreachable or not a valid sqlite file",
		).WithContext("path", f.dsPath)
	}

	f.datasource = Datasource(f.dsPath)

	return f.datasource, nil
}

func updateFilesystemDatasourcePath(f *filesystemDatasource) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return errs.Wrap(err, DatasourceCode, "failed to get user home directory")
	}

	if override, exists := os.LookupEnv(KnoxRootEnvVar); exists {
		home = override
	}

	f.dsPath = filepath.Join(home, DefaultDirName, DefaultFileName)
	return nil
}
