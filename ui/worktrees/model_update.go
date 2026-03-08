package worktrees

import (
	"fmt"

	"gx/ui/components"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return cmdLoadWorktrees(m.repo)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m = m.resized()
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		switch m.mode {
		case modeError:
			return m.handleErrorKey(msg)
		case modeDelete:
			return m.handleDeleteKey(msg)
		case modeTrack:
			return m.handleTrackKey(msg)
		case modeRename:
			return m.handleRenameKey(msg)
		case modeClone:
			return m.handleCloneKey(msg)
		case modeNew:
			return m.handleNewKey(msg)
		case modeYank:
			return m.handleYankKey(msg)
		}
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			m = m.resized()
			return m, nil
		case key.Matches(msg, keys.New) && !m.spinnerActive:
			m = m.enterNewMode()
			return m, nil
		case key.Matches(msg, keys.Delete) && len(m.worktrees) > 0 && !m.spinnerActive:
			m.mode = modeDelete
			m.statusMsg = ""
			return m, nil
		case key.Matches(msg, keys.Rename) && len(m.worktrees) > 0 && !m.spinnerActive:
			m = m.enterRenameMode()
			return m, nil
		case key.Matches(msg, keys.Clone) && len(m.worktrees) > 0 && !m.spinnerActive:
			m = m.enterCloneMode()
			return m, nil
		case key.Matches(msg, keys.Yank) && len(m.worktrees) > 0 && !m.spinnerActive:
			return m.enterYankMode()
		case key.Matches(msg, keys.Paste) && m.clipboard != nil && len(m.worktrees) > 0 && !m.spinnerActive:
			wt := m.selectedWorktree()
			if wt != nil {
				m.statusMsg = "Pasting…"
				return m, cmdPaste(*m.clipboard, *wt)
			}
		case key.Matches(msg, keys.Pull) && len(m.worktrees) > 0 && !m.spinnerActive:
			wt := m.selectedWorktree()
			if wt != nil {
				m.spinnerActive = true
				m.spinnerLabel = "Pulling " + wt.Name + "…"
				return m, tea.Batch(cmdPull(*wt), m.spinner.Tick)
			}
		case key.Matches(msg, keys.Push) && len(m.worktrees) > 0 && !m.spinnerActive:
			wt := m.selectedWorktree()
			if wt != nil {
				if wt.Branch == "" {
					return m.showError("cannot push: worktree is in detached HEAD state"), nil
				}
				m.spinnerActive = true
				m.spinnerLabel = "Pushing " + wt.Name + "…"
				return m, tea.Batch(cmdPush(m.repo, *wt), m.spinner.Tick)
			}
		case key.Matches(msg, keys.Track) && len(m.worktrees) > 0 && !m.spinnerActive && m.sidebarUpstream == "":
			wt := m.selectedWorktree()
			if wt != nil {
				if wt.Branch == "" {
					return m.showError("cannot track: worktree is in detached HEAD state"), nil
				}
				m.mode = modeTrack
				m.statusMsg = ""
			}
		}

	case deleteResultMsg:
		if msg.err != nil {
			return m.showError(msg.err.Error()), nil
		}
		m.statusMsg = ""
		return m, cmdLoadWorktrees(m.repo)

	case renameResultMsg:
		if msg.err != nil {
			return m.showError(msg.err.Error()), nil
		}
		m.statusMsg = ""
		return m, cmdLoadWorktrees(m.repo)

	case cloneResultMsg:
		if msg.err != nil {
			return m.showError(msg.err.Error()), nil
		}
		m.statusMsg = ""
		return m, cmdLoadWorktrees(m.repo)

	case newResultMsg:
		if msg.err != nil {
			return m.showError(msg.err.Error()), nil
		}
		m.statusMsg = ""
		return m, cmdLoadWorktrees(m.repo)

	case yankDataMsg:
		if m.mode != modeYank || msg.worktreePath != m.yankSource.Path {
			return m, nil
		}
		if msg.err != nil {
			return m.showError(msg.err.Error()), nil
		}
		m.yankLoading = false
		m.yankChecklist = components.NewChecklist(changesToChecklistItems(msg.changes))
		return m, nil

	case clearStatusMsg:
		if msg.gen == m.statusGen {
			m.statusMsg = ""
		}
		return m, nil

	case spinner.TickMsg:
		if m.spinnerActive {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil

	case pullResultMsg:
		m.spinnerActive = false
		if msg.err != nil {
			return m.showError(msg.err.Error()), nil
		}
		m.statusGen++
		m.statusMsg = "Pulled"
		cmds = append(cmds, cmdClearStatus(m.statusGen))
		if wt := m.selectedWorktree(); wt != nil && wt.Branch != "" {
			cmds = append(cmds, cmdLoadSyncStatus(m.repo, wt.Branch))
		}
		return m, tea.Batch(cmds...)

	case pushResultMsg:
		m.spinnerActive = false
		if msg.err != nil {
			return m.showError(msg.err.Error()), nil
		}
		m.statusGen++
		m.statusMsg = "Pushed"
		cmds = append(cmds, cmdClearStatus(m.statusGen))
		if wt := m.selectedWorktree(); wt != nil && wt.Branch != "" {
			cmds = append(cmds, cmdLoadSyncStatus(m.repo, wt.Branch))
		}
		return m, tea.Batch(cmds...)

	case trackResultMsg:
		m.spinnerActive = false
		if msg.err != nil {
			return m.showError(msg.err.Error()), nil
		}
		m.statusGen++
		m.statusMsg = "Tracking remote branch"
		cmds = append(cmds, cmdClearStatus(m.statusGen))
		if wt := m.selectedWorktree(); wt != nil && wt.Branch != "" {
			m.sidebarLoading = true
			m.viewport.SetContent(m.sidebarContent())
			cmds = append(cmds, cmdLoadSyncStatus(m.repo, wt.Branch), cmdLoadSidebarData(m.repo, *wt))
		}
		return m, tea.Batch(cmds...)

	case pasteResultMsg:
		if msg.err != nil {
			return m.showError(msg.err.Error()), nil
		}
		m.clipboard = nil
		m.statusGen++
		m.statusMsg = fmt.Sprintf("Pasted %d file(s)", msg.n)
		clearCmd := cmdClearStatus(m.statusGen)
		if wt := m.selectedWorktree(); wt != nil {
			m.sidebarLoading = true
			m.viewport.SetContent(m.sidebarContent())
			return m, tea.Batch(clearCmd, cmdLoadSidebarData(m.repo, *wt))
		}
		return m, clearCmd

	case worktreesLoadedMsg:
		m.loading = false
		if msg.err != nil {
			if m.ready {
				return m.showError(msg.err.Error()), nil
			}
			m.err = msg.err
			return m, nil
		}
		m.worktrees = msg.worktrees
		m.dirties = make(map[string]dirtyState)
		m.table.SetRows(buildRows(m.worktrees, m.statuses, m.dirties, m.table.Cursor(), m.settings.UseNerdFontIcons))

		for i, wt := range m.worktrees {
			if wt.Path == m.activeWorktreePath {
				m.table.SetCursor(i)
				break
			}
		}

		for _, wt := range m.worktrees {
			if wt.Branch != "" {
				cmds = append(cmds, cmdLoadSyncStatus(m.repo, wt.Branch))
			}
			cmds = append(cmds, cmdLoadDirtyStatus(wt))
		}
		if len(m.worktrees) > 0 {
			m.sidebarLoading = true
			m.viewport.SetContent(m.sidebarContent())
			cmds = append(cmds, cmdLoadSidebarData(m.repo, m.worktrees[m.table.Cursor()]))
		}
		return m, tea.Batch(cmds...)

	case syncStatusMsg:
		m.statuses[msg.branch] = msg.status
		m.table.SetRows(buildRows(m.worktrees, m.statuses, m.dirties, m.table.Cursor(), m.settings.UseNerdFontIcons))
		return m, nil

	case dirtyStatusMsg:
		m.dirties[msg.worktreePath] = msg.dirty
		m.table.SetRows(buildRows(m.worktrees, m.statuses, m.dirties, m.table.Cursor(), m.settings.UseNerdFontIcons))
		return m, nil

	case sidebarDataMsg:
		if len(m.worktrees) > 0 && m.worktrees[m.table.Cursor()].Path == msg.worktreePath {
			m.sidebarUpstream = msg.upstream
			m.sidebarAheadCommits = msg.aheadCommits
			m.sidebarBehindCommits = msg.behindCommits
			m.sidebarChanges = msg.changes
			m.sidebarLoading = false
			m.viewport.SetContent(m.sidebarContent())
		}
		return m, nil
	}

	prevCursor := m.table.Cursor()

	var tableCmd tea.Cmd
	m.table, tableCmd = m.table.Update(msg)
	cmds = append(cmds, tableCmd)

	if m.table.Cursor() != prevCursor && len(m.worktrees) > 0 {
		m.table.SetRows(buildRows(m.worktrees, m.statuses, m.dirties, m.table.Cursor(), m.settings.UseNerdFontIcons))
		m.sidebarLoading = true
		m.sidebarUpstream = ""
		m.sidebarAheadCommits = nil
		m.sidebarBehindCommits = nil
		m.sidebarChanges = nil
		m.viewport.SetContent(m.sidebarContent())
		cmds = append(cmds, cmdLoadSidebarData(m.repo, m.worktrees[m.table.Cursor()]))
	}

	var vpCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, vpCmd)

	return m, tea.Batch(cmds...)
}
