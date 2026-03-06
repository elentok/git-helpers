package ui

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	ColorGreen   = lipgloss.Color("2")
	ColorYellow  = lipgloss.Color("3")
	ColorCyan    = lipgloss.Color("6")
	ColorMagenta = lipgloss.Color("5")
	ColorRed     = lipgloss.Color("1")
	ColorGray    = lipgloss.Color("8")
	ColorBorder  = lipgloss.Color("240")
)

// Status styles
var (
	StyleStatusSynced   = lipgloss.NewStyle().Foreground(ColorGreen)
	StyleStatusAhead    = lipgloss.NewStyle().Foreground(ColorMagenta)
	StyleStatusBehind   = lipgloss.NewStyle().Foreground(ColorYellow)
	StyleStatusDiverged = lipgloss.NewStyle().Foreground(ColorRed)
	StyleStatusUnknown  = lipgloss.NewStyle().Foreground(ColorGray)
)

// Dirty-state styles
var (
	StyleDirtyClean     = lipgloss.NewStyle().Foreground(ColorGray)
	StyleDirtyModified  = lipgloss.NewStyle().Foreground(ColorYellow)
	StyleDirtyUntracked = lipgloss.NewStyle().Foreground(ColorCyan)
	StyleDirtyMixed     = lipgloss.NewStyle().Foreground(ColorMagenta)
)

// Text styles
var (
	StyleBold = lipgloss.NewStyle().Bold(true)
	StyleDim  = lipgloss.NewStyle().Foreground(ColorGray)
)
