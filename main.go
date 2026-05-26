package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

type overlayMode int

const (
	overlayNone overlayMode = iota
	overlayHelp
	overlayPattern
)

type GameState struct {
	grid        *Grid
	cursorR     int
	cursorC     int
	running     bool
	speed       int // 0 = manual, 1-15 = generations per second
	generations int

	// Theme
	theme *Theme

	// Pattern color tracking
	nextColorIdx int // cycles through theme.PatternColors

	// Overlay state
	overlay     overlayMode
	ovQuery     string
	ovHighlight int
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

func (s *GameState) speedText() string {
	if s.speed == 0 {
		return "Manual"
	}
	return fmt.Sprintf("%d/s", s.speed)
}

func (s *GameState) speedInterval() time.Duration {
	return time.Second / time.Duration(s.speed)
}

func (s *GameState) nextPatternColor() uint8 {
	idx := uint8(s.nextColorIdx % len(s.theme.PatternColors))
	s.nextColorIdx++
	return idx + 1 // 1-based index into palette
}

func (s *GameState) cellFg(colorIdx uint8) termbox.Attribute {
	switch colorIdx {
	case cellDead:
		return termbox.ColorDefault
	case colorManual:
		return s.theme.ManualCellFg
	default:
		paletteIdx := int(colorIdx-1) % len(s.theme.PatternColors)
		return s.theme.PatternColors[paletteIdx]
	}
}

func (s *GameState) Render() {
	w, h := termbox.Size()
	if h < 5 || w < 10 {
		return
	}

	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		fmt.Println("Failed to clear screen:", err)
		return
	}

	switch s.overlay {
	case overlayHelp:
		s.renderHelp()
	case overlayPattern:
		s.renderPatternOverlay()
	default:
		s.renderGame()
	}

	if err := termbox.Flush(); err != nil {
		fmt.Println("Failed to flush:", err)
		return
	}
}

func (s *GameState) renderGame() {
	w, h := termbox.Size()
	t := s.theme

	// ── Title bar (row 0) ──
	s.drawBar(0, w, t.TitleFg, t.TitleBg, " ")
	title := " ◆ Game of Life ◆ "
	drawStr((w-len(title))/2, 0, title, t.TitleFg, t.TitleBg)

	// ── Info bar (row 1) ──
	s.drawBar(1, w, t.InfoFg, t.InfoBg, " ")
	info := fmt.Sprintf(" Gen: %-6d │ %-8s │ Pop: %-6d │ Grid: %dx%d │ Speed: %-8s ",
		s.generations, statusText(s.running, s.speed), s.Population(), s.grid.Rows(), s.grid.Cols(), s.speedText())
	drawStr(1, 1, info, t.InfoFg, t.InfoBg)

	// ── Grid area (rows 2..h-2, cols 0..w-1) ──
	gridStartY := 2
	gridEndY := h - 1
	gridRows := gridEndY - gridStartY

	for r := 0; r < gridRows && r < s.grid.Rows(); r++ {
		for c := 0; c < w && c < s.grid.Cols(); c++ {
			var ch rune
			var fg, bg termbox.Attribute

			colorIdx := s.grid.Color(r, c)
			if colorIdx != cellDead {
				ch = t.CellChar
				fg = s.cellFg(colorIdx)
				bg = termbox.ColorDefault
			} else {
				ch = ' '
				fg = termbox.ColorDefault
				bg = termbox.ColorDefault
			}

			// Cursor highlight
			if r == s.cursorR && c == s.cursorC {
				ch = '◉'
				fg = termbox.ColorWhite
				bg = termbox.ColorBlue
			}

			termbox.SetCell(c, gridStartY+r, ch, fg, bg)
		}
	}

	// ── Status bar (row h-1) ──
	s.drawBar(h-1, w, t.StatusFg, t.StatusBg, " ")
	status := fmt.Sprintf(" Cursor: (%d,%-3d) │ %d pattern(s) │ %s │ ? for help ",
		s.cursorR, s.cursorC, len(patterns), statusText(s.running, s.speed))
	drawStr(1, h-1, status, t.StatusFg, t.StatusBg)
}

func (s *GameState) drawBar(y, w int, fg, bg termbox.Attribute, ch string) {
	for x := 0; x < w; x++ {
		termbox.SetCell(x, y, rune(ch[0]), fg, bg)
	}
}

