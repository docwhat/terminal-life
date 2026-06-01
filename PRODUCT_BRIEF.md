# Product Brief

> What is being built and why.

**Terminal Game of Life** — a Conway's Game of Life simulation running in the terminal, with themed rendering, pattern placement, and interactive controls. Also targets a WASM build for browser sharing via GitHub Pages.

## Why

- Personal project for fun and experimentation
- Showcase work to others (terminal demo + shareable browser link)
- Experiment with agentic AI for code quality and simplification

## Audience

Primarily the author. Secondary: friends and colleagues pointed to the GitHub Pages WASM demo.

## Reference

Inspired by [sachaos/go-life](https://github.com/sachaos/go-life) but with richer theming, pattern placement, and interactive controls.

## MVP

1. A runnable terminal app with a working Game of Life simulation
2. A WASM build running on GitHub Pages that friends can open in a browser

## Current Status

**V1 core is complete.** Remaining: WASM build, GitHub Pages deployment, GitHub Actions release cycle, versioning, and ongoing polish/refactoring.

## Nice-to-Have (Post-V1)

- Configurable rules (beyond B3/S23)
- Infinite or wrapping playing fields
- Multiple "on" states for cells (extended coloring schemes)
- Other Life variants from Wikipedia
