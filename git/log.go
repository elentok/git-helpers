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

// CommitsBetween returns commits reachable from toRef but not fromRef, ordered
// newest first.
func CommitsBetween(repo Repo, fromRef, toRef string) ([]Commit, error) {
	return commitsBetween(repo, fromRef, toRef)
}

func commitsBetween(repo Repo, fromRef, toRef string) ([]Commit, error) {
	mergeBase, _, err := run(repo.Root, []string{"merge-base", fromRef, toRef})
	if err != nil {
		// No merge base (e.g. orphan branch) - return empty rather than error
		return nil, nil
	}

	out, _, err := run(repo.Root, []string{"log", "--pretty=format:%h\t%s", mergeBase + ".." + toRef})
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
