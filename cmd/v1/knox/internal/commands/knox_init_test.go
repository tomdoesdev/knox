package commands

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tomdoesdev/knox/internal/v1/constants"
	"github.com/tomdoesdev/knox/internal/v1/project"
	errors2 "github.com/tomdoesdev/knox/pkg/errs"
	"github.com/urfave/cli/v3"
)

// containsIndentation checks if the JSON string contains proper indentation
func containsIndentation(jsonStr string) bool {
	// Check if the JSON contains indentation (spaces for formatting)
	lines := strings.Split(jsonStr, "\n")
	if len(lines) < 2 {
		return false
	}

	// Check if any line starts with spaces (indicating indentation)
	for _, line := range lines {
		if len(line) > 0 && line[0] == ' ' {
			return true
		}
	}

	return false
}

func TestNewInitCommand(t *testing.T) {
	cmd := NewInitCommand()

	if cmd.Name != "init" {
		t.Errorf("Expected command name 'init', got '%s'", cmd.Name)
	}

	if cmd.Usage != "initialise new knox project vault" {
		t.Errorf("Expected usage 'initialise new knox project vault', got '%s'", cmd.Usage)
	}

	if cmd.Action == nil {
		t.Error("Expected Action to be set")
	}
}

func TestInitCommand_Success(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func(dir string) {
		err := os.Chdir(dir)
		if err != nil {
			t.Fatal(err)
		}
	}(originalWd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create and run init command
	cmd := NewInitCommand()

	// Create a mock cli.Command context
	app := &cli.Command{}
	ctx := context.Background()

	err = cmd.Action(ctx, app)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify knox.json was created
	configPath := filepath.Join(tmpDir, constants.DefaultProjectFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("knox.json file was not created")
	}

	// Verify file content
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	var config project.Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		t.Fatal(err)
	}

	if config.ProjectID == "" {
		t.Error("ProjectID should not be empty")
	}

	if len(config.ProjectID) != 12 {
		t.Errorf("Expected ProjectID length 12, got %d", len(config.ProjectID))
	}
}

func TestInitCommand_ProjectExists(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func(dir string) {
		err := os.Chdir(dir)
		if err != nil {
			t.Fatal(err)
		}
	}(originalWd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create existing knox.json file
	configPath := filepath.Join(tmpDir, constants.DefaultProjectFileName)
	existingConfig := `{"project_id": "existing123"}`
	err = os.WriteFile(configPath, []byte(existingConfig), 0660)
	if err != nil {
		t.Fatal(err)
	}

	// Create and run init command
	cmd := NewInitCommand()

	// Create a mock cli.Command context
	app := &cli.Command{}
	ctx := context.Background()

	err = cmd.Action(ctx, app)
	if err == nil {
		t.Error("Expected error when project already exists")
	}

	// Verify error is a project exists error
	errors2.AssertErrorCode(t, err, errors2.ProjectExistsCode)
}

func TestInitialise_Success(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func(dir string) {
		err := os.Chdir(dir)
		if err != nil {
			t.Fatal(err)
		}
	}(originalWd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Call initialise function directly
	err = initialise()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify knox.json was created
	configPath := filepath.Join(tmpDir, constants.DefaultProjectFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("knox.json file was not created")
	}

	// Verify file permissions (umask may affect actual permissions)
	fileInfo, err := os.Stat(configPath)
	if err != nil {
		t.Fatal(err)
	}

	// Check that file is readable and writable by owner
	mode := fileInfo.Mode()
	if mode&0600 != 0600 {
		t.Errorf("File should be readable and writable by owner, got mode %v", mode)
	}
}

func TestInitialise_ProjectExists(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func(dir string) {
		err := os.Chdir(dir)
		if err != nil {
			t.Fatal(err)
		}
	}(originalWd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create existing knox.json file
	configPath := filepath.Join(tmpDir, constants.DefaultProjectFileName)
	existingConfig := `{"project_id": "existing123"}`
	err = os.WriteFile(configPath, []byte(existingConfig), 0660)
	if err != nil {
		t.Fatal(err)
	}

	// Call initialise function directly
	err = initialise()
	if err == nil {
		t.Error("Expected error when project already exists")
	}

	// Verify error is a project exists error
	errors2.AssertErrorCode(t, err, errors2.ProjectExistsCode)
}

func TestInitialise_InvalidDirectory(t *testing.T) {
	// Change to a non-existent directory (this should cause getwd to fail)
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func(dir string) {
		err := os.Chdir(dir)
		if err != nil {
			t.Fatal(err)
		}
	}(originalWd)

	// Create a directory and then remove it to simulate an invalid working directory
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Chdir(subDir)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(subDir)
	if err != nil {
		t.Fatal(err)
	}

	// Call initialise function - this should handle the error gracefully
	err = initialise()
	if err == nil {
		t.Error("Expected error when working directory is invalid")
	}
}

func TestInitialise_ValidateJSONFormat(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func(dir string) {
		err := os.Chdir(dir)
		if err != nil {
			t.Fatal(err)
		}
	}(originalWd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Call initialise function
	err = initialise()
	if err != nil {
		t.Fatal(err)
	}

	// Read and validate JSON format
	configPath := filepath.Join(tmpDir, constants.DefaultProjectFileName)
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	// Check that it's valid JSON
	var config map[string]interface{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		t.Fatalf("Generated file is not valid JSON: %v", err)
	}

	// Check that it has the expected structure
	if _, exists := config["project_id"]; !exists {
		t.Error("Generated JSON missing 'project_id' field")
	}

	// Check that the JSON is properly formatted (indented)
	// We'll just verify it's indented by checking for spaces
	if !containsIndentation(string(data)) {
		t.Error("Generated JSON is not properly indented")
	}
}
