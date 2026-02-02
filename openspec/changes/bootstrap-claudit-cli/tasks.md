# Tasks: Bootstrap Claudit CLI

## Phase 1: Foundation

- [ ] **1.1** Initialize Go module and create project structure
  - Create `go.mod` with module path
  - Create directory structure: `cmd/`, `internal/claude/`, `internal/git/`, `internal/storage/`, `internal/web/`
  - Add dependencies:
    - `github.com/spf13/cobra` (CLI framework)
    - `github.com/onsi/ginkgo/v2` (testing framework)
    - `github.com/onsi/gomega` (assertions)
  - Verify: `go build` succeeds

- [ ] **1.2** Set up acceptance test infrastructure
  - Create `tests/acceptance/` directory structure
  - Create `acceptance_suite_test.go` with Ginkgo bootstrap
  - Create `testutil/binary.go` - build and execute claudit binary
  - Create `testutil/git.go` - create/manage temp git repos
  - Create `testutil/fixtures.go` - sample JSONL transcripts
  - Verify: `ginkgo tests/acceptance` runs (empty suite)

- [ ] **1.3** Implement root command with Cobra
  - Create `main.go` entry point
  - Create `cmd/root.go` with version flag
  - Verify: `claudit --version` displays version

- [ ] **1.4** Implement git repository detection
  - Create `internal/git/repo.go`
  - Detect if cwd is inside a git repository
  - Get repository root path
  - Verify: Unit tests pass

- [ ] **1.5** Add acceptance tests for CLI foundation
  - Test `claudit` displays help
  - Test `claudit --version` displays version
  - Verify: `ginkgo tests/acceptance` passes

## Phase 2: Conversation Storage

- [ ] **2.1** Implement JSONL transcript parser
  - Create `internal/claude/transcript.go`
  - Parse user, assistant, system, tool_result entries
  - Handle unknown types gracefully
  - Verify: Unit tests with sample JSONL

- [ ] **2.2** Implement compression and encoding
  - Create `internal/storage/compress.go`
  - Implement gzip compression
  - Implement base64 encoding
  - Implement SHA256 checksum
  - Verify: Round-trip compression test

- [ ] **2.3** Implement storage format
  - Create `internal/storage/format.go`
  - Define `StoredConversation` struct
  - JSON serialization/deserialization
  - Verify: Unit tests

- [ ] **2.4** Implement git notes operations
  - Create `internal/git/notes.go`
  - Add note to commit
  - Read note from commit
  - List commits with notes
  - Verify: Unit tests with temp git repo

- [ ] **2.5** Implement store command
  - Create `cmd/store.go`
  - Read PostToolUse hook JSON from stdin
  - Detect git commit commands
  - Read transcript, compress, store as note
  - Verify: Unit tests with mock stdin

- [ ] **2.6** Add acceptance tests for store command
  - Test: Pipe hook JSON with `git commit` command, verify note created
  - Test: Pipe hook JSON with non-commit command, verify silent exit
  - Test: Verify stored note contains expected metadata
  - Test: Verify transcript can be decompressed from note
  - Verify: `ginkgo tests/acceptance` passes

## Phase 3: Initialization

- [ ] **3.1** Implement Claude hooks configuration
  - Create `internal/claude/hooks.go`
  - Read/write `.claude/settings.local.json`
  - Merge hooks without overwriting existing config
  - Verify: Unit tests

