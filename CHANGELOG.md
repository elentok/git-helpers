# Changelog

## v0.3.2

- Added `gx list-worktrees` command that prints all worktree names, one per line
- Added `gx worktree-abs-path <name>` command that prints the absolute path of the named worktree
- When pushing a branch for the first time, the GitHub PR creation URL is detected and a modal asks whether to open it in the browser (defaults to Yes)
- Fixed `run` to capture stderr even on success (needed for parsing remote push output)

## v0.3.1

- Rebinded pull to `p` and push to `P`, freeing up the old `l` / `s` keys
- After yanking files (pressing `y` and confirming), the app enters a dedicated paste mode where only navigation (`j`/`k`) and `p` to paste (or `esc` to cancel) are active — this is what freed `p` for pull in normal mode
- Refreshes the worktree list after a paste completes

## v0.3.0

- `gx clone-wt` now uses the `.bare` directory trick: clones into `my-repo/.bare/` and writes a `my-repo/.git` file pointing to it, so worktrees live cleanly alongside `.bare/` rather than inside it
- Delete worktree now shows a spinner while the deletion runs and a "Worktree {name} deleted successfully" toast on completion
- Added `gx doctor` command to check a repo for common configuration issues:
  - Verifies the origin fetch refspec is set correctly
  - For `.bare`-style repos: verifies the outer `.git` file points to `.bare`
  - For `.bare`-style repos: verifies each worktree's `.git` file points to the correct location
- Added `gx doctor --fix` to interactively apply fixes with confirmation prompts

## v0.2.1

- Added `U` keybinding to run `git remote update` and refresh all worktree statuses

## v0.2.0

- Added `gx version` command (also `--version`, `-v`) to print the current binary version
- Added `scripts/bump.sh` for bumping the version, creating an annotated git tag

## v0.1.5

- `gx clone-wt` now immediately fixes the fetch refspec after cloning, so remote tracking refs populate correctly on the first fetch
- On startup, the worktrees view checks whether the fetch refspec is misconfigured or remote tracking refs are missing, and offers to fix it automatically
- Delete and track confirmations are now shown as a centred modal with Yes/No buttons instead of a status-bar prompt
- Pull and push now also refresh the sidebar after completing
- Fixed a bug where the `origin/<branch>` fallback could match a bad local branch instead of the remote tracking ref

## v0.1.4

- Added `R` keybinding to refresh the worktree list and all statuses

## v0.1.3

- Added `t` keybinding to set a remote tracking branch for the selected worktree

## v0.1.2

- The sidebar now shows a "no remote tracking branch" note with a hint to press `t` when no upstream is configured

## v0.1.1

- Status column now shows ahead/behind relative to the remote tracking branch instead of the main branch
- Sidebar ahead/behind commit lists now compare against the remote tracking branch instead of main
- Sidebar section headings updated to "Commits ahead of remote" and "Commits behind remote"
