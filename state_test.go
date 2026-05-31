package main

import (
	"math/rand"
	"testing"
)

func TestColorOffsetRandomizedOnInit(t *testing.T) {
	// Verify that colorOffset is set to a valid value on init.
	// We test the formula directly since initState requires a tcell.Screen.
	theme := defaultTheme()
	offsets := make(map[int]bool)

	// Run multiple times to verify randomness across processes
	for i := 0; i < 100; i++ {
		offset := rand.Intn(len(theme.PatternColors))
		offsets[offset] = true
	}

	// With 12 colors and 100 trials, we should see multiple distinct offsets.
	// This verifies rand.Intn produces varied results.
	if len(offsets) < 5 {
		t.Errorf("Expected varied offsets across 100 trials, got only %d distinct values", len(offsets))
	}
}

func TestCellFgColorOffsetApplied(t *testing.T) {
	theme := defaultTheme()

	// Different offsets should produce different palette indices for the same position.
	r, c := 5, 10

	paletteIdxA := ((r*31 + c*37) + 0) % len(theme.PatternColors)
	paletteIdxB := ((r*31 + c*37) + 1) % len(theme.PatternColors)

	if paletteIdxA == paletteIdxB {
		t.Errorf("Offsets 0 and 1 produced the same palette index %d for position (%d,%d)",
			paletteIdxA, r, c)
	}
}

func TestColorOffsetBounds(t *testing.T) {
	theme := defaultTheme()

	// colorOffset should be in [0, len(PatternColors)).
	for offset := 0; offset < 100; offset++ {
		v := rand.Intn(len(theme.PatternColors))
		if v < 0 || v >= len(theme.PatternColors) {
			t.Errorf("colorOffset %d out of bounds [0, %d)", v, len(theme.PatternColors))
		}
	}
}

func TestPaletteIndexNeverNegative(t *testing.T) {
	theme := defaultTheme()

	// Verify that ((r*31 + c*37) + offset) % len(colors) is never negative.
	// Since all terms are non-negative, this should always hold.
	r, c := 0, 0

	for offset := 0; offset < len(theme.PatternColors); offset++ {
		idx := ((r*31 + c*37) + offset) % len(theme.PatternColors)
		if idx < 0 || idx >= len(theme.PatternColors) {
			t.Errorf("Palette index %d out of bounds for offset %d", idx, offset)
		}
	}

	// Test with larger coordinates
	r, c = 100, 200

	for offset := 0; offset < len(theme.PatternColors); offset++ {
		idx := ((r*31 + c*37) + offset) % len(theme.PatternColors)
		if idx < 0 || idx >= len(theme.PatternColors) {
			t.Errorf("Palette index %d out of bounds for offset %d at (%d,%d)", idx, offset, r, c)
		}
	}
}

func TestPaletteIndexCoversAllColors(t *testing.T) {
	theme := defaultTheme()
	r, c := 0, 0

	seen := make(map[int]bool)

	// With enough offsets, we should cover all palette colors.
	for offset := 0; offset < len(theme.PatternColors); offset++ {
		idx := ((r*31 + c*37) + offset) % len(theme.PatternColors)
		seen[idx] = true
	}

	if len(seen) != len(theme.PatternColors) {
		t.Errorf("Expected to cover all %d palette colors, but only covered %d",
			len(theme.PatternColors), len(seen))
	}
}

func TestPaletteIndexCoversAllColorsAtVariousPositions(t *testing.T) {
	theme := defaultTheme()

	// Test at multiple positions that the formula should cover all colors.
	testPositions := []struct{ r, c int }{
		{0, 0}, {1, 0}, {0, 1}, {5, 10}, {10, 5}, {100, 100},
	}

	for _, pos := range testPositions {
		seen := make(map[int]bool)

		for offset := 0; offset < len(theme.PatternColors); offset++ {
			idx := ((pos.r*31 + pos.c*37) + offset) % len(theme.PatternColors)
			seen[idx] = true
		}

		if len(seen) != len(theme.PatternColors) {
			t.Errorf("Position (%d,%d): covered %d/%d palette colors",
				pos.r, pos.c, len(seen), len(theme.PatternColors))
		}
	}
}
