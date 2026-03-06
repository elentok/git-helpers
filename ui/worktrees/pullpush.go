package worktrees

import (
	"gx/git"

	tea "github.com/charmbracelet/bubbletea"
)

type pullResultMsg struct{ err error }
type pushResultMsg struct{ err error }

func cmdPull(wt git.Worktree) tea.Cmd {
	return func() tea.Msg {
		return pullResultMsg{err: git.Pull(wt.Path)}
	}
}

func cmdPush(wt git.Worktree) tea.Cmd {
	return func() tea.Msg {
		return pushResultMsg{err: git.Push(wt.Path)}
	}
}
