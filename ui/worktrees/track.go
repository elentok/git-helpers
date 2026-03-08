package worktrees

import (
	"fmt"

	"gx/git"
	"gx/ui"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type trackResultMsg struct{ err error }

func cmdTrack(repo git.Repo, wt git.Worktree) tea.Cmd {
	return func() tea.Msg {
		remote := git.BranchRemote(repo, wt.Branch)
		return trackResultMsg{err: git.TrackRemote(repo.Root, remote, wt.Branch)}
	}
}

func (m Model) handleTrackKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		m.spinnerActive = true
		m.spinnerLabel = "Tracking " + wt.Name + "…"
		return m, tea.Batch(cmdTrack(m.repo, *wt), m.spinner.Tick)

	case key.Matches(msg, cancel):
		m.mode = modeNormal
		m.statusMsg = ""
	}
	return m, nil
}

func (m Model) trackConfirmView() string {
	wt := m.selectedWorktree()
	if wt == nil {
		return ""
	}
	remote := git.BranchRemote(m.repo, wt.Branch)
	prompt := fmt.Sprintf("  Track %s/%s? ", remote, wt.Branch)
	hint := ui.StyleBold.Render("[y/N]")
	return prompt + hint
}
