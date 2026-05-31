package main

import (
	"os"
	"testing"
)

func TestParseThemeFromArgsHelpExits(t *testing.T) {
	// Save and restore original args
	original := os.Args
	defer func() { os.Args = original }()

	tests := []struct {
		name      string
		args      []string
		shouldErr bool
	}{
		{"--help", []string{"prog", "--help"}, true},
		{"-h", []string{"prog", "-h"}, true},
		{"--list-themes", []string{"prog", "--list-themes"}, true},
		{"-l", []string{"prog", "-l"}, true},
		{"no args", []string{"prog"}, false},
		{"--theme GruvboxDark", []string{"prog", "--theme", "Gruvbox Dark"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			theme, err := parseThemeFromArgs()

			if tt.shouldErr && err == nil {
				t.Error("expected error for exit flag, got nil")
			}

			if !tt.shouldErr && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}

			if tt.shouldErr && theme != nil {
				t.Error("expected nil theme for exit flag")
			}
		})
	}
}
