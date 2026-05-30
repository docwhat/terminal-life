// Package main implements a terminal Game of Life with themed rendering.
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

type overlayMode int

const (
	overlayNone overlayMode = iota
	overlayHelp
	overlayPattern
	overlayTheme
)

// GameState holds all runtime state for the Game of Life.
type GameState struct {
	screen tcell.Screen

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

// Population counts the number of alive cells.
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

// Render draws the current state to the screen.
func (s *GameState) Render() {
	s.screen.Clear()
	w, h := s.screen.Size()

	if h < 5 || w < 10 {
		s.screen.Show()

		return
	}

	switch s.overlay {
	case overlayHelp:
		s.renderHelp()
	case overlayPattern:
		s.renderPatternOverlay()
	case overlayTheme:
		s.renderThemeOverlay()
	default:
		s.renderGame()
	}

	s.screen.Show()
}

// handleResize resizes the grid to fit the new screen dimensions.
func (s *GameState) handleResize(screen tcell.Screen) {
	w, h := screen.Size()
	gridCols := w
	gridRows := h - 4

	s.grid.Resize(gridRows, gridCols)
}

// drawBar draws a full-width bar at row y.
func (s *GameState) drawBar(y, w int, fg, bg tcell.Color, ch string) {
	style := tcell.StyleDefault.Foreground(fg).Background(bg)

	for x := 0; x < w; x++ {
		s.screen.SetContent(x, y, rune(ch[0]), nil, style)
	}
}

// clampHighlight ensures ovHighlight is within bounds.
func (s *GameState) clampHighlight(n int) {
	if n == 0 {
		s.ovHighlight = 0

		return
	}

	if s.ovHighlight >= n {
		s.ovHighlight = n - 1
	}

	if s.ovHighlight < 0 {
		s.ovHighlight = 0
	}
}

// calcScroll calculates the scroll offset for a list.
func (s *GameState) calcScroll(visible int) int {
	scroll := 0

	if s.ovHighlight >= scroll+visible {
		scroll = s.ovHighlight - visible + 1
	}

	if s.ovHighlight < scroll {
		scroll = s.ovHighlight
	}

	return scroll
}

// cellFg returns the foreground color for a cell given its color index and age.
// For non-pattern cells (colorManual), the color fades with age while keeping the same hue.
// For pattern cells, the color is fixed from the pattern palette.
func (s *GameState) cellFg(colorIdx uint8, age int) tcell.Color {
	switch colorIdx {
	case cellDead:
		return s.theme.Background
	case colorManual:
		return s.fadeColor(s.theme.ManualCellFg, age)
	default:
		paletteIdx := int(colorIdx-1) % len(s.theme.PatternColors)

		return s.theme.PatternColors[paletteIdx]
	}
}

// fadeColor returns a faded version of the input color based on age.
// Age 0 is full brightness; higher ages fade toward dim while preserving hue.
func (s *GameState) fadeColor(color tcell.Color, age int) tcell.Color {
	r, g, b := color.RGB()

	h, sat, val := rgbToHSV(float64(r), float64(g), float64(b))

	// Fade parameters: cells fade from bright to dim over maxAge generations.
	const maxAge = 20

	const minVal = 0.3 // minimum brightness (dim)

	// Clamp age to maxAge for the fade curve.

	effectiveAge := age

	if effectiveAge > maxAge {
		effectiveAge = maxAge
	}

	// Linear fade: value goes from 1.0 down to minVal over maxAge generations.
	fade := 1.0 - (float64(effectiveAge)/float64(maxAge))*(1.0-minVal)
	newVal := val * fade

	// Clamp value to [0, 1].
	if newVal > 1.0 {
		newVal = 1.0
	}

	if newVal < 0 {
		newVal = 0
	}

	newR, newG, newB := hsvToRGB(h, sat, newVal)

	return tcell.NewRGBColor(int32(newR), int32(newG), int32(newB))
}

// rgbToHSV converts RGB values (0-255) to HSV (h in degrees, s and v in 0-1).
func rgbToHSV(r, g, b float64) (float64, float64, float64) {
	rNorm := r / 255.0
	gNorm := g / 255.0
	bNorm := b / 255.0

	maxC := rNorm

	if gNorm > maxC {
		maxC = gNorm
	}

	if bNorm > maxC {
		maxC = bNorm
	}

	minC := rNorm

	if gNorm < minC {
		minC = gNorm
	}

	if bNorm < minC {
		minC = bNorm
	}

	d := maxC - minC

	// Hue
	var h float64

	if maxC != minC {
		switch maxC {
		case rNorm:
			h = 60 * (1.0 + (gNorm-bNorm)/d)
		case gNorm:
			h = 60 * (2.0 + (bNorm-rNorm)/d)
		case bNorm:
			h = 60 * (4.0 + (rNorm-gNorm)/d)
		}
	}

	// Value
	v := maxC

	// Saturation
	var s float64

	if maxC != 0 {
		s = d / maxC
	}

	return h, s, v
}

// hsvToRGB converts HSV (h in degrees, s and v in 0-1) to RGB (0-255).
func hsvToRGB(h, s, v float64) (uint8, uint8, uint8) {
	if s == 0 {
		c := uint8(v * 255.0)

		return c, c, c
	}

	// Normalize hue to [0, 360).
	if h < 0 {
		h += 360
	}

	h = h - 360.0*(h/360.0)

	c := v * s
	x := c * (1.0 - absF(h/60.0-2.0*(h/60.0)-1.0))
	m := v - c

	var r, g, b float64

	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return uint8((r + m) * 255.0), uint8((g + m) * 255.0), uint8((b + m) * 255.0)
}

// absF returns the absolute value of a float64.
func absF(x float64) float64 {
	if x < 0 {
		return -x
	}

	return x
}

// speedText returns a human-readable speed string.
func (s *GameState) speedText() string {
	if s.speed == 0 {
		return "Manual"
	}

	return fmt.Sprintf("%d/s", s.speed)
}

// speedInterval returns the tick duration for the current speed.
func (s *GameState) speedInterval() time.Duration {
	return time.Second / time.Duration(s.speed)
}

// nextPatternColor returns the next color index for pattern placement.
func (s *GameState) nextPatternColor() uint8 {
	idx := uint8(s.nextColorIdx % len(s.theme.PatternColors))
	s.nextColorIdx++

	return idx + 1 // 1-based index into palette
}

// handleOverlayEvent processes a key event while an overlay is active.
func (s *GameState) handleOverlayEvent(ev *tcell.EventKey) interface{} {
	switch s.overlay {
	case overlayHelp:
		s.overlay = overlayNone

		return nil
	case overlayPattern:
		return s.handlePatternOverlayEvent(ev)
	case overlayTheme:
		return s.handleThemeOverlayEvent(ev)
	}

	return nil
}

// handlePatternOverlayEvent processes key events in the pattern picker.
func (s *GameState) handlePatternOverlayEvent(ev *tcell.EventKey) interface{} {
	switch {
	case ev.Key() == tcell.KeyEscape:
		s.overlay = overlayNone

		return nil
	case ev.Key() == tcell.KeyEnter:
		s.overlay = overlayNone
		filtered := s.filterPatterns()

		if s.ovHighlight >= 0 && s.ovHighlight < len(filtered) {
			return filtered[s.ovHighlight]
		}

		return nil
	case ev.Key() == tcell.KeyUp:
		if s.ovHighlight > 0 {
			s.ovHighlight--
		}
	case ev.Key() == tcell.KeyDown:
		filtered := s.filterPatterns()

		if s.ovHighlight < len(filtered)-1 {
			s.ovHighlight++
		}
	case ev.Key() == tcell.KeyBackspace:
		if len(s.ovQuery) > 0 {
			s.ovQuery = s.ovQuery[:len(s.ovQuery)-1]
		}
	case ev.Rune() != 0 && ev.Rune() < 128:
		if len(s.ovQuery) < 20 {
			s.ovQuery += string(ev.Rune())
		}
	}

	return nil
}

// handleThemeOverlayEvent processes key events in the theme picker.
func (s *GameState) handleThemeOverlayEvent(ev *tcell.EventKey) interface{} {
	switch {
	case ev.Key() == tcell.KeyEscape:
		s.overlay = overlayNone

		return nil
	case ev.Key() == tcell.KeyEnter:
		s.overlay = overlayNone
		filtered := s.filterThemes()

		if s.ovHighlight >= 0 && s.ovHighlight < len(filtered) {
			return filtered[s.ovHighlight]
		}

		return nil
	case ev.Key() == tcell.KeyUp:
		if s.ovHighlight > 0 {
			s.ovHighlight--
		}
	case ev.Key() == tcell.KeyDown:
		filtered := s.filterThemes()

		if s.ovHighlight < len(filtered)-1 {
			s.ovHighlight++
		}
	case ev.Key() == tcell.KeyBackspace:
		if len(s.ovQuery) > 0 {
			s.ovQuery = s.ovQuery[:len(s.ovQuery)-1]
		}
	case ev.Rune() != 0 && ev.Rune() < 128:
		if len(s.ovQuery) < 20 {
			s.ovQuery += string(ev.Rune())
		}
	}

	return nil
}

// filterPatterns returns patterns matching the current query.
func (s *GameState) filterPatterns() []Pattern {
	var filtered []Pattern

	for _, p := range patterns() {
		if s.ovQuery == "" || fuzzyMatch(s.ovQuery, p.Name) {
			filtered = append(filtered, p)
		}
	}

	return filtered
}

// filterThemes returns themes matching the current query.
func (s *GameState) filterThemes() []*Theme {
	allThemes := builtInThemes()

	var filtered []*Theme

	for _, th := range allThemes {
		if s.ovQuery == "" || fuzzyMatch(s.ovQuery, th.Name) {
			filtered = append(filtered, th)
		}
	}

	return filtered
}

// drawDialog draws a bordered dialog box.
func (s *GameState) drawDialog(x0, y0, w, h int, t *Theme) {
	style := tcell.StyleDefault.Foreground(t.DialogFg).Background(t.DialogBg)

	// Background
	for y := y0; y < y0+h; y++ {
		for x := x0; x < x0+w; x++ {
			s.screen.SetContent(x, y, ' ', nil, style)
		}
	}

	// Border
	for x := x0 + 1; x < x0+w-1; x++ {
		s.screen.SetContent(x, y0, '─', nil, style)
		s.screen.SetContent(x, y0+h-1, '─', nil, style)
	}

	for y := y0 + 1; y < y0+h-1; y++ {
		s.screen.SetContent(x0, y, '│', nil, style)
		s.screen.SetContent(x0+w-1, y, '│', nil, style)
	}

	s.screen.SetContent(x0, y0, '╭', nil, style)
	s.screen.SetContent(x0+w-1, y0, '╮', nil, style)
	s.screen.SetContent(x0, y0+h-1, '╰', nil, style)
	s.screen.SetContent(x0+w-1, y0+h-1, '╯', nil, style)
}

// renderGame draws the main game view.
func (s *GameState) renderGame() {
	w, h := s.screen.Size()
	t := s.theme

	// ── Title bar (row 0) ──
	s.drawBar(0, w, t.TitleFg, t.TitleBg, " ")

	title := " ◆ Game of Life ◆ "
	drawStr(s.screen, (w-len(title))/2, 0, title, t.TitleFg, t.TitleBg)

	// ── Info bar (row 1) ──
	s.drawBar(1, w, t.InfoFg, t.InfoBg, " ")
	info := fmt.Sprintf(" Gen: %-6d │ %-8s │ Pop: %-6d │ Grid: %dx%d │ Speed: %-8s │ Theme: %-20s ",
		s.generations, statusText(s.running, s.speed), s.Population(),
		s.grid.Rows(), s.grid.Cols(), s.speedText(), t.Name)
	drawStr(s.screen, 1, 1, info, t.InfoFg, t.InfoBg)

	// ── Grid area (rows 2..h-2, cols 0..w-1) ──
	gridStartY := 2
	gridEndY := h - 1
	gridRows := gridEndY - gridStartY

	for r := 0; r < gridRows && r < s.grid.Rows(); r++ {
		for c := 0; c < w && c < s.grid.Cols(); c++ {
			colorIdx := s.grid.Color(r, c)
			isDuck := s.grid.IsDuck(r, c)
			age := s.grid.Age(r, c)

			var ch rune

			var fg, bg tcell.Color

			if colorIdx != cellDead {
				if isDuck {
					ch = '🦆'
				} else {
					ch = t.CellChar
				}

				fg = s.cellFg(colorIdx, age)
				bg = t.Background
			} else {
				ch = ' '
				fg = t.Background
				bg = t.Background
			}

			// Cursor highlight
			if r == s.cursorR && c == s.cursorC {
				ch = '◉'
				fg = tcell.ColorWhite
				bg = tcell.ColorBlue
			}

			style := tcell.StyleDefault.Foreground(fg).Background(bg)
			s.screen.SetContent(c, gridStartY+r, ch, nil, style)
		}
	}

	// ── Status bar (row h-1) ──
	s.drawBar(h-1, w, t.StatusFg, t.StatusBg, " ")
	status := fmt.Sprintf(" Cursor: (%d,%-3d) │ %d pattern(s) │ %s │ ? help │ p pattern │ t theme ",
		s.cursorR, s.cursorC, len(patterns()), statusText(s.running, s.speed))
	drawStr(s.screen, 1, h-1, status, t.StatusFg, t.StatusBg)
}

// renderHelp draws the help dialog.
func (s *GameState) renderHelp() {
	w, h := s.screen.Size()
	t := s.theme

	// Dialog dimensions
	dialogW := 58
	dialogH := 19

	if dialogW > w-2 {
		dialogW = w - 2
	}

	if dialogH > h-2 {
		dialogH = h - 2
	}

	x0 := (w - dialogW) / 2
	y0 := (h - dialogH) / 2

	s.drawDialog(x0, y0, dialogW, dialogH, t)

	// Title
	title := " ◆ Controls ◆ "
	titleX := x0 + (dialogW-len(title))/2

	for i, ch := range title {
		s.screen.SetContent(titleX+i, y0+1, ch, nil, tcell.StyleDefault.Foreground(t.DialogFg).Background(t.DialogBg))
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
		"  t           Change theme",
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
			s.screen.SetContent(x0+2+j, y, ch, nil, tcell.StyleDefault.Foreground(t.DialogFg).Background(t.DialogBg))
		}
	}

	// Footer
	footer := " Press any key to close "
	footerX := x0 + (dialogW-len(footer))/2
	drawStr(s.screen, footerX, y0+dialogH-1, footer, t.DialogFg, t.DialogBg)
}

