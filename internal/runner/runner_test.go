package runner

import (
	"testing"
)

func TestEnvRunner_BuildEnvironment_CleanEnvironment(t *testing.T) {
	config := Config{
		InheritEnv: false, // Clean environment (default)
	}
	runner := NewEnvRunner(config)

	envVars := EnvVars{
		"TEST_VAR": "test_value",
		"API_KEY":  "secret123",
	}

	env := runner.buildEnvironment(envVars)

	// Should only contain the template variables
	if len(env) != 2 {
		t.Errorf("expected 2 environment variables, got %d", len(env))
	}

	expectedVars := map[string]bool{
		"TEST_VAR=test_value": true,
		"API_KEY=secret123":   true,
	}

	for _, envVar := range env {
		if !expectedVars[envVar] {
			t.Errorf("unexpected environment variable: %s", envVar)
		}
	}
}

func TestEnvRunner_BuildEnvironment_InheritEnvironment(t *testing.T) {
	config := Config{
		InheritEnv:    true,
		AllowOverride: false,
	}
	runner := NewEnvRunner(config)

	envVars := EnvVars{
		"TEST_VAR": "test_value",
	}

	env := runner.buildEnvironment(envVars)

	// Should contain current environment + template variables
	// (exact count depends on current environment, but should be > 1)
	if len(env) <= 1 {
		t.Errorf("expected more than 1 environment variable when inheriting, got %d", len(env))
	}

	// Should contain our template variable
	found := false
	for _, envVar := range env {
		if envVar == "TEST_VAR=test_value" {
			found = true
			break
		}
	}
	if !found {
		t.Error("template variable TEST_VAR=test_value not found in environment")
	}
}
