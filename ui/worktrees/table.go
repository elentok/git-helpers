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
	nameW := int(float64(width) * 0.35)
	branchW := int(float64(width) * 0.30)
	statusW := width - nameW - branchW
	if statusW < 8 {
		statusW = 8
	}
	t.SetColumns([]table.Column{
		{Title: "Worktree", Width: nameW},
		{Title: "Branch", Width: branchW},
		{Title: "Status", Width: statusW},
	})
	t.SetWidth(width)
	t.SetHeight(height)
}

func buildRows(worktrees []git.Worktree, statuses map[string]git.SyncStatus) []table.Row {
	rows := make([]table.Row, len(worktrees))
	for i, wt := range worktrees {
		rows[i] = table.Row{wt.Name, wt.Branch, statusCell(statuses[wt.Branch])}
	}
	return rows
}

func statusCell(s git.SyncStatus) string {
	switch s.Name {
	case git.StatusSame:
		return ui.StyleStatusSynced.Render("synced")
	case git.StatusAhead:
		return ui.StyleStatusAhead.Render(s.Pretty())
	case git.StatusBehind:
		return ui.StyleStatusBehind.Render(s.Pretty())
	case git.StatusDiverged:
		return ui.StyleStatusDiverged.Render(s.Pretty())
	default:
		return ui.StyleStatusUnknown.Render("—")
	}
}