// renderPatternOverlay draws the pattern picker.
func (s *GameState) renderPatternOverlay() {
	w, h := s.screen.Size()
	t := s.theme

	filtered := s.filterPatterns()
	s.clampHighlight(len(filtered))

	// Header
	header := " ◆ Place Pattern (type to filter, ↑↓ navigate, Enter place, Esc cancel) ◆ "

	s.drawBar(0, w, t.StatusFg, t.StatusBg, " ")
	drawStr(s.screen, intMax(0, (w-len(header))/2), 0, header, t.StatusFg, t.StatusBg)

	// Query line
	queryText := " Filter: " + s.ovQuery + "▌"
	drawStr(s.screen, 2, 2, queryText, t.ManualCellFg, t.Background)

	// Pattern list
	listStart := 4
	listEnd := h - 2
	visible := listEnd - listStart

	if visible <= 0 {
		visible = 1
	}

	scroll := s.calcScroll(visible)

	for i := scroll; i < len(filtered) && (i-scroll)+listStart < listEnd; i++ {
		y := listStart + (i - scroll)

		var fg, bg tcell.Color

		text := "  " + filtered[i].Name

		if i == s.ovHighlight {
			fg = t.StatusBg
			bg = t.ManualCellFg
			text = "▸▸ " + filtered[i].Name
		} else {
			fg = t.DialogFg
			bg = t.Background
		}

		drawStr(s.screen, 2, y, text, fg, bg)
	}

	// Footer
	footer := fmt.Sprintf(" %d pattern(s) │ Cursor: (%d,%d) ", len(filtered), s.cursorR, s.cursorC)
	s.drawBar(h-1, w, t.StatusFg, t.StatusBg, " ")
	drawStr(s.screen, 1, h-1, footer, t.StatusFg, t.StatusBg)
}

