package export

import (
	"encoding/json"
	"io"
	"sort"

	"github.com/oscar.rivas/falias/internal/model"
)

// JSONExporter handles exporting scan results to JSON
type JSONExporter struct {
	pretty bool
}

// NewJSONExporter creates a new JSON exporter
func NewJSONExporter(pretty bool) *JSONExporter {
	return &JSONExporter{
		pretty: pretty,
	}
}

// Export writes the scan result to the given writer as JSON
func (e *JSONExporter) Export(result *model.ScanResult, w io.Writer) error {
	encoder := json.NewEncoder(w)
	if e.pretty {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(result)
}

// ExportAliases exports only the aliases in a simplified format
func (e *JSONExporter) ExportAliases(result *model.ScanResult, w io.Writer) error {
	// Create a simplified structure for output
	type SimpleAlias struct {
		Name     string `json:"name"`
		Value    string `json:"value"`
		Type     string `json:"type"`
		File     string `json:"file"`
		Line     int    `json:"line"`
		Override bool   `json:"overridden,omitempty"`
	}

	aliases := make([]SimpleAlias, 0, len(result.Aliases))

	entries := result.GetAliasesSorted()

	// Sort by name
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	// Convert to simple format
	for _, entry := range entries {
		aliases = append(aliases, SimpleAlias{
			Name:     entry.Name,
			Value:    entry.ActiveValue,
			Type:     string(entry.Type),
			File:     entry.ActiveLocation.FilePath,
			Line:     entry.ActiveLocation.LineNum,
			Override: entry.IsOverridden,
		})
	}

	encoder := json.NewEncoder(w)
	if e.pretty {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(aliases)
}
