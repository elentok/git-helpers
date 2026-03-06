# Convert to go

I wrote this project as a way to write git helpers that are too complicated for bash, but as I
started using it I realized I want something a bit different.

I want to convert this codebase to Go (with the BubbleTea framework) and change from just a CLI to a
TUI.

## Worktrees page

The first feature I want to implement in the Go version is a "Worktrees page" that will show the
list of worktrees in the current repository.

It will show a table of worktrees with the following columns:

1. name
2. attached branch
3. sync status - is it synced with master? if not, is it behind master or after master? (the purpose
   of this is to let me see quickly if I forgot to push or if there are new updates on the remote).

You need to be able to run the new "gx" command from both the bare repo root or a worktree root. If
your run from inside the worktree it will be the active row in the table.

One worktree in the table can be "active", you move between them with the up/down arrow keys or j/k.

### Sidebar

The sidebar shows a preview of the active worktree:

1. list of commits until master/main (whichever is the main branch)
2. list of untracked files and uncommited changes

### Interactions

- d - delete the worktree, including the directory and the branch it's attached to (asks for
  confirmation before doing anything)
- r - rename the worktree and the branch it's attached to
- c - clone the worktree
  - ask the user for a new name, with the old name as the default
  - clone to the new name, keeping everything including untracked files
- y - copy untracked files and uncommited changes to another worktree
  - should show a list of uncommited/untracked files with checkboxes (all checked by default) and
    let the user pick which ones to copy
  - after pressing enter in the selection dialog the files are added to a "clipboard" (virtual),
    there should be some indication on the screen that there are X files in the clipboard
  - the user can then navigate to another worktree and press "p" to paste the files in the clipboard
  - the clipboard empties after pasting
- gpl - run "git pull" inside the worktree
- gps - run "git push" inside the worktree

## Guidelines

- The TUI should use vim bindings:
  - HJKL for movement in addition to regular arrows
  - single key or chained keys to perform actions (e.g. "d" to delete a worktree)
- The current codebase is written the Deno/Javascript-way, the new codebase should be written the
  Go-way.
- Write unit tests

## Instructions

1. Plan the implementation, split it into milestones (you can split those to steps if you prefer)
2. Write the plan to docs/go-migration-plan.md
