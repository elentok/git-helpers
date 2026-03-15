package worktrees

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type tmuxResultMsg struct{ err error }

func cmdTmuxNewSession(name, path string) tea.Cmd {
	return func() tea.Msg {
		if err := exec.Command("tmux", "new-session", "-d", "-s", name, "-c", path).Run(); err != nil {
			return tmuxResultMsg{err: err}
		}
		err := exec.Command("tmux", "switch-client", "-t", name).Run()
		return tmuxResultMsg{err: err}
	}
}

func cmdTmuxNewWindow(name, path string) tea.Cmd {
	return func() tea.Msg {
		err := exec.Command("tmux", "new-window", "-n", name, "-c", path).Run()
		return tmuxResultMsg{err: err}
	}
}
