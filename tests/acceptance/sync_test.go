package acceptance_test

import (
	"encoding/json"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/anthropics/claudit/tests/acceptance/testutil"
)

var _ = Describe("Sync Command", func() {
	var local, remote *testutil.GitRepo

	BeforeEach(func() {
		var err error
		local, remote, err = testutil.NewGitRepoWithRemote()
		Expect(err).NotTo(HaveOccurred())

		// Create initial commit and push
		Expect(local.WriteFile("README.md", "# Test")).To(Succeed())
		Expect(local.Commit("Initial commit")).To(Succeed())
		Expect(local.Run("git", "push", "-u", "origin", "master")).To(Succeed())
	})

	AfterEach(func() {
		if local != nil {
			local.Cleanup()
		}
		if remote != nil {
			remote.Cleanup()
		}
	})

	Describe("claudit sync push", func() {
		It("pushes notes to remote", func() {
			// Create a note on the commit
			head, err := local.GetHead()
			Expect(err).NotTo(HaveOccurred())

			// Store a conversation
			transcriptPath := filepath.Join(local.Path, "transcript.jsonl")
			Expect(os.WriteFile(transcriptPath, []byte(testutil.SampleTranscript()), 0644)).To(Succeed())

			hookInput := testutil.SampleHookInput("session-123", transcriptPath, "git commit -m 'test'")
			_, _, err = testutil.RunClauditInDirWithStdin(local.Path, hookInput, "store")
			Expect(err).NotTo(HaveOccurred())

			// Push notes
			stdout, _, err := testutil.RunClauditInDir(local.Path, "sync", "push")
			Expect(err).NotTo(HaveOccurred())
			Expect(stdout).To(ContainSubstring("Pushed"))

			// Verify remote has the notes ref
			output, err := remote.RunOutput("git", "notes", "--ref", "refs/notes/claude-conversations", "list")
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(ContainSubstring(head))
		})
	})

	Describe("claudit sync pull", func() {
		It("fetches notes from remote", func() {
			head, err := local.GetHead()
			Expect(err).NotTo(HaveOccurred())

			// Store and push a note from local
			transcriptPath := filepath.Join(local.Path, "transcript.jsonl")
			Expect(os.WriteFile(transcriptPath, []byte(testutil.SampleTranscript()), 0644)).To(Succeed())

			hookInput := testutil.SampleHookInput("session-123", transcriptPath, "git commit -m 'test'")
			_, _, err = testutil.RunClauditInDirWithStdin(local.Path, hookInput, "store")
			Expect(err).NotTo(HaveOccurred())
			_, _, err = testutil.RunClauditInDir(local.Path, "sync", "push")
			Expect(err).NotTo(HaveOccurred())

			// Create a clone without notes
			clone, err := testutil.NewGitRepo()
			Expect(err).NotTo(HaveOccurred())
			defer clone.Cleanup()

			Expect(clone.Run("git", "remote", "add", "origin", remote.Path)).To(Succeed())
			Expect(clone.Run("git", "fetch", "origin")).To(Succeed())
			Expect(clone.Run("git", "checkout", "-b", "master", "origin/master")).To(Succeed())

			// Clone should not have notes yet
			Expect(clone.HasNote("refs/notes/claude-conversations", head)).To(BeFalse())

			// Pull notes
			stdout, _, err := testutil.RunClauditInDir(clone.Path, "sync", "pull")
			Expect(err).NotTo(HaveOccurred())
			Expect(stdout).To(ContainSubstring("Fetched"))

			// Now clone should have the note
			Expect(clone.HasNote("refs/notes/claude-conversations", head)).To(BeTrue())
		})
	})

	Describe("notes round-trip", func() {
		It("preserves conversation through push/pull", func() {
			// Store conversation locally
			transcriptPath := filepath.Join(local.Path, "transcript.jsonl")
			Expect(os.WriteFile(transcriptPath, []byte(testutil.SampleTranscript()), 0644)).To(Succeed())

			hookInput := testutil.SampleHookInput("session-roundtrip", transcriptPath, "git commit -m 'test'")
			_, _, err := testutil.RunClauditInDirWithStdin(local.Path, hookInput, "store")
			Expect(err).NotTo(HaveOccurred())

			head, err := local.GetHead()
			Expect(err).NotTo(HaveOccurred())

			// Get original note content
			originalNote, err := local.GetNote("refs/notes/claude-conversations", head)
			Expect(err).NotTo(HaveOccurred())

			// Push to remote
			_, _, err = testutil.RunClauditInDir(local.Path, "sync", "push")
			Expect(err).NotTo(HaveOccurred())

			// Create clone and pull
			clone, err := testutil.NewGitRepo()
			Expect(err).NotTo(HaveOccurred())
			defer clone.Cleanup()

			Expect(clone.Run("git", "remote", "add", "origin", remote.Path)).To(Succeed())
			Expect(clone.Run("git", "fetch", "origin")).To(Succeed())
			Expect(clone.Run("git", "checkout", "-b", "master", "origin/master")).To(Succeed())

			_, _, err = testutil.RunClauditInDir(clone.Path, "sync", "pull")
			Expect(err).NotTo(HaveOccurred())

			// Compare notes
			clonedNote, err := clone.GetNote("refs/notes/claude-conversations", head)
			Expect(err).NotTo(HaveOccurred())

			var original, cloned map[string]interface{}
			Expect(json.Unmarshal([]byte(originalNote), &original)).To(Succeed())
			Expect(json.Unmarshal([]byte(clonedNote), &cloned)).To(Succeed())

			Expect(cloned["session_id"]).To(Equal(original["session_id"]))
			Expect(cloned["checksum"]).To(Equal(original["checksum"]))
			Expect(cloned["transcript"]).To(Equal(original["transcript"]))
		})
	})
})
