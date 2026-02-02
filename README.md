# Claudit

Store and resume Claude Code conversations as Git Notes.

## What is Claudit?

Claudit is a CLI tool that captures your Claude Code conversation history and stores it directly in your Git repository using Git Notes. This allows teams to:

- **Preserve AI context** - Keep the reasoning behind code changes alongside the commits
- **Share knowledge** - Team members can see the AI-assisted development process
- **Resume sessions** - Pick up where you left off on any commit (coming soon)

## Current Status: Milestone 1 Complete

Claudit can now capture and store conversations. The following features are implemented:

### Commands

**`claudit init`** - Set up a repository for conversation capture
```bash
cd your-project
claudit init
```

This configures:
- Claude Code's PostToolUse hook to capture conversations on commit
- Git hooks (pre-push, post-merge, post-checkout) for automatic note syncing

**`claudit store`** - Capture a conversation (called automatically by Claude Code hook)

When you make a git commit during a Claude Code session, the conversation is automatically compressed and stored as a Git Note on that commit.

**`claudit sync`** - Manually sync notes with remote
```bash
claudit sync push              # Push notes to origin
claudit sync pull              # Fetch notes from origin
claudit sync push --remote upstream  # Push to different remote
```

## How It Works

1. **Hook Integration** - When Claude Code runs a git commit, the PostToolUse hook triggers `claudit store`
2. **Compression** - The conversation transcript is gzip compressed and base64 encoded
3. **Storage** - Stored as a Git Note in `refs/notes/claude-conversations` namespace
4. **Sync** - Git hooks automatically sync notes when you push/pull

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

## Quick Start

```bash
# 1. Initialize in your project
cd your-project
claudit init

# 2. Start a Claude Code session and make commits
# Conversations are captured automatically

# 3. Push to share (notes sync automatically with git push)
git push
```

## Development

```bash
make build       # Build binary
make test        # Run unit tests
make acceptance  # Run acceptance tests
make all         # Run all tests
```

## Roadmap

- **Milestone 1** (Complete) - Store conversations as Git Notes
- **Milestone 2** - Resume sessions from any commit
- **Milestone 3** - Web visualization of commit graph with conversations

## License

See LICENSE file for details.
