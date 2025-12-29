# falias - Shell Alias Discovery TUI

A production-ready terminal UI application that discovers and displays all shell aliases defined across your bash/zsh configuration files, recursively following sourced files.

## Features

- **Complete Discovery**: Finds all aliases across `.bashrc`, `.zshrc`, and sourced files
- **Recursive Scanning**: Follows `source` and `.` commands to discover aliases in imported files
- **Smart Path Resolution**: Handles `~`, `$HOME`, `${HOME}`, and `$XDG_CONFIG_HOME` expansions
- **Themeable Interface**: 6 built-in themes with live preview (default, light, dark, high-contrast, nord, gruvbox)
- **Override Detection**: Shows when aliases are redefined with full definition history
- **Loop Prevention**: Detects and prevents infinite include loops
- **Read-Only**: Never modifies your configuration files
- **Safe Parsing**: No shell execution - pure static analysis
- **Beautiful TUI**: Built with Bubble Tea for a smooth terminal experience
- **Search & Filter**: Fuzzy search and multiple view modes
- **Copy to Clipboard**: Quick copy of alias names, values, or full definitions

## Installation

### Homebrew (Recommended)

```bash
brew tap OscarRG/tap
brew install falias
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/OscarRG/falias.git
cd falias

# Install dependencies
go mod download

# Build the binary
go build -o bin/falias ./cmd/falias

# Optionally, install to your PATH
sudo cp bin/falias /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/oscar.rivas/falias/cmd/falias@latest
```

## Quick Start

```bash
# Auto-detect shell and scan
falias

# Force specific shell
falias --shell zsh

# Use a specific theme
falias --theme nord

# List available themes
falias --list-themes

# Export to JSON
falias --json

# Show debug info
falias --debug
```

## Usage

### Command-Line Flags

```
falias [flags]

Flags:
  --shell bash|zsh    Force shell type (auto-detect if not specified)
  --root <path>       Override starting file (default: ~/.bashrc or ~/.zshrc)
  --theme <name>      Set color theme (default, light, dark, high-contrast, nord, gruvbox)
  --list-themes       List available themes and exit
  --json              Export aliases as JSON and exit
  --debug             Show includes graph and unresolved paths
  --help              Show help
  --version           Show version
```

### Configuration

falias stores your preferences in `~/.config/falias/config.yaml`:

```yaml
theme: nord # Your selected theme
```

You can edit this file manually or change themes from within the TUI.

### Keyboard Shortcuts (TUI Mode)

| Key             | Action                                            |
| --------------- | ------------------------------------------------- |
| `↑/↓` or `j/k`  | Navigate list                                     |
| `/`             | Focus search bar                                  |
| `Enter`         | View alias details                                |
| `c`             | Copy alias value to clipboard                     |
| `n`             | Copy alias name to clipboard                      |
| `p`             | Copy full alias definition                        |
| `t`             | Toggle view mode (All/By File/Overridden/Globals) |
| `T`             | Open theme picker with live preview               |
| `r`             | Rescan configuration files                        |
| `h` or `?`      | Show help                                         |
| `q` or `Ctrl+C` | Quit                                              |
| `Esc`           | Close modal or unfocus search                     |

## Examples

### Basic Usage

```bash
# Start the TUI
falias

# Search for git aliases
# (Press '/', type 'git')

# Change theme interactively
# (Press 'T', use arrow keys to preview, Enter to save)
```

### JSON Export

```bash
# Export all aliases as JSON
falias --json > aliases.json

# Pretty-print with jq
falias --json | jq '.'

# Filter for specific aliases
falias --json | jq '.[] | select(.name | startswith("git"))'
```

### Debug Mode

```bash
# Show which files are scanned and includes graph
falias --debug
```

### Custom Root File

```bash
# Scan a specific file
falias --root ~/.config/bash/aliases.sh

# Scan from a different location
falias --root /etc/bash.bashrc
```

## How It Works

### Alias Discovery

falias uses static parsing (no shell execution) to discover aliases:

1. Starts from shell rc files (`.bashrc`, `.zshrc`, etc.)
2. Parses each file line-by-line for:
   - Alias definitions: `alias name='value'`
   - Global aliases (zsh): `alias -g name='value'`
   - Source statements: `source file` or `. file`
3. Resolves sourced file paths safely
4. Recursively scans included files
5. Tracks all definitions in parse order

### Path Resolution

Safely resolves common shell path patterns:

- `~` → User home directory
- `$HOME`, `${HOME}` → User home directory
- `$XDG_CONFIG_HOME` → XDG config directory or `~/.config`
- Quoted paths: `"..."` or `'...'`

**Does NOT evaluate**:

- Command substitutions: `$(...)`, `` `...` ``
- Complex expressions
- Arbitrary shell code

### Loop Prevention

- Maintains a visited set of canonical file paths
- Enforces maximum include depth (25 levels)
- Shows warnings if depth limit is reached

### Override Detection

When an alias is defined multiple times:

- The last definition becomes active
- All previous definitions are preserved
- Full history is shown in details view
- Badge indicates it's been overridden

## Project Structure

```
falias/
├── cmd/
│   └── falias/
│       └── main.go              # CLI entry point
├── internal/
│   ├── model/
│   │   └── types.go             # Data structures
│   ├── config/
│   │   ├── config.go            # Configuration management
│   │   └── themes.go            # Theme definitions
│   ├── resolve/
│   │   └── path.go              # Path expansion
│   ├── parser/
│   │   └── alias.go             # Alias parsing
│   ├── scanner/
│   │   ├── scanner.go           # Main scanning logic
│   │   ├── include.go           # Source/include parsing
│   │   └── file.go              # File reading
│   ├── export/
│   │   └── json.go              # JSON export
│   └── ui/
│       ├── model.go             # Bubble Tea model
│       ├── update.go            # Update logic
│       ├── view.go              # View rendering
│       ├── keys.go              # Keybindings
│       └── styles.go            # Lipgloss styles
├── go.mod
├── go.sum
└── README.md
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/parser
go test ./internal/resolve
go test ./internal/scanner
```

### Building

```bash
# Build for current platform
go build -o bin/falias ./cmd/falias

# Build with version info
go build -ldflags "-X main.version=1.0.0" -o bin/falias ./cmd/falias

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o falias-linux ./cmd/falias
GOOS=darwin GOARCH=arm64 go build -o falias-macos-arm ./cmd/falias
```

## Limitations

- **Read-only**: Does not modify configuration files (by design)
- **Static analysis**: Cannot resolve complex shell expressions
- **No execution**: Will not run scripts or commands to discover dynamic aliases
- **File-based**: Only discovers file-based aliases, not runtime-defined ones

## Troubleshooting

### No aliases found

```bash
# Check which files are being scanned
falias --debug

# Manually specify your rc file
falias --root ~/.bashrc
```

### Unresolved paths

If falias shows "unresolved paths" in debug mode:

- The path contains variables it can't resolve
- The path uses command substitution
- The file is conditionally included but doesn't exist

These are shown for informational purposes and won't cause failures.

### Permission denied

If you see "permission denied" errors:

- Some configuration files may have restrictive permissions
- falias will skip these files and continue scanning

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

MIT License - See LICENSE file for details

## Acknowledgments

Built with:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [clipboard](https://github.com/atotto/clipboard) - Cross-platform clipboard
- [GoReleaser](https://goreleaser.com/) - Release automation
