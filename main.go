package main

import (
	"fmt"
	"strings"
	"time"

	"game-of-life/gol"
	"github.com/nsf/termbox-go"
)

const sidebarWidth = 20

type GameState struct {
	grid        *gol.Grid
	cursorR     int
	cursorC     int
	running     bool
	generations int
	interval    time.Duration // tick interval
}

func (s *GameState) Population() int {
	pop := 0
	for r := 0; r < s.grid.Rows(); r++ {
		for c := 0; c < s.grid.Cols(); c++ {
			if s.grid.Cells(r, c) {
				pop++
			}
		}
	}
	return pop
}

func (s *GameState) Render() {
	w, h := termbox.Size()
	if h < 5 || w < sidebarWidth+10 {
		return
	}

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// ── Title bar (row 0) ──
	titleFg, titleBg := termbox.ColorWhite, termbox.ColorBlue
	for x := 0; x < w; x++ {
		termbox.SetCell(x, 0, ' ', titleFg, titleBg)
	}
	title := " ● Game of Life ● "
	drawStr((w-len(title))/2, 0, title, termbox.ColorYellow, titleBg)

	// ── Info bar (row 1) ──
	infoFg, infoBg := termbox.ColorBlack, termbox.ColorBlue
	for x := 0; x < w; x++ {
		termbox.SetCell(x, 1, ' ', infoFg, infoBg)
	}
	info := fmt.Sprintf(" Gen: %-6d | %-8s | Pop: %-6d | Speed: %.1fs/gen ",
		s.generations, statusText(s.running), s.Population(), s.interval.Seconds())
	drawStr(1, 1, info, termbox.ColorWhite, infoBg)

	// ── Sidebar (right side) ──
	gridW := w - sidebarWidth
	sbX := gridW
	sbFg, sbBg := termbox.ColorBlack, termbox.ColorCyan
	for x := sbX; x < w; x++ {
		termbox.SetCell(x, 1, ' ', sbFg, sbBg)
	}

	// Sidebar header
	sbTitle := "  Controls  "
	for x := sbX; x < w; x++ {
		termbox.SetCell(x, 2, ' ', sbFg, sbBg)
	}
	drawStr(sbX+1, 2, sbTitle, termbox.ColorWhite, sbBg)

	// Sidebar content
	sbLines := []string{
		"",
		"  Space   Pause/Resume",
		"  Enter   Toggle cell",
		"  ↑↓←→    Move cursor",
		"  c       Clear grid",
		"  r       Randomize",
		"  + / -   Speed up/down",
		"  p       Place pattern",
		"  q / Esc Quit",
		"",
		"  Grid: ",
		fmt.Sprintf("  %d x %d", s.grid.Rows(), s.grid.Cols()),
	}

	startY := 3
	for i, line := range sbLines {
		y := startY + i
		if y >= h-1 {
			break
		}
		// Clear sidebar row
		for x := sbX; x < w; x++ {
			termbox.SetCell(x, y, ' ', sbFg, sbBg)
		}
		if i >= len(sbLines)-3 {
			drawStr(sbX+1, y, line, termbox.ColorYellow, sbBg)
		} else {
			drawStr(sbX+1, y, line, termbox.ColorWhite, sbBg)
		}
	}

	// ── Grid area (rows 2..h-2, cols 0..gridW-1) ──
	gridStartY := 3
	gridEndY := h - 1
	gridRows := gridEndY - gridStartY

	for r := 0; r < gridRows && r < s.grid.Rows(); r++ {
		for c := 0; c < gridW && c < s.grid.Cols(); c++ {
			var ch rune
			var fg, bg termbox.Attribute

			if s.grid.Cells(r, c) {
				ch = '●'
				fg = termbox.ColorWhite
				bg = termbox.ColorBlue
			} else {
				ch = ' '
				fg = termbox.ColorDefault
				bg = termbox.ColorDefault
			}

			if r == s.cursorR && c == s.cursorC {
				fg = termbox.ColorBlack
				bg = termbox.ColorYellow
			}

			termbox.SetCell(c, gridStartY+r, ch, fg, bg)
		}
	}

	// ── Status bar (row h-1) ──
	statusFg, statusBg := termbox.ColorBlack, termbox.ColorCyan
	for x := 0; x < w; x++ {
		termbox.SetCell(x, h-1, ' ', statusFg, statusBg)
	}
	status := fmt.Sprintf(" Cursor: (%d,%-3d) | %d pattern(s) available | %s",
		s.cursorR, s.cursorC, len(patterns), statusText(s.running))
	drawStr(1, h-1, status, termbox.ColorWhite, statusBg)

	termbox.Flush()
}

