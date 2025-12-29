package scanner

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/oscar.rivas/falias/internal/model"
	"github.com/oscar.rivas/falias/internal/parser"
	"github.com/oscar.rivas/falias/internal/resolve"
)

const (
	maxDepth = 25 // Maximum include depth to prevent infinite loops
)

// Scanner orchestrates the scanning of shell configuration files
type Scanner struct {
	pathResolver  *resolve.PathResolver
	aliasParser   *parser.AliasParser
	includeParser *IncludeParser
	fileReader    *FileReader
}

// NewScanner creates a new scanner
func NewScanner() *Scanner {
	return &Scanner{
		pathResolver:  resolve.NewPathResolver(),
		aliasParser:   parser.NewAliasParser(),
		includeParser: NewIncludeParser(),
		fileReader:    NewFileReader(),
	}
}

// ScanShellFiles scans shell configuration files starting from the given root paths
func (s *Scanner) ScanShellFiles(shell string, rootPaths []string) (*model.ScanResult, error) {
	result := model.NewScanResult(shell, rootPaths)

	// Track visited files to prevent loops
	visited := make(map[string]bool)

	// Process each root file
	for _, rootPath := range rootPaths {
		expanded, ok := s.pathResolver.ResolvePath(rootPath)
		if !ok {
			result.UnresolvedPaths = append(result.UnresolvedPaths, rootPath)
			continue
		}

		if err := s.scanFile(expanded, result, visited, 0); err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Error scanning %s: %v", expanded, err))
		}
	}

	return result, nil
}

// scanFile recursively scans a single file
func (s *Scanner) scanFile(filePath string, result *model.ScanResult, visited map[string]bool, depth int) error {
	// Check depth limit
	if depth >= maxDepth {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Max include depth (%d) reached at %s", maxDepth, filePath))
		return nil
	}

	// Canonicalize path
	canonPath, err := resolve.Canonicalize(filePath)
	if err != nil {
		canonPath = filePath // Use original if canonicalization fails
	}

	// Check if already visited
	if visited[canonPath] {
		return nil // Skip, already processed
	}
	visited[canonPath] = true

	// Create source file entry
	sourceFile := &model.SourceFile{
		Path:     canonPath,
		Exists:   FileExists(canonPath),
		Readable: FileReadable(canonPath),
		Aliases:  make([]model.AliasDefinition, 0),
		Includes: make([]string, 0),
	}

	// Store the file entry
	result.Files[canonPath] = sourceFile

	// If file doesn't exist or isn't readable, return early
	if !sourceFile.Exists {
		sourceFile.Error = "file does not exist"
		return nil
	}

	if !sourceFile.Readable {
		sourceFile.Error = "permission denied"
		return nil
	}

	// Read file lines
	lines, err := s.fileReader.ReadLines(canonPath)
	if err != nil {
		sourceFile.Error = err.Error()
		return err
	}

	// Parse each line
	for lineNum, line := range lines {
		lineNumber := lineNum + 1 // Line numbers start at 1

		// Try to parse as alias
		if parser.IsAliasLine(line) {
			if aliasDef, ok := s.aliasParser.ParseLine(line, canonPath, lineNumber); ok {
				sourceFile.Aliases = append(sourceFile.Aliases, *aliasDef)
				result.AddAlias(*aliasDef)
			}
		}

		// Try to parse as include
		if IsIncludeLine(line) {
			includes := s.includeParser.ParseLine(line)
			for _, inc := range includes {
				sourceFile.Includes = append(sourceFile.Includes, inc.Path)

				// Try to resolve the included path
				resolvedPath, ok := s.pathResolver.ResolvePath(inc.Path)
				if !ok {
					result.UnresolvedPaths = append(result.UnresolvedPaths, inc.Path)
					continue
				}

				// Recursively scan the included file
				if err := s.scanFile(resolvedPath, result, visited, depth+1); err != nil {
					result.Warnings = append(result.Warnings, fmt.Sprintf("Error scanning included file %s: %v", resolvedPath, err))
				}

				// Mark file as conditional if needed
				if inc.Conditional {
					if sf, exists := result.Files[resolvedPath]; exists {
						sf.Conditional = true
					}
				}
			}
		}

		// Try to parse variable assignments for path resolution
		if varName, value, ok := resolve.ParseVariableAssignment(line); ok {
			// Update resolver with discovered variables
			s.pathResolver.SetVariable(varName, value)
		}
	}

	return nil
}

// DetectShell attempts to detect the user's shell
func DetectShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "bash" // Default to bash
	}

	// Extract shell name from path
	shellName := filepath.Base(shell)
	if shellName == "zsh" {
		return "zsh"
	}

	return "bash"
}

// GetDefaultRootFiles returns the default root files for a given shell
func GetDefaultRootFiles(shell string) []string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "~"
	}

	switch shell {
	case "zsh":
		return []string{
			filepath.Join(homeDir, ".zshrc"),
			filepath.Join(homeDir, ".zprofile"),
		}
	case "bash":
		return []string{
			filepath.Join(homeDir, ".bashrc"),
			filepath.Join(homeDir, ".bash_profile"),
			filepath.Join(homeDir, ".profile"),
		}
	default:
		return []string{
			filepath.Join(homeDir, ".bashrc"),
		}
	}
}

// FilterExistingFiles filters a list of files to only those that exist
func FilterExistingFiles(files []string) []string {
	existing := make([]string, 0)
	for _, file := range files {
		if FileExists(file) {
			existing = append(existing, file)
		}
	}
	return existing
}
