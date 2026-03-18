package confirm

import (
	"io"
	"os"
	"strings"

	"gx/ui"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type doneMsg struct{}

type model struct {
	prompt    string
	choiceYes bool
	done      bool
}

// Run renders a small styled confirmation UI and returns true when accepted.
func Run(prompt string) (bool, error) {
	return run(prompt, os.Stdin, os.Stdout)
}

func run(prompt string, in io.Reader, out io.Writer) (bool, error) {
	m := model{prompt: prompt}
	p := tea.NewProgram(m, tea.WithInput(in), tea.WithOutput(out))
	finalModel, err := p.Run()
	if err != nil {
		return false, err
	}
	fm := finalModel.(model)
	return fm.done && fm.choiceYes, nil
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "left", "h":
			m.choiceYes = true
		case "right", "l":
			m.choiceYes = false
		case "y":
			m.choiceYes = true
			m.done = true
			return m, tea.Quit
		case "n":
			m.choiceYes = false
			m.done = true
			return m, tea.Quit
		case "enter":
			m.done = true
			return m, tea.Quit
		case "ctrl+c", "esc", "q":
			m.choiceYes = false
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	body := lipgloss.NewStyle().Padding(1, 2).Render(m.prompt)
	hint := ui.StyleDim.Render("left/right: choose  y/n: quick select  enter: confirm")
	yes := optionLabel("Yes", m.choiceYes)
	no := optionLabel("No", !m.choiceYes)
	return tea.NewView(strings.Join([]string{
		body,
		"  " + yes + "   " + no,
		"  " + hint,
		"",
	}, "\n"))
}

func optionLabel(label string, selected bool) string {
	s := lipgloss.NewStyle().Padding(0, 1)
	if selected {
		s = s.Foreground(ui.ColorGreen).Bold(true)
		return s.Render("> " + label + " <")
	} else {
		s = s.Foreground(ui.ColorGray)
		return s.Render("  " + label + "  ")
	}
}
