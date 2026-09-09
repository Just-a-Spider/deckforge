# ADR-0008: CWD-First Workspace Resolution and In-TUI Theme Management

## Status
Accepted

## Context
Running DeckForge CLI and companion TUI in day-to-day developer workflows revealed three operational friction points:
1. When running `./bin/deckforge` or `deckforge` inside an arbitrary directory containing presentation decks, the runtime defaulted to a stale cached `ActiveRoot` in `~/.config/deckforge/config.json`, rather than prioritizing the current working directory (`cwd`).
2. Pressing `e` (Edit) on a deck inside the TUI passed the deck directory path (`d.Path`) to `$EDITOR` (e.g. `nvim`, `nano`), crashing with `"is a directory and cannot be opened"`. DeckForge lacked an inline mechanism to switch a presentation's theme on demand without exiting to a text editor.
3. Legacy theme directories (`engine/themes/`) persisted from early development iterations, conflicting with the standardized `/themes/<name>/` architecture established in ADR-0007 and causing path fragmentation.

## Decision
1. **CWD-First Workspace Resolution**:
   - Update `cmd/root.go` so invocation without arguments or with `deckforge tui` immediately resolves `os.Getwd()` to an absolute path and passes it to `runTUI(cwd)`.
   - Update `internal/tui/app.go` to normalize relative roots to absolute paths via `filepath.Abs` and persist `cfg.ActiveRoot = absPath` on launch.
2. **In-TUI Theme Picker Modal & Editor Target Fix**:
   - Update `internal/tui/views_decks.go` so `case "e"` opens `filepath.Join(d.Path, "deck.json")` instead of the directory path.
   - Implement an interactive Theme Picker Modal overlay on `t` keypress in `DecksView`, populated via `ThemeManager.ListThemes(d.Path)`.
   - On theme selection (`Enter`), invoke `workspace.SetDeckTheme`, update in-memory state (`d.Config.Theme`, `d.ThemeName`), automatically recompile the presentation via `Compiler.Build`, and display real-time visual feedback.
   - Synchronize highlighted deck selection to `ThemeStudioView.ActiveDeck` when switching tabs.
3. **Legacy Theme Purging & Path Standardization**:
   - Purge `engine/themes/` from the repository root.
   - Update all theme creation, cloning, and segment-editing paths in `internal/tui/views_theme.go` to target `<workspace>/themes/<name>/` strictly.

## Consequences
- **Positive**: Running `deckforge` in any directory immediately scopes the workspace to that directory; deck authors can swap themes and preview recompilation inside the TUI with one keystroke (`t`); external editors open `deck.json` reliably without directory errors; theme resolution is strictly 3-tiered (workspace, user global, embedded builtin).
- **Negative**: Adds interactive modal state management to `DecksView`.
