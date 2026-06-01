# Build Plan

> A list of ordered tasks. Each single task will be tackled by one AI Agent run.

## Completed (V1 Core)

- [x] Terminal rendering with tcell (grid, title bar, info bar, status bar)
- [x] Simulation engine (B3/S23 rules, color inheritance, age tracking)
- [x] 12 built-in themes with full color palettes
- [x] Pattern placement (18 patterns including gliders, oscillators, Gosper Gun)
- [x] Interactive controls (pause, step, speed, cursor, toggle cells, clear, randomize)
- [x] Overlay system (help, pattern picker with fuzzy filter, theme picker with fuzzy filter)
- [x] Terminal resize handling
- [x] Clean exit (Ctrl+C, q, Escape)
- [x] CLI flags (--theme, --list-themes, --help)
- [x] Duck mutations (1-in-10,000 chance, 🦆 rendering)
- [x] Cell age fade effect (HSV-based, bright → dim over 20 generations)
- [x] Unit tests (grid, state, theme, args)

## Remaining Tasks

### Task 1: WASM Build

**Goal:** Produce a WASM binary that runs in a browser.

**Steps:**

1. Research tcell WASM compatibility (does it work with `GOOS=js GOARCH=wasm`?)
2. If tcell doesn't support WASM, implement a minimal HTML/JS fallback renderer
3. Create `wasm/` directory with `wasm.wasm`, `wasm_index.html`, and `wasm.js` loader
4. Verify the game runs in a browser (Chrome, Firefox)

**Acceptance criteria:**

- `GOOS=js GOARCH=wasm go build -o wasm/wasm.wasm` succeeds
- Opening `wasm_index.html` in a browser shows the game running
- Keyboard controls work (arrow keys, space, enter, etc.)

**Tests:** Build succeeds (compile-time check). Manual browser verification.

**Blocks:** Task 2 (GitHub Pages)

---

### Task 2: GitHub Pages Deployment

**Goal:** Deploy the WASM build to GitHub Pages so it's shareable via URL.

**Steps:**

1. Create `.github/workflows/deploy-pages.yml`
2. Configure workflow to build WASM and push to `gh-pages` branch (or use Pages action)
3. Enable GitHub Pages in repo settings (if not automated)
4. Verify the live URL loads and the game works

**Acceptance criteria:**

- Workflow triggers on push to main (or manually)
- Game is accessible at `https://docwhat.github.io/terminal-life/` (or similar)
- Game runs in browser at the live URL

**Tests:** Manual verification of live URL.

**Blocks:** None (blocks nothing)

---

### Task 3: GitHub Actions Release Cycle

**Goal:** Automated releases producing binaries for multiple platforms.

**Steps:**

1. Create `.github/workflows/release.yml`
2. Configure Go version, build matrix (linux/amd64, darwin/amd64, darwin/arm64, windows/amd64)
3. Tagging triggers release (e.g., `git tag v0.2.0 && git push --tags`)
4. Artifacts include tarballs with binary

**Acceptance criteria:**

- Creating a tag triggers the workflow
- Release page has binaries for all platforms
- Binaries run correctly

**Tests:** Create a test tag, verify release artifacts.

**Blocks:** Task 4 (versioning)

---

### Task 4: Versioning

**Goal:** Add a visible version string to the app.

**Steps:**

1. Define a version variable (e.g., `var Version = "dev"`)
2. Set version via ldflags during release builds (`-X main.Version=v0.2.0`)
3. Display version in the info bar or help overlay
4. Add `--version` / `-v` flag

**Acceptance criteria:**

- `--version` prints the version and exits
- Info bar or help shows the version string
- Release builds show the tagged version

**Tests:** Unit test for `--version` flag. Manual verification of info bar.

**Blocks:** None

---

### Task 5: Polish Pass (Ongoing)

**Goal:** Improve code quality, readability, and maintainability.

**Areas to explore:**

- Simplify rendering functions (reduce duplication between overlay renderers)
- Extract overlay picker logic into a reusable component
- Review variable naming and function organization
- Add missing Go docs
- Run `trunk check` and `trunk fmt`
- Consider refactoring the event loop for clarity

**Acceptance criteria:** Subjective — code is cleaner, more readable, and easier to modify.

**Tests:** Existing test suite passes. `trunk check` passes.

**Blocks:** None (ongoing)

---

## Escalation Criteria

Send a task back to the architect when:

- Design ambiguity (multiple valid approaches, need a decision)
- Unexpected dependencies or breaking changes
- Test failures that can't be resolved in one pass
- Scope creep (task touches more than intended)
- WASM compatibility requires a major architectural change
