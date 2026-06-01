# User Stories

> Who uses it? What tasks are users expecting to accomplish?

## Primary User: The Author

### Playing

- As a player, I can start the game and watch patterns evolve automatically.
- As a player, I can pause/resume the simulation.
- As a player, I can step through generations one at a time (manual mode).
- As a player, I can adjust the speed (1–15 generations per second).

### Interacting with the Grid

- As a player, I can move a cursor around the grid with arrow keys.
- As a player, I can toggle cells on/off at the cursor position.
- As a player, I can clear the entire grid.
- As a player, I can randomize the grid.
- As a player, I can place named patterns (gliders, oscillators, guns, etc.) at the cursor.

### Theming

- As a player, I can browse and switch between 12 built-in themes (Gruvbox, Monokai, Dracula, Nord, etc.).
- As a player, I can filter themes by name.
- As a player, I can select a theme via CLI flag (`--theme "Gruvbox Dark"`).
- As a player, I can list all available themes (`--list-themes`).

### Information

- As a player, I can see the current generation count, population, grid size, speed, and active theme.
- As a player, I can access a help overlay with all keybindings.
- As a player, I can see the cursor position.

### Cleanup

- As a player, I can quit cleanly (Ctrl+C, `q`, or Escape) and have the terminal restored to normal.

## Secondary User: Friends (Browser)

- As a visitor, I can open the GitHub Pages link and see the game running.
- As a visitor, I can interact with the game using keyboard controls.
