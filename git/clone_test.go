package git_test

import (
	"path/filepath"
	"testing"

	"gx/git"
	"gx/testutil"
)

func TestCloneBareAndBootstrapWorktree(t *testing.T) {
	src := testutil.TempRepo(t)
	cwd := t.TempDir()

	repoRoot, err := git.CloneBare(src, "", cwd)
	if err != nil {
		t.Fatalf("CloneBare: %v", err)
	}

	repo, err := git.FindRepo(repoRoot)
	if err != nil {
		t.Fatalf("FindRepo: %v", err)
	}
	if !repo.IsBare {
		t.Fatal("cloned repo is not bare")
	}

	mainBranch := git.RemoteDefaultBranch(repoRoot)
	if mainBranch != "main" && mainBranch != "master" {
		t.Fatalf("unexpected default branch: %q", mainBranch)
	}

	wtPath := filepath.Join(repoRoot, mainBranch)
	if err := git.AddWorktreeFromRemote(*repo, wtPath, mainBranch, "origin/"+mainBranch); err != nil {
		t.Fatalf("AddWorktreeFromRemote: %v", err)
	}

	branch, err := git.CurrentBranch(wtPath)
	if err != nil {
		t.Fatalf("CurrentBranch: %v", err)
	}
	if branch != mainBranch {
		t.Fatalf("CurrentBranch = %q, want %q", branch, mainBranch)
	}
}
