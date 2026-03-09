package worktrees_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gx/git"
	"gx/testutil"
	"gx/ui/worktrees"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

const (
	termWidth  = 120
	termHeight = 40
	loadWait   = 5 * time.Second
	actionWait = 3 * time.Second
)

func startTUI(t *testing.T, repoDir string) (git.Repo, *teatest.TestModel) {
	t.Helper()
	repo, err := git.FindRepo(repoDir)
	if err != nil {
		t.Fatalf("FindRepo: %v", err)
	}
	m := worktrees.New(*repo, "")
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(termWidth, termHeight))
	return *repo, tm
}

func waitForText(t *testing.T, tm *teatest.TestModel, text string, timeout time.Duration) {
	t.Helper()
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte(text))
	}, teatest.WithDuration(timeout))
}

func quit(t *testing.T, tm *teatest.TestModel) {
	t.Helper()
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}

// ── delete ────────────────────────────────────────────────────────────────────

func TestDeleteConfirmationAppearsAndCancels(t *testing.T) {
	repoDir := testutil.TempBareRepoWithWorktrees(t, "feature-a")
	_, tm := startTUI(t, repoDir)

	waitForText(t, tm, "feature-a", loadWait)

	// Enter delete mode
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	waitForText(t, tm, "Delete", actionWait)

	// Cancel with esc — should return to normal without crashing
	tm.Send(tea.KeyMsg{Type: tea.KeyEsc})

	quit(t, tm)
}

func TestDeleteWorktree(t *testing.T) {
	repoDir := testutil.TempBareRepoWithWorktrees(t, "feature-a", "feature-b")
	repo, tm := startTUI(t, repoDir)

	waitForText(t, tm, "feature-a", loadWait)

	// Delete the selected (first) worktree
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	waitForText(t, tm, "Delete", actionWait)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

	// Wait until git actually has only 1 worktree left
	teatest.WaitFor(t, tm.Output(), func(_ []byte) bool {
		wts, err := git.ListWorktrees(repo)
		return err == nil && len(wts) == 1
	}, teatest.WithDuration(loadWait))

	quit(t, tm)
}

func TestDeleteWorktree_ShowsToastAfterDeletion(t *testing.T) {
	repoDir := testutil.TempBareRepoWithWorktrees(t, "feature-a", "feature-b")
	_, tm := startTUI(t, repoDir)

	waitForText(t, tm, "feature-a", loadWait)

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	waitForText(t, tm, "Delete", actionWait)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

	// The toast proves spinnerActive was cleared — if the spinner stays stuck
	// the model never re-renders status messages and this will time out.
	waitForText(t, tm, "feature-a deleted successfully", loadWait)

	quit(t, tm)
}

func TestDeleteCancelWithN(t *testing.T) {
	repoDir := testutil.TempBareRepoWithWorktrees(t, "feature-a")
	repo, tm := startTUI(t, repoDir)

	waitForText(t, tm, "feature-a", loadWait)

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	waitForText(t, tm, "Delete", actionWait)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

	// Worktree should still exist
	wts, err := git.ListWorktrees(repo)
	if err != nil {
		t.Fatalf("ListWorktrees: %v", err)
	}
	if len(wts) != 1 {
		t.Errorf("expected 1 worktree after cancel, got %d", len(wts))
	}

	quit(t, tm)
}

// ── clone ─────────────────────────────────────────────────────────────────────

func TestCloneInputAppearsAndCancels(t *testing.T) {
	repoDir := testutil.TempBareRepoWithWorktrees(t, "feature-a")
	_, tm := startTUI(t, repoDir)

	waitForText(t, tm, "feature-a", loadWait)

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	waitForText(t, tm, "Clone", actionWait)

	tm.Send(tea.KeyMsg{Type: tea.KeyEsc})

	quit(t, tm)
}

