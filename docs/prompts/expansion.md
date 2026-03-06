# Expansion

This project is meant to be more than just a worktree manager:

- I want to use commands for it (please research what's the best practice for that in Go)
- The worktree manager TUI should move to the "gx worktrees" command (shorthand: "gx wt")
- I want to add a few more commands:
  - "gx clone-wt" - same as "git clone", but:
    - clones as a bare repo and creates a worktree for the primary branch (main/master)
  - "gx push" - pushes the current branch, if it can't push ask the user for
    confirmation to "push --force" (see ~/.dotfiles/core/git/scripts/git-psme)

Please plan how to implement this.
