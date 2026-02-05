# claudit

Automatically save Claude Code conversations as Git Notes on every commit.

```bash
claudit init
```

Now work with `claude` as you would normally. Whenever you or Claude Code commit, the conversation since the last commit will be attached to that commit as a Git Note.

## Why?

In order to understand _how_ and _why_ a commit was made, we need to see the conversation that led to it.

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

**View in your browser:**

```bash
claudit serve
```

## How It Works

Claudit uses [Git Notes](https://git-scm.com/docs/git-notes) to attach conversations to commits. When you run `claudit init`, it sets up hooks so:

1. When Claude makes a commit, the conversation is saved automatically
2. When you make a commit during a Claude session, it's saved too
3. When you push/pull, conversations sync with the remote

No extra steps needed during your normal workflow.

## Commands

| Command                   | Description                     |
| ------------------------- | ------------------------------- |
| `claudit init`            | Set up hooks in your repo       |
| `claudit list`            | Show commits with conversations |
| `claudit resume <commit>` | Resume a saved session          |
| `claudit serve`           | Open web UI                     |
| `claudit sync push/pull`  | Manual sync (usually automatic) |

## Requirements

- Git
- Claude Code CLI (for resume)
- Go 1.21+ (for building from source)

## License

See LICENSE file.
