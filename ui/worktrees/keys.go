package worktrees

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Delete key.Binding
	Rename key.Binding
	Clone  key.Binding
	Yank   key.Binding
	Paste  key.Binding
	Pull   key.Binding
	Push   key.Binding
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
	Paste: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "paste files"),
	),
	// Pull and Push are chained sequences (gpl / gps) handled manually;
	// they have no WithKeys so key.Matches never fires, but they appear in help.
	Pull: key.NewBinding(
		key.WithHelp("gpl", "pull"),
	),
	Push: key.NewBinding(
		key.WithHelp("gps", "push"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// ShortHelp and FullHelp implement help.KeyMap so this can be wired to
// bubbles/help in a later milestone.

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Delete, k.Rename, k.Clone, k.Yank, k.Paste, k.Pull, k.Push, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Delete, k.Rename, k.Clone, k.Yank, k.Paste, k.Pull, k.Push, k.Quit},
	}
}
