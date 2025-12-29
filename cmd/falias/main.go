package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oscar.rivas/falias/internal/config"
	"github.com/oscar.rivas/falias/internal/export"
	"github.com/oscar.rivas/falias/internal/scanner"
	"github.com/oscar.rivas/falias/internal/ui"
)

var (
	shellFlag      = flag.String("shell", "", "Force shell type (bash or zsh, auto-detect if not specified)")
	rootFlag       = flag.String("root", "", "Override starting file (default: ~/.bashrc or ~/.zshrc)")
	jsonFlag       = flag.Bool("json", false, "Export aliases as JSON and exit")
	debugFlag      = flag.Bool("debug", false, "Show includes graph and unresolved paths")
	helpFlag       = flag.Bool("help", false, "Show help")
	versionFlag    = flag.Bool("version", false, "Show version")
	themeFlag      = flag.String("theme", "", "Set the color theme and save to config")
	listThemesFlag = flag.Bool("list-themes", false, "List available themes")
)

const version = "1.0.0"

func main() {
	flag.Parse()

	if *helpFlag {
		showHelp()
		os.Exit(0)
	}

	if *versionFlag {
		fmt.Printf("falias v%s\n", version)
		os.Exit(0)
	}

	if *listThemesFlag {
		listThemes()
		os.Exit(0)
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to load config: %v\n", err)
		fmt.Fprintf(os.Stderr, "Using default configuration.\n")
		cfg = config.DefaultConfig()
	}

	// Set theme if flag is provided
	if *themeFlag != "" {
		if err := cfg.SetTheme(*themeFlag); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Theme set to '%s'\n", *themeFlag)
		configPath, _ := config.GetConfigPath()
		fmt.Printf("Saved to: %s\n", configPath)
		os.Exit(0)
	}

	// Detect or use specified shell
	shell := *shellFlag
	if shell == "" {
		shell = scanner.DetectShell()
	}

	// Validate shell
	if shell != "bash" && shell != "zsh" {
		fmt.Fprintf(os.Stderr, "Error: Invalid shell '%s'. Must be 'bash' or 'zsh'.\n", shell)
		os.Exit(1)
	}

	// Get root files
	var rootFiles []string
	if *rootFlag != "" {
		rootFiles = []string{*rootFlag}
	} else {
		allRoots := scanner.GetDefaultRootFiles(shell)
		rootFiles = scanner.FilterExistingFiles(allRoots)

		if len(rootFiles) == 0 {
			fmt.Fprintf(os.Stderr, "Error: No configuration files found for %s.\n", shell)
			fmt.Fprintf(os.Stderr, "Tried: %v\n", allRoots)
			os.Exit(1)
		}
	}

	// JSON export mode
	if *jsonFlag {
		exportJSON(shell, rootFiles)
		return
	}

	// Debug mode
	if *debugFlag {
		showDebug(shell, rootFiles)
		return
	}

	theme := config.GetTheme(cfg.Theme)
	runTUI(shell, rootFiles, theme)
}

// runTUI starts the Bubble Tea TUI
func runTUI(shell string, rootFiles []string, theme config.Theme) {
	m := ui.NewModel(shell, rootFiles, theme)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// exportJSON exports aliases as JSON
func exportJSON(shell string, rootFiles []string) {
	s := scanner.NewScanner()
	result, err := s.ScanShellFiles(shell, rootFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning: %v\n", err)
		os.Exit(1)
	}

	exporter := export.NewJSONExporter(true)
	if err := exporter.ExportAliases(result, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error exporting: %v\n", err)
		os.Exit(1)
	}
}

// showDebug shows debug information
func showDebug(shell string, rootFiles []string) {
	s := scanner.NewScanner()
	result, err := s.ScanShellFiles(shell, rootFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Shell: %s\n", result.Shell)
	fmt.Printf("Root Files: %v\n\n", result.RootFiles)

	fmt.Printf("Files Scanned (%d):\n", len(result.Files))
	for path, file := range result.Files {
		status := "OK"
		if !file.Exists {
			status = "MISSING"
		} else if !file.Readable {
			status = "PERMISSION DENIED"
		}
		conditional := ""
		if file.Conditional {
			conditional = " [CONDITIONAL]"
		}
		fmt.Printf("  %s - %s%s\n", path, status, conditional)
		if len(file.Includes) > 0 {
			fmt.Printf("    Includes: %v\n", file.Includes)
		}
		if len(file.Aliases) > 0 {
			fmt.Printf("    Aliases: %d\n", len(file.Aliases))
		}
	}

	fmt.Printf("\nAliases Found: %d\n", len(result.Aliases))

	if len(result.UnresolvedPaths) > 0 {
		fmt.Printf("\nUnresolved Paths (%d):\n", len(result.UnresolvedPaths))
		for _, path := range result.UnresolvedPaths {
			fmt.Printf("  %s\n", path)
		}
	}

	if len(result.Warnings) > 0 {
		fmt.Printf("\nWarnings (%d):\n", len(result.Warnings))
		for _, warning := range result.Warnings {
			fmt.Printf("  %s\n", warning)
		}
	}
}

// listThemes lists all available themes
func listThemes() {
	cfg, _ := config.Load()
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	fmt.Println("Available themes:\n")

	themes := config.GetAvailableThemes()
	for _, name := range themes {
		theme := config.GetTheme(name)
		prefix := "  "
		if name == cfg.Theme {
			prefix = "* "
		}
		fmt.Printf("%s%-15s - %s\n", prefix, name, theme.Description)
	}

	fmt.Printf("\nCurrent theme: %s\n", cfg.Theme)
	configPath, _ := config.GetConfigPath()
	fmt.Printf("Config file: %s\n", configPath)
	fmt.Printf("\nTo change theme: falias --theme <name>\n")
}

// showHelp displays help information
func showHelp() {
	fmt.Printf(`falias v%s - Shell Alias Discovery TUI

USAGE:
  falias [flags]

FLAGS:
  --shell bash|zsh    Force shell type (auto-detect if not specified)
  --root <path>       Override starting file (default: ~/.bashrc or ~/.zshrc)
  --json              Export aliases as JSON and exit
  --debug             Show includes graph and unresolved paths
  --theme <name>      Set color theme (use --list-themes to see options)
  --list-themes       List available color themes
  --help              Show this help
  --version           Show version

KEYBOARD SHORTCUTS (in TUI):
  ↑/↓ or j/k          Navigate list
  /                   Focus search bar
  Enter               View alias details
  c                   Copy alias value to clipboard
  n                   Copy alias name to clipboard
  p                   Copy full alias definition
  t                   Toggle view mode (All/By File/Overridden/Globals)
  r                   Rescan configuration files
  h or ?              Show help
  q or Ctrl+C         Quit
  Esc                 Close modal or unfocus search

EXAMPLES:
  falias                      # Auto-detect shell and scan
  falias --shell zsh          # Force zsh
  falias --root ~/.zshrc      # Use specific file
  falias --json               # Export as JSON
  falias --json | jq '.'      # Pretty-print JSON
  falias --debug              # Show debug info
  falias --list-themes        # List available themes
  falias --theme gruvbox      # Set theme to gruvbox

CONFIGURATION:
  Config file: ~/.config/falias/config.yaml
  Theme setting is persistent across sessions

DOCUMENTATION:
  See README.md for detailed information about how falias works.

`, version)
}
