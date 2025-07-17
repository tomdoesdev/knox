package internal

import (
	"fmt"
	"io"

	"github.com/tomdoesdev/knox/internal/secrets"
	"github.com/tomdoesdev/knox/kit/fs"
)

type Workspace struct {
	ProjectID       string
	VaultFilePath   fs.FilePath
	ProjectFilePath fs.FilePath
}
type Knox struct {
	Store     secrets.SecretStore
	Workspace *Workspace
}

type KnoxOptions struct {
	SecretStore secrets.SecretStore
	Workspace   *Workspace
}

func NewKnox(opts *KnoxOptions) (*Knox, error) {
	if opts == nil {
		return nil, fmt.Errorf("options cannot be nil")
	}

	if opts.SecretStore == nil {
		return nil, fmt.Errorf("secret store cannot be nil")
	}

	return &Knox{
		Store:     opts.SecretStore,
		Workspace: opts.Workspace,
	}, nil
}

func (k *Knox) Close() error {
	var closeErr error
	if k.Store != nil {
		if closer, ok := k.Store.(io.Closer); ok {
			closeErr = closer.Close()
		}
	}
	return closeErr
}
