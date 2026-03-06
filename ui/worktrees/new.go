package worktrees

import (
	"path/filepath"
	"strings"

	"gx/git"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type newResultMsg struct{ err error }

func cmdNewWorktree(repo git.Repo, newName string) tea.Cmd {
	return func() tea.Msg {
		newPath := filepath.Join(repo.Root, newName)
		return newResultMsg{err: git.AddWorktree(repo, newName, newPath, repo.MainBranch)}
	}
}

func newWorktreeInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "worktree-name"
	ti.Focus()
	return ti
}

func (m Model) enterNewMode() Model {
	m.mode = modeNew
	m.textInput = newWorktreeInput()
	m.statusMsg = ""
	return m
}

func (m Model) handleNewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	cancel := key.NewBinding(key.WithKeys("esc"))
	submit := key.NewBinding(key.WithKeys("enter"))

	switch {
	case key.Matches(msg, cancel):
		m.mode = modeNormal
		m.statusMsg = ""
		return m, nil
	case key.Matches(msg, submit):
		newName := strings.TrimSpace(m.textInput.Value())
		if newName == "" {
			m.mode = modeNormal
			m.statusMsg = ""
			return m, nil
		}
		m.mode = modeNormal
		m.statusMsg = "Creating…"
		return m, cmdNewWorktree(m.repo, newName)
	}

	var tiCmd tea.Cmd
	m.textInput, tiCmd = m.textInput.Update(msg)
	return m, tiCmd
}

func (m Model) newView() string {
	return "  New worktree: " + m.textInput.View()
}
