# Claudit

Save your Claude Code conversations with your commits. Resume them later. Share them with your team.

## Why?

When you're coding with Claude, the conversation is as valuable as the code. Claudit stores that conversation right in your Git repo, so you can:

- **Resume where you left off** — Pick up any past session, even on a different machine
- **Share context** — Your team sees not just the code, but the reasoning behind it
- **Browse history** — Web UI to explore conversations across your commit history

## Quick Start

```bash
# Install
go install github.com/anthropics/claudit@latest
# Or: make build

# Set up your repo (one time)
cd your-project
claudit init

# That's it. Now when you work with Claude and commit, conversations are saved automatically.
```

## Usage

**See what conversations you have:**
```bash
claudit list
```

**Resume a past session:**
```bash
claudit resume abc123    # By commit SHA
claudit resume HEAD~3    # By git ref
```

**Browse in your browser:**
```bash
claudit serve
```

**Sync with your team:**
```bash
git push   # Conversations sync automatically
git pull   # Fetches conversations from teammates
```

## How It Works

Claudit uses [Git Notes](https://git-scm.com/docs/git-notes) to attach conversations to commits. When you run `claudit init`, it sets up hooks so:

1. When Claude makes a commit, the conversation is saved automatically
2. When you make a commit during a Claude session, it's saved too
3. When you push/pull, conversations sync with the remote

No extra steps needed during your normal workflow.

## Commands

| Command | Description |
|---------|-------------|
| `claudit init` | Set up hooks in your repo |
| `claudit list` | Show commits with conversations |
| `claudit resume <commit>` | Resume a saved session |
| `claudit serve` | Open web UI |
| `claudit sync push/pull` | Manual sync (usually automatic) |

## Requirements

- Git
- Claude Code CLI (for resume)
- Go 1.21+ (for building from source)

## License

See LICENSE file.
