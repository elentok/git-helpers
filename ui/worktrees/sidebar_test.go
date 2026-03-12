package worktrees

import (
	"strings"
	"testing"

	"gx/git"
)

func TestRenderSidebarContent_IncludesBehindSection(t *testing.T) {
	wt := &git.Worktree{Name: "feature-a"}
	ahead := []git.Commit{{Hash: "abc1234", Subject: "ahead commit"}}
	behind := []git.Commit{{Hash: "def5678", Subject: "behind commit"}}

	out := renderSidebarContent(wt, "origin/feature-a", ahead, behind, nil, false, nil, false, false)
	if !strings.Contains(out, "Commits ahead of remote") {
		t.Fatal("missing ahead section")
	}
	if !strings.Contains(out, "Commits behind remote") {
		t.Fatal("missing behind section")
	}
	if !strings.Contains(out, "behind commit") {
		t.Fatal("missing behind commit entry")
	}
}

func TestRenderSidebarContent_NoUpstream(t *testing.T) {
	wt := &git.Worktree{Name: "feature-a"}
	out := renderSidebarContent(wt, "", nil, nil, nil, false, nil, false, false)
	if !strings.Contains(out, "no remote tracking branch") {
		t.Fatal("missing no-tracking note")
	}
	if strings.Contains(out, "Commits ahead") {
		t.Fatal("should not show ahead section when no upstream")
	}
}

func TestRenderSidebarContent_UsesNerdFontIcons(t *testing.T) {
	wt := &git.Worktree{Name: "feature-a"}
	out := renderSidebarContent(wt, "origin/feature-a", nil, nil, nil, false, nil, false, true)
	if !strings.Contains(out, "󰙅 Worktree") {
		t.Fatal("missing nerd-font worktree title")
	}
	if !strings.Contains(out, " Commits ahead of remote") {
		t.Fatal("missing nerd-font ahead title")
	}
}

func TestRenderSidebarContent_BehindMain(t *testing.T) {
	wt := &git.Worktree{Name: "feature-a", Branch: "feature-a"}
	behindMain := []git.Commit{
		{Hash: "aaa1111", Subject: "main commit 1"},
		{Hash: "bbb2222", Subject: "main commit 2"},
	}

	out := renderSidebarContent(wt, "origin/feature-a", nil, nil, behindMain, false, nil, false, false)
	if !strings.Contains(out, "Commits behind main") {
		t.Fatal("missing behind-main section")
	}
	if !strings.Contains(out, "main commit 1") {
		t.Fatal("missing behind-main commit entry")
	}
}

func TestRenderSidebarContent_RebasedOnMain(t *testing.T) {
	wt := &git.Worktree{Name: "feature-a", Branch: "feature-a"}
	out := renderSidebarContent(wt, "origin/feature-a", nil, nil, []git.Commit{}, false, nil, false, false)
	if !strings.Contains(out, "rebased on main") {
		t.Fatal("expected 'rebased on main' indicator")
	}
}

func TestRenderSidebarContent_MainBranchHidesSection(t *testing.T) {
	wt := &git.Worktree{Name: "main", Branch: "main"}
	out := renderSidebarContent(wt, "origin/main", nil, nil, nil, true, nil, false, false)
	if strings.Contains(out, "Commits behind main") {
		t.Fatal("main branch should not show behind-main section")
	}
}
