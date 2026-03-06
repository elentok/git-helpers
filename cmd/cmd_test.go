package cmd

import (
	"bytes"
	"errors"
	"io"
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

func TestExecute_Init(t *testing.T) {
	var stdout bytes.Buffer
	called := false
	d := deps{
		stdout: &stdout,
		stderr: bytes.NewBuffer(nil),
		initConfig: func() (string, error) {
			called = true
			return "/tmp/gx/config.json", nil
		},
	}

	if err := execute([]string{"init"}, d); err != nil {
		t.Fatalf("execute init: %v", err)
	}
	if !called {
		t.Fatal("expected initConfig to be called")
	}
	if !strings.Contains(stdout.String(), "Created config file at /tmp/gx/config.json") {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
}

func TestExecute_EditConfig_RequiresEditor(t *testing.T) {
	d := deps{
		stdout: bytes.NewBuffer(nil),
		stderr: bytes.NewBuffer(nil),
		initConfig: func() (string, error) {
			return "/tmp/gx/config.json", nil
		},
		getenv: func(string) string { return "" },
	}

	err := execute([]string{"edit-config"}, d)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "$EDITOR is not set") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecute_EditConfig_RunsEditor(t *testing.T) {
	var stdout bytes.Buffer
	var gotEditor, gotPath string
	d := deps{
		stdout: &stdout,
		stderr: bytes.NewBuffer(nil),
		initConfig: func() (string, error) {
			return "/tmp/gx/config.json", nil
		},
		getenv: func(k string) string {
			if k == "EDITOR" {
				return "vim"
			}
			return ""
		},
		runEditor: func(editor, path string, _ io.Reader, _, _ io.Writer) error {
			gotEditor = editor
			gotPath = path
			return nil
		},
	}

	if err := execute([]string{"edit-config"}, d); err != nil {
		t.Fatalf("execute edit-config: %v", err)
	}
	if gotEditor != "vim" {
		t.Fatalf("editor = %q, want %q", gotEditor, "vim")
	}
	if gotPath == "" {
		t.Fatal("expected non-empty config path")
	}
}
