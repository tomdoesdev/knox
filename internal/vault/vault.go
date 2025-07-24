package vault

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tomdoesdev/knox/kit/errs"
)

type Vault struct {
	datasource Datasource
	db         *sql.DB
}

func (v *Vault) Close() error {
	return v.db.Close()
}

func Open(dsp DatasourceProvider) (*Vault, error) {
	datasource, err := dsp.Datasource()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", datasource.String())
	if err != nil {
		return nil, errs.Wrap(err, ErrVaultConnectionFailed.Code, "failed to open database connection").
			WithContext("datasource", datasource.String())
	}

	err = db.Ping()
	if err != nil {
		return nil, errs.Wrap(err, ErrVaultConnectionFailed.Code, "failed to ping database").
			WithContext("datasource", datasource.String())
	}

	return &Vault{
		datasource: datasource,
		db:         db,
	}, nil
}
