package resolve

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// Regex to match $VAR or ${VAR} patterns
	varPattern = regexp.MustCompile(`\$\{?([A-Z_][A-Z0-9_]*)\}?`)
)

// PathResolver handles expansion of shell path variables
type PathResolver struct {
	homeDir        string
	xdgConfigHome  string
	customVars     map[string]string
}

// NewPathResolver creates a new path resolver with environment defaults
func NewPathResolver() *PathResolver {
	homeDir, _ := os.UserHomeDir()
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(homeDir, ".config")
	}

	return &PathResolver{
		homeDir:       homeDir,
		xdgConfigHome: xdgConfigHome,
		customVars:    make(map[string]string),
	}
}

// SetVariable sets a custom variable for path expansion
func (r *PathResolver) SetVariable(name, value string) {
	r.customVars[name] = value
}

// ResolvePath attempts to resolve a shell path to an absolute path
// Returns the resolved path and a boolean indicating success
func (r *PathResolver) ResolvePath(path string) (string, bool) {
	if path == "" {
		return "", false
	}

	// Remove surrounding quotes
	path = removeQuotes(path)

	// Handle tilde expansion
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(r.homeDir, path[2:])
	} else if path == "~" {
		path = r.homeDir
	}

	// Expand variables
	expanded, ok := r.expandVariables(path)
	if !ok {
		return "", false
	}

	// Clean the path
	expanded = filepath.Clean(expanded)

	// Make absolute if not already
	if !filepath.IsAbs(expanded) {
		return "", false
	}

	return expanded, true
}

// expandVariables expands $VAR and ${VAR} patterns in the path
func (r *PathResolver) expandVariables(path string) (string, bool) {
	result := path
	foundUnknown := false

	// Find all variable references
	matches := varPattern.FindAllStringSubmatch(path, -1)
	for _, match := range matches {
		fullMatch := match[0]
		varName := match[1]

		var replacement string
		var found bool

		// Check known variables
		switch varName {
		case "HOME":
			replacement = r.homeDir
			found = true
		case "XDG_CONFIG_HOME":
			replacement = r.xdgConfigHome
			found = true
		default:
			// Check custom variables
			if val, ok := r.customVars[varName]; ok {
				replacement = val
				found = true
			} else {
				// Unknown variable
				foundUnknown = true
			}
		}

		if found {
			result = strings.Replace(result, fullMatch, replacement, 1)
		}
	}

	// If we found unknown variables, return failure
	if foundUnknown {
		return "", false
	}

	return result, true
}

// removeQuotes removes surrounding single or double quotes
func removeQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// Canonicalize converts a path to its canonical absolute form
func Canonicalize(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// Evaluate symlinks
	canonPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		// If the file doesn't exist, EvalSymlinks fails
		// Return the absolute path anyway
		return absPath, nil
	}

	return canonPath, nil
}

// ParseVariableAssignment attempts to parse a variable assignment like VAR="value"
// Returns variable name, value, and success boolean
func ParseVariableAssignment(line string) (string, string, bool) {
	line = strings.TrimSpace(line)

	// Match VAR=value, VAR="value", or VAR='value'
	re := regexp.MustCompile(`^([A-Z_][A-Z0-9_]*)=(.+)$`)
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return "", "", false
	}

	varName := matches[1]
	value := removeQuotes(matches[2])

	return varName, value, true
}