// renderThemeOverlay draws the theme picker.
func (s *GameState) renderThemeOverlay() {
	w, h := s.screen.Size()
	t := s.theme

	filtered := s.filterThemes()
	s.clampHighlight(len(filtered))

	// Header
	header := " ◆ Choose Theme (type to filter, ↑↓ navigate, Enter select, Esc cancel) ◆ "

	s.drawBar(0, w, t.StatusFg, t.StatusBg, " ")
	drawStr(s.screen, intMax(0, (w-len(header))/2), 0, header, t.StatusFg, t.StatusBg)

	// Query line
	queryText := " Filter: " + s.ovQuery + "▌"
	drawStr(s.screen, 2, 2, queryText, t.ManualCellFg, t.Background)

	// Theme list
	listStart := 4
	listEnd := h - 2
	visible := listEnd - listStart

	if visible <= 0 {
		visible = 1
	}

	scroll := s.calcScroll(visible)

	for i := scroll; i < len(filtered) && (i-scroll)+listStart < listEnd; i++ {
		y := listStart + (i - scroll)

		var fg, bg tcell.Color

		text := "  " + filtered[i].Name

		if i == s.ovHighlight {
			fg = t.StatusBg
			bg = t.ManualCellFg
			text = "▸▸ " + filtered[i].Name
		} else {
			fg = t.DialogFg
			bg = t.Background
		}

		drawStr(s.screen, 2, y, text, fg, bg)
	}

	// Footer
	footer := fmt.Sprintf(" %d theme(s) │ Current: %s ", len(filtered), t.Name)
	s.drawBar(h-1, w, t.StatusFg, t.StatusBg, " ")
	drawStr(s.screen, 1, h-1, footer, t.StatusFg, t.StatusBg)
}

