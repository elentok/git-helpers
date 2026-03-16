package worktrees

import (
	"strings"

	"gx/git"
	"gx/ui"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

var styleMainBranch = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))

// sortedWorktrees returns a copy of wts with the main branch worktree first.
func sortedWorktrees(wts []git.Worktree, mainBranch string) []git.Worktree {
	if mainBranch == "" {
		return wts
	}
	out := make([]git.Worktree, 0, len(wts))
	var main *git.Worktree
	for i := range wts {
		if wts[i].Branch == mainBranch {
			main = &wts[i]
		} else {
			out = append(out, wts[i])
		}
	}
	if main != nil {
		out = append([]git.Worktree{*main}, out...)
	}
	return out
}

// tableStyles holds the styles configured in newTable so our custom renderer
// can use them without needing access to the unexported table.Model.styles field.
var tableStyles table.Styles

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
	tableStyles = s

	return t
}

func resizeTable(t *table.Model, width, height int) {
	// Account for default table cell left/right padding (2 chars per column)
	// plus inter-column spaces to avoid overflow/wrapping.
	const (
		cols       = 5
		separators = cols - 1
		padding    = cols * 2
	)
	usable := width - separators - padding
	if usable < 20 {
		usable = 20
	}

	branchW := int(float64(usable) * 0.25)
	dirtyW := 5
	baseW := 4
	statusW := int(float64(usable) * 0.20)
	if statusW < 8 {
		statusW = 8
	}
	nameW := usable - branchW - dirtyW - baseW - statusW
	if nameW < 8 {
		nameW = 8
	}
	t.SetColumns([]table.Column{
		{Title: "Worktree", Width: nameW},
		{Title: "Branch", Width: branchW},
		{Title: "Dirty", Width: dirtyW},
		{Title: "Base", Width: baseW},
		{Title: "Status", Width: statusW},
	})
	t.SetWidth(width)
	t.SetHeight(height)
}

// tableView renders the table using ansi.Truncate instead of the
// runewidth.Truncate used internally by bubbles/table. This allows cell values
// to contain arbitrary ANSI escape sequences (e.g. lipgloss highlights) without
// column-alignment corruption, because ansi.Truncate is ANSI-aware and will
// never cut through an escape sequence.
func tableView(t table.Model) string {
	return headersView(t) + "\n" + rowsView(t)
}

func headersView(t table.Model) string {
	cols := t.Columns()
	s := make([]string, 0, len(cols))
	for _, col := range cols {
		if col.Width <= 0 {
			continue
		}
		style := lipgloss.NewStyle().Width(col.Width).MaxWidth(col.Width).Inline(true)
		cell := style.Render(ansi.Truncate(col.Title, col.Width, "…"))
		s = append(s, tableStyles.Header.Render(cell))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, s...)
}

func rowsView(t table.Model) string {
	rows := t.Rows()
	cols := t.Columns()
	cursor := t.Cursor()
	height := t.Height()

	start := cursor - height/2
	if start > len(rows)-height {
		start = len(rows) - height
	}
	if start < 0 {
		start = 0
	}
	end := start + height
	if end > len(rows) {
		end = len(rows)
	}

	rendered := make([]string, 0, end-start)
	for i := start; i < end; i++ {
		rendered = append(rendered, renderRow(rows[i], cols, i == cursor))
	}
	return strings.Join(rendered, "\n")
}

func renderRow(row table.Row, cols []table.Column, selected bool) string {
	cells := make([]string, 0, len(cols))
	for i, col := range cols {
		if col.Width <= 0 {
			continue
		}
		value := ""
		if i < len(row) {
			value = row[i]
		}
		style := lipgloss.NewStyle().Width(col.Width).MaxWidth(col.Width).Inline(true)
		cell := tableStyles.Cell.Render(style.Render(ansi.Truncate(value, col.Width, "…")))
		cells = append(cells, cell)
	}
	rowStr := lipgloss.JoinHorizontal(lipgloss.Top, cells...)
	if selected {
		return tableStyles.Selected.Render(rowStr)
	}
	return rowStr
}

// buildRows builds the table rows, applying search highlighting when a query
// is active. Since tableView uses ansi.Truncate (which is ANSI-aware), cell
// values may contain arbitrary lipgloss styles without any pre-truncation.
func (m Model) buildRows() []table.Row {
	ic := icons(m.settings.UseNerdFontIcons)
	rows := make([]table.Row, len(m.worktrees))
	for i, wt := range m.worktrees {
		isSelected := i == m.table.Cursor()
		isMain := wt.Branch == m.repo.MainBranch
		nameCol := worktreeCell(wt.Name, ic, isMain, isSelected)
		branchCol := branchCell(wt.Branch, ic, isMain, isSelected)
		if m.searchQuery != "" && !isSelected {
			nameCol = highlightMatch(nameCol, m.searchQuery)
			branchCol = highlightMatch(branchCol, m.searchQuery)
		}
		rows[i] = table.Row{
			nameCol,
			branchCol,
			dirtyCell(m.dirties[wt.Path], isSelected),
			baseCell(m.baseStatus[wt.Branch], wt.Branch == m.repo.MainBranch, isSelected),
			statusCell(m.statuses[wt.Branch], isSelected, m.settings.UseNerdFontIcons),
		}
	}
	return rows
}

var styleSearchHighlight = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true)

// highlightMatch wraps the first occurrence of query (case-insensitive) in
// text with a yellow bold lipgloss style.
func highlightMatch(text, query string) string {
	lower := strings.ToLower(text)
	lq := strings.ToLower(query)
	idx := strings.Index(lower, lq)
	if idx < 0 {
		return text
	}
	return text[:idx] + styleSearchHighlight.Render(text[idx:idx+len(query)]) + text[idx+len(query):]
}

func worktreeCell(name string, ic uiIcons, isMain, isSelected bool) string {
	prefix := ic.worktreePrefix
	if isMain && ic.mainPrefix != "" {
		prefix = ic.mainPrefix
	}
	text := prefix + name
	if isMain && !isSelected {
		return styleMainBranch.Render(text)
	}
	return text
}

func branchCell(name string, ic uiIcons, isMain, isSelected bool) string {
	text := name
	if ic.branchPrefix != "" && name != "" {
		text = ic.branchPrefix + name
	}
	if isMain && !isSelected {
		return styleMainBranch.Render(text)
	}
	return text
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

func baseCell(rebased *bool, isMainBranch bool, selected bool) string {
	if isMainBranch {
		if selected {
			return "—"
		}
		return ui.StyleDim.Render("—")
	}
	if rebased == nil {
		return "" // not yet loaded
	}
	if *rebased {
		if selected {
			return "✓"
		}
		return ui.StyleStatusSynced.Render("✓")
	}
	if selected {
		return "✗"
	}
	return ui.StyleStatusDiverged.Render("✗")
}

func statusCell(s git.SyncStatus, selected bool, useNerdFontIcons bool) string {
	label := "—"
	switch s.Name {
	case git.StatusSame:
		label = "synced"
	case git.StatusAhead, git.StatusBehind, git.StatusDiverged:
		label = s.Pretty()
	}
	if useNerdFontIcons {
		label = strings.ReplaceAll(label, "ahead", "\uf062")
		label = strings.ReplaceAll(label, "behind", "\uf063")
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
