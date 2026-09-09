package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Brand Colors
	crimson     = lipgloss.Color("#f43f5e") // High-contrast vibrant rose/crimson for foregrounds & highlights
	crimsonDark = lipgloss.Color("#820024") // Deep burgundy for tab & badge backgrounds
	navy        = lipgloss.Color("#0f172a")
	slate       = lipgloss.Color("#334155")
	muted       = lipgloss.Color("#64748b")
	faint       = lipgloss.Color("#94a3b8")
	light       = lipgloss.Color("#f8f9fc")
	white       = lipgloss.Color("#ffffff")
	lime        = lipgloss.Color("#10b981")
	blue        = lipgloss.Color("#3b82f6")
	orange      = lipgloss.Color("#ff5722")

	// Base Layout Styles - Zero wasted vertical space
	DocStyle = lipgloss.NewStyle().
			Padding(0, 1)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Background(crimsonDark).
			Padding(0, 2)

	RootPathStyle = lipgloss.NewStyle().
			Foreground(faint).
			Italic(true)

	// Tab Styles
	ActiveTab = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Background(crimsonDark).
			Padding(0, 2).
			MarginRight(1)

	InactiveTab = lipgloss.NewStyle().
			Foreground(slate).
			Background(light).
			Padding(0, 2).
			MarginRight(1)

	// List & Card Styles - Compact padding, zero bottom margins
	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(muted).
			Padding(0, 1)

	ActiveCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(crimson).
			Padding(0, 1)

	ItemTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white)

	ItemDesc = lipgloss.NewStyle().
			Foreground(faint)

	SelectedItem = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Background(crimsonDark).
			Padding(0, 1)

	BadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Background(slate).
			Padding(0, 1).
			MarginLeft(1)

	ActiveBadgeStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(white).
				Background(lime).
				Padding(0, 1).
				MarginLeft(1)

	// Keybindings & Status Bar
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(faint).
			Background(navy).
			Padding(0, 1)

	KeyHelpStyle = lipgloss.NewStyle().
			Foreground(muted)

	KeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(crimson)

	SuccessStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lime)

	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ef4444"))
)
