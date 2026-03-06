package git

import (
	"os"
	"path/filepath"
	"strings"
)

// CloneBare clones repoURL as a bare repository. If targetDir is empty,
// git's default naming is mirrored from the repository URL.
func CloneBare(repoURL, targetDir, cwd string) (string, error) {
	before := map[string]struct{}{}
	if targetDir == "" {
		entries, _ := os.ReadDir(cwd)
		for _, e := range entries {
			before[e.Name()] = struct{}{}
		}
	}

	args := []string{"clone", "--bare", repoURL}
	if targetDir != "" {
		args = append(args, targetDir)
	}
	if _, err := run(cwd, args); err != nil {
		return "", err
	}

	name := targetDir
	if name == "" {
		entries, _ := os.ReadDir(cwd)
		for _, e := range entries {
			if _, ok := before[e.Name()]; !ok {
				name = e.Name()
				break
			}
		}
		if name == "" {
			name = inferCloneDirFromURL(repoURL)
		}
	}
	if filepath.IsAbs(name) {
		return name, nil
	}
	return filepath.Join(cwd, name), nil
}

func inferCloneDirFromURL(repoURL string) string {
	s := strings.TrimSuffix(repoURL, "/")
	i := strings.LastIndexAny(s, "/:")
	name := s
	if i >= 0 && i+1 < len(s) {
		name = s[i+1:]
	}
	return strings.TrimSuffix(name, ".git")
}

// AddWorktreeFromRemote adds a worktree for branch. It first tries to check out
// an existing local branch, then falls back to creating a tracking branch from
// remoteBranch (for freshly cloned bare repositories).
func AddWorktreeFromRemote(repo Repo, worktreePath, branch, remoteBranch string) error {
	if _, err := run(repo.Root, []string{"worktree", "add", worktreePath, branch}); err == nil {
		return nil
	}
	_, err := run(repo.Root, []string{"worktree", "add", "--track", "-b", branch, worktreePath, remoteBranch})
	return err
}
