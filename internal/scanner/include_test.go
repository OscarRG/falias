package scanner

import (
	"testing"
)

func TestParseInclude(t *testing.T) {
	parser := NewIncludeParser()

	tests := []struct {
		name            string
		line            string
		wantPaths       []string
		wantConditional []bool
	}{
		{
			name:            "simple source",
			line:            "source ~/.bashrc",
			wantPaths:       []string{"~/.bashrc"},
			wantConditional: []bool{false},
		},
		{
			name:            "dot command",
			line:            ". ~/.bash_profile",
			wantPaths:       []string{"~/.bash_profile"},
			wantConditional: []bool{false},
		},
		{
			name:            "quoted path",
			line:            `source "$HOME/.bashrc"`,
			wantPaths:       []string{`"$HOME/.bashrc"`},
			wantConditional: []bool{false},
		},
		{
			name:            "conditional source bash",
			line:            "if [ -f ~/.bashrc ]; then source ~/.bashrc; fi",
			wantPaths:       []string{"~/.bashrc"},
			wantConditional: []bool{true},
		},
		{
			name:            "conditional source zsh double bracket",
			line:            "[[ -f ~/.zsh/aliases.zsh ]] && source ~/.zsh/aliases.zsh",
			wantPaths:       []string{"~/.zsh/aliases.zsh"},
			wantConditional: []bool{true},
		},
		{
			name:            "with leading whitespace",
			line:            "  source ~/.bashrc",
			wantPaths:       []string{"~/.bashrc"},
			wantConditional: []bool{false},
		},
		{
			name:            "with inline comment",
			line:            "source ~/.bashrc # load config",
			wantPaths:       []string{"~/.bashrc"},
			wantConditional: []bool{false},
		},
		{
			name:      "comment line",
			line:      "# source ~/.bashrc",
			wantPaths: []string{},
		},
		{
			name:      "not a source line",
			line:      "alias foo=bar",
			wantPaths: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := parser.ParseLine(tt.line)

			if len(results) != len(tt.wantPaths) {
				t.Errorf("ParseLine() returned %d results, want %d", len(results), len(tt.wantPaths))
				return
			}

			for i, result := range results {
				if result.Path != tt.wantPaths[i] {
					t.Errorf("ParseLine()[%d].Path = %v, want %v", i, result.Path, tt.wantPaths[i])
				}

				if len(tt.wantConditional) > i && result.Conditional != tt.wantConditional[i] {
					t.Errorf("ParseLine()[%d].Conditional = %v, want %v", i, result.Conditional, tt.wantConditional[i])
				}
			}
		})
	}
}

func TestIsIncludeLine(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"source ~/.bashrc", true},
		{". ~/.bashrc", true},
		{"  source ~/.bashrc", true},
		{"if [ -f ~/.bashrc ]; then source ~/.bashrc; fi", true},
		{"# source ~/.bashrc", false},
		{"alias foo=bar", false},
		{"", false},
		{"export FOO=bar", false},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			if got := IsIncludeLine(tt.line); got != tt.want {
				t.Errorf("IsIncludeLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
