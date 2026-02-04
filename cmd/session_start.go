package cmd

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/DanielJonesEB/claudit/internal/session"
	"github.com/spf13/cobra"
)

// SessionStartInput represents the SessionStart hook JSON input from Claude Code
type SessionStartInput struct {
	SessionID      string `json:"session_id"`
	TranscriptPath string `json:"transcript_path"`
	Cwd            string `json:"cwd"`
}

var sessionStartCmd = &cobra.Command{
	Use:   "session-start",
	Short: "Handle Claude Code SessionStart hook",
	Long: `Reads SessionStart hook JSON from stdin and records the active session.

This command is designed to be called by Claude Code's SessionStart hook.`,
	RunE: runSessionStart,
}

func init() {
	rootCmd.AddCommand(sessionStartCmd)
}

func runSessionStart(cmd *cobra.Command, args []string) error {
	// Read hook input from stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		logWarning("failed to read stdin: %v", err)
		return nil // Exit silently to not disrupt workflow
	}

	var hook SessionStartInput
	if err := json.Unmarshal(input, &hook); err != nil {
		logWarning("failed to parse hook JSON: %v", err)
		return nil // Exit silently
	}

	// Validate required fields
	if hook.SessionID == "" || hook.TranscriptPath == "" {
		logWarning("missing required fields in hook data")
		return nil
	}

	// Create active session record
	activeSession := &session.ActiveSession{
		SessionID:      hook.SessionID,
		TranscriptPath: hook.TranscriptPath,
		StartedAt:      time.Now().UTC().Format(time.RFC3339),
		ProjectPath:    hook.Cwd,
	}

	// Write active session file
	if err := session.WriteActiveSession(activeSession); err != nil {
		// Log but don't fail - don't disrupt user's workflow
		logWarning("failed to write active session: %v", err)
		return nil
	}

	logInfo("session started: %s", hook.SessionID[:8])
	return nil
}
