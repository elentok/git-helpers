package worktrees

import (
	"gx/ui"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// enterConfirm switches to confirm mode with the given prompt and the command
// to run if the user selects Yes. spinnerLabel, if non-empty, activates the
// spinner while the command runs.
func (m Model) enterConfirm(prompt string, cmd tea.Cmd, spinnerLabel string) Model {
	m.mode = modeConfirm
	m.confirmPrompt = prompt
	m.confirmYes = false
	m.confirmCmd = cmd
	m.confirmSpinnerLabel = spinnerLabel
	m.statusMsg = ""
	return m
}

func (m Model) handleConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	toggleLeft := key.NewBinding(key.WithKeys("left", "h"))
	toggleRight := key.NewBinding(key.WithKeys("right", "l"))
	yes := key.NewBinding(key.WithKeys("y"))
	no := key.NewBinding(key.WithKeys("n"))
	cancel := key.NewBinding(key.WithKeys("esc", "q"))
	submit := key.NewBinding(key.WithKeys("enter"))

	switch {
	case key.Matches(msg, toggleLeft):
		m.confirmYes = true
	case key.Matches(msg, toggleRight):
		m.confirmYes = false
	case key.Matches(msg, yes):
		return m.runConfirmed()
	case key.Matches(msg, no), key.Matches(msg, cancel):
		m.mode = modeNormal
	case key.Matches(msg, submit):
		if m.confirmYes {
			return m.runConfirmed()
		}
		m.mode = modeNormal
	}
	return m, nil
}

func (m Model) runConfirmed() (tea.Model, tea.Cmd) {
	m.mode = modeNormal
	if m.confirmSpinnerLabel != "" {
		m.spinnerActive = true
		m.spinnerLabel = m.confirmSpinnerLabel
		return m, tea.Batch(m.confirmCmd, m.spinner.Tick)
	}
	return m, m.confirmCmd
}

func (m Model) confirmModalView() string {
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorBorder).
		Padding(1, 2)

	yesLabel := confirmButton("Yes", m.confirmYes)
	noLabel := confirmButton("No", !m.confirmYes)
	hint := ui.StyleDim.Render("←/→ or h/l: choose  y/n: quick select  enter: confirm")

	inner := lipgloss.JoinVertical(lipgloss.Left,
		m.confirmPrompt,
		"",
		"  "+yesLabel+"   "+noLabel,
		"",
		hint,
	)
	modal := borderStyle.Render(inner)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
}

func confirmButton(label string, selected bool) string {
	s := lipgloss.NewStyle().Padding(0, 1)
	if selected {
		return s.Foreground(ui.ColorGreen).Bold(true).Render("> " + label + " <")
	}
	return s.Foreground(ui.ColorGray).Render("  " + label + "  ")
}
