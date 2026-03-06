# Go Migration Plan

Migrate git-helpers from a Deno/TypeScript CLI to a Go TUI using BubbleTea. The first feature
is a Worktrees page that replaces the current `status` command with an interactive, always-visible
worktree dashboard.

## Architecture Overview

```
gx/
├── main.go                  # Entry point, initialize app
├── go.mod
├── git/                     # Git operations (pure library, no TUI dependency)
│   ├── repo.go              # Repo discovery, bare repo detection
│   ├── worktree.go          # Worktree list/add/remove/move
│   ├── branch.go            # Branch create/delete/rename, local & remote
│   ├── remote.go            # Remote list/update/prune
│   ├── status.go            # Sync status (ahead/behind/diverged), uncommitted changes
│   ├── log.go               # Commit log retrieval
│   ├── run.go               # Low-level git command executor
│   ├── repo_test.go
│   ├── worktree_test.go
│   ├── branch_test.go
│   └── status_test.go
├── ui/                      # TUI layer (BubbleTea)
│   ├── app.go               # Root model, page routing, global keybindings
│   ├── worktrees/            # Worktrees page
│   │   ├── model.go         # Page model, Update, View
│   │   ├── table.go         # Worktree table component
│   │   ├── sidebar.go       # Sidebar: commits + changed files
│   │   ├── keys.go          # Key map definitions
│   │   ├── delete.go        # Delete confirmation dialog
│   │   ├── rename.go        # Rename text input dialog
│   │   ├── clone.go         # Clone dialog (name input)
│   │   ├── yank.go          # File selection (checkboxes) for yank
│   │   └── paste.go         # Paste action
│   ├── styles.go            # Shared lipgloss styles
│   └── components/          # Reusable components
│       ├── confirm.go       # Yes/no confirmation dialog
│       ├── textinput.go     # Single-line text input
│       └── checklist.go     # Multi-select checkbox list
└── testutil/                # Test helpers
    ├── repo.go              # Create temp bare repos with worktrees
    └── commit.go            # Create dummy commits
```

### Key design decisions

- **`git/` has no TUI imports.** It returns Go structs; the UI layer formats them. This makes
  the git layer independently testable.
- **Each interaction (delete, rename, clone, yank) is a sub-model** within `ui/worktrees/`. The
  page model delegates to the active sub-model when one is open.
- **Clipboard is page-level state**, not global. It holds a list of file paths and source worktree.
  The sidebar or a status bar shows "N files in clipboard" when non-empty.

---

## Milestone 1: Project scaffold + git library

Set up the Go module, implement the git operations layer, and verify with tests. No TUI yet.

### Steps

1. `go mod init`, add dependencies (`bubbletea`, `lipgloss`, `bubbles`).
2. Implement `git/run.go` - execute git commands, capture stdout/stderr, return structured errors.
3. Implement `git/repo.go`:
   - `FindRepo(path) (Repo, error)` - walk up to find `.git` or bare repo.
   - `IdentifyDir(path) (DirInfo, error)` - detect bare repo root vs worktree root.
   - Repo struct: `Root string`, `IsBare bool`, `MainBranch string` (detect main vs master).
4. Implement `git/worktree.go`:
   - `ListWorktrees(repo) ([]Worktree, error)` - parse `git worktree list --porcelain`.
   - Worktree struct: `Path string`, `Name string`, `Branch string`, `Head string`.
5. Implement `git/branch.go`:
   - `DeleteLocalBranch`, `DeleteRemoteBranch`, `RenameBranch`, `CreateBranch`.
6. Implement `git/remote.go`:
   - `UpdateRemotes`, `PruneRemote`.
7. Implement `git/status.go`:
   - `SyncStatus(repo, branch) (Status, error)` - ahead/behind/diverged relative to main branch.
   - `UncommittedChanges(worktreePath) ([]Change, error)` - parse `git status --porcelain`.
8. Implement `git/log.go`:
   - `CommitsSinceMain(repo, branch) ([]Commit, error)` - log from HEAD to main branch merge-base.
9. Write unit tests for all of the above using temp repos (`testutil/`).

### Done when

- `go test ./git/...` passes.
- Can programmatically list worktrees, get sync status, and get commit logs from a real bare repo.

---

## Milestone 2: Basic TUI - worktree table

Render the worktrees page with the table, navigation, and sidebar. Read-only, no interactions yet.

### Steps

1. Implement `ui/app.go`:
   - Root BubbleTea model that holds the current page.
   - On init: detect repo from cwd, load worktrees, determine active worktree.
   - Global quit: `q`, `ctrl+c`.
2. Implement `ui/worktrees/table.go`:
   - Table with columns: Name, Branch, Sync Status.
   - Sync status shows: "synced", "N behind", "N ahead", "diverged".
   - Active row = worktree matching cwd (if launched from inside a worktree).
   - Navigate with `j`/`k` or up/down arrows.
