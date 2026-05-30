package main

import "math/rand"

const (
	cellDead    uint8 = 0
	colorManual uint8 = 0xFF // special: manually toggled cell
)

// Grid represents the Game of Life grid with colored cells.
type Grid struct {
	rows  int
	cols  int
	cells [][]uint8 // 0 = dead, 1..N = pattern color index, 0xFF = manual
	ducks [][]bool  // true if cell is a duck (persists for life of cell)
}

// NewGrid creates a new grid with the given dimensions.
func NewGrid(rows, cols int) *Grid {
	cells := make([][]uint8, rows)
	ducks := make([][]bool, rows)

	for i := range cells {
		cells[i] = make([]uint8, cols)
		ducks[i] = make([]bool, cols)
	}

	return &Grid{
		rows:  rows,
		cols:  cols,
		cells: cells,
		ducks: ducks,
	}
}

// Get returns whether the cell at (r, c) is alive.
func (g *Grid) Get(r, c int) bool {
	if r < 0 || r >= g.rows || c < 0 || c >= g.cols {
		return false
	}

	return g.cells[r][c] != cellDead
}

// Cells returns whether the cell at (r, c) is alive.
func (g *Grid) Cells(r, c int) bool {
	return g.Get(r, c)
}

// Color returns the color index for a cell (0 if dead).
func (g *Grid) Color(r, c int) uint8 {
	if r < 0 || r >= g.rows || c < 0 || c >= g.cols {
		return cellDead
	}

	return g.cells[r][c]
}

// IsDuck returns whether the cell at (r, c) is a duck.
func (g *Grid) IsDuck(r, c int) bool {
	if r < 0 || r >= g.rows || c < 0 || c >= g.cols {
		return false
	}

	return g.ducks[r][c]
}

// Set sets a cell alive with the manual color.
func (g *Grid) Set(r, c int) {
	if r >= 0 && r < g.rows && c >= 0 && c < g.cols {
		g.cells[r][c] = colorManual
	}
}

// SetColor sets a cell alive with a specific color index.
func (g *Grid) SetColor(r, c int, color uint8) {
	if r >= 0 && r < g.rows && c >= 0 && c < g.cols {
		g.cells[r][c] = color
	}
}

// Toggle toggles a cell between alive and dead.
func (g *Grid) Toggle(r, c int) {
	if r >= 0 && r < g.rows && c >= 0 && c < g.cols {
		if g.cells[r][c] == cellDead {
			g.cells[r][c] = colorManual

			// 1-in-10000 chance to be a duck
			if rand.Intn(10000) == 0 {
				g.ducks[r][c] = true
			}
		} else {
			g.cells[r][c] = cellDead
			g.ducks[r][c] = false
		}
	}
}

// CountNeighbors returns the number of alive neighbors for cell (r, c).
func (g *Grid) CountNeighbors(r, c int) int {
	count := 0

	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}

			if g.Get(r+dr, c+dc) {
				count++
			}
		}
	}

	return count
}

// Reset clears all cells.
func (g *Grid) Reset() {
	for i := range g.cells {
		for j := range g.cells[i] {
			g.cells[i][j] = cellDead
			g.ducks[i][j] = false
		}
	}
}

// Cols returns the number of columns.
func (g *Grid) Cols() int { return g.cols }

// Evolve advances the grid by one generation.
func (g *Grid) Evolve() {
	g.Tick()
}

// Randomize fills the grid with random alive/dead cells.
func (g *Grid) Randomize() {
	for r := range g.cells {
		for c := range g.cells[r] {
			g.ducks[r][c] = false

			if rand.Intn(2) == 1 {
				g.cells[r][c] = colorManual

				// 1-in-10000 chance for a duck
				if rand.Intn(10000) == 0 {
					g.ducks[r][c] = true
				}
			} else {
				g.cells[r][c] = cellDead
			}
		}
	}
}

// Resize changes the grid dimensions, preserving existing cells.
func (g *Grid) Resize(rows, cols int) {
	newCells := make([][]uint8, rows)
	newDucks := make([][]bool, rows)

	for i := range newCells {
		newCells[i] = make([]uint8, cols)
		newDucks[i] = make([]bool, cols)

		if i < len(g.cells) {
			copy(newCells[i], g.cells[i])
			copy(newDucks[i], g.ducks[i])
		}
	}

	g.rows = rows
	g.cols = cols
	g.cells = newCells
	g.ducks = newDucks
}

// Rows returns the number of rows.
func (g *Grid) Rows() int { return g.rows }

// Tick advances the grid by one generation.
func (g *Grid) Tick() {
	next := make([][]uint8, g.rows)
	nextDucks := make([][]bool, g.rows)

	for i := range next {
		next[i] = make([]uint8, g.cols)
		nextDucks[i] = make([]bool, g.cols)
	}

	for r := 0; r < g.rows; r++ {
		for c := 0; c < g.cols; c++ {
			n := g.CountNeighbors(r, c)

			if g.cells[r][c] != cellDead {
				// Survival: keep existing color and duck status
				if n == 2 || n == 3 {
					next[r][c] = g.cells[r][c]
					nextDucks[r][c] = g.ducks[r][c]
				}
			} else if n == 3 {
				// Birth: inherit color from neighbors
				next[r][c] = g.dominantNeighborColor(r, c)
				// New cells are not ducks
				nextDucks[r][c] = false
			}
		}
	}

	copy(g.cells, next)
	copy(g.ducks, nextDucks)
}

// dominantNeighborColor picks the most common alive-neighbor color.
func (g *Grid) dominantNeighborColor(r, c int) uint8 {
	colorCount := make(map[uint8]int)

	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}

			nr, nc := r+dr, c+dc
			if nr >= 0 && nr < g.rows && nc >= 0 && nc < g.cols {
				col := g.cells[nr][nc]
				if col != cellDead {
					colorCount[col]++
				}
			}
		}
	}

	bestColor := colorManual
	bestCount := 0

	for col, cnt := range colorCount {
		if cnt > bestCount {
			bestCount = cnt
			bestColor = col
		}
	}

	return bestColor
}
