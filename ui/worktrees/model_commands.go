package worktrees

import (
	"time"

	"gx/git"

	tea "github.com/charmbracelet/bubbletea"
)

func cmdClearStatus(gen int) tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
		return clearStatusMsg{gen: gen}
	})
}

func cmdLoadWorktrees(repo git.Repo) tea.Cmd {
	return func() tea.Msg {
		wts, err := git.ListWorktrees(repo)
		return worktreesLoadedMsg{worktrees: wts, err: err}
	}
}

func cmdLoadSyncStatus(repo git.Repo, branch string) tea.Cmd {
	return func() tea.Msg {
		status, _ := git.WorktreeSyncStatus(repo, branch)
		return syncStatusMsg{branch: branch, status: status}
	}
}

func cmdLoadDirtyStatus(wt git.Worktree) tea.Cmd {
	return func() tea.Msg {
		changes, _ := git.UncommittedChanges(wt.Path)
		return dirtyStatusMsg{
			worktreePath: wt.Path,
			dirty:        dirtyStateFromChanges(changes),
		}
	}
}

func cmdLoadSidebarData(repo git.Repo, wt git.Worktree) tea.Cmd {
	return func() tea.Msg {
		aheadCommits, _ := git.CommitsSinceMain(repo, wt.Branch)
		behindCommits, _ := git.CommitsBehindMain(repo, wt.Branch)
		changes, _ := git.UncommittedChanges(wt.Path)
		return sidebarDataMsg{
			worktreePath:  wt.Path,
			aheadCommits:  aheadCommits,
			behindCommits: behindCommits,
			changes:       changes,
		}
	}
}
