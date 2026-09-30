# ADR-0016: Universal Workspace Anchoring, Upward Theme Traversal, and Embedded Node CDP Runner

## Status
Accepted

## Context
Following the implementation of the Chromedp CDP and Node CDP export engines (ADR-0015), real-world multi-workspace testing in consumer repositories (e.g. `SysMon/.deckforge/decks/sysmon-overview`) revealed two critical failures when executing DeckForge from global installations (`~/.local/bin/deckforge` or `~/go/bin/deckforge`):
1. **Silent Fallback and Color Mismatch**: When running `deckforge export .deckforge/decks/sysmon-overview`, `NewThemeManager` was initialized with the relative deck path (`.deckforge/decks/sysmon-overview`). Its resolution logic only looked in `.deckforge/decks/sysmon-overview/.deckforge/themes/`, missing the workspace root's `.deckforge/themes/sysmon-hud`. Because `ResolveTheme` silently fell back to `academic-crimson`, the deck was compiled with light surfaces and red accents instead of dark telemetry palettes (`#080c12` canvas and amber/cyan HUD accents).
2. **Missing Companion Runner**: `deckforge export --engine node` failed in consumer repositories with `Export error: scripts/export-pdf.mjs runner not found` because the companion script was only stored as a loose file in the development repository.

## Decision
1. **Universal Workspace Anchoring (`internal/workspace/root.go`)**:
   - Implemented `DetectWorkspaceRoot(targetPath string) string` with upward recursive filesystem traversal.
   - Automatically detects `.deckforge/`, `deckforge.json`, or `decks/` & `themes/` co-presence. If `targetPath` is within or equal to `.deckforge/`, it properly resolves the enclosing project root.
2. **Upward Tree Traversal for Themes and Components**:
   - Updated `ThemeManager.ListThemes(deckPath)` and `ComponentManager.ListComponents(deckPath)` to recursively collect all ancestor directories from root down to `deckPath`, scanning each level for `.deckforge/themes` (or components) and legacy `themes/` folders.
   - Added current working directory (`cwd`) scanning as a deterministic baseline.
   - Eliminated silent fallback in `ResolveTheme`: any fallback now prints a clear diagnostic warning to `os.Stderr` displaying searched scopes and available themes.
3. **Embedded Portable Node Runner**:
   - Embedded `export-pdf.mjs` directly into the `internal/exporter` package using Go's `//go:embed`.
   - In `exportWithNode`, DeckForge checks local repo overrides first; if absent, it extracts the embedded script to `os.UserCacheDir()/deckforge/export-pdf.mjs` and executes via `node`. DeckForge remains a single zero-dependency binary.
4. **Command Updates**:
   - Updated `runBuild`, `runServe`, `runExport`, `runList`, `runTheme`, and `runComponent` in `cmd/root.go` to anchor managers with `workspace.DetectWorkspaceRoot`.

## Consequences
- **Positive**:
  - Full theme color fidelity is preserved regardless of how deeply nested the deck path argument is.
  - Portable `--engine node` export works globally without requiring external script files.
  - Zero silent failures: developers and agents receive immediate diagnostic feedback if a theme is unresolvable.
- **Negative**:
  - Binary size increases by ~6.8 KB due to embedded Node script.
