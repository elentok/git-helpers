package worktrees

import (
	"strings"

	"gx/git"
)

func (m Model) splitWidth() (tableW, sidebarW int) {
	if m.useStackedLayout() {
		return m.width, m.width
	}
	tableW = int(float64(m.width) * 0.55)
	sidebarW = m.width - tableW
	return
}

func (m Model) splitHeight(total int) (tableH, sidebarH int) {
	if !m.useStackedLayout() {
		return total, total
	}
	tableH = int(float64(total) * 0.58)
	if tableH < 8 {
		tableH = 8
	}
	if tableH > total-6 {
		tableH = total - 6
	}
	if tableH < 1 {
		tableH = 1
	}
	sidebarH = total - tableH
	if sidebarH < 1 {
		sidebarH = 1
	}
	return
}

func (m Model) useStackedLayout() bool {
	return m.width <= 100
}

func (m Model) helpLineCount() int {
	v := m.help.View(keys)
	if v == "" {
		return 1
	}
	return strings.Count(v, "\n") + 1
}

func (m Model) contentHeight() int {
	h := m.height - m.helpLineCount()
	if h < 4 {
		return 4
	}
	return h
}

func (m Model) resized() Model {
	m.help.Width = m.width
	tableW, sidebarW := m.splitWidth()
	h := m.contentHeight()
	tableH, sidebarH := m.splitHeight(h)

	tableInnerW := tableW - 2
	tableInnerH := tableH - 2
	if tableInnerW < 1 {
		tableInnerW = 1
	}
	if tableInnerH < 1 {
		tableInnerH = 1
	}
	resizeTable(&m.table, tableInnerW, tableInnerH)

	vpW := sidebarW - 2
	vpH := sidebarH - 2
	if vpW < 1 {
		vpW = 1
	}
	if vpH < 1 {
		vpH = 1
	}
	m.viewport.Width = vpW
	m.viewport.Height = vpH
	m.viewport.SetContent(m.sidebarContent())

	return m
}

func (m Model) sidebarContent() string {
	var wt *git.Worktree
	if len(m.worktrees) > 0 {
		w := m.worktrees[m.table.Cursor()]
		wt = &w
	}
	return renderSidebarContent(
		wt,
		m.sidebarAheadCommits,
		m.sidebarBehindCommits,
		m.sidebarChanges,
		m.sidebarLoading,
		m.settings.UseNerdFontIcons,
	)
}
