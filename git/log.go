package git

import (
	"strings"
	"time"
)

// Commit is a single git commit with abbreviated hash and subject line.
type Commit struct {
	Hash    string
	Subject string
	Date    time.Time
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

// HeadCommit returns the latest commit on the given branch.
func HeadCommit(repoRoot, branch string) (Commit, error) {
	out, _, err := run(repoRoot, []string{"log", "-1", "--pretty=format:%h\t%ci\t%s", branch})
	if err != nil || out == "" {
		return Commit{}, err
	}
	hash, rest, _ := strings.Cut(out, "\t")
	dateStr, subject, _ := strings.Cut(rest, "\t")
	date, _ := time.Parse("2006-01-02 15:04:05 -0700", dateStr)
	return Commit{Hash: hash, Subject: subject, Date: date}, nil
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
