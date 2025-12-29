package ui

import "github.com/charmbracelet/bubbles/key"

// keyMap defines all keybindings for the application
type keyMap struct {
	Up        key.Binding
	Down      key.Binding
	Enter     key.Binding
	Search    key.Binding
	Copy      key.Binding
	CopyName  key.Binding
	CopyFull  key.Binding
	Toggle      key.Binding
	ThemePicker key.Binding
	Rescan      key.Binding
	Help        key.Binding
	Quit        key.Binding
	Escape      key.Binding
}

// defaultKeyMap returns the default keybindings
func defaultKeyMap() keyMap {
	return keyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("⏎", "details"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		Copy: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "copy value"),
		),
		CopyName: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "copy name"),
		),
		CopyFull: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "copy full"),
		),
		Toggle: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "toggle view"),
		),
		ThemePicker: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "theme picker"),
		),
		Rescan: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "rescan"),
		),
		Help: key.NewBinding(
			key.WithKeys("h", "?"),
			key.WithHelp("h/?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "close/unfocus"),
		),
	}
}
