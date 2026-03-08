# gx

A collection of git helper (worktree management, etc...)

## Features

- Browse all linked worktrees in a table with sync status (ahead / behind / diverged)
- Sidebar showing commits ahead of main and uncommitted file changes
- Create, rename, clone, and delete worktrees interactively
- Yank files from one worktree and paste them into another
- Pull and push the selected worktree's branch
- Scrollable error modal for any git failures
- See [Changelog](./CHANGELOG.md)

![screenshot](./docs/screenshot.png)

## Requirements

- Go 1.21+
- Git

## Installation

Using homebrew:

```sh
brew tap elentok/stuff
brew install --cask gx
```

Using `go install`:

```sh
go install github.com/elentok/gx@latest
```

```sh
make install
```

## Usage

Run from inside any git repository or bare repo:

```sh
gx
```

If launched from inside a worktree, the cursor starts on that worktree.

You can also run the TUI explicitly:

```sh
gx worktrees
gx wt
```

Clone as bare + bootstrap initial worktree:

```sh
gx clone-wt <repo-url> [directory]
```

Push current worktree branch, with styled force-with-lease confirmation on rejection:

```sh
gx push
```

Create an initial config file with defaults:

```sh
gx init
```

Edit config in `$EDITOR`:

```sh
gx edit-config
```

## Configuration

Optional config file:

```sh
~/.config/gx/config.json
```

Example:

```json
{
  "use-nerdfont-icons": true
}
```

## Key bindings

| Key            | Action                                             |
| -------------- | -------------------------------------------------- |
| `j` / `↓`      | Move down                                          |
| `k` / `↑`      | Move up                                            |
| `d`            | Delete selected worktree (and its branch)          |
| `r`            | Rename selected worktree and branch                |
| `c`            | Clone selected worktree (copies uncommitted files) |
| `y`            | Yank files from selected worktree into clipboard   |
| `p`            | Paste yanked files into selected worktree          |
| `l`            | Pull selected worktree's branch                    |
| `s`            | Push selected worktree's branch                    |
| `?`            | Toggle full help                                   |
| `q` / `Ctrl+C` | Quit                                               |

## Development

```sh
make test   # run all tests
make run    # run without building
```