// drawStr draws a string at the given position with foreground and background colors.
func drawStr(screen tcell.Screen, x, y int, text string, fg, bg tcell.Color) {
	style := tcell.StyleDefault.Foreground(fg).Background(bg)

	for i, ch := range text {
		screen.SetContent(x+i, y, ch, nil, style)
	}
}

// statusText returns a human-readable status string.
func statusText(running bool, speed int) string {
	if speed == 0 {
		return "▶ Manual"
	}

	if running {
		return "▶ Running"
	}

	return "❚❚ Paused"
}

// handleKeyEvent processes a key event in the main game view.
func handleKeyEvent(ev *tcell.EventKey, state *GameState) bool {
	switch {
	case ev.Key() == tcell.KeyEscape || ev.Rune() == 'q':
		return true
	case ev.Key() == tcell.KeyF5 || ev.Rune() == 'c':
		state.grid.Reset()
		state.generations = 0
	case ev.Key() == tcell.KeyF8 || ev.Rune() == 'r':
		state.grid.Randomize()
		state.generations = 0
	case ev.Rune() == '+':
		if state.speed < 15 {
			state.speed++
		}
	case ev.Rune() == '-':
		if state.speed > 0 {
			state.speed--
		}
	case ev.Rune() == ' ':
		if state.speed == 0 {
			state.grid.Evolve()
			state.generations++
		} else {
			state.running = !state.running
		}
	case ev.Key() == tcell.KeyUp:
		if state.cursorR > 0 {
			state.cursorR--
		}
	case ev.Key() == tcell.KeyDown:
		if state.cursorR < state.grid.Rows()-1 {
			state.cursorR++
		}
	case ev.Key() == tcell.KeyLeft:
		if state.cursorC > 0 {
			state.cursorC--
		}
	case ev.Key() == tcell.KeyRight:
		if state.cursorC < state.grid.Cols()-1 {
			state.cursorC++
		}
	case ev.Key() == tcell.KeyEnter:
		state.grid.Toggle(state.cursorR, state.cursorC)
	case ev.Rune() == 'p':
		state.overlay = overlayPattern
		state.ovQuery = ""
		state.ovHighlight = 0
	case ev.Rune() == 't':
		state.overlay = overlayTheme
		state.ovQuery = ""
		state.ovHighlight = 0
	case ev.Rune() == '?' || ev.Rune() == 'h':
		state.overlay = overlayHelp
	}

	return false
}

