package worktrees

import (
	"fmt"

	"gx/git"

	tea "github.com/charmbracelet/bubbletea"
)

type trackResultMsg struct{ err error }

func cmdTrack(repo git.Repo, wt git.Worktree) tea.Cmd {
	return func() tea.Msg {
		remote := git.BranchRemote(repo, wt.Branch)
		return trackResultMsg{err: git.TrackRemote(repo.Root, remote, wt.Branch)}
	}
}

func (m Model) enterTrackConfirm() Model {
	wt := m.selectedWorktree()
	if wt == nil {
		return m
	}
	remote := git.BranchRemote(m.repo, wt.Branch)
	prompt := fmt.Sprintf("Track %s/%s?\n\nThis will run:\n  git branch --set-upstream-to=%s/%s %s",
		remote, wt.Branch, remote, wt.Branch, wt.Branch)
	return m.enterConfirm(prompt, cmdTrack(m.repo, *wt), "Tracking "+wt.Name+"…")
}
