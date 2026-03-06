package git_test

import (
	"path/filepath"
	"testing"

	"gx/git"
	"gx/testutil"
)

func TestCommitsSinceMain_hasCommits(t *testing.T) {
	repoDir := testutil.TempBareRepoWithWorktrees(t, "feature")
	repo, _ := git.FindRepo(repoDir)

	commits, err := git.CommitsSinceMain(*repo, "feature")
	if err != nil {
		t.Fatalf("CommitsSinceMain: %v", err)
	}
	if len(commits) != 1 {
		t.Fatalf("got %d commits, want 1", len(commits))
	}
	if commits[0].Subject != "add feature" {
		t.Errorf("Subject = %q, want %q", commits[0].Subject, "add feature")
	}
	if commits[0].Hash == "" {
		t.Error("Hash is empty")
	}
}

func TestCommitsSinceMain_noCommits(t *testing.T) {
	repoDir := testutil.TempBareRepo(t)
	repo, _ := git.FindRepo(repoDir)

	commits, err := git.CommitsSinceMain(*repo, "main")
	if err != nil {
		t.Fatalf("CommitsSinceMain: %v", err)
	}
	if len(commits) != 0 {
		t.Errorf("got %d commits for main vs main, want 0", len(commits))
	}
}

func TestCommitsSinceMain_multipleCommits(t *testing.T) {
	repoDir := testutil.TempBareRepoWithWorktrees(t, "feature")
	wtDir := filepath.Join(repoDir, "feature")
	repo, _ := git.FindRepo(repoDir)

	// Add a second commit to the feature branch
	testutil.WriteFile(t, wtDir, "extra.txt", "extra")
	testutil.CommitAll(t, wtDir, "second commit")

	commits, err := git.CommitsSinceMain(*repo, "feature")
	if err != nil {
		t.Fatalf("CommitsSinceMain: %v", err)
	}
	if len(commits) != 2 {
		t.Fatalf("got %d commits, want 2", len(commits))
	}
}
