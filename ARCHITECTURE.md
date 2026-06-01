# Architecture

> App structure, data model, risks, and any other decisions made about the project.

## Language & Dependencies

- **Go 1.26.3** (module: `docwhat.org/terminal-life`)
- **tcell v2** — terminal rendering and input handling
- **go-colorful** (indirect) — color manipulation (HSV conversion)

## File Structure

Flat layout (no subdirectories):

| File          | Responsibility                                                           |
| ------------- | ------------------------------------------------------------------------ |
| `main.go`     | Entry point, screen init, event loop, rendering, overlays, key handling  |
| `grid.go`     | Grid data structure, simulation logic (B3/S23), cell colors, ages, ducks |
| `patterns.go` | Named pattern definitions and placement logic                            |
| `theme.go`    | Theme definitions (12 built-in themes) and lookup                        |
| `*_test.go`   | Unit tests for each module                                               |

## Data Model

### Grid

- `cells [][]uint8` — 0 = dead, 1–N = pattern color index, 0xFF = manual toggle
- `ducks [][]bool` — 1-in-10,000 chance per cell; ducks render as 🦆 and persist for the cell's life
- `ages [][]int` — generations alive; used for fade effect on manual cells (bright → dim over 20 gens)

### GameState

- `grid *Grid` — the simulation grid
- `cursorR, cursorC int` — cursor position
- `running bool` — simulation paused/running
- `speed int` — 0 = manual step, 1–15 = generations per second
- `generations int` — generation counter (resets on clear/randomize/pattern placement)
- `theme *Theme` — current color theme
- `nextColorIdx int` — cycles through pattern palette for each placed pattern
- `colorOffset int` — random offset applied on grid reset/theme change for palette variety
- `overlay overlayMode` — current overlay (none, help, pattern picker, theme picker)
- `ovQuery, ovHighlight` — overlay filter query and selection index

### Theme

Each theme defines: UI chrome colors (title, info, status, dialog), cell character, manual cell color, pattern color palette (12 colors), and background.

## Simulation Rules

Standard Conway B3/S23:

- **Birth:** dead cell with exactly 3 alive neighbors → born (inherits dominant neighbor color)
- **Survival:** alive cell with 2 or 3 neighbors → survives (keeps color, increments age)
- **Death:** all other cases → dies

## Key Design Decisions

1. **Flat file structure** — project is small enough that directories add unnecessary complexity.
2. **No persistent state** — everything resets on launch; no save/load.
3. **Color inheritance on birth** — new cells inherit the most common color from alive neighbors, creating visually coherent regions.
4. **Fade effect** — manually toggled cells fade from bright to dim over 20 generations (HSV-based), providing a visual age indicator.
5. **Duck mutations** — 1-in-10,000 chance for a cell to become a duck (🦆), rendered persistently for the cell's lifetime.
6. **Fuzzy search** — pattern and theme overlays support fuzzy matching (case-insensitive subsequence).

## Failure Modes

- Terminal too small (< 5 rows or < 10 cols) — screen clears, game doesn't render
- No Unicode support — box-drawing characters and emoji won't render correctly
- No color support — terminal falls back to monochrome
- WASM build issues — browser may not support required features
- Clean exit always attempted (defer screen.Fini()) unless segfault

## Risks

- **WASM compatibility** — tcell may not fully support WASM; may need a separate rendering path
- **Terminal edge cases** — resize during overlay, interrupt during render
- **AI polish pass** — agentic AI may introduce regressions; test suite is the safety net
