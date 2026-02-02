package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// IsInsideWorkTree returns true if the current directory is inside a git repository
func IsInsideWorkTree() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}

// GetRepoRoot returns the root directory of the git repository
func GetRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetCurrentBranch returns the name of the current branch
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetHeadCommit returns the SHA of HEAD
func GetHeadCommit() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// EnsureGitDir returns the path to the .git directory, handling worktrees
func EnsureGitDir() (string, error) {
	root, err := GetRepoRoot()
	if err != nil {
		return "", err
	}

	gitDir := filepath.Join(root, ".git")
	info, err := os.Stat(gitDir)
	if err != nil {
		return "", err
	}

	// If .git is a file, this is a worktree - read the actual git dir path
	if !info.IsDir() {
		content, err := os.ReadFile(gitDir)
		if err != nil {
			return "", err
		}
		// Format: "gitdir: /path/to/actual/.git/worktrees/name"
		parts := strings.SplitN(string(content), ": ", 2)
		if len(parts) == 2 {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return gitDir, nil
}