// printThemeList prints all available theme names to stdout.
func printThemeList() {
	fmt.Println("Available themes:")

	for _, t := range builtInThemes() {
		fmt.Printf("  %s\n", t.Name)
	}
}

// parseThemeFromArgs returns the theme selected via CLI flags, or nil.
func parseThemeFromArgs() *Theme {
	for i, arg := range os.Args[1:] {
		switch arg {
		case "--list-themes", "-l":
			printThemeList()

			return nil
		case "--help", "-h":
			fmt.Println("Usage: terminal-life [OPTIONS]")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  -t, --theme NAME   Select theme by name")
			fmt.Println("  -l, --list-themes  List all available themes")
			fmt.Println("  -h, --help         Show this help message")

			return nil
		case "--theme", "-t":
			if i+1 < len(os.Args[1:]) {
				themeName := os.Args[2+i]
				if t := findTheme(themeName); t != nil {
					return t
				}

				fmt.Fprintf(os.Stderr, "Unknown theme: %s\n", themeName)
				fmt.Fprintln(os.Stderr, "Use --list-themes to see available themes.")
				os.Exit(1)
			}

			fmt.Fprintln(os.Stderr, "Flag --theme requires a theme name.")
			os.Exit(1)
		}
	}

	return nil
}

// intMax returns the larger of a and b.
func intMax(a, b int) int {
	if a > b {
		return a
	}

	return b
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

func main() {
	theme := parseThemeFromArgs()
	if theme == nil {
		theme = defaultTheme()
	}

	screen := initScreen(theme)
	defer screen.Fini()

	state := initState(screen, theme)

	tick := time.NewTicker(state.speedInterval())
	defer tick.Stop()

	state.Render()

	runEventLoop(state, screen, tick)
}

// initScreen creates and initializes a tcell screen with the given theme.
func initScreen(theme *Theme) tcell.Screen {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Println("Failed to create screen:", err)
		os.Exit(1)
	}

	if err := screen.Init(); err != nil {
		fmt.Println("Failed to init screen:", err)
		os.Exit(1)
	}

	screen.SetStyle(tcell.StyleDefault.Background(theme.Background).Foreground(theme.Background))

	return screen
}

