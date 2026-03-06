package testutil

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// evalDir resolves symlinks in a directory path. On macOS, t.TempDir() returns
// /var/... which is a symlink to /private/var/..., while git resolves the real path.
func evalDir(t *testing.T, dir string) string {
	t.Helper()
	real, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatalf("EvalSymlinks(%s): %v", dir, err)
	}
	return real
}

// TempRepo creates a regular git repo with one initial commit on "main".
func TempRepo(t *testing.T) string {
	t.Helper()
	dir := evalDir(t, t.TempDir())
	mustGit(t, dir, "init", "--initial-branch=main")
	configUser(t, dir)
	WriteFile(t, dir, "README.md", "# test")
	mustGit(t, dir, "add", ".")
	mustGit(t, dir, "commit", "-m", "initial")
	return dir
}

// TempBareRepo creates a bare git repo by cloning a regular repo.
// The bare repo has one commit on "main".
func TempBareRepo(t *testing.T) string {
	t.Helper()
	src := TempRepo(t)
	bare := evalDir(t, t.TempDir())
	// Remove the empty TempDir so git clone can create it cleanly
	os.RemoveAll(bare)
	mustRun(t, ".", "git", "clone", "--bare", src, bare)
	return bare
}

// TempBareRepoWithWorktrees creates a bare repo with linked worktrees.
// Each name results in a branch and a worktree directory under the bare repo.
func TempBareRepoWithWorktrees(t *testing.T, names ...string) string {
	t.Helper()
	repoDir := TempBareRepo(t)
	for _, name := range names {
		wtDir := filepath.Join(repoDir, name)
		mustGit(t, repoDir, "worktree", "add", "-b", name, wtDir)
		configUser(t, wtDir)
		WriteFile(t, wtDir, "file.txt", name)
		mustGit(t, wtDir, "add", ".")
		mustGit(t, wtDir, "commit", "-m", "add "+name)
	}
	return repoDir
}

// WriteFile writes content to a file inside dir, creating it if needed.
func WriteFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile %s: %v", path, err)
	}
}

func configUser(t *testing.T, dir string) {
	t.Helper()
	mustGit(t, dir, "config", "user.email", "test@test.com")
	mustGit(t, dir, "config", "user.name", "Test")
}

// Mkdir creates a directory, failing the test if it can't.
func Mkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("Mkdir %s: %v", path, err)
	}
}

// CommitAll stages all changes in dir and creates a commit with the given message.
func CommitAll(t *testing.T, dir, message string) {
	t.Helper()
	mustGit(t, dir, "add", ".")
	mustGit(t, dir, "commit", "-m", message)
}

func mustGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	mustRun(t, dir, "git", args...)
}

func mustRun(t *testing.T, dir, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run %s %v in %s: %v\n%s", name, args, dir, err, out)
	}
}
