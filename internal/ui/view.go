package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/oscar.rivas/falias/internal/config"
	"github.com/oscar.rivas/falias/internal/model"
)

// View renders the TUI
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	if m.scanning {
		return m.renderScanning()
	}

	if m.showHelp {
		return m.renderHelp()
	}

	if m.showDetails {
		return m.renderDetails()
	}

	if m.showThemePicker {
		return m.renderThemePicker()
	}

	return m.renderMain()
}

// renderScanning renders the scanning state
func (m Model) renderScanning() string {
	var s strings.Builder

	s.WriteString("\n\n")
	s.WriteString(m.styles.SpinnerStyle.Render(m.spinner.View()))
	s.WriteString(" Scanning aliases...\n\n")
	s.WriteString(m.styles.MutedStyle.Render(fmt.Sprintf("Shell: %s\n", m.shell)))
	for _, file := range m.rootFiles {
		s.WriteString(m.styles.MutedStyle.Render(fmt.Sprintf("  %s\n", file)))
	}

	return s.String()
}

// renderMain renders the main interface
func (m Model) renderMain() string {
	header := m.renderHeader()
	search := m.renderSearch()
	list := m.renderList()
	footer := m.renderFooter()
	status := m.renderStatus()

	// Calculate available height for list
	headerHeight := lipgloss.Height(header)
	searchHeight := lipgloss.Height(search)
	footerHeight := lipgloss.Height(footer)
	statusHeight := lipgloss.Height(status)

	listHeight := m.height - headerHeight - searchHeight - footerHeight - statusHeight - 4

	// Assemble the view
	var s strings.Builder
	s.WriteString(header)
	s.WriteString("\n")
	s.WriteString(search)
	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Height(listHeight).Render(list))
	s.WriteString("\n")
	s.WriteString(footer)
	if status != "" {
		s.WriteString("\n")
		s.WriteString(status)
	}

	return s.String()
}

// renderHeader renders the header
func (m Model) renderHeader() string {
	title := m.styles.HeaderStyle.Render("falias")

	shellInfo := ""
	if m.scanResult != nil && len(m.scanResult.RootFiles) > 0 {
		shellInfo = m.styles.ShellInfoStyle.Render(fmt.Sprintf("(%s: %s)",
			m.scanResult.Shell,
			filepath.Base(m.scanResult.RootFiles[0])))
	}

	count := ""
	if m.scanResult != nil {
		count = m.styles.CountStyle.Render(fmt.Sprintf("%d/%d aliases",
			len(m.displayedAliases),
			len(m.allAliases)))
	}

	viewMode := m.styles.ShellInfoStyle.Render(fmt.Sprintf("[%s]", m.viewMode.String()))

	left := lipgloss.JoinHorizontal(lipgloss.Left, title, shellInfo, viewMode)
	right := count

	// Calculate spacing
	spacer := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if spacer < 0 {
		spacer = 0
	}

	return lipgloss.JoinHorizontal(lipgloss.Left,
		left,
		strings.Repeat(" ", spacer),
		right)
}

// renderSearch renders the search bar
func (m Model) renderSearch() string {
	label := m.styles.SearchLabelStyle.Render("Search: ")
	input := m.searchInput.View()

	return label + input
}

// renderList renders the alias list
func (m Model) renderList() string {
	if len(m.displayedAliases) == 0 {
		return m.styles.MutedStyle.Render("No aliases found")
	}

	var s strings.Builder

	// Calculate how many items we can show
	maxItems := m.height - 10 // Rough calculation
	if maxItems < 5 {
		maxItems = 5
	}

	// Calculate scroll window
	start := m.cursor - maxItems/2
	if start < 0 {
		start = 0
	}
	end := start + maxItems
	if end > len(m.displayedAliases) {
		end = len(m.displayedAliases)
		start = end - maxItems
		if start < 0 {
			start = 0
		}
	}

	for i := start; i < end; i++ {
		alias := m.displayedAliases[i]
		line := m.renderListItem(alias, i == m.cursor)
		s.WriteString(line)
		s.WriteString("\n")
	}

	return s.String()
}

