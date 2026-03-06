package worktrees

type uiIcons struct {
	worktreePrefix string
	branchPrefix   string
	worktreeTitle  string
	aheadTitle     string
	behindTitle    string
	changesTitle   string
}

func icons(useNerdFont bool) uiIcons {
	if !useNerdFont {
		return uiIcons{
			worktreeTitle: "Worktree",
			aheadTitle:    "Commits ahead of main",
			behindTitle:   "Commits behind main",
			changesTitle:  "Changes",
		}
	}
	return uiIcons{
		worktreePrefix: "󰉋 ",
		branchPrefix:   " ",
		worktreeTitle:  "󰙅 Worktree",
		aheadTitle:     " Commits ahead of main",
		behindTitle:    " Commits behind main",
		changesTitle:   "󰈔 Changes",
	}
}
