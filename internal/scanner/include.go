package scanner

import (
	"regexp"
	"strings"
)

var (
	// Matches: source <path> or . <path>
	sourcePattern = regexp.MustCompile(`^\s*(source|\.)\s+(.+)`)

	// Matches: if [ -f <path> ]; then source <path>; fi
	// Also matches zsh syntax: [[ -f <path> ]] && source <path>
	conditionalPattern = regexp.MustCompile(`(?:if\s+)?\[\[?\s+-f\s+(.+?)\s+\]\]?.*(?:source|\.)\s+(.+?)(?:\s*;|\s*$)`)
)

// IncludeParser handles parsing of source/include statements
type IncludeParser struct{}

// NewIncludeParser creates a new include parser
func NewIncludeParser() *IncludeParser {
	return &IncludeParser{}
}

// ParseResult represents the result of parsing an include line
type ParseResult struct {
	Path        string
	Conditional bool
}

// ParseLine attempts to parse a source/include statement from a line
// Returns the path(s) and whether they were found
func (p *IncludeParser) ParseLine(line string) []ParseResult {
	results := make([]ParseResult, 0)

	// Remove leading/trailing whitespace
	line = strings.TrimSpace(line)

	// Skip comments and empty lines
	if line == "" || strings.HasPrefix(line, "#") {
		return results
	}

	// Check for conditional include first
	if matches := conditionalPattern.FindStringSubmatch(line); matches != nil {
		// Extract the path from the source command (not the test)
		path := extractPath(matches[2])
		if path != "" {
			results = append(results, ParseResult{
				Path:        path,
				Conditional: true,
			})
		}
		return results
	}

	// Check for simple source/dot command
	if matches := sourcePattern.FindStringSubmatch(line); matches != nil {
		path := extractPath(matches[2])
		if path != "" {
			results = append(results, ParseResult{
				Path:        path,
				Conditional: false,
			})
		}
		return results
	}

	return results
}

// extractPath extracts a file path from a source argument
// Handles quotes, removes inline comments, etc.
func extractPath(s string) string {
	s = strings.TrimSpace(s)

	// Remove inline comment (but be careful with quotes)
	s = removeTrailingComment(s)

	// Remove trailing semicolon
	s = strings.TrimRight(s, ";")
	s = strings.TrimSpace(s)

	return s
}

// removeTrailingComment removes a trailing comment from a string
func removeTrailingComment(s string) string {
	inSingleQuote := false
	inDoubleQuote := false
	escaped := false

	for i := 0; i < len(s); i++ {
		ch := s[i]

		if escaped {
			escaped = false
			continue
		}

		if ch == '\\' {
			escaped = true
			continue
		}

		if ch == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			continue
		}

		if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			continue
		}

		if ch == '#' && !inSingleQuote && !inDoubleQuote {
			return strings.TrimSpace(s[:i])
		}
	}

	return s
}

// IsIncludeLine quickly checks if a line might contain a source/include
func IsIncludeLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	// Ignore comment lines
	if strings.HasPrefix(trimmed, "#") {
		return false
	}
	return strings.HasPrefix(trimmed, "source ") ||
		strings.HasPrefix(trimmed, ". ") ||
		strings.Contains(trimmed, "source ") // For conditionals
}