func (s *GameState) renderHelp() {
	w, h := termbox.Size()
	t := s.theme

	// Dialog dimensions
	dialogW := 56
	dialogH := 18
	if dialogW > w-2 {
		dialogW = w - 2
	}
	if dialogH > h-2 {
		dialogH = h - 2
	}
	x0 := (w - dialogW) / 2
	y0 := (h - dialogH) / 2

	// Draw dialog background
	for y := y0; y < y0+dialogH; y++ {
		for x := x0; x < x0+dialogW; x++ {
			termbox.SetCell(x, y, ' ', t.DialogFg, t.DialogBg)
		}
	}

	// Border
	for x := x0; x < x0+dialogW; x++ {
		termbox.SetCell(x, y0, '─', t.DialogFg, t.DialogBg)
		termbox.SetCell(x, y0+dialogH-1, '─', t.DialogFg, t.DialogBg)
	}
	for y := y0; y < y0+dialogH; y++ {
		termbox.SetCell(x0, y, '│', t.DialogFg, t.DialogBg)
		termbox.SetCell(x0+dialogW-1, y, '│', t.DialogFg, t.DialogBg)
	}
	termbox.SetCell(x0, y0, '╭', t.DialogFg, t.DialogBg)
	termbox.SetCell(x0+dialogW-1, y0, '╮', t.DialogFg, t.DialogBg)
	termbox.SetCell(x0, y0+dialogH-1, '╰', t.DialogFg, t.DialogBg)
	termbox.SetCell(x0+dialogW-1, y0+dialogH-1, '╯', t.DialogFg, t.DialogBg)

	// Title
	title := " ◆ Controls ◆ "
	titleX := x0 + (dialogW-len(title))/2
	for i, ch := range title {
		termbox.SetCell(titleX+i, y0+1, ch, t.DialogFg, t.DialogBg)
	}

	// Help lines
	helpLines := []string{
		"",
		"  Space       Advance gen (manual) / Pause",
		"  Enter       Toggle cell at cursor",
		"  ↑ ↓ ← →    Move cursor",
		"  c           Clear grid",
		"  r           Randomize",
		"  + / -       Speed up / down (1-15/s)",
		"  p           Place pattern",
		"  ? / h       Show this help",
		"  q / Esc     Quit",
		"",
	}

	for i, line := range helpLines {
		y := y0 + 3 + i
		if y >= y0+dialogH-1 {
			break
		}
		for j, ch := range line {
			termbox.SetCell(x0+2+j, y, ch, t.DialogFg, t.DialogBg)
		}
	}

	// Footer
	footer := " Press any key to close "
	footerX := x0 + (dialogW-len(footer))/2
	drawStr(footerX, y0+dialogH-1, footer, t.DialogFg, t.DialogBg)
}

func (s *GameState) renderPatternOverlay() {
	w, h := termbox.Size()
	t := s.theme

	// Filter patterns
	var filtered []Pattern
	for _, p := range patterns {
		if s.ovQuery == "" || fuzzyMatch(s.ovQuery, p.Name) {
			filtered = append(filtered, p)
		}
	}
	if s.ovHighlight >= len(filtered) {
		s.ovHighlight = len(filtered) - 1
	}
	if s.ovHighlight < 0 {
		s.ovHighlight = 0
	}

	// Header
	header := " ◆ Place Pattern (type to filter, ↑↓ navigate, Enter place, Esc cancel) ◆ "
	s.drawBar(0, w, t.StatusFg, t.StatusBg, " ")
	drawStr(max(0, (w-len(header))/2), 0, header, t.StatusFg, t.StatusBg)

	// Query line
	queryText := " Filter: " + s.ovQuery + "▌"
	drawStr(2, 2, queryText, t.ManualCellFg, termbox.ColorDefault)

	// Pattern list
	listStart := 4
	listEnd := h - 2
	visible := listEnd - listStart
	if visible <= 0 {
		visible = 1
	}

	// Calculate scroll offset
	scroll := 0
	if s.ovHighlight >= scroll+visible {
		scroll = s.ovHighlight - visible + 1
	}
	if s.ovHighlight < scroll {
		scroll = s.ovHighlight
	}

	for i := scroll; i < len(filtered) && (i-scroll)+listStart < listEnd; i++ {
		y := listStart + (i - scroll)
		var fg, bg termbox.Attribute
		text := "  " + filtered[i].Name
		if i == s.ovHighlight {
			fg = t.StatusBg
			bg = t.ManualCellFg
			text = "▸▸ " + filtered[i].Name
		} else {
			fg = t.DialogFg
			bg = termbox.ColorDefault
		}
		drawStr(2, y, text, fg, bg)
	}

	// Footer
	footer := fmt.Sprintf(" %d pattern(s) │ Cursor: (%d,%d) ", len(filtered), s.cursorR, s.cursorC)
	s.drawBar(h-1, w, t.StatusFg, t.StatusBg, " ")
	drawStr(1, h-1, footer, t.StatusFg, t.StatusBg)
}

