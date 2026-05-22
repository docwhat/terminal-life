package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/nsf/termbox-go"
)

// Theme holds terminal capability flags and color attributes.
type Theme struct {
	TrueColor bool
	DarkBg    bool

	// UI chrome
	TitleFg, TitleBg   termbox.Attribute
	InfoFg, InfoBg     termbox.Attribute
	StatusFg, StatusBg termbox.Attribute
	DialogFg, DialogBg termbox.Attribute

	// Cell rendering
	CellChar      rune                // character used for alive cells
	ManualCellFg  termbox.Attribute   // color for manually toggled cells
	PatternColors []termbox.Attribute // palette cycled per pattern placement
}

// NewTheme detects terminal capabilities and builds a color theme.
func NewTheme() *Theme {
	t := &Theme{
		DarkBg:    true, // default assumption
		TrueColor: false,
		CellChar:  '●',
	}

	// ── Truecolor detection ──
	colorterm := os.Getenv("COLORTERM")
	t.TrueColor = strings.Contains(colorterm, "24bit") || strings.Contains(colorterm, "truecolor")

	// ── Background brightness detection ──
	t.DarkBg = detectDarkBackground()

	// ── Build palette ──
	if t.TrueColor {
		t.buildTruecolorPalette()
	} else {
		t.build256Palette()
	}

	return t
}

// detectDarkBackground checks environment hints for a light terminal theme.
func detectDarkBackground() bool {
	// Check COLORFGBG (e.g. "15;0" = white on black, "0;15" = black on white)
	if v := os.Getenv("COLORFGBG"); v != "" {
		parts := strings.Split(v, ";")
		if len(parts) >= 2 {
			bg, err := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err == nil && bg >= 8 && bg != 255 {
				return false // light background
			}
		}
	}

	// Check TERM_PROGRAM environment for known light-theme terminals
	// (most terminals default to dark, so we stay conservative)
	return true
}

func rgb(r, g, b byte) termbox.Attribute {
	return termbox.RGBToAttribute(r, g, b)
}

func (t *Theme) buildTruecolorPalette() {
	// Vibrant, distinct colors for pattern tracking
	t.PatternColors = []termbox.Attribute{
		rgb(0, 180, 216),   // cyan
		rgb(45, 198, 83),   // green
		rgb(255, 183, 3),   // yellow
		rgb(251, 133, 0),   // orange
		rgb(239, 71, 111),  // red-pink
		rgb(247, 37, 133),  // pink
		rgb(114, 9, 183),   // purple
		rgb(67, 97, 238),   // blue
		rgb(76, 201, 240),  // teal
		rgb(6, 214, 160),   // lime
		rgb(255, 107, 107), // coral
		rgb(255, 209, 102), // gold
		rgb(162, 155, 254), // lavender
		rgb(253, 121, 168), // rose
		rgb(99, 230, 190),  // mint
		rgb(255, 159, 243), // peach
	}
	t.ManualCellFg = rgb(255, 255, 255) // white for manual cells

	if t.DarkBg {
		t.TitleFg = rgb(255, 255, 255)
		t.TitleBg = rgb(30, 30, 46)
		t.InfoFg = rgb(205, 214, 244)
		t.InfoBg = rgb(30, 30, 46)
		t.StatusFg = rgb(205, 214, 244)
		t.StatusBg = rgb(49, 50, 68)
		t.DialogFg = rgb(205, 214, 244)
		t.DialogBg = rgb(30, 30, 46)
	} else {
		t.TitleFg = rgb(30, 30, 30)
		t.TitleBg = rgb(220, 220, 230)
		t.InfoFg = rgb(30, 30, 30)
		t.InfoBg = rgb(220, 220, 230)
		t.StatusFg = rgb(30, 30, 30)
		t.StatusBg = rgb(200, 200, 215)
		t.DialogFg = rgb(30, 30, 30)
		t.DialogBg = rgb(220, 220, 230)
	}
}

func (t *Theme) build256Palette() {
	// Good 256-color fallbacks
	t.PatternColors = []termbox.Attribute{
		termbox.ColorCyan,
		termbox.ColorGreen,
		termbox.ColorYellow,
		termbox.Attribute(209), // orange
		termbox.ColorRed,
		termbox.Attribute(206), // pink
		termbox.ColorMagenta,
		termbox.ColorBlue,
		termbox.Attribute(81),  // teal
		termbox.Attribute(48),  // lime
		termbox.Attribute(203), // coral
		termbox.Attribute(220), // gold
		termbox.Attribute(153), // lavender
		termbox.Attribute(212), // rose
		termbox.Attribute(79),  // mint
		termbox.Attribute(218), // peach
	}
	t.ManualCellFg = termbox.ColorWhite

	if t.DarkBg {
		t.TitleFg = termbox.ColorWhite
		t.TitleBg = termbox.ColorDefault
		t.InfoFg = termbox.ColorWhite
		t.InfoBg = termbox.ColorDefault
		t.StatusFg = termbox.ColorWhite
		t.StatusBg = termbox.ColorDefault
		t.DialogFg = termbox.ColorWhite
		t.DialogBg = termbox.ColorDefault
	} else {
		t.TitleFg = termbox.ColorBlack
		t.TitleBg = termbox.ColorDefault
		t.InfoFg = termbox.ColorBlack
		t.InfoBg = termbox.ColorDefault
		t.StatusFg = termbox.ColorBlack
		t.StatusBg = termbox.ColorDefault
		t.DialogFg = termbox.ColorBlack
		t.DialogBg = termbox.ColorDefault
	}
}