func TestCloneWorktree(t *testing.T) {
	repoDir := testutil.TempBareRepoWithWorktrees(t, "feature-a")
	_, tm := startTUI(t, repoDir)

	// Add an untracked file to the source worktree before starting the TUI
	wtDir := filepath.Join(repoDir, "feature-a")
	testutil.WriteFile(t, wtDir, "untracked.txt", "hello from untracked")

	waitForText(t, tm, "feature-a", loadWait)

	// Open clone input (pre-filled with "feature-a")
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	waitForText(t, tm, "Clone", actionWait)

	// Clear pre-filled value and type new name
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlU})
	tm.Type("feature-copy")
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

	// Wait until the untracked file appears in the clone. Waiting for the file
	// (rather than just the worktree in git's list) avoids a race where git
	// reports the worktree as existing before cmdClone has finished copying files.
	clonedFile := filepath.Join(repoDir, "feature-copy", "untracked.txt")
	teatest.WaitFor(t, tm.Output(), func(_ []byte) bool {
		_, err := os.ReadFile(clonedFile)
		return err == nil
	}, teatest.WithDuration(loadWait))

	data, err := os.ReadFile(clonedFile)
	if err != nil {
		t.Fatalf("untracked.txt missing in clone: %v", err)
	}
	if string(data) != "hello from untracked" {
		t.Errorf("untracked.txt content = %q, want %q", string(data), "hello from untracked")
	}

	quit(t, tm)
}

// ── new ───────────────────────────────────────────────────────────────────────

func TestNewInputAppearsAndCancels(t *testing.T) {
	repoDir := testutil.TempBareRepoWithWorktrees(t, "feature-a")
	_, tm := startTUI(t, repoDir)

	waitForText(t, tm, "feature-a", loadWait)

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	waitForText(t, tm, "New worktree", actionWait)

	tm.Send(tea.KeyMsg{Type: tea.KeyEsc})
	quit(t, tm)
}

func TestNewWorktree(t *testing.T) {
	repoDir := testutil.TempBareRepoWithWorktrees(t, "feature-a")
	repo, tm := startTUI(t, repoDir)

	waitForText(t, tm, "feature-a", loadWait)

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	waitForText(t, tm, "New worktree", actionWait)
	tm.Type("feature-new")
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

	teatest.WaitFor(t, tm.Output(), func(_ []byte) bool {
		wts, err := git.ListWorktrees(repo)
		if err != nil {
			return false
		}
		for _, wt := range wts {
			if wt.Name == "feature-new" && wt.Branch == "feature-new" {
				return true
			}
		}
		return false
	}, teatest.WithDuration(loadWait))

	quit(t, tm)
}

// ── yank / paste ──────────────────────────────────────────────────────────────

func TestYankModalAppearsAndCancels(t *testing.T) {
	repoDir := testutil.TempBareRepoWithWorktrees(t, "feature-a")
	_, tm := startTUI(t, repoDir)

	waitForText(t, tm, "feature-a", loadWait)

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	waitForText(t, tm, "Yank files from", actionWait)

	tm.Send(tea.KeyMsg{Type: tea.KeyEsc})

	quit(t, tm)
}

func TestYankAndPaste(t *testing.T) {
	repoDir := testutil.TempBareRepoWithWorktrees(t, "feature-a", "feature-b")

	// Add an untracked file to feature-a before starting the TUI
	wtDir := filepath.Join(repoDir, "feature-a")
	testutil.WriteFile(t, wtDir, "shared.txt", "hello from feature-a")

	_, tm := startTUI(t, repoDir)
	waitForText(t, tm, "feature-a", loadWait)

	// Yank from feature-a (cursor is on row 0)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	waitForText(t, tm, "Yank files from", actionWait)
	// Confirm with all items checked
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

	// Clipboard indicator should appear
	waitForText(t, tm, "feature-a", actionWait)

	// Navigate to feature-b and paste
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})

	// Wait until the pasted file appears in feature-b
	teatest.WaitFor(t, tm.Output(), func(_ []byte) bool {
		data, err := os.ReadFile(filepath.Join(repoDir, "feature-b", "shared.txt"))
		return err == nil && string(data) == "hello from feature-a"
	}, teatest.WithDuration(loadWait))

	quit(t, tm)
}

