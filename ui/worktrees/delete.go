package worktrees

import (
	"fmt"

	"gx/git"
	"gx/ui"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// deleteResultMsg is sent when a delete operation completes.
type deleteResultMsg struct{ err error }

// cmdDelete removes the worktree directory and force-deletes its branch.
func cmdDelete(repo git.Repo, wt git.Worktree) tea.Cmd {
	return func() tea.Msg {
		if err := git.RemoveWorktree(repo, wt.Name, true); err != nil {
			return deleteResultMsg{err: err}
		}
		if wt.Branch != "" {
			if err := git.DeleteLocalBranch(repo, wt.Branch, true); err != nil {
				return deleteResultMsg{err: err}
			}
		}
		return deleteResultMsg{}
	}
}

// handleDeleteKey handles key events while in delete-confirmation mode.
func (m Model) handleDeleteKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	confirm := key.NewBinding(key.WithKeys("y"))
	cancel := key.NewBinding(key.WithKeys("n", "esc", "q"))

	switch {
	case key.Matches(msg, confirm):
		wt := m.selectedWorktree()
		if wt == nil {
			m.mode = modeNormal
			return m, nil
		}
		m.mode = modeNormal
		m.statusMsg = fmt.Sprintf("Deleting '%s'…", wt.Name)
		return m, cmdDelete(m.repo, *wt)

	case key.Matches(msg, cancel):
		m.mode = modeNormal
		m.statusMsg = ""
	}
	return m, nil
}

// deleteConfirmView returns the one-line status bar text for delete mode.
func (m Model) deleteConfirmView() string {
	wt := m.selectedWorktree()
	if wt == nil {
		return ""
	}
	prompt := fmt.Sprintf("  Delete '%s' (branch: %s)? ", wt.Name, wt.Branch)
	hint := ui.StyleBold.Render("[y/N]")
	return prompt + hint
}
