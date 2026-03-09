package worktrees

import (
	"fmt"

	"gx/git"

	tea "github.com/charmbracelet/bubbletea"
)

type pullResultMsg struct{ err error }
type pushResultMsg struct{ err error }
type forcePushResultMsg struct{ err error }

func cmdPull(wt git.Worktree) tea.Cmd {
	return func() tea.Msg {
		return pullResultMsg{err: git.Pull(wt.Path)}
	}
}

func cmdPush(repo git.Repo, wt git.Worktree) tea.Cmd {
	return func() tea.Msg {
		remote := git.BranchRemote(repo, wt.Branch)
		return pushResultMsg{err: git.PushBranch(wt.Path, remote, wt.Branch)}
	}
}

func cmdForcePush(repo git.Repo, wt git.Worktree) tea.Cmd {
	return func() tea.Msg {
		remote := git.BranchRemote(repo, wt.Branch)
		return forcePushResultMsg{err: git.PushBranchForce(wt.Path, remote, wt.Branch)}
	}
}

func forcePushPrompt(wt git.Worktree) string {
	return fmt.Sprintf("Push rejected for %s. Force push?", wt.Branch)
}
