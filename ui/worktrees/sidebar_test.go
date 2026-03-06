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

	out := renderSidebarContent(wt, ahead, behind, nil, false)
	if !strings.Contains(out, "Commits ahead of main") {
		t.Fatal("missing ahead section")
	}
	if !strings.Contains(out, "Commits behind main") {
		t.Fatal("missing behind section")
	}
	if !strings.Contains(out, "behind commit") {
		t.Fatal("missing behind commit entry")
	}
}
