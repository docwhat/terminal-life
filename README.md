# Game of Life written in Go

> This is Conway's game of life, implemented in Go for the terminal.

## Usage

By default, the grid that'd created matches your terminal window size.

You can use the cursor keys (keyboard arrows) to select a cell. You can use the enter key to toggle the cell's live state.

You can use the space bar to pause the game.

## Life's Rules

1. Any live cell with fewer than two live neighbours dies, as if by underpopulation.
1. Any live cell with two or three live neighbours lives on to the next generation.
1. Any live cell with more than three live neighbours dies, as if by overpopulation.
1. Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

## Development

Use `trunk` (from <https://trunk.io>) to ensure everything is formatted and linted correctly.
