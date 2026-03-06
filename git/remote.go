package git

import "strings"

// ListRemotes returns the names of all configured remotes.
func ListRemotes(repo Repo) ([]string, error) {
	out, err := run(repo.Root, []string{"remote"})
	if err != nil {
		return nil, err
	}
	var remotes []string
	for _, r := range strings.Split(out, "\n") {
		if r != "" {
			remotes = append(remotes, r)
		}
	}
	return remotes, nil
}

// UpdateRemotes fetches updates from all remotes.
func UpdateRemotes(repo Repo) error {
	_, err := run(repo.Root, []string{"remote", "update"})
	return err
}

// PruneRemote removes remote-tracking references for deleted remote branches.
func PruneRemote(repo Repo, remote string) error {
	_, err := run(repo.Root, []string{"remote", "prune", remote})
	return err
}

// Pull fetches and integrates changes from the remote into the worktree.
func Pull(worktreePath string) error {
	_, err := run(worktreePath, []string{"pull"})
	return err
}

// BranchRemote returns the remote configured for branch (e.g. "origin"),
// falling back to "origin" if none is set.
func BranchRemote(repo Repo, branch string) string {
	remote := runAllowFail(repo.Root, []string{"config", "branch." + branch + ".remote"})
	if remote == "" {
		return "origin"
	}
	return remote
}

// Push uploads local branch commits to the remote using an explicit
// "git push <remote> <branch>" invocation.
func Push(worktreePath, remote, branch string) error {
	_, err := run(worktreePath, []string{"push", remote, branch})
	return err
}

// PushBranch pushes branch to remote.
func PushBranch(worktreePath, remote, branch string) error {
	return Push(worktreePath, remote, branch)
}

// PushBranchForceWithLease force-pushes branch using --force-with-lease.
func PushBranchForceWithLease(worktreePath, remote, branch string) error {
	_, err := run(worktreePath, []string{"push", "--force-with-lease", remote, branch})
	return err
}

// IsNonFastForwardPushError returns true when err matches a rejected push that
// can be resolved with a force push.
func IsNonFastForwardPushError(err error) bool {
	runErr, ok := err.(*RunError)
	if !ok {
		return false
	}

	s := strings.ToLower(runErr.Stdout + "\n" + runErr.Stderr)
	return strings.Contains(s, "non-fast-forward") ||
		strings.Contains(s, "[rejected]") ||
		strings.Contains(s, "fetch first") ||
		strings.Contains(s, "failed to push some refs")
}

// PruneAllRemotes prunes all configured remotes.
func PruneAllRemotes(repo Repo) error {
	remotes, err := ListRemotes(repo)
	if err != nil {
		return err
	}
	for _, remote := range remotes {
		if err := PruneRemote(repo, remote); err != nil {
			return err
		}
	}
	return nil
}
