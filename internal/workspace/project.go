package workspace

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/tomdoesdev/knox/internal"
	"github.com/tomdoesdev/knox/kit/errs"
)

// Project represents a workspace project configuration
type Project struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	CreatedAt   time.Time         `json:"created_at"`
	SecretMap   map[string]string `json:"secret_map"`
}

// SecretReference represents a parsed secret reference in format "secret@vault/collection"
type SecretReference struct {
	Secret     string
	Vault      string
	Collection string
}

// NewProject creates a new project with the given name and description
func NewProject(name, description string) *Project {
	return &Project{
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		SecretMap:   make(map[string]string),
	}
}

// AddSecret adds a secret mapping to the project
func (p *Project) AddSecret(logicalName, secretPath string) {
	p.SecretMap[logicalName] = secretPath
}

// RemoveSecret removes a secret mapping from the project
func (p *Project) RemoveSecret(logicalName string) {
	delete(p.SecretMap, logicalName)
}

// GetSecret returns the secret path for a logical name
func (p *Project) GetSecret(logicalName string) (string, bool) {
	path, exists := p.SecretMap[logicalName]
	return path, exists
}

// ListSecrets returns all logical secret names in the project
func (p *Project) ListSecrets() []string {
	secrets := make([]string, 0, len(p.SecretMap))
	for name := range p.SecretMap {
		secrets = append(secrets, name)
	}
	return secrets
}

// ToJSON serializes the project to JSON bytes
func (p *Project) ToJSON() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

// FromJSON deserializes a project from JSON bytes
func FromJSON(data []byte) (*Project, error) {
	var project Project
	err := json.Unmarshal(data, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// ParseSecretReference parses a secret reference in format "secret@vault/collection"
func ParseSecretReference(ref string) (*SecretReference, error) {
	parts := strings.Split(ref, "@")
	if len(parts) != 2 {
		return nil, errs.New(internal.SecretInvalidCode, "invalid secret reference format, expected 'secret@vault/collection'").WithContext("reference", ref)
	}

	secret := strings.TrimSpace(parts[0])
	location := strings.TrimSpace(parts[1])

	if secret == "" {
		return nil, errs.New(internal.SecretInvalidCode, "secret name cannot be empty").WithContext("reference", ref)
	}

	if location == "" {
		return nil, errs.New(internal.SecretInvalidCode, "vault/collection location cannot be empty").WithContext("reference", ref)
	}

	// Parse vault/collection from location
	locationParts := strings.Split(location, "/")
	if len(locationParts) != 2 {
		return nil, errs.New(internal.SecretInvalidCode, "invalid location format, expected 'vault/collection'").WithContext("reference", ref).WithContext("location", location)
	}

	vault := strings.TrimSpace(locationParts[0])
	collection := strings.TrimSpace(locationParts[1])

	if vault == "" {
		return nil, errs.New(internal.SecretInvalidCode, "vault name cannot be empty").WithContext("reference", ref)
	}

	if collection == "" {
		return nil, errs.New(internal.SecretInvalidCode, "collection name cannot be empty").WithContext("reference", ref)
	}

	return &SecretReference{
		Secret:     secret,
		Vault:      vault,
		Collection: collection,
	}, nil
}

// ValidateName validates a project name
func ValidateName(name string) error {
	if name == "" {
		return errs.New(internal.ProjectInvalidCode, "project name cannot be empty")
	}

	// Project names should be valid filenames and reasonable identifiers
	validName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validName.MatchString(name) {
		return errs.New(internal.ProjectInvalidCode, "project name must contain only letters, numbers, underscores, and hyphens").WithContext("name", name)
	}

	if len(name) > 50 {
		return errs.New(internal.ProjectInvalidCode, "project name must be 50 characters or less").WithContext("name", name)
	}

	return nil
}

// Validate validates the project structure
func (p *Project) Validate() error {
	if err := ValidateName(p.Name); err != nil {
		return err
	}

	if p.SecretMap == nil {
		return errs.New(internal.ProjectInvalidCode, "project must have a secret map")
	}

	// Validate all secret references
	for logicalName, secretRef := range p.SecretMap {
		if logicalName == "" {
			return errs.New(internal.ProjectInvalidCode, "logical secret name cannot be empty")
		}

		_, err := ParseSecretReference(secretRef)
		if err != nil {
			return errs.Wrap(err, internal.ProjectInvalidCode, fmt.Sprintf("invalid secret reference for '%s'", logicalName))
		}
	}

	return nil
}

// ValidateWithVaults validates the project against available vaults
func (p *Project) ValidateWithVaults(availableVaults []string) error {
	if err := p.Validate(); err != nil {
		return err
	}

	vaultSet := make(map[string]bool)
	for _, vault := range availableVaults {
		vaultSet[vault] = true
	}

	// Check all secret references point to available vaults
	for logicalName, secretRef := range p.SecretMap {
		ref, err := ParseSecretReference(secretRef)
		if err != nil {
			return err // Already validated above, but being safe
		}

		if !vaultSet[ref.Vault] {
			return errs.New(internal.ProjectInvalidCode, fmt.Sprintf("secret '%s' references unknown vault '%s'", logicalName, ref.Vault))
		}
	}

	return nil
}
