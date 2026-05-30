package main

import (
	"testing"
)

func TestDefaultTheme(t *testing.T) {
	theme := defaultTheme()
	if theme == nil {
		t.Fatal("defaultTheme() should not return nil")
	}

	themes := builtInThemes()
	if theme.Name != themes[0].Name {
		t.Errorf("defaultTheme() should return first built-in theme, got %s", theme.Name)
	}
}

func TestFindTheme(t *testing.T) {
	tests := []struct {
		name  string
		found bool
	}{
		{"Dracula", true},
		{"Monokai", true},
		{"Gruvbox Dark", true},
		{"Nonexistent", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme := findTheme(tt.name)

			if tt.found && theme == nil {
				t.Errorf("findTheme(%q) should have found a theme", tt.name)
			}

			if !tt.found && theme != nil {
				t.Errorf("findTheme(%q) should have returned nil", tt.name)
			}
		})
	}
}

func TestBuiltInThemesCount(t *testing.T) {
	themes := builtInThemes()
	if len(themes) == 0 {
		t.Fatal("builtInThemes() should return at least one theme")
	}
}
