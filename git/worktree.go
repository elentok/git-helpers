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
	return parseWorktreePorcelain(out, repo.LinkedWorktreeDir()), nil
}

func parseWorktreePorcelain(out, worktreeDir string) []Worktree {
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
			name := strings.TrimPrefix(after, worktreeDir+"/")
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

// AddWorktree creates a new linked worktree at newPath with a new branch newName,
// starting at fromRef (branch name, tag, or commit hash). fromRef may be empty,
// in which case git uses the current HEAD of the repo.
func AddWorktree(repo Repo, newName, newPath, fromRef string) error {
	args := []string{"worktree", "add", "-b", newName, newPath}
	if fromRef != "" {
		args = append(args, fromRef)
	}
	_, err := run(repo.Root, args)
	return err
}