- [ ] **3.2** Implement git hooks installation
  - Add to `internal/git/repo.go` or new file
  - Install pre-push → `claudit sync push`
  - Install post-merge, post-checkout → `claudit sync pull`
  - Handle existing hooks (append, don't overwrite)
  - Verify: Hooks are executable and call claudit commands

- [ ] **3.3** Implement init command
  - Create `cmd/init.go`
  - Configure Claude hooks
  - Install git hooks
  - Display success message
  - Verify: Command runs without error

- [ ] **3.4** Add acceptance tests for init command
  - Test: `claudit init` in git repo creates `.claude/settings.local.json`
  - Test: Verify PostToolUse hook configuration is correct
  - Test: Verify git hooks are installed and executable
  - Test: `claudit init` outside git repo fails with error
  - Test: `claudit init` preserves existing settings
  - Verify: `ginkgo tests/acceptance` passes

- [ ] **3.5** Implement sync command
  - Create `cmd/sync.go`
  - `sync push` - push notes to origin
  - `sync pull` - fetch notes from origin
  - Verify: Commands execute git push/fetch for notes ref

- [ ] **3.6** Add acceptance tests for sync and git hooks
  - Create `testutil/remote.go` - create local bare repos as remotes
  - Test: `claudit sync push` pushes notes to bare repo remote
  - Test: `claudit sync pull` fetches notes from bare repo remote
  - Test: Verify notes round-trip through push/pull
  - Test: Git hooks invoke `claudit sync` commands (not inline bash)
  - Verify: `ginkgo tests/acceptance` passes

## Phase 4: Session Resume

- [ ] **4.1** Implement session file management
  - Create `internal/claude/session.go`
  - Compute encoded project path
  - Write JSONL to Claude's location
  - Update sessions-index.json
  - Verify: Unit tests

- [ ] **4.2** Set up isolated Claude test environment
  - Create `testutil/claude_env.go`
  - Helper to create temp HOME directory for test isolation
  - Pre-populate `$HOME/.claude/projects/` structure for tests
  - Helper to run Claude CLI with isolated HOME
  - Skip helper when `CLAUDIT_SKIP_CLAUDE_TESTS=1` is set
  - Verify: Claude CLI runs in isolated environment without touching real config

- [ ] **4.3** Implement resume command
  - Create `cmd/resume.go`
  - Resolve commit reference
  - Read and decompress conversation
  - Verify checksum (warn on mismatch)
  - Restore session files
  - Check for uncommitted changes (prompt)
  - Checkout commit
  - Launch `claude --resume`
  - Verify: Command structure in place

- [ ] **4.4** Add acceptance tests for resume command
  - Test: `claudit resume <sha>` restores transcript to Claude location
  - Test: `claudit resume <sha>` calls `claude --resume <session-id>`
  - Test: `claudit resume` with branch name resolves correctly
  - Test: `claudit resume` with relative ref (HEAD~1) works
  - Test: `claudit resume` on commit without conversation shows error
  - Test: `claudit resume` warns about uncommitted changes
  - Verify: `ginkgo tests/acceptance` passes

- [ ] **4.5** Implement list command
  - Add `cmd/list.go`
  - List commits with conversations
  - Display SHA, date, message preview
  - Verify: Command runs

- [ ] **4.6** Add acceptance tests for list command
  - Test: `claudit list` shows commits with conversations
  - Test: `claudit list` in repo with no conversations shows empty
  - Test: Output format includes SHA, date, message
  - Verify: `ginkgo tests/acceptance` passes

## Phase 5: Web Visualization

- [ ] **5.1** Set up embedded static assets
  - Create `internal/web/static/` directory
  - Create basic HTML template
  - Configure Go embed
  - Verify: Assets compile into binary

- [ ] **5.2** Implement HTTP server foundation
  - Create `internal/web/server.go`
  - Localhost binding, configurable port
  - Serve embedded assets
  - Verify: Server starts and serves index page

- [ ] **5.3** Implement commits API
  - Create `internal/web/handlers.go`
  - `GET /api/commits` - list with pagination
  - `GET /api/commits/:sha` - full conversation
  - Verify: Handlers return expected JSON

- [ ] **5.4** Implement graph API
  - Add `GET /api/graph` endpoint
  - Return commit graph structure
  - Mark commits with conversations
  - Verify: Handler returns expected JSON

- [ ] **5.5** Implement resume API
  - Add `POST /api/resume/:sha` endpoint
  - Check for uncommitted changes
  - Trigger resume flow
  - Verify: Handler structure in place

- [ ] **5.6** Implement serve command
  - Create `cmd/serve.go`
  - Start server with port flag
  - Display URL in terminal
  - Verify: `claudit serve` starts server

- [ ] **5.7** Add acceptance tests for serve command and APIs
  - Test: `claudit serve` starts server on default port
  - Test: `claudit serve --port 3000` uses custom port
  - Test: `GET /api/commits` returns commit list with has_conversation flag
  - Test: `GET /api/commits/:sha` returns full conversation
  - Test: `GET /api/commits/:sha` returns 404 for commit without conversation
  - Test: `GET /api/graph` returns graph structure
  - Test: `POST /api/resume/:sha` triggers resume (with Claude mock)
  - Verify: `ginkgo tests/acceptance` passes

- [ ] **5.8** Build commit graph UI
  - Create SVG-based graph renderer
  - Display branches and commits
  - Highlight commits with conversations
  - Implement scroll/navigation
  - Verify: Visual inspection

- [ ] **5.9** Build conversation viewer UI
  - Create message display components
  - Style user vs assistant messages
  - Render markdown content
  - Collapsible tool uses
  - Verify: Visual inspection

- [ ] **5.10** Integrate resume button
  - Add "Resume Session" button
  - Call resume API
  - Display status feedback
  - Verify: Visual inspection

## Phase 6: Polish

- [ ] **6.1** Create Makefile
  - `make build` - build binary
  - `make test` - run unit tests
  - `make acceptance` - run acceptance tests
  - `make all` - build + all tests
  - Verify: All targets work

- [ ] **6.2** Add error handling and logging
  - Consistent error messages
  - Debug logging flag (`--debug`)
  - Verify: Errors are user-friendly

- [ ] **6.3** Write README
  - Installation instructions
  - Quick start guide
  - Command reference
  - Verify: Documentation is complete

- [ ] **6.4** Final acceptance test review
  - Ensure all scenarios from specs have corresponding tests
  - Add any missing edge case tests
  - Verify: Full test coverage of specified behavior

## Dependencies

```
1.1 ─► 1.2 ─► 1.3 ─► 1.4 ─► 1.5
              │
              └─► 2.1 ─┬─► 2.5 ─► 2.6
                  2.2 ─┤
                  2.3 ─┤
                  2.4 ─┘
                        │
              ┌─────────┘
              │
              ├─► 3.1 ─┬─► 3.3 ─► 3.4
              │   3.2 ─┘
              │         │
              │         └─► 3.5 ─► 3.6
              │
              └─► 4.1 ─► 4.2 ─► 4.3 ─► 4.4 ─► 4.5 ─► 4.6
                               │
              ┌────────────────┘
              │
              └─► 5.1 ─► 5.2 ─► 5.3 ─► 5.4 ─► 5.5 ─► 5.6 ─► 5.7
                                 │
                                 ├─► 5.8 ─┬─► 5.10
                                 └─► 5.9 ─┘

6.1, 6.2, 6.3, 6.4 can run after Phase 5
```

## External Dependencies

| Dependency | Required At | Purpose | Notes |
|------------|-------------|---------|-------|
| Git CLI | Runtime + Tests | All git operations | Must be installed on system |
| Go 1.21+ | Build | Compilation | For generics, embed improvements |
| Ginkgo CLI | Tests | Run acceptance tests | `go install github.com/onsi/ginkgo/v2/ginkgo@latest` |
| Claude CLI | Runtime + Tests | Resume sessions | Real CLI with HOME isolation; skip with `CLAUDIT_SKIP_CLAUDE_TESTS=1` |
| Anthropic API Key | Tests | Claude CLI auth | Required for resume integration tests |

### Network Requirements

- **Remote repos**: Tests use local bare git repos (no network)
- **Claude CLI**: Makes API calls to Anthropic (requires network + valid API key)
- **Web server**: Tests make HTTP requests to localhost

### Test Isolation

Tests override `HOME` environment variable to isolate Claude config:
```go
cmd.Env = append(os.Environ(), "HOME="+tempHome)
```
This prevents tests from polluting real `~/.claude/` directory.

## Parallelization Opportunities

- **Phase 2**: Tasks 2.1, 2.2, 2.3, 2.4 can be developed in parallel
- **Phase 3**: Tasks 3.1 and 3.2 can be developed in parallel
- **Phase 5**: Tasks 5.8 and 5.9 can be developed in parallel (UI components)
- **Phase 6**: All tasks can run in parallel
