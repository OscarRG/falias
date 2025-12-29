package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/oscar.rivas/falias/internal/config"
)

// Styles holds all the lipgloss styles for the UI
type Styles struct {
	// Header styles
	HeaderStyle     lipgloss.Style
	ShellInfoStyle  lipgloss.Style
	CountStyle      lipgloss.Style

	// Search bar styles
	SearchLabelStyle lipgloss.Style
	SearchInputStyle lipgloss.Style

	// List styles
	ListItemStyle     lipgloss.Style
	SelectedItemStyle lipgloss.Style
	AliasNameStyle    lipgloss.Style
	AliasValueStyle   lipgloss.Style
	AliasFileStyle    lipgloss.Style

	// Badge styles
	GlobalBadgeStyle      lipgloss.Style
	OverriddenBadgeStyle  lipgloss.Style
	MissingBadgeStyle     lipgloss.Style
	ConditionalBadgeStyle lipgloss.Style

	// Modal styles
	ModalBoxStyle    lipgloss.Style
	ModalTitleStyle  lipgloss.Style
	ModalLabelStyle  lipgloss.Style
	ModalValueStyle  lipgloss.Style
	ModalActiveStyle lipgloss.Style

	// Footer styles
	FooterStyle lipgloss.Style
	KeyStyle    lipgloss.Style
	HelpStyle   lipgloss.Style

	// Help modal styles
	HelpBoxStyle    lipgloss.Style
	HelpHeaderStyle lipgloss.Style
	HelpKeyStyle    lipgloss.Style
	HelpDescStyle   lipgloss.Style

	// Status messages
	StatusStyle      lipgloss.Style
	ErrorStatusStyle lipgloss.Style

	// Spinner style
	SpinnerStyle lipgloss.Style

	// Muted text style
	MutedStyle lipgloss.Style
}

// NewStyles creates a new Styles instance based on the theme
func NewStyles(theme config.Theme) Styles {
	return Styles{
		// Header styles
		HeaderStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary).
			Padding(0, 1),

		ShellInfoStyle: lipgloss.NewStyle().
			Foreground(theme.Secondary).
			Padding(0, 1),

		CountStyle: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Align(lipgloss.Right),

		// Search bar styles
		SearchLabelStyle: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true),

		SearchInputStyle: lipgloss.NewStyle().
			Foreground(theme.Foreground),

		// List styles
		ListItemStyle: lipgloss.NewStyle().
			Padding(0, 2),

		SelectedItemStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("231")).
			Foreground(lipgloss.Color("0")).
			Bold(true).
			Padding(0, 2),

		AliasNameStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Success),

		AliasValueStyle: lipgloss.NewStyle().
			Foreground(theme.Foreground),

		AliasFileStyle: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Italic(true),

		// Badge styles
		GlobalBadgeStyle: lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true).
			Background(theme.Secondary).
			Foreground(lipgloss.Color("0")),

		OverriddenBadgeStyle: lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true).
			Background(theme.Warning).
			Foreground(lipgloss.Color("0")),

		MissingBadgeStyle: lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true).
			Background(theme.Error).
			Foreground(lipgloss.Color("255")),

		ConditionalBadgeStyle: lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true).
			Background(theme.Muted).
			Foreground(lipgloss.Color("255")),

		// Modal styles
		ModalBoxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Primary).
			Padding(1, 2).
			Width(70),

		ModalTitleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary).
			Underline(true),

		ModalLabelStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Secondary),

		ModalValueStyle: lipgloss.NewStyle().
			Foreground(theme.Foreground),

		ModalActiveStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Success),

		// Footer styles
		FooterStyle: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Padding(0, 1),

		KeyStyle: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true),

		HelpStyle: lipgloss.NewStyle().
			Foreground(theme.Muted),

		// Help modal styles
		HelpBoxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Primary).
			Padding(1, 2).
			Width(60),

		HelpHeaderStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary).
			Underline(true).
			Align(lipgloss.Center),

		HelpKeyStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Success).
			Width(15),

		HelpDescStyle: lipgloss.NewStyle().
			Foreground(theme.Foreground),

		// Status messages
		StatusStyle: lipgloss.NewStyle().
			Foreground(theme.Success).
			Bold(true).
			Padding(0, 1),

		ErrorStatusStyle: lipgloss.NewStyle().
			Foreground(theme.Error).
			Bold(true).
			Padding(0, 1),

		// Spinner style
		SpinnerStyle: lipgloss.NewStyle().
			Foreground(theme.Primary),

		// Muted text style
		MutedStyle: lipgloss.NewStyle().
			Foreground(theme.Muted),
	}
}
