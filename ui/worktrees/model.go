package worktrees

import (
	"gx/git"
	"gx/ui"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mode int

const (
	modeNormal mode = iota
	modeDelete
	modeRename
	modeClone
	modeError
)

// ── messages ─────────────────────────────────────────────────────────────────

type worktreesLoadedMsg struct {
	worktrees []git.Worktree
	err       error
}

type syncStatusMsg struct {
	branch string
	status git.SyncStatus
}

type sidebarDataMsg struct {
	worktreePath string
	commits      []git.Commit
	changes      []git.Change
}

// ── commands ──────────────────────────────────────────────────────────────────

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

func cmdLoadSidebarData(repo git.Repo, wt git.Worktree) tea.Cmd {
	return func() tea.Msg {
		commits, _ := git.CommitsSinceMain(repo, wt.Branch)
		changes, _ := git.UncommittedChanges(wt.Path)
		return sidebarDataMsg{worktreePath: wt.Path, commits: commits, changes: changes}
	}
}

// ── model ─────────────────────────────────────────────────────────────────────

// Model is the BubbleTea model for the worktrees page.
type Model struct {
	repo               git.Repo
	activeWorktreePath string // path of the worktree the user launched from

	worktrees []git.Worktree
	statuses  map[string]git.SyncStatus

	table    table.Model
	viewport viewport.Model

	sidebarCommits []git.Commit
	sidebarChanges []git.Change
	sidebarLoading bool

	mode         mode
	textInput    textinput.Model // shared by rename and clone modes
	statusMsg    string
	errorViewport viewport.Model

	width  int
	height int
	ready  bool // true once we've received the first WindowSizeMsg

	loading bool
	err     error
}

// New creates a new worktrees page model. activeWorktreePath is the path of the
// worktree the user is currently in (empty if launched from the bare repo root).
func New(repo git.Repo, activeWorktreePath string) Model {
	return Model{
		repo:               repo,
		activeWorktreePath: activeWorktreePath,
		statuses:           make(map[string]git.SyncStatus),
		table:              newTable(),
		loading:            true,
	}
}

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
		// Mode-specific handling first
		switch m.mode {
		case modeError:
			return m.handleErrorKey(msg)
		case modeDelete:
			return m.handleDeleteKey(msg)
		case modeRename:
			return m.handleRenameKey(msg)
		case modeClone:
			return m.handleCloneKey(msg)
		}
		// Normal mode
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Delete) && len(m.worktrees) > 0:
			m.mode = modeDelete
			m.statusMsg = ""
			return m, nil
		case key.Matches(msg, keys.Rename) && len(m.worktrees) > 0:
			m = m.enterRenameMode()
			return m, nil
		case key.Matches(msg, keys.Clone) && len(m.worktrees) > 0:
			m = m.enterCloneMode()
			return m, nil
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
		m.table.SetRows(buildRows(m.worktrees, m.statuses))

		// Position cursor on the active worktree
		for i, wt := range m.worktrees {
			if wt.Path == m.activeWorktreePath {
				m.table.SetCursor(i)
				break
			}
		}

		// Load sync status for all worktrees and initial sidebar data in parallel
		for _, wt := range m.worktrees {
			if wt.Branch != "" {
				cmds = append(cmds, cmdLoadSyncStatus(m.repo, wt.Branch))
			}
		}
		if len(m.worktrees) > 0 {
			m.sidebarLoading = true
			m.viewport.SetContent(m.sidebarContent())
			cmds = append(cmds, cmdLoadSidebarData(m.repo, m.worktrees[m.table.Cursor()]))
		}
		return m, tea.Batch(cmds...)

	case syncStatusMsg:
		m.statuses[msg.branch] = msg.status
		m.table.SetRows(buildRows(m.worktrees, m.statuses))
		return m, nil

	case sidebarDataMsg:
		// Discard stale results (user may have moved the cursor)
		if len(m.worktrees) > 0 && m.worktrees[m.table.Cursor()].Path == msg.worktreePath {
			m.sidebarCommits = msg.commits
			m.sidebarChanges = msg.changes
			m.sidebarLoading = false
			m.viewport.SetContent(m.sidebarContent())
		}
		return m, nil
	}

	// Pass remaining messages to table and viewport; detect cursor changes
	prevCursor := m.table.Cursor()

	var tableCmd tea.Cmd
	m.table, tableCmd = m.table.Update(msg)
	cmds = append(cmds, tableCmd)

	if m.table.Cursor() != prevCursor && len(m.worktrees) > 0 {
		m.sidebarLoading = true
		m.sidebarCommits = nil
		m.sidebarChanges = nil
		m.viewport.SetContent(m.sidebarContent())
		cmds = append(cmds, cmdLoadSidebarData(m.repo, m.worktrees[m.table.Cursor()]))
	}

	var vpCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, vpCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing…"
	}
	if m.err != nil {
		return "\n  Error: " + m.err.Error()
	}
	if m.mode == modeError {
		return m.errorModalView()
	}

	_, sidebarW := m.splitWidth()
	innerSidebarW := sidebarW - 2 // rounded border adds 1 on each side
	innerSidebarH := m.contentHeight() - 2

	tableView := m.table.View()

	sidebarView := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorBorder).
		Width(innerSidebarW).
		Height(innerSidebarH).
		Render(m.viewport.View())

	content := lipgloss.JoinHorizontal(lipgloss.Top, tableView, sidebarView)
	return lipgloss.JoinVertical(lipgloss.Left, content, m.statusBarView())
}