// renderListItem renders a single list item
func (m Model) renderListItem(alias *model.AliasEntry, selected bool) string {
	// Name
	name := m.styles.AliasNameStyle.Render(alias.Name)

	// Value (truncated)
	maxValueLen := 40
	value := alias.ActiveValue
	if len(value) > maxValueLen {
		value = value[:maxValueLen-3] + "..."
	}
	valueStr := m.styles.AliasValueStyle.Render(value)

	// Badges
	var badges []string
	if alias.Type == model.AliasTypeGlobal {
		badges = append(badges, m.styles.GlobalBadgeStyle.Render("global"))
	}
	if alias.IsOverridden {
		badges = append(badges, m.styles.OverriddenBadgeStyle.Render("overridden"))
	}

	// File info
	file := m.styles.AliasFileStyle.Render(filepath.Base(alias.ActiveLocation.FilePath))

	// Combine
	content := fmt.Sprintf("%-20s %-45s %s %s",
		name,
		valueStr,
		strings.Join(badges, " "),
		file)

	if selected {
		return m.styles.SelectedItemStyle.Render("▸ " + content)
	}
	return m.styles.ListItemStyle.Render("  " + content)
}

// renderFooter renders the footer with keybindings
func (m Model) renderFooter() string {
	keys := []string{
		m.styles.KeyStyle.Render("↑/↓") + ":nav",
		m.styles.KeyStyle.Render("/") + ":search",
		m.styles.KeyStyle.Render("⏎") + ":details",
		m.styles.KeyStyle.Render("c") + ":copy",
		m.styles.KeyStyle.Render("t") + ":toggle",
		m.styles.KeyStyle.Render("T") + ":theme",
		m.styles.KeyStyle.Render("q") + ":quit",
		m.styles.KeyStyle.Render("h") + ":?",
	}

	return m.styles.FooterStyle.Render(strings.Join(keys, " "))
}

// renderStatus renders status/error messages
func (m Model) renderStatus() string {
	if m.errorMessage != "" {
		msg := m.errorMessage
		m.errorMessage = "" // Clear after showing
		return m.styles.ErrorStatusStyle.Render(msg)
	}

	if m.statusMessage != "" {
		msg := m.statusMessage
		m.statusMessage = "" // Clear after showing
		return m.styles.StatusStyle.Render(msg)
	}

	return ""
}

// renderDetails renders the details modal
func (m Model) renderDetails() string {
	alias := m.getCurrentAlias()
	if alias == nil {
		return "No alias selected"
	}

	var content strings.Builder

	// Title
	content.WriteString(m.styles.ModalTitleStyle.Render("Alias Details"))
	content.WriteString("\n\n")

	// Name
	content.WriteString(m.styles.ModalLabelStyle.Render("Name: "))
	content.WriteString(m.styles.ModalValueStyle.Render(alias.Name))
	content.WriteString("\n")

	// Type
	content.WriteString(m.styles.ModalLabelStyle.Render("Type: "))
	content.WriteString(m.styles.ModalValueStyle.Render(string(alias.Type)))
	content.WriteString("\n")

	// Value
	content.WriteString(m.styles.ModalLabelStyle.Render("Value: "))
	content.WriteString(m.styles.ModalValueStyle.Render(alias.ActiveValue))
	content.WriteString("\n\n")

	// Definitions
	if len(alias.Definitions) == 1 {
		content.WriteString(m.styles.ModalLabelStyle.Render("Defined in:"))
		content.WriteString("\n")
		loc := alias.ActiveLocation
		content.WriteString("  " + m.styles.ModalValueStyle.Render(fmt.Sprintf("%s:%d", loc.FilePath, loc.LineNum)))
		content.WriteString("\n\n")
		content.WriteString(m.styles.HelpStyle.Render(fmt.Sprintf("Open: code -g %s:%d", loc.FilePath, loc.LineNum)))
	} else {
		content.WriteString(m.styles.ModalLabelStyle.Render("Definition History:"))
		content.WriteString("\n")
		for i, def := range alias.Definitions {
			isActive := i == len(alias.Definitions)-1
			marker := fmt.Sprintf("%d.", i+1)
			valueStr := def.Value
			if isActive {
				marker = m.styles.ModalActiveStyle.Render(marker)
				valueStr = m.styles.ModalActiveStyle.Render(valueStr + " [ACTIVE]")
			}
			content.WriteString(fmt.Sprintf("  %s %s\n", marker, valueStr))
			content.WriteString(fmt.Sprintf("     %s\n",
				m.styles.MutedStyle.Render(fmt.Sprintf("%s:%d", def.Location.FilePath, def.Location.LineNum))))
		}
	}

	content.WriteString("\n")
	content.WriteString(m.styles.HelpStyle.Render("[ESC to close]"))

	box := m.styles.ModalBoxStyle.Render(content.String())

	// Center the modal
	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		box)
}