3. Implement `ui/worktrees/sidebar.go`:
   - When a worktree is selected, show:
     - Commits since main branch (abbreviated hash + message).
     - Uncommitted changes (modified/added/deleted files).
   - Sidebar updates on row change.
4. Implement `ui/worktrees/model.go`:
   - Compose table + sidebar side by side.
   - Handle window resize with `tea.WindowSizeMsg`.
5. Implement `ui/styles.go` - lipgloss styles for table, sidebar, status bar.
6. Implement `main.go`:
   - Parse cwd, find repo, launch BubbleTea program with worktrees page.
   - Wire up `bin/gx` to run the Go binary (or just `go run .` during dev).

### Done when

- Running `gx` from a bare repo or worktree shows the table with real data.
- Can navigate rows with j/k, sidebar updates.
- Active worktree is highlighted when launched from inside one.

---

## Milestone 3: Delete and rename

Implement the `d` and `r` interactions.

### Steps

1. Implement `ui/components/confirm.go` - reusable yes/no dialog.
2. Implement `ui/components/textinput.go` - single-line input with default value.
3. Implement `ui/worktrees/delete.go`:
   - `d` on a worktree opens confirmation: "Delete worktree X and branch Y? [y/N]".
   - On confirm: call `git/worktree.Remove` + `git/branch.DeleteLocalBranch` +
     `git/branch.DeleteRemoteBranch`.
   - Refresh table after deletion.
   - Show error inline if deletion fails.
4. Implement `ui/worktrees/rename.go`:
   - `r` opens text input pre-filled with current name.
   - On submit: call `git worktree move`, rename directory, rename branch.
   - Port the rename logic from the existing `rename-worktree.ts` (gitdir/`.git` file fixups).
   - Refresh table after rename.
5. Implement `ui/worktrees/keys.go` - centralized key bindings using `bubbles/key`.

### Done when

- Can delete a worktree with `d`, confirm, see it removed from the table.
- Can rename a worktree with `r`, type new name, see it updated.

---

## Milestone 4: Clone worktree

Implement the `c` interaction.

### Steps

1. Implement `ui/worktrees/clone.go`:
   - `c` opens text input pre-filled with current worktree name.
   - On submit: create new worktree as a copy.
   - Clone strategy: `git worktree add` with new branch, then `cp -r` working tree files
     (including untracked) from source to destination.
   - Refresh table after clone.

### Done when

- Can clone a worktree with `c`, specify name, see new worktree in table.
- Cloned worktree has all files including untracked ones from the source.

---

## Milestone 5: Yank and paste

Implement the `y` (yank files) and `p` (paste files) interactions.

### Steps

1. Implement `ui/components/checklist.go` - multi-select checkbox list with toggle and toggle-all.
2. Implement `ui/worktrees/yank.go`:
   - `y` shows a checklist of uncommitted + untracked files in the selected worktree.
   - All items checked by default.
   - Navigate with j/k, toggle with space, confirm with enter.
   - Store selected file paths + source worktree path in page-level clipboard state.
3. Update `ui/worktrees/model.go`:
   - Show clipboard indicator in status bar: "N files in clipboard" (or empty).
4. Implement `ui/worktrees/paste.go`:
   - `p` copies files from clipboard source to current worktree destination.
   - Preserve relative paths. Create directories as needed.
   - Clear clipboard after paste.
   - Show success/error message.

### Done when

- Can yank files from one worktree, navigate to another, paste them.
- Clipboard indicator shows file count, clears after paste.

---

## Milestone 6: Git pull and push

Implement `gpl` and `gps` chained key interactions.

### Steps

1. Add chained key support to the key handling in `ui/worktrees/model.go`:
   - Track a key buffer with a short timeout (e.g. 500ms).
   - `g` starts a chain, `gp` continues, `gpl` triggers pull, `gps` triggers push.
2. Implement pull/push as async commands (`tea.Cmd`) that run `git pull`/`git push` in the
   selected worktree directory.
3. Show a spinner or status message while the operation is running.
4. Refresh sync status after completion.

### Done when

- `gpl` runs git pull in the selected worktree, updates sync status.
- `gps` runs git push in the selected worktree, updates sync status.
- Status bar shows progress and result.

---

## Milestone 7: Polish and ship

### Steps

1. Error handling pass: ensure all git errors surface as user-visible messages, not panics.
2. Handle edge cases: empty repo, no worktrees, detached HEAD, worktree with no branch.
3. Add a help bar at the bottom showing available keybindings.
4. Update `bin/gx` to build and run the Go binary.
5. Add a `Makefile` or `go install` instructions.
6. Update `README.md`.