func drawStr(x, y int, text string, fg, bg termbox.Attribute) {
	for i, ch := range text {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}

func statusText(running bool) string {
	if running {
		return "▶ Running"
	}
	return "❚❚ Paused"
}

func handleKeyEvent(ev termbox.Event, state *GameState) bool {
	switch {
	case ev.Key == termbox.KeyEsc || ev.Ch == 'q':
		return true
	case ev.Key == termbox.KeyF5 || ev.Ch == 'c':
		state.grid.Reset()
		state.generations = 0
	case ev.Key == termbox.KeyF8 || ev.Ch == 'r':
		state.grid.Randomize()
		state.generations = 0
	case ev.Ch == '+':
		if state.interval > 50*time.Millisecond {
			state.interval /= 2
			if state.interval < 100*time.Millisecond {
				state.interval = 100 * time.Millisecond
			}
		}
	case ev.Ch == '-':
		if state.interval < 5*time.Second {
			state.interval *= 2
			if state.interval > 5*time.Second {
				state.interval = 5 * time.Second
			}
		}
	case ev.Key == termbox.KeySpace:
		state.running = !state.running
	case ev.Key == termbox.KeyArrowUp:
		if state.cursorR > 0 {
			state.cursorR--
		}
	case ev.Key == termbox.KeyArrowDown:
		if state.cursorR < state.grid.Rows()-1 {
			state.cursorR++
		}
	case ev.Key == termbox.KeyArrowLeft:
		if state.cursorC > 0 {
			state.cursorC--
		}
	case ev.Key == termbox.KeyArrowRight:
		if state.cursorC < state.grid.Cols()-1 {
			state.cursorC++
		}
	case ev.Key == termbox.KeyEnter:
		state.grid.Toggle(state.cursorR, state.cursorC)
	case ev.Ch == 'p':
		pattern := patternOverlay(state)
		if pattern != nil {
			PlacePattern(state.grid, state.cursorR, state.cursorC, *pattern)
			state.generations = 0
		}
	}
	return false
}

func main() {
	if err := termbox.Init(); err != nil {
		fmt.Println("Failed to init termbox:", err)
		return
	}
	defer termbox.Close()

	w, h := termbox.Size()
	gridW := w - sidebarWidth
	if gridW < 10 {
		gridW = w - 4
	}
	gridRows := h - 4
	if gridRows < 4 {
		gridRows = 4
	}

	grid := gol.NewGrid(gridRows, gridW)
	grid.Randomize()

	state := &GameState{
		grid:     grid,
		running:  true,
		interval: 1 * time.Second,
	}

	tick := time.NewTicker(state.interval)
	defer tick.Stop()

	// Run PollEvent in a goroutine so the main loop can select on tick
	eventCh := make(chan termbox.Event, 1)
	go func() {
		for {
			eventCh <- termbox.PollEvent()
		}
	}()

	state.Render()

	for {
		select {
		case <-tick.C:
			if state.running {
				state.grid.Evolve()
				state.generations++
				state.Render()
			}
		case ev := <-eventCh:
			switch ev.Type {
			case termbox.EventKey:
				if handleKeyEvent(ev, state) {
					return
				}
				if ev.Ch == '+' || ev.Ch == '-' {
					tick.Stop()
					tick = time.NewTicker(state.interval)
				}
			case termbox.EventResize:
				w, h := termbox.Size()
				gridW := w - sidebarWidth
				gridRows := h - 4
				if gridW < 10 {
					gridW = w - 4
				}
				if gridRows < 4 {
					gridRows = 4
				}
				state.grid.Resize(gridRows, gridW)
			case termbox.EventError:
				return
			}

			state.Render()
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
		w, h := termbox.Size()

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
		header := " Place Pattern (type to filter, ↑↓ navigate, Enter place, Esc cancel) "
		for x := 0; x < w; x++ {
			termbox.SetCell(x, 0, ' ', termbox.ColorBlack, termbox.ColorCyan)
		}
		drawStr(1, 0, header, termbox.ColorWhite, termbox.ColorCyan)

		// Query line
		queryText := " Filter: " + query + "_"
		drawStr(2, 2, queryText, termbox.ColorYellow, termbox.ColorDefault)

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
		footer := fmt.Sprintf(" %d pattern(s) | Cursor: (%d,%d) ", len(filtered), state.cursorR, state.cursorC)
		for x := 0; x < w; x++ {
			termbox.SetCell(x, h-1, ' ', termbox.ColorBlack, termbox.ColorCyan)
		}
		drawStr(1, h-1, footer, termbox.ColorWhite, termbox.ColorCyan)

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
