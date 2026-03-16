# Stashify

Add the "gx stashify" command, e.g.

```
gx stashify git rebase master
```

If the current repo (or worktree) is dirty it will stash the changes, run the
command and then pop the stages.

It will only pop automatically if the command was successful, otherwise it will
show a message, showing the error and confirm with the user if they want to pop.

This should also be usable from the worktrees UI for the following two scenarios:

- Pressing "p" to pull a dirty worktree:
  - Will ask the user if they wish to stash the changes first
  - If they do it will use stashify (but in an interactive way)
  - If they don't it should show "Pull aborted (dirty worktree)"

- Pressing "b" to will rebase the worktree on top of the main branch:
  - It will ask for confirmation first
  - If the worktree is dirty it will ask the user if they wish to stash
  - If they do it will use stashify
  - If not it will show "Rebase aborted (dirty worktree)"

## Plan

### Git Layer — `git/stash.go` (new)

Three functions following the pattern of `Pull` in `git/remote.go:88`:

- `Stash(dir string) (string, error)` — `git stash`
- `StashPop(dir string) (string, error)` — `git stash pop`
- `Rebase(dir string, onto string) (string, error)` — `git rebase <onto>`

### CLI — `cmd/stashify.go` (new)

`runStashify(args []string, d deps) error`:

1. Check dirty via `git.UncommittedChanges(cwd)`
2. If dirty: `git.Stash(cwd)` (with spinner)
3. Run the command via `exec.Command`
4. Success → auto-pop
5. Failure → print error, `d.confirmForce("Pop stash anyway?")` to prompt user

Register in `cmd/cmd.go` switch + usage.

### UI — Confirm Modal Cancel Message

Add `confirmCancelMsg string` field to Model (`model_state.go`).

In `model_confirm_modal.go`:

- `enterConfirmWithCancel(prompt, cmd, spinnerLabel, cancelMsg)` — sets cancel msg
- `enterConfirm` clears `confirmCancelMsg = ""`
- `handleConfirmKey` cancel branches: if `confirmCancelMsg != ""`, show it as status

### UI — Pull with Stash (`p` key)

Modify pull key handler (`model_update.go:89`):

- Check `m.dirties[wt.Path]` synchronously
- **Dirty**: `enterConfirmWithCancel("Stash changes before pulling?", cmdStashPull, spinner, "Pull aborted (dirty worktree)")`
- **Clean**: existing direct pull behavior

New types in `pullpush.go`:

- `stashPullResultMsg{err, log, stashed, wtPath}` — returned by `cmdStashPull`
- `stashPopResultMsg{err, wtPath, opLabel}` — returned by `cmdStashPop`

Handlers:

- `stashPullResultMsg`: success+stashed → `cmdStashPop` to auto-pop; failure+stashed → confirm "Pop stash?"
- `stashPopResultMsg`: success → "Pulled (stash restored)", reload dirty/sync; failure → showError

### UI — Rebase on Main (`b` key, new)

Add `Rebase` key binding (`b`) in `keys.go`.

Key handler guards: branch not empty, not main, MainBranch detected.

- **Clean**: `enterConfirm("Rebase X on main?", cmdRebase(repo, wt, false), "Rebasing…")`
- **Dirty**: `enterConfirm("Rebase X on main?", cmdRebasePreflight(repo, wt), "")` — two-step confirm

Two-step pattern:

- `cmdRebasePreflight(repo, wt)` returns `rebasePreflightMsg` immediately
- Handler shows second confirm: `enterConfirmWithCancel("Stash changes before rebasing?", cmdRebase(repo, wt, true), spinner, "Rebase aborted (dirty worktree)")`

New types:

- `rebasePreflightMsg{repo, wt}`
- `rebaseResultMsg{err, log, stashed, wtPath}`
- `cmdRebase(repo, wt, stash bool)` — optionally stash → rebase → return result

Handler: same pattern as stashPull — auto-pop on success, confirm pop on failure. Also reload base status for all worktrees.

### Files Changed

| File                                  | Change                                                    |
| ------------------------------------- | --------------------------------------------------------- |
| `git/stash.go`                        | New — `Stash`, `StashPop`, `Rebase`                       |
| `cmd/stashify.go`                     | New — `runStashify`                                       |
| `cmd/cmd.go`                          | Add case + usage                                          |
| `ui/worktrees/model_state.go`         | Add `confirmCancelMsg` field                              |
| `ui/worktrees/model_confirm_modal.go` | `enterConfirmWithCancel`, cancel-message handling          |
| `ui/worktrees/keys.go`               | Add `Rebase` key binding (`b`)                            |
| `ui/worktrees/pullpush.go`           | New msg types + commands for stash-aware pull/rebase       |
| `ui/worktrees/model_update.go`       | Update pull handler, add rebase handler, new msg handlers  |
