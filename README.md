# Claudit

Store, resume, and visualize Claude Code conversations as Git Notes.

## What is Claudit?

Claudit is a CLI tool that captures your Claude Code conversation history and stores it directly in your Git repository using Git Notes. This allows teams to:

- **Preserve AI context** - Keep the reasoning behind code changes alongside the commits
- **Share knowledge** - Team members can see the AI-assisted development process
- **Resume sessions** - Pick up where you left off on any commit
- **Visualize history** - Browse conversations in a web interface

## Status: All Milestones Complete

Claudit is feature-complete with full conversation storage, session resume, and web visualization.

## Commands

### `claudit init` - Set up repository
```bash
cd your-project
claudit init                                    # Interactive ref selection
claudit init --notes-ref=refs/notes/commits     # Use default ref (recommended)
claudit init --notes-ref=refs/notes/claude-conversations  # Use custom ref
```

Configures:
- Git notes ref selection (default or custom namespace)
- Git settings for notes visibility (displayRef, rewriteRef)
- Claude Code hooks:
  - PostToolUse hook to capture conversations when Claude commits
  - SessionStart/SessionEnd hooks to track active sessions
- Git hooks:
  - pre-push, post-merge, post-checkout for automatic note syncing
  - post-commit for capturing manual commits during active sessions

**Notes ref options:**
- `refs/notes/commits` (default) - Standard git notes ref, works with `git notes show HEAD` and `git log`
- `refs/notes/claude-conversations` (custom) - Separate namespace, requires `--ref` flag for manual git commands

### `claudit store` - Capture conversation (automatic)

Called automatically by Claude Code hook when you commit. Compresses and stores the conversation as a Git Note.

**Two capture modes:**
- **PostToolUse hook** - Triggers when Claude Code runs `git commit`
- **post-commit hook** - Captures conversations for manual commits made during active sessions

The dual-hook approach ensures conversations are captured whether Claude makes the commit or you do.

```bash
claudit store           # Called automatically by PostToolUse hook
claudit store --manual  # Called by post-commit git hook for manual commits
```

**Manual commit capture** works by tracking the active session:
- When a Claude session starts, `claudit session-start` records the session state
- Manual commits during the session trigger `claudit store --manual`
- Both hooks can fire for the same commit - duplicates are detected and skipped (idempotent)

### `claudit list` - Show commits with conversations
```bash
claudit list
```

Lists all commits that have stored conversations, showing:
- Commit SHA
- Date
- Message preview
- Number of messages

### `claudit resume <commit>` - Resume a session
```bash
claudit resume abc123       # Resume from short SHA
claudit resume HEAD~2       # Resume from relative ref
claudit resume feature-branch  # Resume from branch name
claudit resume abc123 --force  # Skip uncommitted changes warning
```

Restores the conversation and launches Claude Code to continue the session.

### `claudit serve` - Web visualization
```bash
claudit serve              # Start on port 8080, open browser
claudit serve --port 3000  # Custom port
claudit serve --no-browser # Don't open browser automatically
```

Opens a web interface showing:
- Commit list with conversation indicators
- Full conversation viewer with collapsible tool uses
- Resume button to continue any session

### `claudit sync` - Manual sync
```bash
claudit sync push              # Push notes to origin
claudit sync pull              # Fetch notes from origin
claudit sync push --remote upstream  # Push to different remote
```

## How It Works

1. **Hook Integration** - Conversations are captured automatically:
   - When Claude Code runs `git commit`, the PostToolUse hook triggers `claudit store`
   - For manual commits during a Claude session, the post-commit git hook triggers `claudit store --manual`
2. **Session Tracking** - SessionStart/SessionEnd hooks track the active Claude session in `.claudit/active-session.json`
3. **Compression** - The conversation transcript is gzip compressed and base64 encoded
4. **Storage** - Stored as a Git Note on your configured ref (default: `refs/notes/commits`)
5. **Duplicate Prevention** - If both hooks fire for the same commit, the second one detects the duplicate and skips
6. **Sync** - Git hooks automatically sync notes when you push/pull
7. **Resume** - Restores session files to Claude's expected location and launches with `--resume`

### Storage Format

Each note contains:
```json
{
  "version": 1,
  "session_id": "uuid",
  "timestamp": "2024-01-15T10:30:00Z",
  "project_path": "/path/to/repo",
  "git_branch": "feature-branch",
  "message_count": 42,
  "checksum": "sha256:abc123...",
  "transcript": "<compressed conversation>"
}
```

## Installation

```bash
# Build from source
make build

# The binary is created at ./claudit
```

### Requirements

- Git CLI
- Go 1.21+ (for building)
- Claude CLI (for resume functionality)

## Quick Start

```bash
# 1. Initialize in your project
cd your-project
claudit init

# 2. Start a Claude Code session and make commits
# Conversations are captured automatically

# 3. Push to share (notes sync automatically with git push)
git push

# 4. View your conversation history
claudit list

# 5. Browse in web UI
claudit serve

# 6. Resume any session
claudit resume <commit>
```

## Development

```bash
make build       # Build binary
make test        # Run unit tests
make acceptance  # Run acceptance tests
make all         # Run all tests
```

## Test Suite

42 acceptance tests covering:
- CLI foundation (help, version)
- Init command (hooks setup)
- Store command (note creation, compression)
- Sync command (push/pull, git hooks)
- Resume command (session restore, checkout)
- List command (conversation listing)
- Serve command (web server basics)

## Milestones

- **Milestone 1** ✅ - Store conversations as Git Notes
- **Milestone 2** ✅ - Resume sessions from any commit
- **Milestone 3** ✅ - Web visualization of commit graph with conversations

## License

See LICENSE file for details.
