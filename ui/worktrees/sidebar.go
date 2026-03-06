package worktrees

import (
	"strings"

	"gx/git"
	"gx/ui"

	"github.com/charmbracelet/lipgloss"
)

func renderSidebarContent(wt *git.Worktree, commits []git.Commit, changes []git.Change, loading bool) string {
	if wt == nil {
		return ui.StyleDim.Render("  no worktree selected")
	}
	if loading {
		return ui.StyleDim.Render("  loading…")
	}

	var b strings.Builder

	b.WriteString(ui.StyleBold.Render("Commits ahead of main"))
	b.WriteString("\n\n")
	if len(commits) == 0 {
		b.WriteString(ui.StyleDim.Render("  none") + "\n")
	} else {
		for _, c := range commits {
			b.WriteString("  ")
			b.WriteString(ui.StyleDim.Render(c.Hash))
			b.WriteString("  ")
			b.WriteString(c.Subject)
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(ui.StyleBold.Render("Changes"))
	b.WriteString("\n\n")
	if len(changes) == 0 {
		b.WriteString(ui.StyleDim.Render("  clean") + "\n")
	} else {
		for _, c := range changes {
			b.WriteString("  ")
			b.WriteString(changeKindStyle(c.Kind).Render(string(c.Kind)))
			b.WriteString("  ")
			b.WriteString(c.Path)
			b.WriteString("\n")
		}
	}

	return b.String()
}

func changeKindStyle(k git.ChangeKind) lipgloss.Style {
	switch k {
	case git.ChangeAdded:
		return ui.StyleStatusSynced
	case git.ChangeDeleted:
		return ui.StyleStatusDiverged
	case git.ChangeModified:
		return ui.StyleStatusBehind
	default:
		return ui.StyleDim
	}
}
