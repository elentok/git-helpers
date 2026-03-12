package worktrees

import (
	"strings"

	"gx/git"
	"gx/ui"

	"github.com/charmbracelet/lipgloss"
)

func renderSidebarContent(wt *git.Worktree, upstream string, aheadCommits, behindCommits, behindMainCommits []git.Commit, isMainBranch bool, changes []git.Change, loading bool, useNerdFontIcons bool) string {
	if wt == nil {
		return ui.StyleDim.Render("  no worktree selected")
	}
	if loading {
		return ui.StyleDim.Render("  loading…")
	}

	var b strings.Builder
	ic := icons(useNerdFontIcons)

	b.WriteString(ui.StyleBold.Render(ic.worktreeTitle))
	b.WriteString("\n\n")
	b.WriteString("  ")
	b.WriteString(wt.Name)
	b.WriteString("\n\n")

	if upstream == "" {
		b.WriteString(ui.StyleDim.Render("  no remote tracking branch") + "\n")
		b.WriteString(ui.StyleDim.Render("  press t to track origin/<branch>") + "\n")
	} else {
		b.WriteString(ui.StyleBold.Render(ic.aheadTitle))
		b.WriteString("\n\n")
		if len(aheadCommits) == 0 {
			b.WriteString(ui.StyleDim.Render("  none") + "\n")
		} else {
			for _, c := range aheadCommits {
				b.WriteString("  ")
				b.WriteString(ui.StyleDim.Render(c.Hash))
				b.WriteString("  ")
				b.WriteString(c.Subject)
				b.WriteString("\n")
			}
		}

		b.WriteString("\n")
		b.WriteString(ui.StyleBold.Render(ic.behindTitle))
		b.WriteString("\n\n")
		if len(behindCommits) == 0 {
			b.WriteString(ui.StyleDim.Render("  none") + "\n")
		} else {
			for _, c := range behindCommits {
				b.WriteString("  ")
				b.WriteString(ui.StyleDim.Render(c.Hash))
				b.WriteString("  ")
				b.WriteString(c.Subject)
				b.WriteString("\n")
			}
		}
	}

	if !isMainBranch {
		b.WriteString("\n")
		b.WriteString(ui.StyleBold.Render(ic.behindMainTitle))
		b.WriteString("\n\n")
		switch {
		case behindMainCommits == nil:
			b.WriteString(ui.StyleDim.Render("  loading…") + "\n")
		case len(behindMainCommits) == 0:
			b.WriteString(ui.StyleStatusSynced.Render("  ✓ rebased on main") + "\n")
		default:
			for _, c := range behindMainCommits {
				b.WriteString("  ")
				b.WriteString(ui.StyleDim.Render(c.Hash))
				b.WriteString("  ")
				b.WriteString(c.Subject)
				b.WriteString("\n")
			}
		}
	}

	b.WriteString("\n")
	b.WriteString(ui.StyleBold.Render(ic.changesTitle))
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
