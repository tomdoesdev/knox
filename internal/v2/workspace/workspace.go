package workspace

import "time"

type LinkedVault struct {
	ID        int       `json:"id"`
	Alias     string    `json:"alias"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
}

type LinkedProject struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	VaultID     int       `json:"vault_id"`
	ProjectName *string   `json:"project_name,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// FindWorkspace finds the nearest .knox directory, traversing up the directory tree until it finds it.
func FindWorkspace() {

}

func CreateWorkspace() {

}

func IsWorkspace() bool {
	return true
}
