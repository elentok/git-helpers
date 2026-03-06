package worktrees

import (
	"fmt"
	"path/filepath"

	"gx/git"

	tea "github.com/charmbracelet/bubbletea"
)

// pasteResultMsg is sent when a paste operation completes.
type pasteResultMsg struct {
	n   int // number of files pasted
	err error
}

// cmdPaste copies clipboard files from their source into the destination worktree.
func cmdPaste(cb clipboardState, dst git.Worktree) tea.Cmd {
	return func() tea.Msg {
		for _, relPath := range cb.files {
			src := filepath.Join(cb.srcPath, relPath)
			dstPath := filepath.Join(dst.Path, relPath)
			if err := copyPath(src, dstPath); err != nil {
				return pasteResultMsg{err: fmt.Errorf("copy %s: %w", relPath, err)}
			}
		}
		return pasteResultMsg{n: len(cb.files)}
	}
}
