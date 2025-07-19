package template

import (
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	secrets2 "github.com/tomdoesdev/knox/internal/v1/secrets"
)

func TestIntegration_TemplateProcessingWithSecretStore(t *testing.T) {
	tmpDir := t.TempDir()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	vaultPath := filepath.Join(tmpDir, "test.db")
	projectID := "test-project"

	store, err := secrets2.NewFileSecretStore(vaultPath, projectID, secrets2.NewNoOpEncryptionHandler())
	if err != nil {
		t.Fatalf("failed to create secret store: %v", err)
	}
	defer store.Close()

	err = store.WriteSecret("API_KEY", "secret123")
	if err != nil {
		t.Fatalf("failed to write secret: %v", err)
	}

	err = store.WriteSecret("DB_PASSWORD", "password456")
	if err != nil {
		t.Fatalf("failed to write secret: %v", err)
	}

	templateContent := `# Application configuration
API_KEY={{.Secret "API_KEY"}}
DB_PASSWORD={{.Secret "DB_PASSWORD"}}
DEBUG={{.Default "DEBUG_MODE" "false"}}
ENV={{.Default "ENVIRONMENT" "development"}}
PATH={{.Env "PATH"}}`

	templatePath := filepath.Join(tmpDir, ".env.template")
	err = os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("failed to write template file: %v", err)
	}

	os.Setenv("PATH", "/usr/local/bin:/usr/bin:/bin")
	defer os.Unsetenv("PATH")

	processor := NewProcessor(store)
	result, err := processor.ProcessFile("")
	if err != nil {
		t.Fatalf("failed to process template: %v", err)
	}

	expected := map[string]string{
		"API_KEY":     "secret123",
		"DB_PASSWORD": "password456",
		"DEBUG":       "false",
		"ENV":         "development",
		"PATH":        "/usr/local/bin:/usr/bin:/bin",
	}

	if len(result) != len(expected) {
		t.Errorf("expected %d variables, got %d", len(expected), len(result))
	}

	for key, expectedValue := range expected {
		if actualValue, exists := result[key]; !exists {
			t.Errorf("expected key %q not found", key)
		} else if actualValue != expectedValue {
			t.Errorf("for key %q: expected %q, got %q", key, expectedValue, actualValue)
		}
	}
}

func TestIntegration_ProcessFileOptional_WithoutTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	vaultPath := filepath.Join(tmpDir, "test.db")
	projectID := "test-project"

	store, err := secrets2.NewFileSecretStore(vaultPath, projectID, secrets2.NewNoOpEncryptionHandler())
	if err != nil {
		t.Fatalf("failed to create secret store: %v", err)
	}
	defer store.Close()

	processor := NewProcessor(store)
	result, err := processor.ProcessFileOptional("")
	if err != nil {
		t.Errorf("unexpected error for missing template file: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected empty result for missing template file, got %d variables", len(result))
	}
}
