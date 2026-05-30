package main

import "math/rand"

// cellDead represents a dead cell.
const cellDead uint8 = 0

// colorManual is a special color index for manually toggled cells.
const colorManual uint8 = 0xFF

// Grid represents the Game of Life grid with colored cells.
type Grid struct {
	rows  int
	cols  int
	cells [][]uint8 // 0 = dead, 1..N = pattern color index, 0xFF = manual
	ducks [][]bool  // true if cell is a duck (persists for life of cell)
	ages  [][]int   // generations alive (for heat map fading of non-pattern cells)
}

// NewGrid creates a new grid with the given dimensions.
func NewGrid(rows, cols int) *Grid {
	cells := make([][]uint8, rows)
	ducks := make([][]bool, rows)
	ages := make([][]int, rows)

	for i := range cells {
		cells[i] = make([]uint8, cols)
		ducks[i] = make([]bool, cols)
		ages[i] = make([]int, cols)
	}

	return &Grid{
		rows:  rows,
		cols:  cols,
		cells: cells,
		ducks: ducks,
		ages:  ages,
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

// Age returns the number of generations a cell has been alive (0 if dead).
func (g *Grid) Age(r, c int) int {
	if r < 0 || r >= g.rows || c < 0 || c >= g.cols {
		return 0
	}

	return g.ages[r][c]
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
		g.ages[r][c] = 0
	}
}

// SetColor sets a cell alive with a specific color index.
func (g *Grid) SetColor(r, c int, color uint8) {
	if r >= 0 && r < g.rows && c >= 0 && c < g.cols {
		g.cells[r][c] = color
		g.ages[r][c] = 0
	}
}

// Toggle toggles a cell between alive and dead.
func (g *Grid) Toggle(r, c int) {
	if r >= 0 && r < g.rows && c >= 0 && c < g.cols {
		if g.cells[r][c] == cellDead {
			g.cells[r][c] = colorManual
			g.ages[r][c] = 0

			// 1-in-10000 chance to be a duck
			if rand.Intn(10000) == 0 {
				g.ducks[r][c] = true
			}
		} else {
			g.cells[r][c] = cellDead
			g.ducks[r][c] = false
			g.ages[r][c] = 0
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
			g.ages[i][j] = 0
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
			g.ages[r][c] = 0

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
	newAges := make([][]int, rows)

	for i := range newCells {
		newCells[i] = make([]uint8, cols)
		newDucks[i] = make([]bool, cols)
		newAges[i] = make([]int, cols)

		if i < len(g.cells) {
			copy(newCells[i], g.cells[i])
			copy(newDucks[i], g.ducks[i])
			copy(newAges[i], g.ages[i])
		}
	}

	g.rows = rows
	g.cols = cols
	g.cells = newCells
	g.ducks = newDucks
	g.ages = newAges
}

// Rows returns the number of rows.
func (g *Grid) Rows() int { return g.rows }

// Tick advances the grid by one generation.
func (g *Grid) Tick() {
	next := make([][]uint8, g.rows)
	nextDucks := make([][]bool, g.rows)
	nextAges := make([][]int, g.rows)

	for i := range next {
		next[i] = make([]uint8, g.cols)
		nextDucks[i] = make([]bool, g.cols)
		nextAges[i] = make([]int, g.cols)
	}

	for r := 0; r < g.rows; r++ {
		for c := 0; c < g.cols; c++ {
			n := g.CountNeighbors(r, c)

			if g.cells[r][c] != cellDead {
				// Survival: keep existing color, duck status, and increment age
				if n == 2 || n == 3 {
					next[r][c] = g.cells[r][c]
					nextDucks[r][c] = g.ducks[r][c]
					nextAges[r][c] = g.ages[r][c] + 1
				}
			} else if n == 3 {
				// Birth: inherit color from neighbors, age starts at 0
				next[r][c] = g.dominantNeighborColor(r, c)
				// New cells are not ducks
				nextDucks[r][c] = false
				nextAges[r][c] = 0
			}
		}
	}

	copy(g.cells, next)
	copy(g.ducks, nextDucks)
	copy(g.ages, nextAges)
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
