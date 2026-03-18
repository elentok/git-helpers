package worktrees

import (
	"gx/ui"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// enterLogsMode switches the model into logs mode to display the last job output.
func (m Model) enterLogsMode() Model {
	vpW := m.width * 2 / 3
	if vpW < 40 {
		vpW = 40
	}
	if vpW > 100 {
		vpW = 100
	}
	vpH := m.height/2 - 6
	if vpH < 3 {
		vpH = 3
	}
	vp := viewport.New(viewport.WithWidth(vpW-2), viewport.WithHeight(vpH))
	vp.SetContent(m.lastJobLog)
	m.logsViewport = vp
	m.mode = modeLogs
	return m
}

// handleLogsKey scrolls the logs viewport or dismisses it.
func (m Model) handleLogsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc, tea.KeyEnter:
		m.mode = modeNormal
		return m, nil
	}
	if msg.Type == tea.KeyRunes && msg.String() == "q" {
		m.mode = modeNormal
		return m, nil
	}
	var cmd tea.Cmd
	m.logsViewport, cmd = m.logsViewport.Update(msg)
	return m, cmd
}

// logsModalView renders a centred modal with the last job's output.
func (m Model) logsModalView() string {
	titleStyle := lipgloss.NewStyle().Bold(true)
	hintStyle := ui.StyleDim
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorBorder).
		Padding(0, 1).
		Width(m.logsViewport.Width())

	content := m.lastJobLog
	if content == "" {
		content = ui.StyleDim.Render("(no output)")
	}
	m.logsViewport.SetContent(content)

	inner := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(m.lastJobLabel),
		"",
		m.logsViewport.View(),
		"",
		hintStyle.Render("esc / enter / q  to dismiss"),
	)
	return borderStyle.Render(inner)
}
