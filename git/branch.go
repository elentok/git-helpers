package git

import "strings"

// Branch represents a local or remote git branch.
type Branch struct {
	Name       string
	GitName    string // as passed to git commands (e.g. "origin/main" for remote)
	RemoteName string // non-empty for remote branches
	IsRemote   bool
}

// ListBranches returns all local and remote branches, excluding HEAD pointers.
func ListBranches(repo Repo) ([]Branch, error) {
	out, err := run(repo.Root, []string{"branch", "--all"})
	if err != nil {
		return nil, err
	}
	var branches []Branch
	for _, line := range strings.Split(out, "\n") {
		if line == "" || strings.Contains(line, "/HEAD ") || strings.Contains(line, " -> ") {
			continue
		}
		branches = append(branches, parseBranchLine(line))
	}
	return branches, nil
}

func parseBranchLine(line string) Branch {
	// Strip leading *, +, spaces (current/worktree markers)
	line = strings.TrimLeft(line, " *+")
	line = strings.TrimSpace(line)

	if strings.HasPrefix(line, "remotes/") {
		// "remotes/origin/branchname"
		parts := strings.SplitN(strings.TrimPrefix(line, "remotes/"), "/", 2)
		if len(parts) == 2 {
			return Branch{
				Name:       parts[1],
				GitName:    parts[0] + "/" + parts[1],
				RemoteName: parts[0],
				IsRemote:   true,
			}
		}
	}

	return Branch{Name: line, GitName: line}
}

// CurrentBranch returns the short name of the current branch in dir.
func CurrentBranch(dir string) (string, error) {
	return run(dir, []string{"rev-parse", "--abbrev-ref", "HEAD"})
}

// TrackRemote configures branch to track <remote>/<branch>.
func TrackRemote(repoRoot, remote, branch string) error {
	_, err := run(repoRoot, []string{"branch", "--track", remote + "/" + branch})
	return err
}

// UpstreamBranch returns the upstream tracking ref (e.g. "origin/feature") for
// a local branch. It first checks the configured tracking branch, then falls
// back to "origin/<branch>" if that ref exists. Returns "" if neither applies.
func UpstreamBranch(repoRoot, branch string) string {
	if up := runAllowFail(repoRoot, []string{
		"for-each-ref",
		"--format=%(upstream:short)",
		"refs/heads/" + branch,
	}); up != "" {
		return up
	}
	// Fall back to the conventional origin/<branch> ref.
	candidate := "origin/" + branch
	if runAllowFail(repoRoot, []string{"rev-parse", "--verify", candidate}) != "" {
		return candidate
	}
	return ""
}

// CreateBranch creates and checks out a new branch.
func CreateBranch(repo Repo, name string) error {
	_, err := run(repo.Root, []string{"checkout", "-b", name})
	return err
}

// DeleteLocalBranch deletes a local branch. Use force=true to delete unmerged branches.
func DeleteLocalBranch(repo Repo, name string, force bool) error {
	flag := "-d"
	if force {
		flag = "-D"
	}
	_, err := run(repo.Root, []string{"branch", flag, name})
	return err
}

// DeleteRemoteBranch deletes a branch on the remote.
func DeleteRemoteBranch(repo Repo, remoteName, branchName string) error {
	_, err := run(repo.Root, []string{"push", "--delete", remoteName, branchName})
	return err
}

// RenameBranch renames a local branch from oldName to newName.
func RenameBranch(repo Repo, oldName, newName string) error {
	_, err := run(repo.Root, []string{"branch", "-m", oldName, newName})
	return err
}
