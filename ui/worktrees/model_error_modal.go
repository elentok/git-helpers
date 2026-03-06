package worktrees

import (
	"gx/ui"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	vp := viewport.New(vpW-2, vpH)
	vp.SetContent(errMsg)
	m.errorViewport = vp
	m.mode = modeError
	return m
}

// handleErrorKey scrolls the error viewport or dismisses it.
func (m Model) handleErrorKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	m.errorViewport, cmd = m.errorViewport.Update(msg)
	return m, cmd
}

// errorModalView renders a centred modal with the error text.
func (m Model) errorModalView() string {
	titleStyle := lipgloss.NewStyle().Foreground(ui.ColorRed).Bold(true)
	hintStyle := ui.StyleDim
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorRed).
		Padding(0, 1).
		Width(m.errorViewport.Width)

	inner := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Error"),
		"",
		m.errorViewport.View(),
		"",
		hintStyle.Render("esc / enter / q  to dismiss"),
	)
	modal := borderStyle.Render(inner)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
}