// initState creates the GameState with an initialized grid.
func initState(screen tcell.Screen, theme *Theme) *GameState {
	w, h := screen.Size()
	gridCols := w
	gridRows := h - 4

	grid := NewGrid(gridRows, gridCols)
	grid.Randomize()

	// Seed random for duck mutations
	rand.Seed(time.Now().UnixNano())

	return &GameState{
		screen:       screen,
		grid:         grid,
		theme:        theme,
		running:      true,
		speed:        3,
		nextColorIdx: 0,
	}
}

// runEventLoop processes screen events until the user quits.
func runEventLoop(state *GameState, screen tcell.Screen, tick *time.Ticker) {
	// Event channel from PollEvent goroutine
	eventCh := make(chan tcell.Event, 1)

	go func() {
		for {
			eventCh <- screen.PollEvent()
		}
	}()

	for {
		select {
		case <-tick.C:
			handleTick(state)
		case ev := <-eventCh:
			if !handleScreenEvent(ev, state, screen, tick) {
				return
			}
		}
	}
}

// handleTick advances the simulation if running.
func handleTick(state *GameState) {
	if state.running && state.speed > 0 && state.overlay == overlayNone {
		state.grid.Evolve()
		state.generations++
		state.Render()
	}
}

// handleScreenEvent processes a single screen event. Returns false to quit.
func handleScreenEvent(ev tcell.Event, state *GameState, screen tcell.Screen, tick *time.Ticker) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		return handleKeyEventEvent(ev, state, screen, tick)
	case *tcell.EventResize:
		state.handleResize(screen)
		state.Render()
	case *tcell.EventInterrupt:
		return false
	case *tcell.EventError:
		return false
	}

	return true
}

// handleKeyEventEvent processes a key event. Returns false to quit.
func handleKeyEventEvent(ev *tcell.EventKey, state *GameState, screen tcell.Screen, tick *time.Ticker) bool {
	if state.overlay != overlayNone {
		return handleOverlayKeyEvent(ev, state, screen)
	}

	if handleKeyEvent(ev, state) {
		return false
	}

	if ev.Rune() == '+' || ev.Rune() == '-' {
		tick.Stop()

		if state.speed > 0 {
			tick = time.NewTicker(state.speedInterval())
		} else {
			tick = time.NewTicker(time.Hour) // effectively disabled in manual mode
		}
	}

	state.Render()

	return true
}

// handleOverlayKeyEvent processes a key event while an overlay is active.
func handleOverlayKeyEvent(ev *tcell.EventKey, state *GameState, screen tcell.Screen) bool {
	result := state.handleOverlayEvent(ev)

	switch v := result.(type) {
	case Pattern:
		if v.Name != "" {
			colorIdx := state.nextPatternColor()
			PlacePattern(state.grid, state.cursorR, state.cursorC, v, colorIdx)
			state.generations = 0
		}
	case *Theme:
		if v != nil {
			state.theme = v
			screen.SetStyle(tcell.StyleDefault.Background(v.Background).Foreground(v.Background))
		}
	}

	state.Render()

	return true
}
