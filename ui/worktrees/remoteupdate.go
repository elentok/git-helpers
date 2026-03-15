package worktrees

import (
	"gx/git"

	tea "github.com/charmbracelet/bubbletea"
)

type remoteUpdateResultMsg struct {
	err error
	log string
}

func cmdRemoteUpdate(repo git.Repo) tea.Cmd {
	return func() tea.Msg {
		out, err := git.UpdateRemotes(repo)
		return remoteUpdateResultMsg{err: err, log: out}
	}
}
