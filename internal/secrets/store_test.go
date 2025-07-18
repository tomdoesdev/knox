package secrets

import (
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tomdoesdev/knox/pkg/errs"
)

func TestVCSProtection(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create a .git directory to simulate a git repo
	gitDir := filepath.Join(tmpDir, ".git")
	err := os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("failed to create .git directory: %v", err)
	}

	// Try to create a vault in the git repo (should fail)
	vaultPath := filepath.Join(tmpDir, "vault.db")
	_, err = NewFileSecretStore(vaultPath, "test-project", NewNoOpEncryptionHandler())
	if err == nil {
		t.Error("expected error when creating vault in git repo, but got none")
	}

	// Verify it's the VCS protection error
	if !errs.Is(err, StoreVCSWarningCode) {
		t.Errorf("expected VCS warning error, got: %v", err)
	}

	// Try again with force flag (should succeed)
	store, err := NewFileSecretStoreWithOptions(vaultPath, "test-project", NewNoOpEncryptionHandler(), true)
	if err != nil {
		t.Errorf("unexpected error when using force flag: %v", err)
	}
	if store != nil {
		store.Close()
	}
}

func TestVCSProtection_ExistingFile(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create a .git directory to simulate a git repo
	gitDir := filepath.Join(tmpDir, ".git")
	err := os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("failed to create .git directory: %v", err)
	}

	// Create an existing vault file
	vaultPath := filepath.Join(tmpDir, "existing-vault.db")
	existingStore, err := NewFileSecretStoreWithOptions(vaultPath, "test-project", NewNoOpEncryptionHandler(), true)
	if err != nil {
		t.Fatalf("failed to create existing vault: %v", err)
	}
	existingStore.Close()

	// Try to open existing vault (should succeed even in git repo)
	store, err := NewFileSecretStore(vaultPath, "test-project", NewNoOpEncryptionHandler())
	if err != nil {
		t.Errorf("unexpected error when opening existing vault: %v", err)
	}
	if store != nil {
		store.Close()
	}
}

func TestVCSProtection_NoVCS(t *testing.T) {
	// Create a temporary directory without VCS
	tmpDir := t.TempDir()

	// Try to create a vault (should succeed)
	vaultPath := filepath.Join(tmpDir, "vault.db")
	store, err := NewFileSecretStore(vaultPath, "test-project", NewNoOpEncryptionHandler())
	if err != nil {
		t.Errorf("unexpected error when creating vault outside VCS: %v", err)
	}
	if store != nil {
		store.Close()
	}
}

func TestCheckVCSProtection(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create a .git directory to simulate a git repo
	gitDir := filepath.Join(tmpDir, ".git")
	err := os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("failed to create .git directory: %v", err)
	}

	vaultPath := filepath.Join(tmpDir, "vault.db")

	// Test without force (should fail)
	err = checkVCSProtection(vaultPath, false)
	if err == nil {
		t.Error("expected VCS protection error")
	}

	// Test with force (should succeed)
	err = checkVCSProtection(vaultPath, true)
	if err != nil {
		t.Errorf("unexpected error with force flag: %v", err)
	}
}
