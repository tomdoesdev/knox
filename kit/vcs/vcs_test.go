package vcs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectVCS_Git(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create a .git directory
	gitDir := filepath.Join(tmpDir, ".git")
	err := os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("failed to create .git directory: %v", err)
	}

	// Test detection
	vcs, err := DetectVCS(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if vcs == nil {
		t.Fatal("expected VCS info, got nil")
	}

	if vcs.Type != Git {
		t.Errorf("expected Git, got %s", vcs.Type)
	}

	if vcs.Root != tmpDir {
		t.Errorf("expected root %s, got %s", tmpDir, vcs.Root)
	}
}

func TestDetectVCS_Subdirectory(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()

	// Create .git directory at root
	gitDir := filepath.Join(tmpDir, ".git")
	err := os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("failed to create .git directory: %v", err)
	}

	// Create subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	// Test detection from subdirectory
	vcs, err := DetectVCS(subDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if vcs == nil {
		t.Fatal("expected VCS info, got nil")
	}

	if vcs.Type != Git {
		t.Errorf("expected Git, got %s", vcs.Type)
	}

	if vcs.Root != tmpDir {
		t.Errorf("expected root %s, got %s", tmpDir, vcs.Root)
	}
}

func TestDetectVCS_NoVCS(t *testing.T) {
	// Create a temporary directory without any VCS
	tmpDir := t.TempDir()

	// Test detection
	vcs, err := DetectVCS(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if vcs != nil {
		t.Errorf("expected nil, got %+v", vcs)
	}
}

func TestIsUnderGit(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Test without .git
	if IsUnderGit(tmpDir) {
		t.Error("expected false for directory without .git")
	}

	// Create .git directory
	gitDir := filepath.Join(tmpDir, ".git")
	err := os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("failed to create .git directory: %v", err)
	}

	// Test with .git
	if !IsUnderGit(tmpDir) {
		t.Error("expected true for directory with .git")
	}
}

func TestDetectVCS_FileInGitRepo(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create .git directory
	gitDir := filepath.Join(tmpDir, ".git")
	err := os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("failed to create .git directory: %v", err)
	}

	// Create a file in the repo
	filePath := filepath.Join(tmpDir, "test.txt")
	err = os.WriteFile(filePath, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Test detection of file in git repo
	vcs, err := DetectVCS(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if vcs == nil {
		t.Fatal("expected VCS info, got nil")
	}

	if vcs.Type != Git {
		t.Errorf("expected Git, got %s", vcs.Type)
	}

	if vcs.Root != tmpDir {
		t.Errorf("expected root %s, got %s", tmpDir, vcs.Root)
	}
}

func TestGetVCSIgnorePatterns(t *testing.T) {
	patterns := GetVCSIgnorePatterns(Git)
	expected := []string{".git/", ".gitignore", ".gitmodules", ".gitattributes"}

	if len(patterns) != len(expected) {
		t.Errorf("expected %d patterns, got %d", len(expected), len(patterns))
	}

	for i, pattern := range expected {
		if i >= len(patterns) || patterns[i] != pattern {
			t.Errorf("expected pattern %s at index %d", pattern, i)
		}
	}
}

func TestShouldIgnoreFile(t *testing.T) {
	tests := []struct {
		filePath string
		vcsType  VCSType
		expected bool
	}{
		{"/path/to/.git/config", Git, true},
		{"/path/to/.gitignore", Git, true},
		{"/path/to/normal.txt", Git, false},
		{"/path/to/.hg/config", Mercurial, true},
		{"/path/to/.svn/entries", Subversion, true},
		{"/path/to/.bzr/branch", Bazaar, true},
	}

	for _, tt := range tests {
		result := ShouldIgnoreFile(tt.filePath, tt.vcsType)
		if result != tt.expected {
			t.Errorf("ShouldIgnoreFile(%s, %s) = %v, expected %v",
				tt.filePath, tt.vcsType, result, tt.expected)
		}
	}
}