// ── layout helpers ────────────────────────────────────────────────────────────

func (m Model) splitWidth() (tableW, sidebarW int) {
	tableW = int(float64(m.width) * 0.55)
	sidebarW = m.width - tableW
	return
}

func (m Model) contentHeight() int {
	h := m.height - 1 // reserve 1 line for status bar
	if h < 4 {
		return 4
	}
	return h
}

// selectedWorktree returns a pointer to the currently highlighted worktree, or nil.
func (m Model) selectedWorktree() *git.Worktree {
	if len(m.worktrees) == 0 {
		return nil
	}
	w := m.worktrees[m.table.Cursor()]
	return &w
}

// statusBarView renders the 1-line bar at the bottom of the screen.
func (m Model) statusBarView() string {
	switch m.mode {
	case modeError:
		return ""
	case modeDelete:
		return m.deleteConfirmView()
	case modeRename:
		return m.renameView()
	case modeClone:
		return m.cloneView()
	default:
		if m.statusMsg != "" {
			return "  " + m.statusMsg
		}
		return ui.StyleDim.Render("  d delete  r rename  c clone  q quit")
	}
}

func (m Model) resized() Model {
	tableW, sidebarW := m.splitWidth()
	h := m.contentHeight()

	resizeTable(&m.table, tableW, h)

	vpW := sidebarW - 2
	vpH := h - 2
	if vpW < 1 {
		vpW = 1
	}
	if vpH < 1 {
		vpH = 1
	}
	m.viewport.Width = vpW
	m.viewport.Height = vpH
	m.viewport.SetContent(m.sidebarContent())

	return m
}

func (m Model) sidebarContent() string {
	var wt *git.Worktree
	if len(m.worktrees) > 0 {
		w := m.worktrees[m.table.Cursor()]
		wt = &w
	}
	return renderSidebarContent(wt, m.sidebarCommits, m.sidebarChanges, m.sidebarLoading)
}

// ── error modal ───────────────────────────────────────────────────────────────

// showError switches the model into error mode with a scrollable viewport.
func (m Model) showError(errMsg string) Model {
	vpW := m.width * 2 / 3
	if vpW < 40 {
		vpW = 40
	}
	if vpW > 80 {
		vpW = 80
	}
	// overhead: border (2) + title line (1) + blank (1) + blank (1) + hint (1) = 6
	vpH := m.height/2 - 6
	if vpH < 3 {
		vpH = 3
	}
	vp := viewport.New(vpW-2, vpH) // -2 for left/right border
	vp.SetContent(errMsg)
	m.errorViewport = vp
	m.mode = modeError
	return m
}

// handleErrorKey scrolls the error viewport or dismisses it.
func (m Model) handleErrorKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc, tea.KeyEnter:
		m.mode = modeNormal
		return m, nil
	}
	if msg.Type == tea.KeyRunes && msg.String() == "q" {
		m.mode = modeNormal
		return m, nil
	}
	var cmd tea.Cmd
	m.errorViewport, cmd = m.errorViewport.Update(msg)
	return m, cmd
}

// errorModalView renders a centred modal with the error text.
func (m Model) errorModalView() string {
	titleStyle := lipgloss.NewStyle().Foreground(ui.ColorRed).Bold(true)
	hintStyle := ui.StyleDim
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorRed).
		Padding(0, 1).
		Width(m.errorViewport.Width)

	inner := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Error"),
		"",
		m.errorViewport.View(),
		"",
		hintStyle.Render("esc / enter / q  to dismiss"),
	)
	modal := borderStyle.Render(inner)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
}
