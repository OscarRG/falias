package resolve

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePath(t *testing.T) {
	resolver := NewPathResolver()
	homeDir, _ := os.UserHomeDir()

	tests := []struct {
		name    string
		path    string
		wantOk  bool
		wantContains string // Check if result contains this
	}{
		{
			name:    "tilde expansion",
			path:    "~/.bashrc",
			wantOk:  true,
			wantContains: homeDir,
		},
		{
			name:    "HOME variable",
			path:    "$HOME/.bashrc",
			wantOk:  true,
			wantContains: homeDir,
		},
		{
			name:    "HOME with braces",
			path:    "${HOME}/.bashrc",
			wantOk:  true,
			wantContains: homeDir,
		},
		{
			name:    "single quoted path",
			path:    "'~/.bashrc'",
			wantOk:  true,
			wantContains: homeDir,
		},
		{
			name:    "double quoted path",
			path:    `"~/.bashrc"`,
			wantOk:  true,
			wantContains: homeDir,
		},
		{
			name:    "absolute path",
			path:    "/etc/bashrc",
			wantOk:  true,
			wantContains: "/etc",
		},
		{
			name:   "unknown variable",
			path:   "$UNKNOWN_VAR/file",
			wantOk: false,
		},
		{
			name:   "relative path",
			path:   "./file",
			wantOk: false,
		},
		{
			name:   "empty path",
			path:   "",
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := resolver.ResolvePath(tt.path)

			if ok != tt.wantOk {
				t.Errorf("ResolvePath() ok = %v, want %v", ok, tt.wantOk)
				return
			}

			if tt.wantOk && tt.wantContains != "" {
				if !contains(got, tt.wantContains) {
					t.Errorf("ResolvePath() = %v, want to contain %v", got, tt.wantContains)
				}
			}
		})
	}
}

func TestSetVariable(t *testing.T) {
	resolver := NewPathResolver()
	resolver.SetVariable("MYVAR", "/custom/path")

	path, ok := resolver.ResolvePath("$MYVAR/file")
	if !ok {
		t.Errorf("Expected successful resolution after SetVariable")
	}

	if !contains(path, "/custom/path") {
		t.Errorf("ResolvePath() = %v, want to contain /custom/path", path)
	}
}

func TestCanonialize(t *testing.T) {
	// Create a temp file to test with
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(tmpFile, []byte("test"), 0644)

	canon, err := Canonicalize(tmpFile)
	if err != nil {
		t.Errorf("Canonicalize() error = %v", err)
	}

	if !filepath.IsAbs(canon) {
		t.Errorf("Canonicalize() result is not absolute: %v", canon)
	}
}

func TestParseVariableAssignment(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		wantName  string
		wantValue string
		wantOk    bool
	}{
		{
			name:      "simple assignment",
			line:      "FOO=bar",
			wantName:  "FOO",
			wantValue: "bar",
			wantOk:    true,
		},
		{
			name:      "quoted value",
			line:      `FOO="bar"`,
			wantName:  "FOO",
			wantValue: "bar",
			wantOk:    true,
		},
		{
			name:      "single quoted",
			line:      "FOO='bar'",
			wantName:  "FOO",
			wantValue: "bar",
			wantOk:    true,
		},
		{
			name:      "with path",
			line:      `XDG_CONFIG_HOME="$HOME/.config"`,
			wantName:  "XDG_CONFIG_HOME",
			wantValue: "$HOME/.config",
			wantOk:    true,
		},
		{
			name:   "not an assignment",
			line:   "alias foo=bar",
			wantOk: false,
		},
		{
			name:   "export statement",
			line:   "export FOO=bar",
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, value, ok := ParseVariableAssignment(tt.line)

			if ok != tt.wantOk {
				t.Errorf("ParseVariableAssignment() ok = %v, want %v", ok, tt.wantOk)
				return
			}

			if !tt.wantOk {
				return
			}

			if name != tt.wantName {
				t.Errorf("ParseVariableAssignment() name = %v, want %v", name, tt.wantName)
			}

			if value != tt.wantValue {
				t.Errorf("ParseVariableAssignment() value = %v, want %v", value, tt.wantValue)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		findInString(s, substr)))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
