package worktrees

import (
	"os/exec"

	"gx/git"

	tea "github.com/charmbracelet/bubbletea"
)

type lazygitFinishedMsg struct{ err error }

func cmdLazygit(wt git.Worktree) tea.Cmd {
	c := exec.Command("lazygit", "-p", wt.Path)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return lazygitFinishedMsg{err: err}
	})
}
