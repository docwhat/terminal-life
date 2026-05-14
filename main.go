package main

import (
	"fmt"
	"time"

	"game-of-life/gol"
	"github.com/nsf/termbox-go"
)

type GameState struct {
	grid        *gol.Grid
	cursorR     int
	cursorC     int
	running     bool
	generations int
}

func (s *GameState) Render() {
	w, h := termbox.Size()

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Draw header
	drawStr(0, 0, fmt.Sprintf(" Game of Life"), termbox.ColorGreen, termbox.ColorDefault)
	drawStr(0, 1, fmt.Sprintf(" Gen: %d | Status: %s", s.generations, statusText(s.running)), termbox.ColorCyan, termbox.ColorDefault)
	drawStr(0, 2, " ↑↓←→: Move | Enter: Toggle | Space: Pause/Resume", termbox.ColorYellow, termbox.ColorDefault)

	// Draw grid
	startY := 4
	for r := 0; r < h-4 && r < s.grid.Rows(); r++ {
		for c := 0; c < w && c < s.grid.Cols(); c++ {
			var ch rune
			var fg, bg termbox.Attribute
			if s.grid.Cells(r, c) {
				ch = '█'
				fg = termbox.ColorWhite
				bg = termbox.ColorBlue
			} else {
				ch = ' '
				fg = termbox.ColorDefault
				bg = termbox.ColorDefault
			}

			// Highlight cursor position
			if r == s.cursorR && c == s.cursorC {
				fg = termbox.ColorBlack
				bg = termbox.ColorYellow
			}

			termbox.SetCell(c, startY+r, ch, fg, bg)
		}
	}

	termbox.Flush()
}

func drawStr(x, y int, text string, fg, bg termbox.Attribute) {
	for i, ch := range text {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}

func statusText(running bool) string {
	if running {
		return "Running"
	}
	return "Paused"
}

func handleKeyEvent(ev termbox.Event, state *GameState) bool {
	switch {
	case ev.Key == termbox.KeyEsc || ev.Ch == 'q':
		return true
	case ev.Key == termbox.KeyF5:
		state.grid.Reset()
		state.generations = 0
	case ev.Key == termbox.KeyF8:
		state.grid.Randomize()
		state.generations = 0
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
	grid := gol.NewGrid(h-4, w)
	grid.Randomize()

	state := &GameState{
		grid:    grid,
		running: true,
	}

	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()

	state.Render()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if handleKeyEvent(ev, state) {
				return
			}
		case termbox.EventResize:
			w, h := termbox.Size()
			state.grid.Resize(h-4, w)
		case termbox.EventError:
			return
		default:
			select {
			case <-tick.C:
				if state.running {
					state.grid.Evolve()
					state.generations++
				}
			default:
			}
		}

		state.Render()
	}
}
