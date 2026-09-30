# ADR-0014: User-Agnostic Workspaces, Scoped Discovery, and 1:1 Landscape PDF Engine

## Status
Accepted

## Context
When running DeckForge in external consumer repositories (such as `/home/andre/Desktop/Projects/SysMon`), several issues were observed:
1. **Repository Context Leakage**:
   The `deckforge` agent in consumer projects attempted to introspect and extract information from the `deckforge` development repository. This occurred because `models.LoadWorkspaceConfigForRoot("")` prioritized persisted global `activeRoot` over the current working directory (`cwd`), and because `SKILL.md` instructed agents to call `deckforge book log`, which is private to the DeckForge development repository.
2. **Aggressive Multi-Directory Crawling**:
   `deckforge list` printed 15 global theme presets by default and traversed depth-2 subdirectories across arbitrary project folders, polluting the listing with unrelated presentations.
3. **Hardcoded Machine and User Paths**:
   `internal/exporter/exporter.go` hardcoded `/home/andre/.local/bin/google-chrome`. `book/SUMMARY.md` used absolute `file:///home/andre/...` paths that broke portability on any other developer machine.
4. **Vertical Multi-Page PDF Distortion**:
   Chromium headless PDF export defaulted to Letter portrait (8.5" x 11" / 612x792 pts). A 3-slide presentation was severed into 7 vertical pages with white backgrounds and broken formatting due to missing `@page { size: 1920px 1080px; }`, missing `print-color-adjust: exact`, and unsuppressed DOM UI elements.

## Decision
1. **Scoped Workspace Listing (`cmd/root.go`, `internal/workspace/discovery.go`)**:
   - Limit `deckforge list` to presentations within the current workspace root (`.deckforge/decks/`, `./decks/`, or root if it is a deck).
   - Omit the 15-preset global theme dump unless explicitly requested via `deckforge list --themes`.
   - Restrict `ScanDecksInRoot` to immediate child directories (depth 1 only), eliminating runaway depth-2 traversal into foreign repositories.
   - Format paths relative to `cwd`.
2. **User/System-Agnostic Workspace Configuration (`internal/models/workspace.go`)**:
   - Guarantee `targetRoot` defaults to current working directory (`cwd`) when `root` is empty string.
   - Prevent a persisted `ActiveRoot` from one repository from hijacking another repository.
3. **Zero-Engine-Dependency Agent Protocol (`SKILL.md`, `cmd/agent.go`)**:
   - Explicitly instruct agents that DeckForge is a self-contained installed CLI/MCP tool.
   - Forbid agents from searching for, cloning, or reading the `deckforge` engine repository.
   - Mandate runtime introspection via `deckforge agent catalog`, `deckforge agent styling-guide`, and `deckforge agent theme-list`.
   - Remove `deckforge book log` from end-user agent skills.
4. **1:1 1080p Landscape PDF Export Engine (`internal/exporter/exporter.go`, `engine/core/stage.css`)**:
   - Remove all hardcoded user paths. Discover Chrome across `$CHROME_BIN`, `$CHROMIUM_BIN`, `$BROWSER`, PATH lookup, dynamic `os.UserHomeDir()` (`~/.local/bin`), and standard cross-platform OS installation directories (Linux, macOS, Windows).
   - Use `--headless=new` with modern flags (`--disable-gpu`, `--no-pdf-header-footer`, `--run-all-compositor-stages-before-draw`, `--virtual-time-budget=3000`).
   - Add `@page { size: 1920px 1080px; margin: 0; }`.
   - In `@media print`: enforce `print-color-adjust: exact`, preserve `background-color: var(--bg-canvas)`, lock `.slide` to `1920px x 1080px` with `page-break-inside: avoid` and `break-after: page`, and completely hide non-slide DOM elements (`#captureToast`, `.deck-hud`, `.keyboard-help-modal`, `.master-index-drawer`, `.hairlines::before`, `.hairlines::after`).
5. **Portable Living Development Book**:
   - Convert `SUMMARY.md` links to relative paths.
   - Update `cmd/book.go` to resolve both relative paths and `file://` URIs against `bookDir`.

## Consequences
- **Positive**:
  - `deckforge list` is clean, instantaneous, and scoped strictly to the current project.
  - Zero context leakage: agents operate cleanly within consumer projects without wandering into `deckforge`.
  - PDF generation produces exact 1:1 1920x1080 landscape slides (1440x810 pts at 72dpi), matching the HTML presentation slide-for-slide.
  - Complete portability across different developer usernames, Linux distributions, macOS, and Windows.
- **Negative**:
  - Developers wanting global themes in `deckforge list` must pass `--themes` or run `deckforge theme list`.
