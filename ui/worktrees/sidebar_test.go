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

	out := renderSidebarContent(wt, ahead, behind, nil, false, false)
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

func TestRenderSidebarContent_UsesNerdFontIcons(t *testing.T) {
	wt := &git.Worktree{Name: "feature-a"}
	out := renderSidebarContent(wt, nil, nil, nil, false, true)
	if !strings.Contains(out, "󰙅 Worktree") {
		t.Fatal("missing nerd-font worktree title")
	}
	if !strings.Contains(out, " Commits ahead of remote") {
		t.Fatal("missing nerd-font ahead title")
	}
}
