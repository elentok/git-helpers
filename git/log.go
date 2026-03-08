package git

import "strings"

// Commit is a single git commit with abbreviated hash and subject line.
type Commit struct {
	Hash    string
	Subject string
}

// CommitsSinceMain returns commits on branch that are not reachable from the
// repo's main branch, ordered newest first.
func CommitsSinceMain(repo Repo, branch string) ([]Commit, error) {
	return commitsBetween(repo, repo.MainBranch, branch)
}

// CommitsBehindMain returns commits on main that are not reachable from branch,
// ordered newest first.
func CommitsBehindMain(repo Repo, branch string) ([]Commit, error) {
	return commitsBetween(repo, branch, repo.MainBranch)
}

// CommitsSinceUpstream returns commits on branch not reachable from its upstream
// tracking branch, ordered newest first. Returns nil if no upstream is configured.
func CommitsSinceUpstream(repo Repo, branch string) ([]Commit, error) {
	upstream := UpstreamBranch(repo.Root, branch)
	if upstream == "" {
		return nil, nil
	}
	return commitsBetween(repo, upstream, branch)
}

// CommitsBehindUpstream returns commits on the upstream tracking branch not
// reachable from branch, ordered newest first. Returns nil if no upstream is
// configured.
func CommitsBehindUpstream(repo Repo, branch string) ([]Commit, error) {
	upstream := UpstreamBranch(repo.Root, branch)
	if upstream == "" {
		return nil, nil
	}
	return commitsBetween(repo, branch, upstream)
}

func commitsBetween(repo Repo, fromRef, toRef string) ([]Commit, error) {
	mergeBase, err := run(repo.Root, []string{"merge-base", fromRef, toRef})
	if err != nil {
		// No merge base (e.g. orphan branch) - return empty rather than error
		return nil, nil
	}

	out, err := run(repo.Root, []string{"log", "--pretty=format:%h\t%s", mergeBase + ".." + toRef})
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}

	var commits []Commit
	for _, line := range strings.Split(out, "\n") {
		hash, subject, ok := strings.Cut(line, "\t")
		if ok {
			commits = append(commits, Commit{Hash: hash, Subject: subject})
		}
	}
	return commits, nil
}
