package config

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Name        string
	Description string
	Primary     lipgloss.Color
	Secondary   lipgloss.Color
	Success     lipgloss.Color
	Warning     lipgloss.Color
	Error       lipgloss.Color
	Muted       lipgloss.Color
	Highlight   lipgloss.Color
	Background  lipgloss.Color
	Foreground  lipgloss.Color
}

// GetTheme returns the theme configuration for a given theme name
func GetTheme(name string) Theme {
	themes := map[string]Theme{
		"default": {
			Name:        "default",
			Description: "Default balanced theme",
			Primary:     lipgloss.Color("39"),  // Blue
			Secondary:   lipgloss.Color("86"),  // Cyan
			Success:     lipgloss.Color("42"),  // Green
			Warning:     lipgloss.Color("214"), // Orange
			Error:       lipgloss.Color("196"), // Red
			Muted:       lipgloss.Color("246"), // Medium gray
			Highlight:   lipgloss.Color("226"), // Yellow
			Background:  lipgloss.Color(""),    // Default terminal bg
			Foreground:  lipgloss.Color("255"), // White
		},
		"light": {
			Name:        "light",
			Description: "Light theme with high contrast",
			Primary:     lipgloss.Color("27"),  // Dark blue
			Secondary:   lipgloss.Color("30"),  // Dark cyan
			Success:     lipgloss.Color("28"),  // Dark green
			Warning:     lipgloss.Color("166"), // Dark orange
			Error:       lipgloss.Color("160"), // Dark red
			Muted:       lipgloss.Color("240"), // Dark gray
			Highlight:   lipgloss.Color("136"), // Dark yellow
			Background:  lipgloss.Color("231"), // White
			Foreground:  lipgloss.Color("16"),  // Black
		},
		"dark": {
			Name:        "dark",
			Description: "Dark theme with vibrant colors",
			Primary:     lipgloss.Color("75"),  // Bright blue
			Secondary:   lipgloss.Color("87"),  // Bright cyan
			Success:     lipgloss.Color("84"),  // Bright green
			Warning:     lipgloss.Color("222"), // Bright orange
			Error:       lipgloss.Color("204"), // Bright red
			Muted:       lipgloss.Color("243"), // Light gray
			Highlight:   lipgloss.Color("229"), // Bright yellow
			Background:  lipgloss.Color("0"),   // Black
			Foreground:  lipgloss.Color("255"), // White
		},
		"high-contrast": {
			Name:        "high-contrast",
			Description: "High contrast for better readability",
			Primary:     lipgloss.Color("51"),  // Very bright cyan
			Secondary:   lipgloss.Color("45"),  // Very bright cyan
			Success:     lipgloss.Color("46"),  // Very bright green
			Warning:     lipgloss.Color("226"), // Very bright yellow
			Error:       lipgloss.Color("196"), // Bright red
			Muted:       lipgloss.Color("250"), // Very light gray
			Highlight:   lipgloss.Color("11"),  // Bright yellow
			Background:  lipgloss.Color("16"),  // True black
			Foreground:  lipgloss.Color("231"), // True white
		},
		"nord": {
			Name:        "nord",
			Description: "Nord color palette",
			Primary:     lipgloss.Color("111"), // Nord blue
			Secondary:   lipgloss.Color("117"), // Nord frost
			Success:     lipgloss.Color("108"), // Nord green
			Warning:     lipgloss.Color("179"), // Nord yellow
			Error:       lipgloss.Color("131"), // Nord red
			Muted:       lipgloss.Color("246"), // Nord gray
			Highlight:   lipgloss.Color("180"), // Nord yellow light
			Background:  lipgloss.Color("234"), // Nord dark
			Foreground:  lipgloss.Color("253"), // Nord light
		},
		"gruvbox": {
			Name:        "gruvbox",
			Description: "Gruvbox color palette",
			Primary:     lipgloss.Color("109"), // Gruvbox blue
			Secondary:   lipgloss.Color("108"), // Gruvbox aqua
			Success:     lipgloss.Color("142"), // Gruvbox green
			Warning:     lipgloss.Color("214"), // Gruvbox orange
			Error:       lipgloss.Color("167"), // Gruvbox red
			Muted:       lipgloss.Color("245"), // Gruvbox gray
			Highlight:   lipgloss.Color("214"), // Gruvbox yellow
			Background:  lipgloss.Color("235"), // Gruvbox dark
			Foreground:  lipgloss.Color("223"), // Gruvbox light
		},
	}

	theme, ok := themes[name]
	if !ok {
		return themes["default"]
	}

	return theme
}

// ListThemes returns a formatted list of available themes
func ListThemes() string {
	themes := GetAvailableThemes()
	result := "Available themes:\n"
	for _, name := range themes {
		theme := GetTheme(name)
		result += fmt.Sprintf("  - %s: %s\n", name, theme.Description)
	}
	return result
}
