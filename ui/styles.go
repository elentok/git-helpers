package ui

import "charm.land/lipgloss/v2"

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

// Text styles
var (
	StyleBold = lipgloss.NewStyle().Bold(true)
	StyleDim  = lipgloss.NewStyle().Foreground(ColorGray)
)
