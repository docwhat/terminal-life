# Tickets

> A record of tickets (prompts) used per session or thought as being useful.

## Initial Planning

You are an expert software architect and technical product manager.

I have a rough concept for a new coding project, but I need a solid plan
before I ask AI to write code.

To build this plan, we will use an iterative 20-Questions-style specification
framework.

Instead of me giving you all the details at once, ask me questions we can
define the scope and requirements.

You'll ask for clarification when needed.

Once we finish answering the questions, concisely summarize our conversation by
generating the 6 files that will map out our project.

### Files

Each file should already have some minimal information at the top of each file about what each file should contain.

- `PRODUCT_BRIEF.md`
- `USER_STORIES.md`
- `ARCHITECTURE.md`
- `BUILD_PLAN.md` -- This is the ordered list of tasks that need be done one by one with an AI agent.
- `TICKETS.md`
- `TEST_PLAN.md`

### Questions

#### Round 1 - Vision and Scope

1. One-sentence description?
1. Reference to a similar app, game, or tool?
1. Is this just for me or is it for other users?
1. Minimum viable thing that proves the concept?

#### Round 2 - Tech Stack

1. Language and framework?
1. External libraries allowed? If yes, which?
1. File structure (single file or multiple)?
1. Build tools, or "just open the file"?

#### Round 3 - Mechanics and Behavior

1. Core interactions?
1. State: What data persists? What data is deleted after a session?
1. Failure modes?
1. Success conditions?
1. What should the model NOT do?

#### Round 4 - Scope Discipline

1. V1 features?
1. How do we know when we have finished V1?
1. What's in the "If I have time" pile?
1. Will there be a polish pass?

#### Round 5 - Build Discipline

1. Order of operations: what blocks what?
1. Test plan per task?
1. Escalation: What triggers a need for the task to be sent back to the architect mid-build?

## Future Tickets (Templates)

### WASM Build

> "Implement a WASM build for the terminal Game of Life. Research tcell WASM compatibility. If tcell doesn't support GOOS=js, implement a minimal HTML/JS fallback renderer. Output to wasm/wasm.wasm with wasm_index.html and wasm.js loader."

### GitHub Pages

> "Create a GitHub Actions workflow to deploy the WASM build to GitHub Pages. Configure for manual or push-triggered deployment. Target branch: gh-pages."

### Release Cycle

> "Create a GitHub Actions release workflow. Build binaries for linux/amd64, darwin/amd64, darwin/arm64, windows/amd64 on tag. Produce tarball artifacts."

### Versioning

> "Add versioning support: a --version flag, version display in the info bar, and ldflags injection during release builds."

### Polish Pass

> "Review and simplify the codebase. Focus on: reducing duplication in overlay renderers, extracting reusable picker logic, improving Go docs, and running trunk check/fmt."