// ── push ──────────────────────────────────────────────────────────────────────

func TestPushWorktree(t *testing.T) {
	repoDir := testutil.TempBareRepoWithWorktrees(t, "feature-a")
	wtDir := filepath.Join(repoDir, "feature-a")

	// Push the branch to origin so a remote tracking ref exists, then add a commit.
	testutil.PushBranchWithUpstream(t, wtDir, "origin", "feature-a")
	testutil.WriteFile(t, wtDir, "extra.txt", "more content")
	testutil.CommitAll(t, wtDir, "second commit")

	_, tm := startTUI(t, repoDir)
	waitForText(t, tm, "feature-a", loadWait)

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'P'}})

	waitForText(t, tm, "Pushed", loadWait)

	quit(t, tm)
}

func TestPushRejectedShowsForcePushPrompt(t *testing.T) {
	repoDir := testutil.TempBareRepoWithWorktrees(t, "feature-a")
	wtDir := filepath.Join(repoDir, "feature-a")

	// Push to origin then amend the local commit to diverge from remote.
	testutil.PushBranchWithUpstream(t, wtDir, "origin", "feature-a")
	testutil.AmendLastCommit(t, wtDir)

	_, tm := startTUI(t, repoDir)
	waitForText(t, tm, "feature-a", loadWait)

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'P'}})

	// The model should detect the non-fast-forward error and show the confirm modal.
	waitForText(t, tm, "Force push?", loadWait)

	// 'q' in confirm mode cancels the modal; send it to return to normal, then quit.
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	quit(t, tm)
}

func TestPushRejectedForcePushConfirmed(t *testing.T) {
	repoDir := testutil.TempBareRepoWithWorktrees(t, "feature-a")
	wtDir := filepath.Join(repoDir, "feature-a")

	testutil.PushBranchWithUpstream(t, wtDir, "origin", "feature-a")
	testutil.AmendLastCommit(t, wtDir)

	_, tm := startTUI(t, repoDir)
	waitForText(t, tm, "feature-a", loadWait)

	// Trigger push (will be rejected).
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'P'}})
	waitForText(t, tm, "Force push?", loadWait)

	// Confirm force push with 'y'.
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

	waitForText(t, tm, "Force-pushed", loadWait)

	quit(t, tm)
}

// ── rename ────────────────────────────────────────────────────────────────────

func TestRenameInputAppearsAndCancels(t *testing.T) {
	repoDir := testutil.TempBareRepoWithWorktrees(t, "feature-a")
	_, tm := startTUI(t, repoDir)

	waitForText(t, tm, "feature-a", loadWait)

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	waitForText(t, tm, "Rename", actionWait)

	// Cancel with esc
	tm.Send(tea.KeyMsg{Type: tea.KeyEsc})

	quit(t, tm)
}

func TestRenameWorktree(t *testing.T) {
	repoDir := testutil.TempBareRepoWithWorktrees(t, "feature-a")
	repo, tm := startTUI(t, repoDir)

	waitForText(t, tm, "feature-a", loadWait)

	// Open rename input (pre-filled with "feature-a")
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	waitForText(t, tm, "Rename", actionWait)

	// Clear the pre-filled value with ctrl+u then type new name
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlU})
	tm.Type("feature-renamed")
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

	// Wait until git reports the renamed worktree
	teatest.WaitFor(t, tm.Output(), func(_ []byte) bool {
		wts, err := git.ListWorktrees(repo)
		if err != nil {
			return false
		}
		for _, wt := range wts {
			if wt.Name == "feature-renamed" {
				return true
			}
		}
		return false
	}, teatest.WithDuration(loadWait))

	quit(t, tm)
}
