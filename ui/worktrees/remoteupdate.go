package worktrees

import (
	"gx/git"

	tea "github.com/charmbracelet/bubbletea"
)

type remoteUpdateResultMsg struct{ err error }

func cmdRemoteUpdate(repo git.Repo) tea.Cmd {
	return func() tea.Msg {
		return remoteUpdateResultMsg{err: git.UpdateRemotes(repo)}
	}
}
