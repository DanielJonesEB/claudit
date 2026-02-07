package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/DanielJonesEB/claudit/internal/claude"
	"github.com/DanielJonesEB/claudit/internal/cli"
	"github.com/DanielJonesEB/claudit/internal/git"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize claudit in the current repository",
	Long: `Configures the current git repository for conversation capture.

This command:
- Uses refs/notes/claude-conversations for note storage
- Creates/updates .claude/settings.local.json with PostToolUse hook
- Installs git hooks for automatic note syncing
- Configures git settings for notes visibility`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	// Verify we're in a git repository
	if err := git.RequireGitRepo(); err != nil {
		return err
	}

	repoRoot, err := git.GetRepoRoot()
	if err != nil {
		return fmt.Errorf("failed to get repository root: %w", err)
	}

	// Configure git settings for notes visibility
	cli.LogDebug("init: configuring git settings for notes ref %s", git.NotesRef)
	if err := configureGitSettings(git.NotesRef); err != nil {
		return fmt.Errorf("failed to configure git settings: %w", err)
	}

	fmt.Printf("✓ Configured notes ref: %s\n", git.NotesRef)
	fmt.Println("✓ Configured git notes settings (displayRef, rewriteRef)")

	// Configure Claude hooks
	cli.LogDebug("init: configuring Claude hooks")
	claudeDir := filepath.Join(repoRoot, ".claude")
	settings, err := claude.ReadSettings(claudeDir)
	if err != nil {
		return fmt.Errorf("failed to read Claude settings: %w", err)
	}

	claude.AddClauditHook(settings)
	claude.AddSessionHooks(settings)

	if err := claude.WriteSettings(claudeDir, settings); err != nil {
		return fmt.Errorf("failed to write Claude settings: %w", err)
	}

	fmt.Println("✓ Configured Claude hooks (PostToolUse, SessionStart, SessionEnd)")

	// Install git hooks
	cli.LogDebug("init: installing git hooks")
	gitDir, err := git.EnsureGitDir()
	if err != nil {
		return fmt.Errorf("failed to find git directory: %w", err)
	}

	if err := git.InstallAllHooks(gitDir); err != nil {
		return fmt.Errorf("failed to install git hooks: %w", err)
	}

	fmt.Println("✓ Installed git hooks (pre-push, post-merge, post-checkout, post-commit)")

	// Check if claudit is in PATH
	if _, err := exec.LookPath("claudit"); err != nil {
		fmt.Println()
		fmt.Println("⚠ Warning: 'claudit' is not in your PATH.")
		fmt.Println("  The hook will not work until claudit is installed.")
		fmt.Println("  Install with: go install github.com/DanielJonesEB/claudit@latest")
	}

	fmt.Println()
	fmt.Println("Claudit is now configured! Conversations will be stored")
	fmt.Printf("as git notes on %s when commits are made via Claude Code.\n", git.NotesRef)

	return nil
}

// configureGitSettings configures git settings for notes visibility
func configureGitSettings(notesRef string) error {
	// Configure notes.displayRef so git log shows notes
	cmd := exec.Command("git", "config", "notes.displayRef", notesRef)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set notes.displayRef: %w", err)
	}

	// Configure notes.rewriteRef so notes follow commits during rebase
	cmd = exec.Command("git", "config", "notes.rewriteRef", notesRef)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set notes.rewriteRef: %w", err)
	}

	return nil
}
