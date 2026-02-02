package git

import (
	"bytes"
	"os/exec"
	"strings"
)

// NotesRef is the git notes ref used for storing conversations
const NotesRef = "refs/notes/claude-conversations"

// AddNote adds a note to a commit
func AddNote(commitSHA string, content []byte) error {
	cmd := exec.Command("git", "notes", "--ref", NotesRef, "add", "-f", "-m", string(content), commitSHA)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// GetNote retrieves a note from a commit
func GetNote(commitSHA string) ([]byte, error) {
	cmd := exec.Command("git", "notes", "--ref", NotesRef, "show", commitSHA)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return output, nil
}

// HasNote checks if a commit has a conversation note
func HasNote(commitSHA string) bool {
	cmd := exec.Command("git", "notes", "--ref", NotesRef, "show", commitSHA)
	return cmd.Run() == nil
}

// ListCommitsWithNotes returns a list of commit SHAs that have conversation notes
func ListCommitsWithNotes() ([]string, error) {
	cmd := exec.Command("git", "notes", "--ref", NotesRef, "list")
	output, err := cmd.Output()
	if err != nil {
		// No notes exist yet - this is not an error
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil, nil
		}
		return nil, err
	}

	var commits []string
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		// Format: "note_sha commit_sha"
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			commits = append(commits, parts[1])
		}
	}

	return commits, nil
}

// PushNotes pushes notes to the remote
func PushNotes(remote string) error {
	// Use --no-verify to prevent pre-push hook from triggering recursively
	cmd := exec.Command("git", "push", "--no-verify", remote, NotesRef)
	return cmd.Run()
}

// FetchNotes fetches notes from the remote
func FetchNotes(remote string) error {
	cmd := exec.Command("git", "fetch", remote, NotesRef+":"+NotesRef)
	return cmd.Run()
}
