# gx

An interactive TUI for managing git worktrees.

## Features

- Browse all linked worktrees in a table with sync status (ahead / behind / diverged)
- Sidebar showing commits ahead of main and uncommitted file changes
- Create, rename, clone, and delete worktrees interactively
- Yank files from one worktree and paste them into another
- Pull and push the selected worktree's branch
- Scrollable error modal for any git failures

## Requirements

- Go 1.21+
- Git

## Installation

```sh
# Install to $GOPATH/bin
make install

# Or build a local binary
make build
```

## Usage

Run from inside any git repository or bare repo:

```sh
gx
```

If launched from inside a worktree, the cursor starts on that worktree.

## Key bindings

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `d` | Delete selected worktree (and its branch) |
| `r` | Rename selected worktree and branch |
| `c` | Clone selected worktree (copies uncommitted files) |
| `y` | Yank files from selected worktree into clipboard |
| `p` | Paste yanked files into selected worktree |
| `l` | Pull selected worktree's branch |
| `s` | Push selected worktree's branch |
| `?` | Toggle full help |
| `q` / `Ctrl+C` | Quit |

## Development

```sh
make test   # run all tests
make run    # run without building
```
