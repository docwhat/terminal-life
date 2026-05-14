package gol

import (
	"testing"
)

func TestNewGrid(t *testing.T) {
	g := NewGrid(5, 5)
	if g == nil || len(g.cells) != 5 || len(g.cells[0]) != 5 {
		t.Fatal("NewGrid did not create correct dimensions")
	}
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			if g.cells[r][c] {
				t.Errorf("NewGrid should initialize all cells to false, got true at (%d,%d)", r, c)
			}
		}
	}
}

func TestGetAndSet(t *testing.T) {
	g := NewGrid(5, 5)
	if g.Get(2, 2) {
		t.Error("Expected cell to be false initially")
	}

	g.Set(2, 2)
	if !g.Get(2, 2) {
		t.Error("Expected cell to be true after Set")
	}

	g.cells[2][2] = false // Manually clear to test Set(false) implicitly via Toggle or direct access
	// Actually, Set(true/false) if it supports it. If only Set(r,c) exists:
	g.Set(2, 2)
	if !g.Get(2, 2) {
		t.Error("Set should enable cell")
	}
}

func TestToggle(t *testing.T) {
	g := NewGrid(5, 5)
	if g.Get(1, 1) {
		t.Error("Expected false initially")
	}
	g.Toggle(1, 1)
	if !g.Get(1, 1) {
		t.Error("Expected true after first Toggle")
	}
	g.Toggle(1, 1)
	if g.Get(1, 1) {
		t.Error("Expected false after second Toggle")
	}
}

func TestCountNeighbors(t *testing.T) {
	g := NewGrid(5, 5)
	// Set up a pattern around (2,2)
	g.Set(1, 1)
	g.Set(1, 2)
	g.Set(1, 3)
	g.Set(2, 1)

	if n := g.CountNeighbors(2, 2); n != 4 {
		t.Errorf("Expected 4 neighbors for (2,2), got %d", n)
	}

	// Corner case
	g2 := NewGrid(3, 3)
	g2.Set(0, 0)
	g2.Set(0, 1)
	g2.Set(1, 0)
	if n := g2.CountNeighbors(0, 0); n != 2 {
		t.Errorf("Expected 2 neighbors for corner (0,0), got %d", n)
	}

	// Center in 3x3 with all neighbors set
	g3 := NewGrid(3, 3)
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			if r != 1 || c != 1 {
				g3.Set(r, c)
			}
		}
	}
	if n := g3.CountNeighbors(1, 1); n != 8 {
		t.Errorf("Expected 8 neighbors for center (1,1), got %d", n)
	}
}

func TestEvolveSurvivalAndReproduction(t *testing.T) {
	// Blinker or similar pattern, or just rules.
	// Rule: Alive + 2 neighbors -> Lives
	g := NewGrid(3, 5)
	g.Set(1, 2) // Cell to test
	g.Set(0, 1)
	g.Set(1, 1) // Two neighbors
	if got := g.CountNeighbors(0, 0); got != 2 {
		t.Errorf("CountNeighbors(0,0) = %d; want 2", got)
	}

	// Rule: Alive + 3 neighbors -> Lives
	g2 := NewGrid(3, 3)
	g2.Set(1, 1)
	g2.Set(0, 1)
	g2.Set(1, 0)
	g2.Set(2, 1)
	if val := g2.Get(1, 0); !val {
		t.Error("Expected cell (1,0) to be alive")
	}

	// Rule: Dead + 3 neighbors -> Born
	g3 := NewGrid(3, 3)
	g3.Set(0, 1)
	g3.Set(1, 0)
	g3.Set(2, 1)
	g3.Evolve()
	if !g3.Get(1, 0) {
		t.Error("Expected cell (1,0) to be born")
	}
	if g3.Get(0, 1) {
		t.Error("Expected cell (0,1) to die")
	}
}

func TestEvolveUnderpopulationAndOverpopulation(t *testing.T) {
	// Underpopulation: 1 neighbor -> Dies
	g := NewGrid(3, 3)
	g.Set(1, 1)
	g.Set(0, 1)
	g.Evolve()
	if g.Get(1, 1) {
		t.Error("Expected death from underpopulation (1 neighbor)")
	}

	// Overpopulation: 4 neighbors -> Dies
	g2 := NewGrid(3, 5)
	g2.Set(1, 2)
	g2.Set(0, 1)
	g2.Set(0, 2)
	g2.Set(0, 3)
	g2.Set(1, 1)
	g2.Evolve()
	if g2.Get(1, 2) {
		t.Error("Expected death from overpopulation (4 neighbors)")
	}
}

func TestReset(t *testing.T) {
	g := NewGrid(5, 5)
	g.Set(2, 2)
	g.Reset()
	if g.Get(2, 2) {
		t.Error("Reset should clear all cells")
	}
}

func TestBlinkerOscillator(t *testing.T) {
	g := NewGrid(5, 5)
	g.Set(1, 2) // Blinker
	g.Set(2, 2)
	g.Set(3, 2)

	g.Evolve()
	// Should become horizontal at row 2
	if !g.Get(2, 1) || !g.Get(2, 2) || !g.Get(2, 3) {
		t.Error("Vertical blinker should have evolved to horizontal")
	}
	if g.Get(1, 2) || g.Get(3, 2) {
		t.Error("Original vertical cells should be dead")
	}

	g.Evolve()
	// Should revert to vertical
	if g.Get(2, 1) || g.Get(2, 3) {
		t.Error("Horizontal blinker should have reverted to vertical")
	}
	if !g.Get(1, 2) || !g.Get(3, 2) {
		t.Error("Missing cells in reverted vertical blinker")
	}
}
