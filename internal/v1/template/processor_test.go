package template

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type mockSecretStore struct {
	secrets map[string]string
}

func (m *mockSecretStore) ReadSecret(key string) (string, error) {
	if value, exists := m.secrets[key]; exists {
		return value, nil
	}
	return "", fmt.Errorf("secret not found: %s", key)
}

func newMockSecretStore(secrets map[string]string) *mockSecretStore {
	return &mockSecretStore{secrets: secrets}
}

func TestProcessor_ProcessTemplate(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		secrets     map[string]string
		envVars     map[string]string
		expected    map[string]string
		expectError bool
	}{
		{
			name:     "simple secret replacement",
			template: "API_KEY={{.Secret \"API_KEY\"}}",
			secrets:  map[string]string{"API_KEY": "secret123"},
			expected: map[string]string{"API_KEY": "secret123"},
		},
		{
			name:     "environment variable",
			template: "PATH={{.Env \"PATH\"}}",
			envVars:  map[string]string{"PATH": "/usr/bin:/bin"},
			expected: map[string]string{"PATH": "/usr/bin:/bin"},
		},
		{
			name:     "default with existing secret",
			template: "DEBUG={{.Default \"DEBUG\" \"false\"}}",
			secrets:  map[string]string{"DEBUG": "true"},
			expected: map[string]string{"DEBUG": "true"},
		},
		{
			name:     "default with missing secret",
			template: "DEBUG={{.Default \"MISSING\" \"false\"}}",
			secrets:  map[string]string{},
			expected: map[string]string{"DEBUG": "false"},
		},
		{
			name: "multiple variables",
			template: `API_KEY={{.Secret "API_KEY"}}
DEBUG={{.Default "DEBUG" "false"}}
PATH={{.Env "PATH"}}`,
			secrets: map[string]string{"API_KEY": "secret123"},
			envVars: map[string]string{"PATH": "/usr/bin"},
			expected: map[string]string{
				"API_KEY": "secret123",
				"DEBUG":   "false",
				"PATH":    "/usr/bin",
			},
		},
		{
			name: "comments and empty lines",
			template: `# This is a comment
API_KEY={{.Secret "API_KEY"}}

# Another comment
DEBUG=true`,
			secrets:  map[string]string{"API_KEY": "secret123"},
			expected: map[string]string{"API_KEY": "secret123", "DEBUG": "true"},
		},
		{
			name:        "missing secret error",
			template:    "API_KEY={{.Secret \"MISSING\"}}",
			secrets:     map[string]string{},
			expectError: true,
		},
		{
			name:        "invalid template syntax",
			template:    "API_KEY={{.Secret \"API_KEY\"",
			secrets:     map[string]string{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			store := newMockSecretStore(tt.secrets)
			processor := NewProcessor(store)

			result, err := processor.ProcessTemplate(tt.template, "test")

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d variables, got %d", len(tt.expected), len(result))
			}

			for key, expectedValue := range tt.expected {
				if actualValue, exists := result[key]; !exists {
					t.Errorf("expected key %q not found", key)
				} else if actualValue != expectedValue {
					t.Errorf("for key %q: expected %q, got %q", key, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestProcessor_ProcessFile(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, ".env.template")

	templateContent := `API_KEY={{.Secret "API_KEY"}}
DEBUG={{.Default "DEBUG" "false"}}`

	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test template file: %v", err)
	}

	secrets := map[string]string{"API_KEY": "secret123"}
	store := newMockSecretStore(secrets)
	processor := NewProcessor(store)

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	result, err := processor.ProcessFile("")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	expected := map[string]string{
		"API_KEY": "secret123",
		"DEBUG":   "false",
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

func TestProcessor_ProcessFile_NotFound(t *testing.T) {
	store := newMockSecretStore(map[string]string{})
	processor := NewProcessor(store)

	_, err := processor.ProcessFile("nonexistent.template")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestProcessor_ProcessFileOptional(t *testing.T) {
	store := newMockSecretStore(map[string]string{})
	processor := NewProcessor(store)

	result, err := processor.ProcessFileOptional("nonexistent.template")
	if err != nil {
		t.Errorf("unexpected error for optional file processing: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected empty result for missing template file, got %d variables", len(result))
	}
}

func TestContext_Secret(t *testing.T) {
	secrets := map[string]string{"API_KEY": "secret123"}
	store := newMockSecretStore(secrets)
	ctx := NewContext(store)

	value, err := ctx.Secret("API_KEY")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if value != "secret123" {
		t.Errorf("expected 'secret123', got %q", value)
	}

	_, err = ctx.Secret("MISSING")
	if err == nil {
		t.Error("expected error for missing secret")
	}
}

func TestContext_Env(t *testing.T) {
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	ctx := NewContext(nil)
	value := ctx.Env("TEST_VAR")
	if value != "test_value" {
		t.Errorf("expected 'test_value', got %q", value)
	}

	empty := ctx.Env("NONEXISTENT")
	if empty != "" {
		t.Errorf("expected empty string for nonexistent env var, got %q", empty)
	}
}

func TestContext_Default(t *testing.T) {
	secrets := map[string]string{"EXISTING": "value"}
	store := newMockSecretStore(secrets)
	ctx := NewContext(store)

	value := ctx.Default("EXISTING", "fallback")
	if value != "value" {
		t.Errorf("expected 'value', got %q", value)
	}

	fallback := ctx.Default("MISSING", "fallback")
	if fallback != "fallback" {
		t.Errorf("expected 'fallback', got %q", fallback)
	}
}
