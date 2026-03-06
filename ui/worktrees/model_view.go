package worktrees

import (
	"fmt"

	"gx/ui"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing…"
	}
	if m.err != nil {
		return "\n  Error: " + m.err.Error()
	}
	if m.mode == modeYank {
		return m.yankModalView()
	}
	if m.mode == modeError {
		return m.errorModalView()
	}

	h := m.contentHeight()
	tableW, sidebarW := m.splitWidth()
	tableH, sidebarH := m.splitHeight(h)

	innerTableW := tableW - 2
	innerSidebarW := sidebarW - 2
	innerTableH := tableH - 2
	innerSidebarH := sidebarH - 2
	if innerTableW < 1 {
		innerTableW = 1
	}
	if innerSidebarW < 1 {
		innerSidebarW = 1
	}
	if innerTableH < 1 {
		innerTableH = 1
	}
	if innerSidebarH < 1 {
		innerSidebarH = 1
	}

	tableView := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorBorder).
		Width(innerTableW).
		Height(innerTableH).
		Render(m.table.View())

	sidebarView := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorBorder).
		Width(innerSidebarW).
		Height(innerSidebarH).
		Render(m.viewport.View())

	var content string
	if m.useStackedLayout() {
		content = lipgloss.JoinVertical(lipgloss.Left, tableView, sidebarView)
	} else {
		content = lipgloss.JoinHorizontal(lipgloss.Top, tableView, sidebarView)
	}
	return lipgloss.JoinVertical(lipgloss.Left, content, m.statusBarView())
}

// statusBarView renders the 1-line bar at the bottom of the screen.
func (m Model) statusBarView() string {
	switch m.mode {
	case modeError:
		return ""
	case modeDelete:
		return m.deleteConfirmView()
	case modeRename:
		return m.renameView()
	case modeClone:
		return m.cloneView()
	case modeNew:
		return m.newView()
	default:
		if m.spinnerActive {
			return "  " + m.spinner.View() + " " + m.spinnerLabel
		}
		if m.statusMsg != "" {
			return "  " + m.statusMsg
		}
		if m.clipboard != nil {
			prefix := fmt.Sprintf("  %d file(s) from %s  ·  ", len(m.clipboard.files), m.clipboard.srcName)
			m.help.Width = m.width - len(prefix)
			return ui.StyleDim.Render(prefix) + m.help.View(keys)
		}
		return m.help.View(keys)
	}
}
