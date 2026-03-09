package worktrees

import (
	"fmt"
	"path/filepath"

	"gx/git"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// pasteResultMsg is sent when a paste operation completes.
type pasteResultMsg struct {
	n   int // number of files pasted
	err error
}

// handlePasteModeKey handles key events in paste mode (clipboard active, waiting for destination).
// Only navigation (j/k/arrows), paste (p), and cancel (esc) are active.
func (m Model) handlePasteModeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	up := key.NewBinding(key.WithKeys("up", "k"))
	down := key.NewBinding(key.WithKeys("down", "j"))
	doPaste := key.NewBinding(key.WithKeys("p"))
	cancel := key.NewBinding(key.WithKeys("esc", "q"))

	switch {
	case key.Matches(msg, cancel):
		m.clipboard = nil
		m.mode = modeNormal
		return m, nil
	case key.Matches(msg, doPaste):
		if m.clipboard != nil {
			wt := m.selectedWorktree()
			if wt != nil {
				m.statusMsg = "Pasting…"
				m.mode = modeNormal
				return m, cmdPaste(*m.clipboard, *wt)
			}
		}
		m.mode = modeNormal
		return m, nil
	case key.Matches(msg, up), key.Matches(msg, down):
		prevCursor := m.table.Cursor()
		var tableCmd tea.Cmd
		m.table, tableCmd = m.table.Update(msg)
		if m.table.Cursor() != prevCursor && len(m.worktrees) > 0 {
			m.table.SetRows(buildRows(m.worktrees, m.statuses, m.dirties, m.table.Cursor(), m.settings.UseNerdFontIcons))
			m.sidebarLoading = true
			m.sidebarUpstream = ""
			m.sidebarAheadCommits = nil
			m.sidebarBehindCommits = nil
			m.sidebarChanges = nil
			m.viewport.SetContent(m.sidebarContent())
			return m, tea.Batch(tableCmd, cmdLoadSidebarData(m.repo, m.worktrees[m.table.Cursor()]))
		}
		return m, tableCmd
	}
	return m, nil
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
