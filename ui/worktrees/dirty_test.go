package worktrees

import (
	"strings"
	"testing"

	"gx/git"
)

func TestDirtyStateFromChanges(t *testing.T) {
	tests := []struct {
		name    string
		changes []git.Change
		want    dirtyState
	}{
		{name: "clean", changes: nil, want: dirtyState{}},
		{
			name:    "modified only",
			changes: []git.Change{{Kind: git.ChangeModified, Path: "a.txt"}},
			want:    dirtyState{hasModified: true},
		},
		{
			name:    "untracked only",
			changes: []git.Change{{Kind: git.ChangeUntracked, Path: "a.txt"}},
			want:    dirtyState{hasUntracked: true},
		},
		{
			name: "mixed",
			changes: []git.Change{
				{Kind: git.ChangeUntracked, Path: "a.txt"},
				{Kind: git.ChangeModified, Path: "b.txt"},
			},
			want: dirtyState{hasModified: true, hasUntracked: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dirtyStateFromChanges(tt.changes)
			if got != tt.want {
				t.Fatalf("dirtyStateFromChanges() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestDirtyCellSymbols(t *testing.T) {
	tests := []struct {
		name  string
		dirty dirtyState
		want  string
	}{
		{name: "clean", dirty: dirtyState{}, want: "-"},
		{name: "modified", dirty: dirtyState{hasModified: true}, want: "M"},
		{name: "untracked", dirty: dirtyState{hasUntracked: true}, want: "?"},
		{name: "mixed", dirty: dirtyState{hasModified: true, hasUntracked: true}, want: "M?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dirtyCell(tt.dirty, false)
			if !strings.Contains(got, tt.want) {
				t.Fatalf("dirtyCell() = %q, want symbol %q", got, tt.want)
			}
		})
	}
}

func TestDirtyAndStatusCellSelectedArePlain(t *testing.T) {
	if got := dirtyCell(dirtyState{hasModified: true, hasUntracked: true}, true); got != "M?" {
		t.Fatalf("dirtyCell(selected) = %q, want %q", got, "M?")
	}
	if got := statusCell(git.SyncStatus{Name: git.StatusSame}, true); got != "synced" {
		t.Fatalf("statusCell(selected) = %q, want %q", got, "synced")
	}
}
