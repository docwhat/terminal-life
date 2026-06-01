# Test Plan

> A description of how we prove the app works as expected.

## Existing Tests (V1 Core)

### grid_test.go

| Test                                         | What it verifies                                                |
| -------------------------------------------- | --------------------------------------------------------------- |
| `TestNewGrid`                                | Grid initializes with correct dimensions, all cells dead        |
| `TestGetAndSet`                              | Set makes cell alive with `colorManual`                         |
| `TestSetColor`                               | SetColor sets specific color index                              |
| `TestToggle`                                 | Toggle flips between alive/dead                                 |
| `TestCountNeighbors`                         | Neighbor counting at center, corner, and edge                   |
| `TestEvolveSurvivalAndReproduction`          | B3/S23 rules: birth with 3 neighbors, survival with 2-3         |
| `TestEvolveUnderpopulationAndOverpopulation` | Death with 1 neighbor, death with 4+ neighbors                  |
| `TestEvolveColorPreservation`                | Surviving cells keep their color across generations             |
| `TestReset`                                  | Reset clears all cells                                          |
| `TestBlinkerOscillator`                      | Blinker oscillates correctly (vertical → horizontal → vertical) |
| `TestCellAgeOnBirth`                         | New cells start with age 0                                      |
| `TestCellAgeIncrementsOnSurvival`            | Age increments each generation a cell survives                  |
| `TestCellAgeResetsOnDeathAndRebirth`         | Age resets to 0 on rebirth                                      |
| `TestResetClearsAges`                        | Reset clears all ages                                           |
| `TestRandomizeResetsAges`                    | Randomize resets ages for all cells                             |
| `TestResizePreservesAges`                    | Resize preserves existing cell ages                             |
| `TestDominantNeighborColor`                  | Birth inherits most common neighbor color                       |

### state_test.go

| Test                                                | What it verifies                                    |
| --------------------------------------------------- | --------------------------------------------------- |
| `TestColorOffsetRandomizedOnInit`                   | Color offset produces varied values across trials   |
| `TestCellFgColorOffsetApplied`                      | Different offsets produce different palette indices |
| `TestColorOffsetBounds`                             | Offset stays within palette bounds                  |
| `TestPaletteIndexNeverNegative`                     | Formula never produces negative indices             |
| `TestPaletteIndexCoversAllColors`                   | All palette colors reachable from a given position  |
| `TestPaletteIndexCoversAllColorsAtVariousPositions` | Coverage holds across multiple positions            |

### theme_test.go

| Test                     | What it verifies                                |
| ------------------------ | ----------------------------------------------- |
| `TestDefaultTheme`       | Returns first built-in theme                    |
| `TestFindTheme`          | Finds existing themes, returns nil for unknowns |
| `TestBuiltInThemesCount` | At least one theme exists                       |

### args_test.go

| Test                              | What it verifies                                       |
| --------------------------------- | ------------------------------------------------------ |
| `TestParseThemeFromArgsHelpExits` | `--help`, `-h`, `--list-themes`, `-l` return `ErrExit` |

## Remaining Test Gaps

### Core Logic

- [ ] `TestDuckPersistence` — duck status persists across generations until cell dies
- [ ] `TestFadeColor` — fadeColor produces dimmer colors with higher ages
- [ ] `TestHSVConversionRoundtrip` — rgbToHSV → hsvToRGB preserves colors
- [ ] `TestFuzzyMatch` — fuzzy matching works correctly (subsequence, case-insensitive)
- [ ] `TestPlacePattern` — pattern placement respects grid bounds
- [ ] `TestSpeedText` — speedText returns correct strings for all speeds
- [ ] `TestStatusText` — statusText returns correct strings for all states

### Manual Verification (Terminal)

- [ ] Game starts and renders correctly
- [ ] Arrow keys move cursor
- [ ] Enter toggles cells
- [ ] Space pauses/resumes (or steps in manual mode)
- [ ] `+`/`-` adjusts speed
- [ ] `c` clears grid, `r` randomizes
- [ ] `p` opens pattern picker, placement works
- [ ] `t` opens theme picker, switching works
- [ ] `?`/`h` shows help overlay
- [ ] `q`/Escape quits cleanly, terminal restored
- [ ] Terminal resize adapts grid
- [ ] `--theme "Name"` sets theme on launch
- [ ] `--list-themes` lists all themes and exits
- [ ] `--help` shows usage and exits

## Test Strategy Per Task

| Task          | Test Approach                                                                     |
| ------------- | --------------------------------------------------------------------------------- |
| WASM Build    | Compile-time check (`GOOS=js GOARCH=wasm go build`). Manual browser verification. |
| GitHub Pages  | Workflow runs successfully. Live URL loads and game works.                        |
| Release Cycle | Tag triggers workflow. Artifacts download and run.                                |
| Versioning    | Unit test for `--version` flag. Manual check of info bar.                         |
| Polish Pass   | Existing test suite passes. `trunk check` passes.                                 |

## Rule for AI Agents

Each BUILD_PLAN task must include its own acceptance criteria. If a task changes core logic, the AI must write a unit test **before** implementing the fix.
