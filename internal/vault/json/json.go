package json

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/tomdoesdev/knox/internal"
	"github.com/tomdoesdev/knox/internal/vault"
	"github.com/tomdoesdev/knox/kit/errs"
	"github.com/tomdoesdev/knox/kit/fs"
)

const (
	ECodeJSONFileFailure errs.Code = "JSON_FILE_FAILURE"
)

type (
	VaultFile struct {
		path    string
		Id      string            `json:"id"`
		Secrets map[string]string `json:"secrets"`
	}
)

type (
	Vault struct {
		dirty bool
		mutex sync.RWMutex

		VaultFile
	}
)

func (j *Vault) Save() error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if !j.dirty {
		return nil
	}

	data, err := json.MarshalIndent(j.VaultFile, "", " ")
	if err != nil {
		return errs.Wrap(err, vault.ECodeFailedToPersist, "failed to marshal secrets for writing").WithPath(j.path)
	}

	err = fs.WriteFile(j.path, data, internal.WorkspaceFilePermissions)
	if err != nil {
		return errs.Wrap(err, vault.ECodeFailedToPersist, "failed to write secrets json").WithPath(j.path)
	}
	j.dirty = false
	return nil
}

// Delete performs an in-memory deletion of a specified key.
func (j *Vault) Delete(key string) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	delete(j.Secrets, key)
	j.dirty = true
	return nil
}

// Set performs an in-memory set of the specified key,value
func (j *Vault) Set(key string, value string) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.Secrets[key] = value
	j.dirty = true
	return nil
}

// Get returns the in-memory value of the specified key
func (j *Vault) Get(key string) string {
	j.mutex.RLock()
	defer j.mutex.RUnlock()
	return j.Secrets[key]
}

func NewVault(json *VaultFile) *Vault {
	return &Vault{
		dirty:     false,
		mutex:     sync.RWMutex{},
		VaultFile: *json,
	}
}

func OpenVault(path string) (*Vault, error) {
	bytes, err := fs.ReadFile(path)
	if err != nil {
		return nil, errs.Wrap(err, ECodeJSONFileFailure, "failed to open vault json").WithPath(path)
	}

	j := new(VaultFile)
	j.path = path

	err = json.Unmarshal(bytes, j)
	if err != nil {
		return nil, errs.Wrap(err, ECodeJSONFileFailure, "failed to unmarshal json").WithPath(path)
	}
	return NewVault(j), nil
}

func CreateVaultFile(path string, id uuid.UUID) error {
	return createFile(path, id.String())
}

func CreateDefaultVault(path string) error {
	return createFile(path, "default")
}

func createFile(path string, id string) error {
	j := new(VaultFile)
	j.Id = id
	j.Secrets = make(map[string]string)

	fileName := fmt.Sprintf("%s.json", id)

	if fs.IsDir(path) {
		path = filepath.Join(path, fileName)
	}

	bytes, err := json.Marshal(j)
	if err != nil {
		return errs.Wrap(err, ECodeJSONFileFailure, "failed to marshal json").WithPath(path)
	}

	err = fs.WriteFile(path, bytes, internal.WorkspaceFilePermissions)
	if err != nil {
		return errs.Wrap(err, ECodeJSONFileFailure, "failed to create vault json").WithPath(path)
	}

	return nil
}
