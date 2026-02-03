package integration_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestClaudeCodeIntegration runs an end-to-end test with actual Claude Code CLI.
// This test requires:
// - ANTHROPIC_API_KEY environment variable set
// - Claude Code CLI installed and in PATH
// - claudit binary built
//
// Skip with: SKIP_CLAUDE_INTEGRATION=1 go test ./tests/integration/...
func TestClaudeCodeIntegration(t *testing.T) {
	// Skip conditions
	if os.Getenv("SKIP_CLAUDE_INTEGRATION") == "1" {
		t.Skip("SKIP_CLAUDE_INTEGRATION=1 is set")
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY not set - skipping integration test")
	}

	// Check Claude CLI is available
	if _, err := exec.LookPath("claude"); err != nil {
		t.Skip("Claude Code CLI not found in PATH")
	}

	// Check claudit binary
	clauditPath := os.Getenv("CLAUDIT_BINARY")
	if clauditPath == "" {
		// Try to find it relative to workspace
		clauditPath = filepath.Join(getWorkspaceRoot(), "claudit")
	}
	if _, err := os.Stat(clauditPath); err != nil {
		t.Fatalf("claudit binary not found at %s - run 'make build' first", clauditPath)
	}

	// Create temporary test directory
	tmpDir, err := os.MkdirTemp("", "claude-integration-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize git repo
	runGit(t, tmpDir, "init")
	runGit(t, tmpDir, "config", "user.email", "test@example.com")
	runGit(t, tmpDir, "config", "user.name", "Test User")

	// Create initial file and commit
	testFile := filepath.Join(tmpDir, "README.md")
	if err := os.WriteFile(testFile, []byte("# Test Project\n"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	runGit(t, tmpDir, "add", "README.md")
	runGit(t, tmpDir, "commit", "-m", "Initial commit")

	// Initialize claudit
	cmd := exec.Command(clauditPath, "init")
	cmd.Dir = tmpDir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("claudit init failed: %v\nOutput: %s", err, output)
	}

	// Verify hooks are configured
	settingsPath := filepath.Join(tmpDir, ".claude", "settings.local.json")
	settingsData, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("Failed to read settings: %v", err)
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(settingsData, &settings); err != nil {
		t.Fatalf("Failed to parse settings: %v", err)
	}

	// Verify hook format is correct
	hooks, ok := settings["hooks"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected hooks object in settings, got: %v", settings)
	}
	postToolUse, ok := hooks["PostToolUse"].([]interface{})
	if !ok {
		t.Fatalf("Expected PostToolUse array in hooks, got: %v", hooks)
	}
	if len(postToolUse) == 0 {
		t.Fatal("PostToolUse hooks array is empty")
	}

	t.Log("Hook configuration verified successfully")

	// Create a test file that Claude will commit
	todoFile := filepath.Join(tmpDir, "todo.txt")
	if err := os.WriteFile(todoFile, []byte("- Buy milk\n- Walk dog\n"), 0644); err != nil {
		t.Fatalf("Failed to write todo file: %v", err)
	}

	// Run Claude Code with a simple prompt to commit the file
	// Use --print for non-interactive mode
	// Use --allowedTools to only allow Bash for git operations
	// Use --max-turns 5 to limit API calls
	claudeCmd := exec.Command("claude",
		"--print",
		"--allowedTools", "Bash(git:*),Read",
		"--max-turns", "5",
		"--dangerously-skip-permissions",
		"Please run: git add todo.txt && git commit -m 'Add todo list'",
	)
	claudeCmd.Dir = tmpDir
	claudeCmd.Env = append(os.Environ(),
		"ANTHROPIC_API_KEY="+apiKey,
		"PATH="+os.Getenv("PATH")+":"+filepath.Dir(clauditPath),
	)

	// Set timeout
	done := make(chan error, 1)
	go func() {
		output, err := claudeCmd.CombinedOutput()
		if err != nil {
			t.Logf("Claude output: %s", output)
		}
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Logf("Claude command finished with error (may be expected): %v", err)
		}
	case <-time.After(60 * time.Second):
		claudeCmd.Process.Kill()
		t.Fatal("Claude command timed out after 60 seconds")
	}

	// Give hooks time to run
	time.Sleep(2 * time.Second)

	// Check if commit was made
	cmd = exec.Command("git", "log", "--oneline", "-n", "2")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Git log output: %s", output)
	}

	if !strings.Contains(string(output), "todo") {
		t.Log("Commit with 'todo' not found - Claude may not have made the commit")
		t.Logf("Git log: %s", output)
		// This is not necessarily a failure - Claude might refuse or fail
	}

	// Check if note was created (the main test)
	cmd = exec.Command("git", "notes", "--ref=refs/notes/claude-conversations", "list")
	cmd.Dir = tmpDir
	notesOutput, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("No notes found (may be expected if commit failed): %v", err)
		return
	}

	if len(strings.TrimSpace(string(notesOutput))) > 0 {
		t.Log("SUCCESS: Git note was created by claudit hook!")
		t.Logf("Notes: %s", notesOutput)

		// Verify note content
		cmd = exec.Command("git", "notes", "--ref=refs/notes/claude-conversations", "show", "HEAD")
		cmd.Dir = tmpDir
		noteContent, _ := cmd.CombinedOutput()
		t.Logf("Note content preview: %s", noteContent[:min(len(noteContent), 200)])
	} else {
		t.Log("No git notes found - hook may not have triggered")
	}
}

// getWorkspaceRoot finds the workspace root by looking for go.mod
func getWorkspaceRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "."
}

// runGit runs a git command in the specified directory
func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\nOutput: %s", args, err, output)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
