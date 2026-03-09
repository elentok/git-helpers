package worktrees

import (
	"fmt"
	"os/exec"
	"runtime"

	"gx/git"

	tea "github.com/charmbracelet/bubbletea"
)

type pullResultMsg struct{ err error }
type pushResultMsg struct {
	err   error
	prURL string
}
type forcePushResultMsg struct{ err error }
type urlOpenedMsg struct{}

func cmdPull(wt git.Worktree) tea.Cmd {
	return func() tea.Msg {
		return pullResultMsg{err: git.Pull(wt.Path)}
	}
}

func cmdPush(repo git.Repo, wt git.Worktree) tea.Cmd {
	return func() tea.Msg {
		remote := git.BranchRemote(repo, wt.Branch)
		prURL, err := git.PushBranch(wt.Path, remote, wt.Branch)
		return pushResultMsg{err: err, prURL: prURL}
	}
}

func cmdForcePush(repo git.Repo, wt git.Worktree) tea.Cmd {
	return func() tea.Msg {
		remote := git.BranchRemote(repo, wt.Branch)
		return forcePushResultMsg{err: git.PushBranchForce(wt.Path, remote, wt.Branch)}
	}
}

func cmdOpenURL(url string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", url)
		default:
			cmd = exec.Command("xdg-open", url)
		}
		_ = cmd.Start()
		return urlOpenedMsg{}
	}
}

func forcePushPrompt(wt git.Worktree) string {
	return fmt.Sprintf("Push rejected for %s. Force push?", wt.Branch)
}
