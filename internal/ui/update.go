package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oscar.rivas/falias/internal/config"
	"github.com/oscar.rivas/falias/internal/model"
)

// Update handles all messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.ready {
			m.ready = true
		}
		return m, nil

	case scanCompleteMsg:
		m.scanning = false
		m.scanResult = msg.result

		// Build sorted alias list
		m.allAliases = make([]*model.AliasEntry, 0, len(msg.result.Aliases))
		for _, alias := range msg.result.Aliases {
			m.allAliases = append(m.allAliases, alias)
		}

		// Sort by name
		sort.Slice(m.allAliases, func(i, j int) bool {
			return strings.ToLower(m.allAliases[i].Name) < strings.ToLower(m.allAliases[j].Name)
		})

		m.filterAliases()
		m.statusMessage = fmt.Sprintf("Found %d aliases", len(m.allAliases))
		return m, nil

	case scanErrorMsg:
		m.scanning = false
		m.errorMessage = fmt.Sprintf("Error scanning: %v", msg.err)
		return m, nil

	case spinner.TickMsg:
		if m.scanning {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil

	default:
		// Update text input if search is focused
		if m.searchFocused {
			var cmd tea.Cmd
			m.searchInput, cmd = m.searchInput.Update(msg)
			m.filterAliases() // Re-filter when search changes
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// handleKeyPress handles keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Global keys that work in any state
	switch {
	case msg.String() == "ctrl+c":
		return m, tea.Quit
	}

	// Help modal
	if m.showHelp {
		if msg.String() == "esc" || msg.String() == "q" || msg.String() == "h" || msg.String() == "?" {
			m.showHelp = false
		}
		return m, nil
	}

	// Details modal
	if m.showDetails {
		if msg.String() == "esc" || msg.String() == "enter" || msg.String() == "q" {
			m.showDetails = false
		}
		return m, nil
	}

	// Theme picker modal
	if m.showThemePicker {
		switch msg.String() {
		case "esc":
			// Revert to original theme
			m.applyTheme(m.originalTheme)
			m.showThemePicker = false
		case "enter":
			// Save the theme
			cfg, _ := config.Load()
			if cfg == nil {
				cfg = config.DefaultConfig()
			}
			if err := cfg.SetTheme(m.currentThemeName); err == nil {
				m.statusMessage = fmt.Sprintf("Theme '%s' saved", m.currentThemeName)
			} else {
				m.errorMessage = fmt.Sprintf("Failed to save theme: %v", err)
			}
			m.showThemePicker = false
		case "up", "k":
			if m.themeCursor > 0 {
				m.themeCursor--
				m.applyTheme(m.themeList[m.themeCursor])
			}
		case "down", "j":
			if m.themeCursor < len(m.themeList)-1 {
				m.themeCursor++
				m.applyTheme(m.themeList[m.themeCursor])
			}
		}
		return m, nil
	}

	// Search input focused
	if m.searchFocused {
		switch msg.String() {
		case "esc":
			m.searchFocused = false
			m.searchInput.Blur()
			return m, nil
		case "enter":
			m.searchFocused = false
			m.searchInput.Blur()
			return m, nil
		default:
			var cmd tea.Cmd
			m.searchInput, cmd = m.searchInput.Update(msg)
			m.filterAliases()
			return m, cmd
		}
	}

	// Main list navigation
	switch {
	case msg.String() == "q":
		return m, tea.Quit

	case msg.String() == "up" || msg.String() == "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case msg.String() == "down" || msg.String() == "j":
		if m.cursor < len(m.displayedAliases)-1 {
			m.cursor++
		}

	case msg.String() == "/":
		m.searchFocused = true
		m.searchInput.Focus()
		return m, textinput.Blink

	case msg.String() == "enter":
		if len(m.displayedAliases) > 0 {
			m.showDetails = true
		}

	case msg.String() == "c":
		// Copy value
		if alias := m.getCurrentAlias(); alias != nil {
			if err := clipboard.WriteAll(alias.ActiveValue); err == nil {
				m.statusMessage = "Copied value to clipboard"
			} else {
				m.errorMessage = "Failed to copy to clipboard"
			}
		}

	case msg.String() == "n":
		// Copy name
		if alias := m.getCurrentAlias(); alias != nil {
			if err := clipboard.WriteAll(alias.Name); err == nil {
				m.statusMessage = "Copied name to clipboard"
			} else {
				m.errorMessage = "Failed to copy to clipboard"
			}
		}

	case msg.String() == "p":
		// Copy full alias definition
		if alias := m.getCurrentAlias(); alias != nil {
			fullDef := fmt.Sprintf("alias %s='%s'", alias.Name, alias.ActiveValue)
			if err := clipboard.WriteAll(fullDef); err == nil {
				m.statusMessage = "Copied full definition to clipboard"
			} else {
				m.errorMessage = "Failed to copy to clipboard"
			}
		}

	case msg.String() == "t":
		m.cycleViewMode()
		m.statusMessage = fmt.Sprintf("View: %s", m.viewMode.String())

	case msg.String() == "T":
		// Open theme picker
		m.originalTheme = m.currentThemeName
		m.showThemePicker = true

	case msg.String() == "r":
		// Rescan
		m.scanning = true
		m.statusMessage = ""
		m.errorMessage = ""
		cmds = append(cmds, m.spinner.Tick, scanAliasesCmd(m.scanner, m.shell, m.rootFiles))

	case msg.String() == "h" || msg.String() == "?":
		m.showHelp = true
	}

	return m, tea.Batch(cmds...)
}
