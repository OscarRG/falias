package model

// AliasType represents the type of shell alias
type AliasType string

const (
	AliasTypeNormal AliasType = "normal"
	AliasTypeGlobal AliasType = "global"
)

// SourceLocation represents where an alias or include was defined
type SourceLocation struct {
	FilePath string `json:"file_path"`
	LineNum  int    `json:"line_num"`
	RawLine  string `json:"raw_line"`
}

// AliasDefinition represents a single definition of an alias
type AliasDefinition struct {
	Name     string         `json:"name"`
	Value    string         `json:"value"`
	Type     AliasType      `json:"type"`
	Location SourceLocation `json:"location"`
}

// AliasEntry represents an alias with all its definitions
type AliasEntry struct {
	Name           string            `json:"name"`
	Type           AliasType         `json:"type"`
	ActiveValue    string            `json:"active_value"`
	ActiveLocation SourceLocation    `json:"active_location"`
	Definitions    []AliasDefinition `json:"definitions"` // All definitions in parse order
	IsOverridden   bool              `json:"is_overridden"`
}

// AddDefinition adds a new definition to the alias entry
func (e *AliasEntry) AddDefinition(def AliasDefinition) {
	e.Definitions = append(e.Definitions, def)
	e.ActiveValue = def.Value
	e.ActiveLocation = def.Location
	e.Type = def.Type
	if len(e.Definitions) > 1 {
		e.IsOverridden = true
	}
}

// SourceFile represents a parsed shell configuration file
type SourceFile struct {
	Path        string            `json:"path"`
	Exists      bool              `json:"exists"`
	Readable    bool              `json:"readable"`
	Conditional bool              `json:"conditional"` // Was it in a conditional include?
	Aliases     []AliasDefinition `json:"aliases"`
	Includes    []string          `json:"includes"` // Raw paths to sourced files
	Error       string            `json:"error,omitempty"`
}

// ScanResult represents the complete result of scanning shell files
type ScanResult struct {
	Aliases         map[string]*AliasEntry `json:"aliases"`          // Key: alias name
	Files           map[string]*SourceFile `json:"files"`            // Key: absolute path
	UnresolvedPaths []string               `json:"unresolved_paths"` // Paths we couldn't resolve
	Warnings        []string               `json:"warnings"`
	Shell           string                 `json:"shell"`
	RootFiles       []string               `json:"root_files"`
}

// NewScanResult creates a new empty scan result
func NewScanResult(shell string, rootFiles []string) *ScanResult {
	return &ScanResult{
		Aliases:         make(map[string]*AliasEntry),
		Files:           make(map[string]*SourceFile),
		UnresolvedPaths: make([]string, 0),
		Warnings:        make([]string, 0),
		Shell:           shell,
		RootFiles:       rootFiles,
	}
}

// AddAlias adds or updates an alias in the scan result
func (r *ScanResult) AddAlias(def AliasDefinition) {
	if entry, exists := r.Aliases[def.Name]; exists {
		entry.AddDefinition(def)
	} else {
		r.Aliases[def.Name] = &AliasEntry{
			Name:           def.Name,
			Type:           def.Type,
			ActiveValue:    def.Value,
			ActiveLocation: def.Location,
			Definitions:    []AliasDefinition{def},
			IsOverridden:   false,
		}
	}
}

// GetAliasesSorted returns all aliases sorted by name
func (r *ScanResult) GetAliasesSorted() []*AliasEntry {
	aliases := make([]*AliasEntry, 0, len(r.Aliases))
	for _, entry := range r.Aliases {
		aliases = append(aliases, entry)
	}
	return aliases
}
