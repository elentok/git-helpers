package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"gx/testutil"
)

func TestExecute_DefaultRunsWorktrees(t *testing.T) {
	called := 0
	d := deps{
		stdout: bytes.NewBuffer(nil),
		stderr: bytes.NewBuffer(nil),
		runWorktrees: func(_ string) error {
			called++
			return nil
		},
	}

	if err := execute(nil, d); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if called != 1 {
		t.Fatalf("runWorktrees called %d times, want 1", called)
	}
}

func TestExecute_WorktreesAliases(t *testing.T) {
	for _, arg := range []string{"worktrees", "wt"} {
		t.Run(arg, func(t *testing.T) {
			called := 0
			d := deps{
				stdout: bytes.NewBuffer(nil),
				stderr: bytes.NewBuffer(nil),
				runWorktrees: func(_ string) error {
					called++
					return nil
				},
			}
			if err := execute([]string{arg}, d); err != nil {
				t.Fatalf("execute: %v", err)
			}
			if called != 1 {
				t.Fatalf("runWorktrees called %d times, want 1", called)
			}
		})
	}
}

func TestExecute_UnknownCommand(t *testing.T) {
	var stderr bytes.Buffer
	d := deps{
		stdout:       bytes.NewBuffer(nil),
		stderr:       &stderr,
		runWorktrees: func(_ string) error { return nil },
	}
	err := execute([]string{"nope"}, d)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := stderr.String(); got == "" {
		t.Fatal("expected usage on stderr")
	}
}

func TestExecute_RunsPush(t *testing.T) {
	d := deps{
		stdout: bytes.NewBuffer(nil),
		stderr: bytes.NewBuffer(nil),
		getwd: func() (string, error) {
			return "/tmp", errors.New("boom")
		},
	}
	if err := execute([]string{"push"}, d); err == nil {
		t.Fatal("expected propagated error")
	}
}

func TestExecute_PushAllowedInRegularRepo(t *testing.T) {
	repoDir := testutil.TempRepo(t)
	d := deps{
		stdout: bytes.NewBuffer(nil),
		stderr: bytes.NewBuffer(nil),
		getwd: func() (string, error) {
			return repoDir, nil
		},
		confirmForce: func(string) (bool, error) { return false, nil },
	}

	err := execute([]string{"push"}, d)
	if err == nil {
		t.Fatal("expected push failure in test repo without remote")
	}
	if strings.Contains(err.Error(), "must be run from a regular repo or linked worktree") {
		t.Fatalf("regular repo should be allowed, got: %v", err)
	}
}

func TestExecute_PushRejectedInBareRepo(t *testing.T) {
	repoDir := testutil.TempBareRepo(t)
	d := deps{
		stdout: bytes.NewBuffer(nil),
		stderr: bytes.NewBuffer(nil),
		getwd: func() (string, error) {
			return repoDir, nil
		},
		confirmForce: func(string) (bool, error) { return false, nil },
	}

	err := execute([]string{"push"}, d)
	if err == nil {
		t.Fatal("expected error in bare repo")
	}
	if !strings.Contains(err.Error(), "must be run from a regular repo or linked worktree") {
		t.Fatalf("unexpected error: %v", err)
	}
}
