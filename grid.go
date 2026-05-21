package main

import "math/rand"

type Grid struct {
	rows  int
	cols  int
	cells [][]bool
}

func NewGrid(rows, cols int) *Grid {
	cells := make([][]bool, rows)
	for i := range cells {
		cells[i] = make([]bool, cols)
	}
	return &Grid{
		rows:  rows,
		cols:  cols,
		cells: cells,
	}
}

func (g *Grid) Get(r, c int) bool {
	if r < 0 || r >= g.rows || c < 0 || c >= g.cols {
		return false
	}
	return g.cells[r][c]
}

func (g *Grid) Cells(r, c int) bool {
	return g.Get(r, c)
}

func (g *Grid) Set(r, c int) {
	if r >= 0 && r < g.rows && c >= 0 && c < g.cols {
		g.cells[r][c] = true
	}
}

func (g *Grid) Toggle(r, c int) {
	if r >= 0 && r < g.rows && c >= 0 && c < g.cols {
		g.cells[r][c] = !g.cells[r][c]
	}
}

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

func (g *Grid) Reset() {
	for i := range g.cells {
		for j := range g.cells[i] {
			g.cells[i][j] = false
		}
	}
}

func (g *Grid) Cols() int { return g.cols }

func (g *Grid) Evolve() {
	g.Tick()
}

func (g *Grid) Rows() int { return g.rows }

func (g *Grid) Randomize() {
	for r := range g.cells {
		for c := range g.cells[r] {
			g.cells[r][c] = rand.Intn(2) == 1
		}
	}
}

func (g *Grid) Resize(rows, cols int) {
	newCells := make([][]bool, rows)
	for i := range newCells {
		newCells[i] = make([]bool, cols)
		if i < len(g.cells) {
			copy(newCells[i], g.cells[i])
		}
	}
	g.rows = rows
	g.cols = cols
	g.cells = newCells
}

func (g *Grid) Tick() {
	next := make([][]bool, g.rows)
	for i := range next {
		next[i] = make([]bool, g.cols)
	}

	for r := 0; r < g.rows; r++ {
		for c := 0; c < g.cols; c++ {
			n := g.CountNeighbors(r, c)
			if g.cells[r][c] {
				next[r][c] = n == 2 || n == 3
			} else {
				next[r][c] = n == 3
			}
		}
	}

	copy(g.cells, next)
}
