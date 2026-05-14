package main

import (
	"fmt"
	"strings"

	"game-of-life/gol"
	"github.com/nsf/termbox-go"
)

type Pattern struct {
	Name  string
	Cells [][2]int // (row, col) offsets from top-left
}

var patterns = []Pattern{
	{"Block", [][2]int{{0, 0}, {0, 1}, {1, 0}, {1, 1}}},
	{"Beehive", [][2]int{{0, 1}, {0, 2}, {1, 0}, {1, 3}, {2, 1}, {2, 2}}},
	{"Loaf", [][2]int{{0, 1}, {0, 2}, {1, 3}, {2, 0}, {2, 3}, {3, 1}, {3, 2}}},
	{"Boat", [][2]int{{0, 0}, {0, 1}, {1, 0}, {1, 1}, {2, 1}}},
	{"T-block", [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, 2}}},
	{"R-pentomino", [][2]int{{0, 1}, {0, 2}, {1, 0}, {1, 1}, {2, 1}}},
	{"Acorn", [][2]int{{0, 1}, {1, 3}, {2, 0}, {2, 1}, {2, 4}, {2, 5}, {2, 6}}},
	{"Blinker", [][2]int{{0, 0}, {0, 1}, {0, 2}}},
	{"Toad", [][2]int{{0, 1}, {0, 2}, {0, 3}, {1, 0}, {1, 1}, {1, 2}}},
	{"Beacon", [][2]int{{0, 0}, {0, 1}, {1, 0}, {2, 3}, {3, 2}, {3, 3}}},
	{"Clock", [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, 2}, {2, 1}}},
	{"Pulsar", [][2]int{
		{0, 2}, {0, 3}, {0, 4}, {0, 8}, {0, 9}, {0, 10},
		{2, 0}, {2, 5}, {2, 7}, {2, 12},
		{3, 0}, {3, 5}, {3, 7}, {3, 12},
		{4, 0}, {4, 5}, {4, 7}, {4, 12},
		{5, 2}, {5, 3}, {5, 4}, {5, 8}, {5, 9}, {5, 10},
		{7, 2}, {7, 3}, {7, 4}, {7, 8}, {7, 9}, {7, 10},
		{8, 0}, {8, 5}, {8, 7}, {8, 12},
		{9, 0}, {9, 5}, {9, 7}, {9, 12},
		{10, 0}, {10, 5}, {10, 7}, {10, 12},
		{12, 2}, {12, 3}, {12, 4}, {12, 8}, {12, 9}, {12, 10},
	}},
	{"Glider", [][2]int{{0, 1}, {1, 2}, {2, 0}, {2, 1}, {2, 2}}},
	{"LWSS", [][2]int{{0, 1}, {0, 4}, {1, 0}, {2, 0}, {2, 4}, {3, 0}, {3, 1}, {3, 2}, {3, 3}, {3, 4}}},
	{"HWSS", [][2]int{
		{0, 1}, {0, 4},
		{1, 0},
		{2, 0}, {2, 5},
		{3, 0}, {3, 1}, {3, 2}, {3, 3}, {3, 4},
	}},
	{"MWSS", [][2]int{
		{0, 1}, {0, 5},
		{1, 0},
		{2, 0}, {2, 6},
		{3, 0}, {3, 1}, {3, 2}, {3, 3}, {3, 4}, {3, 5}, {3, 6},
	}},
	{"Diehard", [][2]int{{0, 6}, {1, 0}, {1, 1}, {2, 1}, {2, 5}, {2, 6}, {2, 7}}},
	{"Pentadecathlon", [][2]int{
		{0, 1},
		{1, 0}, {1, 2},
		{2, 1}, {3, 1}, {4, 1}, {5, 1},
		{6, 0}, {6, 2},
		{7, 1},
		{8, 1}, {9, 1}, {10, 1},
	}},
	{"Gosper Gun", [][2]int{
		{0, 24},
		{1, 22}, {1, 24},
		{2, 12}, {2, 13}, {2, 20}, {2, 21}, {2, 34}, {2, 35},
		{3, 11}, {3, 15}, {3, 20}, {3, 21}, {3, 34}, {3, 35},
		{4, 0}, {4, 1}, {4, 10}, {4, 16}, {4, 20}, {4, 21},
		{5, 0}, {5, 1}, {5, 10}, {5, 14}, {5, 16}, {5, 17}, {5, 22}, {5, 24},
		{6, 10}, {6, 16}, {6, 24},
		{7, 11}, {7, 15},
		{8, 12}, {8, 13},
	}},
}

