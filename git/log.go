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
	mergeBase, err := run(repo.Root, []string{"merge-base", branch, repo.MainBranch})
	if err != nil {
		// No merge base (e.g. orphan branch) - return empty rather than error
		return nil, nil
	}

	out, err := run(repo.Root, []string{
		"log",
		"--pretty=format:%h\t%s",
		mergeBase + ".." + branch,
	})
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
