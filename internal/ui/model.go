package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oscar.rivas/falias/internal/config"
	"github.com/oscar.rivas/falias/internal/model"
	"github.com/oscar.rivas/falias/internal/scanner"
)

// ViewMode represents the current view filtering mode
type ViewMode int

const (
	ViewAll ViewMode = iota
	ViewByFile
	ViewOverridden
	ViewGlobals
)

func (v ViewMode) String() string {
	switch v {
	case ViewAll:
		return "All"
	case ViewByFile:
		return "By File"
	case ViewOverridden:
		return "Overridden"
	case ViewGlobals:
		return "Globals"
	default:
		return "All"
	}
}

// Model represents the Bubble Tea model for the TUI
type Model struct {
	// Data
	scanResult *model.ScanResult
	scanner    *scanner.Scanner

	// UI State
	ready           bool
	scanning        bool
	width           int
	height          int
	cursor          int
	viewMode        ViewMode
	searchFocused   bool
	showDetails     bool
	showHelp        bool
	showThemePicker bool
	statusMessage   string
	errorMessage    string

	// Theme selection
	themeList          []string
	themeCursor        int
	originalTheme      string
	currentThemeName   string

	// Filtered/displayed aliases
	displayedAliases []*model.AliasEntry
	allAliases       []*model.AliasEntry

	// Components
	searchInput textinput.Model
	spinner     spinner.Model
	keys        keyMap
	styles      Styles

	// Config
	shell     string
	rootFiles []string
}

// NewModel creates a new TUI model with a theme
func NewModel(shell string, rootFiles []string, theme config.Theme) Model {
	// Create text input for search
	ti := textinput.New()
	ti.Placeholder = "Search aliases..."
	ti.CharLimit = 100

	// Create spinner
	s := spinner.New()
	s.Spinner = spinner.Dot

	// Create styles from theme
	styles := NewStyles(theme)
	s.Style = styles.SpinnerStyle

	return Model{
		scanner:          scanner.NewScanner(),
		shell:            shell,
		rootFiles:        rootFiles,
		searchInput:      ti,
		spinner:          s,
		keys:             defaultKeyMap(),
		styles:           styles,
		viewMode:         ViewAll,
		displayedAliases: make([]*model.AliasEntry, 0),
		allAliases:       make([]*model.AliasEntry, 0),
		scanning:         true,
		themeList:        config.GetAvailableThemes(),
		currentThemeName: theme.Name,
		themeCursor:      findThemeIndex(theme.Name, config.GetAvailableThemes()),
	}
}

// findThemeIndex finds the index of a theme in the theme list
func findThemeIndex(themeName string, themes []string) int {
	for i, name := range themes {
		if name == themeName {
			return i
		}
	}
	return 0
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		scanAliasesCmd(m.scanner, m.shell, m.rootFiles),
	)
}

// scanAliasesCmd creates a command to scan aliases
func scanAliasesCmd(s *scanner.Scanner, shell string, rootFiles []string) tea.Cmd {
	return func() tea.Msg {
		result, err := s.ScanShellFiles(shell, rootFiles)
		if err != nil {
			return scanErrorMsg{err: err}
		}
		return scanCompleteMsg{result: result}
	}
}

// scanCompleteMsg is sent when scanning completes
type scanCompleteMsg struct {
	result *model.ScanResult
}

// scanErrorMsg is sent when scanning fails
type scanErrorMsg struct {
	err error
}

// windowSizeMsg is sent when the terminal is resized
type windowSizeMsg struct {
	width  int
	height int
}

// filterAliases filters the alias list based on search and view mode
func (m *Model) filterAliases() {
	if m.scanResult == nil {
		m.displayedAliases = make([]*model.AliasEntry, 0)
		return
	}

	// Start with all aliases
	filtered := m.allAliases

	// Apply view mode filter
	switch m.viewMode {
	case ViewOverridden:
		temp := make([]*model.AliasEntry, 0)
		for _, alias := range filtered {
			if alias.IsOverridden {
				temp = append(temp, alias)
			}
		}
		filtered = temp

	case ViewGlobals:
		temp := make([]*model.AliasEntry, 0)
		for _, alias := range filtered {
			if alias.Type == model.AliasTypeGlobal {
				temp = append(temp, alias)
			}
		}
		filtered = temp

	case ViewByFile:
		// TODO: Implement grouping by file
		// For now, just show all
		break

	case ViewAll:
		// No additional filtering
		break
	}

	// Apply search filter
	searchTerm := strings.ToLower(strings.TrimSpace(m.searchInput.Value()))
	if searchTerm != "" {
		temp := make([]*model.AliasEntry, 0)
		for _, alias := range filtered {
			if strings.Contains(strings.ToLower(alias.Name), searchTerm) ||
				strings.Contains(strings.ToLower(alias.ActiveValue), searchTerm) {
				temp = append(temp, alias)
			}
		}
		filtered = temp
	}

	m.displayedAliases = filtered

	// Adjust cursor if needed
	if m.cursor >= len(m.displayedAliases) {
		m.cursor = len(m.displayedAliases) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

// getCurrentAlias returns the currently selected alias
func (m *Model) getCurrentAlias() *model.AliasEntry {
	if len(m.displayedAliases) == 0 || m.cursor >= len(m.displayedAliases) {
		return nil
	}
	return m.displayedAliases[m.cursor]
}

// cycleViewMode cycles to the next view mode
func (m *Model) cycleViewMode() {
	m.viewMode = (m.viewMode + 1) % 4
	m.filterAliases()
	m.cursor = 0
}

// applyTheme applies a theme to the model
func (m *Model) applyTheme(themeName string) {
	theme := config.GetTheme(themeName)
	m.styles = NewStyles(theme)
	m.spinner.Style = m.styles.SpinnerStyle
	m.currentThemeName = themeName
}
