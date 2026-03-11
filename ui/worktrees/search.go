package worktrees

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// enterSearchMode transitions the model into search mode with an empty query.
func (m Model) enterSearchMode() Model {
	ti := textinput.New()
	ti.Prompt = ""
	ti.Focus()
	m.mode = modeSearch
	m.textInput = ti
	m.searchQuery = ""
	m.searchMatches = nil
	m.searchCursor = 0
	m.statusMsg = ""
	return m
}

// exitSearchMode clears search state and returns to normal mode.
func (m Model) exitSearchMode() Model {
	m.mode = modeNormal
	m.searchQuery = ""
	m.searchMatches = nil
	m.searchCursor = 0
	m.table.SetRows(m.buildRows())
	return m
}

// recomputeSearchMatches rebuilds the searchMatches slice from the current query.
func (m Model) recomputeSearchMatches() Model {
	q := strings.ToLower(m.searchQuery)
	m.searchMatches = nil
	if q == "" {
		return m
	}
	for i, wt := range m.worktrees {
		if strings.Contains(strings.ToLower(wt.Name), q) || strings.Contains(strings.ToLower(wt.Branch), q) {
			m.searchMatches = append(m.searchMatches, i)
		}
	}
	return m
}

// handleSearchKey handles key events while in search mode.
func (m Model) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	nextMatch := key.NewBinding(key.WithKeys("ctrl+n"))
	prevMatch := key.NewBinding(key.WithKeys("ctrl+p"))
	exit := key.NewBinding(key.WithKeys("esc", "enter"))

	switch {
	case key.Matches(msg, exit):
		m = m.exitSearchMode()
		return m, nil

	case key.Matches(msg, nextMatch):
		if len(m.searchMatches) > 0 {
			m.searchCursor = (m.searchCursor + 1) % len(m.searchMatches)
			return m.jumpToSearchCursor()
		}
		return m, nil

	case key.Matches(msg, prevMatch):
		if len(m.searchMatches) > 0 {
			m.searchCursor = (m.searchCursor - 1 + len(m.searchMatches)) % len(m.searchMatches)
			return m.jumpToSearchCursor()
		}
		return m, nil
	}

	var tiCmd tea.Cmd
	m.textInput, tiCmd = m.textInput.Update(msg)
	m.searchQuery = m.textInput.Value()
	m = m.recomputeSearchMatches()
	if len(m.searchMatches) > 0 {
		m.searchCursor = 0
		return m.jumpToSearchCursor()
	}
	m.table.SetRows(m.buildRows())
	return m, tiCmd
}

// jumpToSearchCursor moves the table cursor to the current search match and
// returns the sidebar-reload command.
func (m Model) jumpToSearchCursor() (Model, tea.Cmd) {
	idx := m.searchMatches[m.searchCursor]
	m.table.SetCursor(idx)
	m.table.SetRows(m.buildRows())
	m.sidebarLoading = true
	m.sidebarUpstream = ""
	m.sidebarAheadCommits = nil
	m.sidebarBehindCommits = nil
	m.sidebarChanges = nil
	m.viewport.SetContent(m.sidebarContent())
	return m, cmdLoadSidebarData(m.repo, m.worktrees[idx])
}

// searchView returns the one-line status bar text for search mode.
func (m Model) searchView() string {
	query := m.textInput.View()
	if m.searchQuery != "" && len(m.searchMatches) == 0 {
		return "  Search: " + query + "  no matches"
	}
	if len(m.searchMatches) > 0 {
		return fmt.Sprintf("  Search: %s  %d/%d", query, m.searchCursor+1, len(m.searchMatches))
	}
	return "  Search: " + query
}
