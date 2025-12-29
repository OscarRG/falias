package parser

import (
	"regexp"
	"strings"

	"github.com/oscar.rivas/falias/internal/model"
)

var (
	// Regex patterns for alias detection
	// Matches: alias name=value or alias name='value' or alias name="value"
	// Allow letters, numbers, underscores, dots, dashes in alias names
	normalAliasPattern = regexp.MustCompile(`^\s*alias\s+([a-zA-Z_.\-][a-zA-Z0-9_.\-]*)=(.+)`)
	// Matches: alias -g name=value
	globalAliasPattern = regexp.MustCompile(`^\s*alias\s+-g\s+([a-zA-Z_.\-][a-zA-Z0-9_.\-]*)=(.+)`)
)

// AliasParser handles parsing of alias definitions from shell script lines
type AliasParser struct{}

// NewAliasParser creates a new alias parser
func NewAliasParser() *AliasParser {
	return &AliasParser{}
}

// ParseLine attempts to parse an alias definition from a line
// Returns the alias definition and true if successful, or nil and false otherwise
func (p *AliasParser) ParseLine(line string, filePath string, lineNum int) (*model.AliasDefinition, bool) {
	// Remove inline comments (but respect quotes)
	cleaned := removeInlineComment(line)
	if cleaned == "" {
		return nil, false
	}

	// Try global alias first (more specific pattern)
	if matches := globalAliasPattern.FindStringSubmatch(cleaned); matches != nil {
		return &model.AliasDefinition{
			Name:  matches[1],
			Value: extractValue(matches[2]),
			Type:  model.AliasTypeGlobal,
			Location: model.SourceLocation{
				FilePath: filePath,
				LineNum:  lineNum,
				RawLine:  line,
			},
		}, true
	}

	// Try normal alias
	if matches := normalAliasPattern.FindStringSubmatch(cleaned); matches != nil {
		return &model.AliasDefinition{
			Name:  matches[1],
			Value: extractValue(matches[2]),
			Type:  model.AliasTypeNormal,
			Location: model.SourceLocation{
				FilePath: filePath,
				LineNum:  lineNum,
				RawLine:  line,
			},
		}, true
	}

	return nil, false
}

// extractValue extracts the value from an alias assignment, removing quotes and inline comments
func extractValue(s string) string {
	s = strings.TrimSpace(s)

	// If it starts with a quote, find the matching closing quote
	if len(s) > 0 && (s[0] == '\'' || s[0] == '"') {
		quote := s[0]
		// Find the matching closing quote
		for i := 1; i < len(s); i++ {
			if s[i] == quote {
				// Check if it's escaped
				if i > 0 && s[i-1] == '\\' {
					continue
				}
				return s[1:i]
			}
		}
		// No closing quote found, return as-is without the opening quote
		return s[1:]
	}

	// No quotes, remove any trailing comment
	// Split on # but only if not escaped
	parts := splitOnUnescapedChar(s, '#')
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}

	return s
}

// removeInlineComment removes comments from a line, respecting quotes
func removeInlineComment(line string) string {
	inSingleQuote := false
	inDoubleQuote := false
	escaped := false

	var result strings.Builder

	for i := 0; i < len(line); i++ {
		ch := line[i]

		if escaped {
			result.WriteByte(ch)
			escaped = false
			continue
		}

		if ch == '\\' {
			escaped = true
			result.WriteByte(ch)
			continue
		}

		if ch == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			result.WriteByte(ch)
			continue
		}

		if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			result.WriteByte(ch)
			continue
		}

		if ch == '#' && !inSingleQuote && !inDoubleQuote {
			break
		}

		result.WriteByte(ch)
	}

	return strings.TrimSpace(result.String())
}

// splitOnUnescapedChar splits a string on a character that is not escaped
func splitOnUnescapedChar(s string, delim byte) []string {
	var parts []string
	var current strings.Builder
	escaped := false

	for i := 0; i < len(s); i++ {
		ch := s[i]

		if escaped {
			current.WriteByte(ch)
			escaped = false
			continue
		}

		if ch == '\\' {
			escaped = true
			current.WriteByte(ch)
			continue
		}

		if ch == delim {
			parts = append(parts, current.String())
			current.Reset()
			continue
		}

		current.WriteByte(ch)
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

// IsAliasLine quickly checks if a line might contain an alias definition
func IsAliasLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "alias ")
}
