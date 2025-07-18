package vcs

import (
	"os"
	"path/filepath"
	"strings"
)

// VCSType represents different version control systems
type VCSType string

const (
	Git        VCSType = "git"
	Mercurial  VCSType = "hg"
	Subversion VCSType = "svn"
	Bazaar     VCSType = "bzr"
)

// VCSInfo contains information about detected version control
type VCSInfo struct {
	Type VCSType
	Root string
}

// DetectVCS detects if a path is under version control
func DetectVCS(path string) (*VCSInfo, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// Check if path is a file, use its directory instead
	if info, err := os.Stat(absPath); err == nil && !info.IsDir() {
		absPath = filepath.Dir(absPath)
	}

	// Walk up the directory tree looking for VCS directories
	currentPath := absPath
	for {
		// Check for Git
		if gitDir := filepath.Join(currentPath, ".git"); isVCSDir(gitDir) {
			return &VCSInfo{Type: Git, Root: currentPath}, nil
		}

		// Check for Mercurial
		if hgDir := filepath.Join(currentPath, ".hg"); isVCSDir(hgDir) {
			return &VCSInfo{Type: Mercurial, Root: currentPath}, nil
		}

		// Check for Subversion
		if svnDir := filepath.Join(currentPath, ".svn"); isVCSDir(svnDir) {
			return &VCSInfo{Type: Subversion, Root: currentPath}, nil
		}

		// Check for Bazaar
		if bzrDir := filepath.Join(currentPath, ".bzr"); isVCSDir(bzrDir) {
			return &VCSInfo{Type: Bazaar, Root: currentPath}, nil
		}

		// Move up one directory
		parent := filepath.Dir(currentPath)
		if parent == currentPath {
			// Reached root directory
			break
		}
		currentPath = parent
	}

	return nil, nil // No VCS found
}

// IsUnderVCS checks if a path is under version control
func IsUnderVCS(path string) bool {
	vcs, err := DetectVCS(path)
	return err == nil && vcs != nil
}

// IsUnderGit checks if a path is under Git version control
func IsUnderGit(path string) bool {
	vcs, err := DetectVCS(path)
	return err == nil && vcs != nil && vcs.Type == Git
}

// GetVCSIgnorePatterns returns common patterns that should be ignored for each VCS
func GetVCSIgnorePatterns(vcsType VCSType) []string {
	switch vcsType {
	case Git:
		return []string{".git/", ".gitignore", ".gitmodules", ".gitattributes"}
	case Mercurial:
		return []string{".hg/", ".hgignore", ".hgtags"}
	case Subversion:
		return []string{".svn/"}
	case Bazaar:
		return []string{".bzr/", ".bzrignore"}
	default:
		return []string{}
	}
}

// ShouldIgnoreFile checks if a file should be ignored based on VCS patterns
func ShouldIgnoreFile(filePath string, vcsType VCSType) bool {
	patterns := GetVCSIgnorePatterns(vcsType)
	for _, pattern := range patterns {
		if strings.Contains(filePath, pattern) {
			return true
		}
	}
	return false
}

// isVCSDir checks if a directory exists and is likely a VCS directory
func isVCSDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
