package vault

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/tomdoesdev/knox/kit/fs"
	"github.com/tomdoesdev/knox/pkg/errs"
)

const (
	DatasourceIntegrityCode errs.ErrorCode = "invalid_datasource"
)

var (
	ErrDatasourceUnreachable = errs.New(DatasourceIntegrityCode, "datasource database unreachable")
)

const (
	FileSystem DatasourceType = "file_system"
	InMemory   DatasourceType = "in_memory"
)

type DatasourceType string
type Datasource string

func (d Datasource) String() string {
	return string(d)
}

func IsVault(d string) (bool, error) {
	slog.Debug("checking isVault")
	slog.Debug("datasource is", "datasource", fmt.Sprintf("%v", d))
	var result string

	if d == ":memory:" {
		return true, nil
	}

	if d == "" {
		return false, nil
	}

	if !fs.IsFile(d) {
		return false, nil
	}

	db, err := sql.Open("sqlite3", d)
	if err != nil {
		return false, errs.Wrap(err,
			DatasourceIntegrityCode,
			"failed to open database").WithContext("d", d)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='secrets'").Scan(&result)
	if err != nil {
		return false, err
	}

	return result == "secrets", nil

}

type DatasourceProvider interface {
	Datasource() (Datasource, error)
	SourceType() DatasourceType
}

type filesystemDatasource struct {
	datasource Datasource
	dsPath     string
}

func NewFileSystemDatasource() (DatasourceProvider, error) {
	dsp := &filesystemDatasource{}

	err := updateFilesystemDatasourcePath(dsp)
	if err != nil {
		return nil, errs.Wrap(err, DatasourceIntegrityCode, "failed to get datasource path")
	}

	if !fs.IsDir(filepath.Dir(dsp.dsPath)) {
		slog.Debug("creating directories")
		err := os.MkdirAll(filepath.Dir(dsp.dsPath), 0700)
		if err != nil {
			return nil, errs.Wrap(err, DatasourceIntegrityCode, "failed to create datasource directory")
		}
	}

	if !fs.IsFile(dsp.dsPath) {
		slog.Debug("creating file")
		err := createVault(dsp.dsPath)
		if err != nil {
			return nil, err
		}
	}

	return dsp, nil
}

func (f *filesystemDatasource) SourceType() DatasourceType {
	return FileSystem
}

func (f *filesystemDatasource) Datasource() (Datasource, error) {
	if f.datasource != "" {
		return f.datasource, nil
	}

	if isVault, err := IsVault(f.dsPath); err != nil || !isVault {
		return "", errs.Wrap(err,
			DatasourceIntegrityCode,
			"datasource is unreachable or not a valid sqlite file",
		).WithContext("path", f.dsPath)
	}

	f.datasource = Datasource(f.dsPath)

	return f.datasource, nil
}

func updateFilesystemDatasourcePath(f *filesystemDatasource) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return errs.Wrap(err, DatasourceIntegrityCode, "failed to get user home directory")
	}

	if override, exists := os.LookupEnv(KnoxRootEnvVar); exists {
		home = override
	}

	f.dsPath = filepath.Join(home, DefaultDirName, DefaultFileName)
	slog.Debug("updated fs datasource path", "path", f.dsPath)

	return nil
}
