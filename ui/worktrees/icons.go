package worktrees

type uiIcons struct {
	worktreePrefix string
	mainPrefix     string
	branchPrefix   string
	worktreeTitle  string
	aheadTitle     string
	behindTitle    string
	baseTitle      string
	rebasedYes     string
	rebasedNo      string
	changesTitle   string
}

func icons(useNerdFont bool) uiIcons {
	if !useNerdFont {
		return uiIcons{
			worktreeTitle: "Worktree",
			aheadTitle:    "Commits ahead of remote",
			behindTitle:   "Commits behind remote",
			baseTitle:     "Base",
			rebasedYes:    "✓",
			rebasedNo:     "✗",
			changesTitle:  "Changes",
		}
	}
	return uiIcons{
		worktreePrefix: "󰉋 ",
		mainPrefix:     "󰋜 ",
		branchPrefix:   " ",
		worktreeTitle:  "󰙅 Worktree",
		aheadTitle:     " Commits ahead of remote",
		behindTitle:    " Commits behind remote",
		baseTitle:      " Base",
		rebasedYes:     "\uf00c", // fa-check
		rebasedNo:      "\uf00d", // fa-times
		changesTitle:   "󰈔 Changes",
	}
}
