package main

import (
	"fmt"
	"os"

	"gx/git"
	"gx/ui/worktrees"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	repo, err := git.FindRepo(cwd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	// Detect which worktree the user launched from, if any.
	var activeWorktreePath string
	if info, err := git.IdentifyDir(cwd); err == nil {
		activeWorktreePath = info.WorktreeRoot
	}

	m := worktrees.New(*repo, activeWorktreePath)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
