package vault

import (
	"github.com/tomdoesdev/knox/kit/errkit"
)

type Options struct {
	Path string
	Name string
}
type Vault interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Del(key string) error
	Clear() error
}

func New(opts *Options) Vault {
	//storePath := path.Join(opts.Path, "knox.db")
	//
	//if fskit.Exists(storePath) == false {
	//
	//}
	return nil
}

type vaultImpl struct{}

func (s *vaultImpl) Set(key, value string) error {
	return errkit.ErrNotImplemented
}
func (s *vaultImpl) Get(key string) (string, error) {
	return "", errkit.ErrNotImplemented
}

func (s *vaultImpl) Del(key string) error {
	return errkit.ErrNotImplemented
}
func (s *vaultImpl) Clear() error {
	return errkit.ErrNotImplemented
}
