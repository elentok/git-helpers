package worktrees

import (
	"gx/git"
	"gx/ui"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func newTable() table.Model {
	t := table.New(table.WithFocused(true))

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(ui.ColorBorder).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}

func resizeTable(t *table.Model, width, height int) {
	// Account for default table cell left/right padding (2 chars per column)
	// plus inter-column spaces to avoid overflow/wrapping.
	const (
		cols       = 4
		separators = cols - 1
		padding    = cols * 2
	)
	usable := width - separators - padding
	if usable < 20 {
		usable = 20
	}

	nameW := int(float64(usable) * 0.31)
	branchW := int(float64(usable) * 0.25)
	dirtyW := 5
	statusW := usable - nameW - branchW - dirtyW
	if statusW < 8 {
		statusW = 8
	}
	t.SetColumns([]table.Column{
		{Title: "Worktree", Width: nameW},
		{Title: "Branch", Width: branchW},
		{Title: "Dirty", Width: dirtyW},
		{Title: "Status", Width: statusW},
	})
	t.SetWidth(width)
	t.SetHeight(height)
}

func buildRows(worktrees []git.Worktree, statuses map[string]git.SyncStatus, dirties map[string]dirtyState, selected int) []table.Row {
	rows := make([]table.Row, len(worktrees))
	for i, wt := range worktrees {
		isSelected := i == selected
		rows[i] = table.Row{
			wt.Name,
			wt.Branch,
			dirtyCell(dirties[wt.Path], isSelected),
			statusCell(statuses[wt.Branch], isSelected),
		}
	}
	return rows
}

func dirtyCell(d dirtyState, selected bool) string {
	symbol := "-"
	switch {
	case d.hasModified && d.hasUntracked:
		symbol = "M?"
	case d.hasModified:
		symbol = "M"
	case d.hasUntracked:
		symbol = "?"
	}
	return symbol
}

func statusCell(s git.SyncStatus, selected bool) string {
	label := "—"
	switch s.Name {
	case git.StatusSame:
		label = "synced"
	case git.StatusAhead, git.StatusBehind, git.StatusDiverged:
		label = s.Pretty()
	}
	if selected {
		return label
	}
	switch s.Name {
	case git.StatusSame:
		return ui.StyleStatusSynced.Render(label)
	case git.StatusAhead:
		return ui.StyleStatusAhead.Render(label)
	case git.StatusBehind:
		return ui.StyleStatusBehind.Render(label)
	case git.StatusDiverged:
		return ui.StyleStatusDiverged.Render(label)
	default:
		return ui.StyleStatusUnknown.Render(label)
	}
}
