package worktrees

import (
	"gx/ui"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// showError switches the model into error mode with a scrollable viewport.
func (m Model) showError(errMsg string) Model {
	vpW := m.width * 2 / 3
	if vpW < 40 {
		vpW = 40
	}
	if vpW > 80 {
		vpW = 80
	}
	vpH := m.height/2 - 6
	if vpH < 3 {
		vpH = 3
	}
	vp := viewport.New(viewport.WithWidth(vpW-2), viewport.WithHeight(vpH))
	vp.SetContent(errMsg)
	m.errorViewport = vp
	m.mode = modeError
	return m
}

// handleErrorKey scrolls the error viewport or dismisses it.
func (m Model) handleErrorKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "enter", "q":
		m.mode = modeNormal
		return m, nil
	case "o":
		if m.lastJobLog != "" {
			return m.enterLogsMode(), nil
		}
	}
	var cmd tea.Cmd
	m.errorViewport, cmd = m.errorViewport.Update(msg)
	return m, cmd
}

func (m Model) errorHint() string {
	hint := "esc / enter / q  dismiss"
	if m.lastJobLog != "" {
		hint += "  ·  o  view output"
	}
	return ui.StyleDim.Render(hint)
}

// errorModalView renders a centred modal with the error text.
func (m Model) errorModalView() string {
	titleStyle := lipgloss.NewStyle().Foreground(ui.ColorRed).Bold(true)
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorRed).
		Padding(0, 1).
		Width(m.errorViewport.Width())

	inner := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Error"),
		"",
		m.errorViewport.View(),
		"",
		m.errorHint(),
	)
	return borderStyle.Render(inner)
}
