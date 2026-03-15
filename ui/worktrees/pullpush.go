package worktrees

import (
	"fmt"
	"os/exec"
	"runtime"

	"gx/git"

	tea "github.com/charmbracelet/bubbletea"
)

type pullResultMsg struct {
	err error
	log string
}
type pushResultMsg struct {
	err   error
	prURL string
	log   string
}
type forcePushResultMsg struct {
	err error
	log string
}
type urlOpenedMsg struct{}

func cmdPull(wt git.Worktree) tea.Cmd {
	return func() tea.Msg {
		out, err := git.Pull(wt.Path)
		return pullResultMsg{err: err, log: out}
	}
}

func cmdPush(repo git.Repo, wt git.Worktree) tea.Cmd {
	return func() tea.Msg {
		remote := git.BranchRemote(repo, wt.Branch)
		prURL, out, err := git.PushBranch(wt.Path, remote, wt.Branch)
		return pushResultMsg{err: err, prURL: prURL, log: out}
	}
}

func cmdForcePush(repo git.Repo, wt git.Worktree) tea.Cmd {
	return func() tea.Msg {
		remote := git.BranchRemote(repo, wt.Branch)
		out, err := git.PushBranchForce(wt.Path, remote, wt.Branch)
		return forcePushResultMsg{err: err, log: out}
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