func (s *GameState) handleOverlayEvent(ev termbox.Event) *Pattern {
	switch s.overlay {
	case overlayHelp:
		// Any key closes help
		s.overlay = overlayNone
		return nil

	case overlayPattern:
		switch ev.Type {
		case termbox.EventKey:
			switch {
			case ev.Key == termbox.KeyEsc:
				s.overlay = overlayNone
				return nil
			case ev.Key == termbox.KeyEnter:
				s.overlay = overlayNone
				var filtered []Pattern
				for _, p := range patterns {
					if s.ovQuery == "" || fuzzyMatch(s.ovQuery, p.Name) {
						filtered = append(filtered, p)
					}
				}
				if s.ovHighlight >= 0 && s.ovHighlight < len(filtered) {
					return &filtered[s.ovHighlight]
				}
				return nil
			case ev.Key == termbox.KeyArrowUp:
				if s.ovHighlight > 0 {
					s.ovHighlight--
				}
			case ev.Key == termbox.KeyArrowDown:
				var filtered []Pattern
				for _, p := range patterns {
					if s.ovQuery == "" || fuzzyMatch(s.ovQuery, p.Name) {
						filtered = append(filtered, p)
					}
				}
				if s.ovHighlight < len(filtered)-1 {
					s.ovHighlight++
				}
			case ev.Key == termbox.KeyBackspace:
				if len(s.ovQuery) > 0 {
					s.ovQuery = s.ovQuery[:len(s.ovQuery)-1]
				}
			case ev.Ch != 0 && ev.Ch < 128:
				if len(s.ovQuery) < 20 {
					s.ovQuery += string(ev.Ch)
				}
			}
		case termbox.EventError:
			s.overlay = overlayNone
			return nil
		}
	}
	return nil
}

func drawStr(x, y int, text string, fg, bg termbox.Attribute) {
	for i, ch := range text {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}

func statusText(running bool, speed int) string {
	if speed == 0 {
		return "▶ Manual"
	}
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
		if state.speed < 15 {
			state.speed++
		}
	case ev.Ch == '-':
		if state.speed > 0 {
			state.speed--
		}
	case ev.Key == termbox.KeySpace:
		if state.speed == 0 {
			// Manual mode: advance one generation
			state.grid.Evolve()
			state.generations++
		} else {
			state.running = !state.running
		}
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
		state.overlay = overlayPattern
		state.ovQuery = ""
		state.ovHighlight = 0
	case ev.Ch == '?' || ev.Ch == 'h':
		state.overlay = overlayHelp
	}
	return false
}

func main() {
	if err := termbox.Init(); err != nil {
		fmt.Println("Failed to init termbox:", err)
		return
	}
	defer termbox.Close()

	// Detect terminal capabilities and build theme
	theme := NewTheme()

	// Set output mode based on detected capabilities
	if theme.TrueColor {
		termbox.SetOutputMode(termbox.OutputRGB)
	} else {
		termbox.SetOutputMode(termbox.Output256)
	}

	w, h := termbox.Size()
	gridCols := w
	gridRows := h - 4

	grid := NewGrid(gridRows, gridCols)

	state := &GameState{
		grid:         grid,
		theme:        theme,
		running:      true,
		speed:        3, // default 3 generations/second
		nextColorIdx: 0,
	}

	tick := time.NewTicker(state.speedInterval())
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
			if state.running && state.speed > 0 && state.overlay == overlayNone {
				state.grid.Evolve()
				state.generations++
				state.Render()
			}
		case ev := <-eventCh:
			switch ev.Type {
			case termbox.EventKey:
				if state.overlay != overlayNone {
					pattern := state.handleOverlayEvent(ev)
					if pattern != nil {
						colorIdx := state.nextPatternColor()
						PlacePattern(state.grid, state.cursorR, state.cursorC, *pattern, colorIdx)
						state.generations = 0
					}
					state.Render()
					continue
				}
				if handleKeyEvent(ev, state) {
					return
				}
				if ev.Ch == '+' || ev.Ch == '-' {
					tick.Stop()
					if state.speed > 0 {
						tick = time.NewTicker(state.speedInterval())
					} else {
						tick = time.NewTicker(time.Hour) // effectively disabled in manual mode
					}
				}
			case termbox.EventResize:
				w, h := termbox.Size()
				gridCols := w
				gridRows := h - 4
				state.grid.Resize(gridRows, gridCols)
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
