package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gx/git"
	"gx/testutil"
)

func setupDotBareForCmd(t *testing.T) (outerDir string) {
	t.Helper()
	src := testutil.TempRepo(t)
	cwd := t.TempDir()

	raw, err := git.CloneBare(src, "", cwd)
	if err != nil {
		t.Fatalf("CloneBare: %v", err)
	}
	outerDir, err = filepath.EvalSymlinks(raw)
	if err != nil {
		t.Fatalf("EvalSymlinks: %v", err)
	}

	repo, err := git.FindRepo(outerDir)
	if err != nil {
		t.Fatalf("FindRepo: %v", err)
	}
	if err := git.UpdateRemotes(*repo); err != nil {
		t.Fatalf("UpdateRemotes: %v", err)
	}

	branch := repo.MainBranch
	wtPath := filepath.Join(outerDir, branch)
	if err := git.AddWorktreeFromRemote(*repo, wtPath, branch, "origin/"+branch); err != nil {
		t.Fatalf("AddWorktreeFromRemote: %v", err)
	}

	return outerDir
}

func TestDoctor_NoIssues(t *testing.T) {
	outerDir := setupDotBareForCmd(t)

	var stdout bytes.Buffer
	d := deps{
		stdout: &stdout,
		stderr: bytes.NewBuffer(nil),
		getwd:  func() (string, error) { return outerDir, nil },
	}

	if err := execute([]string{"doctor"}, d); err != nil {
		t.Fatalf("doctor: %v", err)
	}
	if !strings.Contains(stdout.String(), "No issues found") {
		t.Errorf("expected 'No issues found', got: %q", stdout.String())
	}
}

func TestDoctor_ReportsIssue(t *testing.T) {
	outerDir := setupDotBareForCmd(t)

	// Corrupt the outer .git file.
	gitFile := filepath.Join(outerDir, ".git")
	if err := os.WriteFile(gitFile, []byte("gitdir: wrong\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	var stdout bytes.Buffer
	d := deps{
		stdout: &stdout,
		stderr: bytes.NewBuffer(nil),
		getwd:  func() (string, error) { return outerDir, nil },
	}

	if err := execute([]string{"doctor"}, d); err != nil {
		t.Fatalf("doctor: %v", err)
	}
	out := stdout.String()
	if !strings.Contains(out, "gitdir: wrong") {
		t.Errorf("expected issue description in output, got: %q", out)
	}
	if !strings.Contains(out, "--fix") {
		t.Errorf("expected '--fix' hint in output, got: %q", out)
	}
}

func TestDoctor_Fix_AppliesWhenConfirmed(t *testing.T) {
	outerDir := setupDotBareForCmd(t)

	gitFile := filepath.Join(outerDir, ".git")
	if err := os.WriteFile(gitFile, []byte("gitdir: wrong\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	var stdout bytes.Buffer
	d := deps{
		stdout:       &stdout,
		stderr:       bytes.NewBuffer(nil),
		getwd:        func() (string, error) { return outerDir, nil },
		confirmForce: func(string) (bool, error) { return true, nil },
	}

	if err := execute([]string{"doctor", "--fix"}, d); err != nil {
		t.Fatalf("doctor --fix: %v", err)
	}

	data, err := os.ReadFile(gitFile)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(data) != "gitdir: ./.bare\n" {
		t.Errorf(".git content = %q, want %q", string(data), "gitdir: ./.bare\n")
	}
	if !strings.Contains(stdout.String(), "Fixed") {
		t.Errorf("expected 'Fixed' in output, got: %q", stdout.String())
	}
}

func TestDoctor_Fix_SkipsWhenDeclined(t *testing.T) {
	outerDir := setupDotBareForCmd(t)

	gitFile := filepath.Join(outerDir, ".git")
	original := []byte("gitdir: wrong\n")
	if err := os.WriteFile(gitFile, original, 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	var stdout bytes.Buffer
	d := deps{
		stdout:       &stdout,
		stderr:       bytes.NewBuffer(nil),
		getwd:        func() (string, error) { return outerDir, nil },
		confirmForce: func(string) (bool, error) { return false, nil },
	}

	if err := execute([]string{"doctor", "--fix"}, d); err != nil {
		t.Fatalf("doctor --fix: %v", err)
	}

	data, _ := os.ReadFile(gitFile)
	if string(data) != string(original) {
		t.Errorf(".git file was modified despite declining fix")
	}
	if !strings.Contains(stdout.String(), "Skipped") {
		t.Errorf("expected 'Skipped' in output, got: %q", stdout.String())
	}
}