// renderHelp renders the help modal
func (m Model) renderHelp() string {
	var content strings.Builder

	content.WriteString(m.styles.HelpHeaderStyle.Render("Keyboard Shortcuts"))
	content.WriteString("\n\n")

	shortcuts := []struct {
		key  string
		desc string
	}{
		{"↑/↓ or j/k", "Navigate list"},
		{"/", "Focus search bar"},
		{"Enter", "View alias details"},
		{"c", "Copy alias value to clipboard"},
		{"n", "Copy alias name to clipboard"},
		{"p", "Copy full alias definition"},
		{"t", "Toggle view mode (All/By File/Overridden/Globals)"},
		{"r", "Rescan configuration files"},
		{"h or ?", "Show this help"},
		{"q or Ctrl+C", "Quit"},
		{"Esc", "Close modal or unfocus search"},
	}

	for _, sc := range shortcuts {
		content.WriteString(m.styles.HelpKeyStyle.Render(sc.key))
		content.WriteString(" ")
		content.WriteString(m.styles.HelpDescStyle.Render(sc.desc))
		content.WriteString("\n")
	}

	content.WriteString("\n")
	content.WriteString(m.styles.HelpStyle.Render("[Press any key to close]"))

	box := m.styles.HelpBoxStyle.Render(content.String())

	// Center the modal
	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		box)
}

// renderThemePicker renders the theme picker modal
func (m Model) renderThemePicker() string {
	var content strings.Builder

	content.WriteString(m.styles.HelpHeaderStyle.Render("Theme Picker"))
	content.WriteString("\n\n")
	content.WriteString(m.styles.HelpStyle.Render("Navigate with ↑/↓ to preview, Enter to save, Esc to cancel"))
	content.WriteString("\n\n")

	// Render theme list
	for i, themeName := range m.themeList {
		theme := config.GetTheme(themeName)

		isSelected := i == m.themeCursor

		selectionStyle := lipgloss.NewStyle().
			Background(lipgloss.Color("231")).
			Foreground(lipgloss.Color("0")).
			Bold(true).
			Padding(0, 1)

		prefix := "  "
		if isSelected {
			prefix = "▸ "
		}

		themeLine := fmt.Sprintf("%-15s", themeName)
		description := theme.Description

		if isSelected {
			// High contrast: bright yellow background with black text
			themeLine = selectionStyle.Render(prefix + themeLine)
			line := fmt.Sprintf("%s - %s", themeLine, description)
			content.WriteString(line)
		} else {
			// Normal style
			themeLine = m.styles.HelpKeyStyle.Render(themeName)
			line := fmt.Sprintf("%s%-15s - %s", prefix, themeLine, description)
			content.WriteString(line)
		}

		content.WriteString("\n")
	}

	content.WriteString("\n")
	content.WriteString(m.styles.HelpStyle.Render("Current selection previewed in real-time"))

	box := m.styles.HelpBoxStyle.Copy().Width(70).Render(content.String())

	// Center the modal
	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		box)
}