// PlacePattern places a pattern centered at (r, c) on the grid.
func PlacePattern(g *gol.Grid, r, c int, p Pattern) {
	for _, cell := range p.Cells {
		gr := r + cell[0]
		gc := c + cell[1]
		if gr >= 0 && gr < g.Rows() && gc >= 0 && gc < g.Cols() {
			g.Set(gr, gc)
		}
	}
}

// fuzzyMatch returns true if query matches name (case-insensitive subsequence).
func fuzzyMatch(query, name string) bool {
	query = strings.ToLower(query)
	name = strings.ToLower(name)
	i := 0
	for j := 0; i < len(query) && j < len(name); j++ {
		if query[i] == name[j] {
			i++
		}
	}
	return i == len(query)
}

// patternOverlay shows a searchable pattern list. Returns the selected pattern or nil.
func patternOverlay(state *GameState) *Pattern {
	var filtered []Pattern
	highlight := 0
	query := ""

	for {
		_, h := termbox.Size()

		// Filter patterns
		filtered = nil
		for _, p := range patterns {
			if query == "" || fuzzyMatch(query, p.Name) {
				filtered = append(filtered, p)
			}
		}
		if highlight >= len(filtered) {
			highlight = len(filtered) - 1
		}
		if highlight < 0 {
			highlight = 0
		}

		// Draw overlay
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		// Header
		drawStr(0, 0, " Place Pattern (type to filter, ↑↓ navigate, Enter place, Esc cancel)", termbox.ColorGreen, termbox.ColorDefault)

		// Query line
		queryText := "Filter: " + query + "_"
		drawStr(0, 2, queryText, termbox.ColorYellow, termbox.ColorDefault)

		// Pattern list
		listStart := 4
		listEnd := h - 2
		visible := listEnd - listStart
		if visible <= 0 {
			visible = 1
		}

		// Calculate scroll offset
		scroll := 0
		if highlight >= scroll+visible {
			scroll = highlight - visible + 1
		}
		if highlight < scroll {
			scroll = highlight
		}

		for i := scroll; i < len(filtered) && (i-scroll)+listStart < listEnd; i++ {
			y := listStart + (i - scroll)
			var fg, bg termbox.Attribute
			text := "  " + filtered[i].Name
			if i == highlight {
				fg = termbox.ColorBlack
				bg = termbox.ColorYellow
				text = ">> " + filtered[i].Name
			} else {
				fg = termbox.ColorWhite
				bg = termbox.ColorDefault
			}
			drawStr(2, y, text, fg, bg)
		}

		// Footer
		drawStr(0, h-1, fmt.Sprintf(" %d pattern(s) | Cursor: (%d,%d)", len(filtered), state.cursorR, state.cursorC), termbox.ColorCyan, termbox.ColorDefault)

		termbox.Flush()

		// Handle input
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			switch {
			case ev.Key == termbox.KeyEsc:
				return nil
			case ev.Key == termbox.KeyEnter:
				if len(filtered) > 0 {
					return &filtered[highlight]
				}
				return nil
			case ev.Key == termbox.KeyArrowUp:
				if highlight > 0 {
					highlight--
				}
			case ev.Key == termbox.KeyArrowDown:
				if highlight < len(filtered)-1 {
					highlight++
				}
			case ev.Key == termbox.KeyBackspace:
				if len(query) > 0 {
					query = query[:len(query)-1]
				}
			case ev.Ch != 0 && ev.Ch < 128:
				if len(query) < 20 {
					query += string(ev.Ch)
				}
			}
		case termbox.EventError:
			return nil
		}
	}
}
