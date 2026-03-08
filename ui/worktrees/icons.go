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
			aheadTitle:    "Commits ahead of remote",
			behindTitle:   "Commits behind remote",
			changesTitle:  "Changes",
		}
	}
	return uiIcons{
		worktreePrefix: "󰉋 ",
		branchPrefix:   " ",
		worktreeTitle:  "󰙅 Worktree",
		aheadTitle:     " Commits ahead of remote",
		behindTitle:    " Commits behind remote",
		changesTitle:   "󰈔 Changes",
	}
}
