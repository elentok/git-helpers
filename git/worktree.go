package git

import (
	"strings"
)

// Worktree represents a git worktree.
type Worktree struct {
	Path       string
	Name       string // path relative to the repo root
	Branch     string // short branch name, empty if detached
	Head       string // commit hash
	IsDetached bool
	IsBare     bool
}

// ListWorktrees returns all linked worktrees for the repo (excludes the bare root).
func ListWorktrees(repo Repo) ([]Worktree, error) {
	out, err := run(repo.Root, []string{"worktree", "list", "--porcelain"})
	if err != nil {
		return nil, err
	}
	return parseWorktreePorcelain(out, repo.Root), nil
}

func parseWorktreePorcelain(out, repoRoot string) []Worktree {
	var worktrees []Worktree
	var cur *Worktree

	flush := func() {
		if cur != nil && !cur.IsBare {
			worktrees = append(worktrees, *cur)
		}
		cur = nil
	}

	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			flush()
			continue
		}
		if after, ok := strings.CutPrefix(line, "worktree "); ok {
			flush()
			name := strings.TrimPrefix(after, repoRoot+"/")
			cur = &Worktree{Path: after, Name: name}
		} else if cur == nil {
			continue
		} else if after, ok := strings.CutPrefix(line, "HEAD "); ok {
			cur.Head = after
		} else if after, ok := strings.CutPrefix(line, "branch "); ok {
			cur.Branch = strings.TrimPrefix(after, "refs/heads/")
		} else if line == "detached" {
			cur.IsDetached = true
		} else if line == "bare" {
			cur.IsBare = true
		}
	}
	flush()

	return worktrees
}

// RemoveWorktree removes a worktree by name or path.
func RemoveWorktree(repo Repo, name string, force bool) error {
	args := []string{"worktree", "remove"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, name)
	_, err := run(repo.Root, args)
	return err
}

// MoveWorktree moves a worktree from one path to another.
func MoveWorktree(repo Repo, from, to string) error {
	_, err := run(repo.Root, []string{"worktree", "move", from, to})
	return err
}
