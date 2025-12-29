package parser

import (
	"testing"

	"github.com/oscar.rivas/falias/internal/model"
)

func TestParseAliasVariations(t *testing.T) {
	parser := NewAliasParser()

	tests := []struct {
		name      string
		line      string
		wantName  string
		wantValue string
		wantType  model.AliasType
		wantOk    bool
	}{
		{
			name:      "simple single quote",
			line:      "alias ll='ls -la'",
			wantName:  "ll",
			wantValue: "ls -la",
			wantType:  model.AliasTypeNormal,
			wantOk:    true,
		},
		{
			name:      "simple double quote",
			line:      `alias gs="git status"`,
			wantName:  "gs",
			wantValue: "git status",
			wantType:  model.AliasTypeNormal,
			wantOk:    true,
		},
		{
			name:      "no quotes",
			line:      "alias ..=cd ..",
			wantName:  "..",
			wantValue: "cd ..",
			wantType:  model.AliasTypeNormal,
			wantOk:    true,
		},
		{
			name:      "with inline comment",
			line:      "alias grep='grep --color=auto' # colorize grep",
			wantName:  "grep",
			wantValue: "grep --color=auto",
			wantType:  model.AliasTypeNormal,
			wantOk:    true,
		},
		{
			name:      "global alias",
			line:      "alias -g G='| grep'",
			wantName:  "G",
			wantValue: "| grep",
			wantType:  model.AliasTypeGlobal,
			wantOk:    true,
		},
		{
			name:      "with leading whitespace",
			line:      "  alias foo='bar'",
			wantName:  "foo",
			wantValue: "bar",
			wantType:  model.AliasTypeNormal,
			wantOk:    true,
		},
		{
			name:      "empty value",
			line:      "alias empty=''",
			wantName:  "empty",
			wantValue: "",
			wantType:  model.AliasTypeNormal,
			wantOk:    true,
		},
		{
			name:      "hash in value",
			line:      "alias test='echo # not a comment'",
			wantName:  "test",
			wantValue: "echo # not a comment",
			wantType:  model.AliasTypeNormal,
			wantOk:    true,
		},
		{
			name:   "comment line",
			line:   "# alias foo='bar'",
			wantOk: false,
		},
		{
			name:   "not an alias",
			line:   "export FOO=bar",
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			def, ok := parser.ParseLine(tt.line, "/test/file", 1)

			if ok != tt.wantOk {
				t.Errorf("ParseLine() ok = %v, want %v", ok, tt.wantOk)
				return
			}

			if !tt.wantOk {
				return
			}

			if def.Name != tt.wantName {
				t.Errorf("ParseLine() name = %v, want %v", def.Name, tt.wantName)
			}

			if def.Value != tt.wantValue {
				t.Errorf("ParseLine() value = %v, want %v", def.Value, tt.wantValue)
			}

			if def.Type != tt.wantType {
				t.Errorf("ParseLine() type = %v, want %v", def.Type, tt.wantType)
			}
		})
	}
}

func TestIsAliasLine(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"alias foo='bar'", true},
		{"  alias foo='bar'", true},
		{"alias -g G='grep'", true},
		{"# alias foo='bar'", false},
		{"export FOO=bar", false},
		{"", false},
		{"source ~/.bashrc", false},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			if got := IsAliasLine(tt.line); got != tt.want {
				t.Errorf("IsAliasLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
