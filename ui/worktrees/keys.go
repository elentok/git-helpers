package worktrees

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	New    key.Binding
	Delete key.Binding
	Rename key.Binding
	Clone  key.Binding
	Yank   key.Binding
	Pull   key.Binding
	Push   key.Binding
	Lazygit       key.Binding
	Track         key.Binding
	Refresh       key.Binding
	RemoteUpdate  key.Binding
	Help    key.Binding
	Quit   key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new worktree"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename"),
	),
	Clone: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "clone"),
	),
	Yank: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "yank files"),
	),
	Pull: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "pull"),
	),
	Push: key.NewBinding(
		key.WithKeys("P"),
		key.WithHelp("P", "push"),
	),
	Lazygit: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "lazygit"),
	),
	Track: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "track"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "refresh"),
	),
	RemoteUpdate: key.NewBinding(
		key.WithKeys("U"),
		key.WithHelp("U", "remote update"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// ShortHelp and FullHelp implement help.KeyMap so this can be wired to
// bubbles/help in a later milestone.

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.New, k.Delete, k.Rename, k.Clone, k.Yank, k.Pull, k.Push, k.Lazygit, k.Track, k.Refresh, k.Quit, k.Help}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.New, k.Delete, k.Rename, k.Clone},
		{k.Yank},
		{k.Pull, k.Push, k.Lazygit, k.Track, k.Refresh, k.RemoteUpdate, k.Help, k.Quit},
	}
}
