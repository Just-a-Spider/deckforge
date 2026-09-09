# ADR-0009: Typography Cascade Precedence, Universal Master Index, Vertical Budget Primitives, and Agent Metadata Mutation

## Status
Accepted

## Context
Autonomous agent stress-testing on external codebases (e.g. `SysMon`) and browser verification revealed five critical functional and presentation flaws:
1. **Broken Font Loading**: Google Web Fonts declared in themes were ignored by browsers because the compiler injected theme CSS with `@import` *after* base styles, violating the W3C CSS specification where `@import` must precede all other CSS rules.
2. **Frozen Slide Counter**: The presentation HUD counter remained permanently stuck at `01 / %02d` because `controller.js` queried `#deckPageNum` while `compiler.go` generated `#hudCounter`.
3. **Unstyled & Hardcoded Master Index**: `master_index.js` contained hardcoded Spanish strings ("Índice Maestro de la Cátedra", "Buscar por tema...") and emitted `.index-card` markup lacking corresponding CSS definitions in `hud.css`, rendering unstyled text.
4. **Agent Friction & Monolithic "AI Look" Collisions**: Autonomous agents lacked programmatic deck inspection and metadata mutation tools, forcing them to grep Go source code (`LoadDeck`) and edit `deck.json` via bash `cat`. Furthermore, lack of layout budget guidance led to repetitive left-bordered lime cards and vertical collisions between multi-row grid cards and summary blocks.
5. **TUI Theme Selector Modal Glitches**: ANSI color codes interfered with string padding, producing jagged badge borders, and unwindowed lists caused visual overflow when scrolling themes.

## Decision
1. **Typography Cascade Precedence**:
   - Update `internal/compiler/compiler.go` to inject `<link rel="preconnect">` and `<link rel="stylesheet">` tags directly in `<head>` for theme Google Fonts.
   - Reorder `<style>` generation so theme CSS (including any fallback `@import`) is emitted first, strictly obeying W3C CSS cascade specifications.
2. **Slide Counter Unification**:
   - Update `engine/core/controller.js` and embedded assets to look up `document.getElementById('deckPageNum') || document.getElementById('hudCounter')`.
3. **Universal Master Index & Styling**:
   - Purge hardcoded Spanish strings from `master_index.js`; substitute universal English ("Table of Contents", "Presentation Index", "Search slides by title, tag, or topic...").
   - Extract deck title dynamically from `document.title`.
   - Add rich styles in `engine/core/hud.css` and embedded assets for `.index-module-banner`, `.index-card`, `.index-card.active-slide`, `.index-card-num`, `.index-card-title`, `.index-card-tag`, `.index-card-badge`, and `.index-empty-state`.
4. **Vertical Budget Primitives & Diverse Card Archetypes**:
   - Enforce an explicit 860px vertical canvas content budget in `internal/models/styling_guide.go` (1080px canvas − 80px topbar − 60px HUD − margins = 860px usable height).
   - Add `.slide-content-area` with `max-height: 860px` flex distribution, `.row-fill`, and `.row-auto` in `engine/core/components.css`.
   - Introduce diverse card styles: `.panel-bordered`, `.panel-top-accent`, `.panel-elevated`, and `.panel-outline` to break left-stripe monoculture.
5. **Programmatic Agent Metadata Inspection & Mutation**:
   - Implement `deckforge agent deck-info <deck-path>` and `deckforge agent deck-meta <deck-path> [--title] [--subtitle] [--theme]` in `cmd/agent.go`.
   - Expose corresponding native MCP tools `deckforge_get_deck_info` and `deckforge_set_deck_meta` in `internal/server/mcp.go`.
6. **TUI Theme Selector Modal Alignment & Windowing**:
   - Strip ANSI escape sequences before computing padding in `internal/tui/views_decks.go`.
   - Implement 7-item windowed scrolling with a position indicator `(X/total)` and hide background deck list during modal activation.

## Consequences
- **Positive**: Google Fonts fetch in parallel and apply deterministically; slide counter reflects real-time slide navigation; Master Index drawer is modern, responsive, and localized in universal English; agents can query and update presentation metadata with structured zero-bash commands; 1080p canvas vertical overflow is prevented by clear layout budget rules.
- **Negative**: Slides with dense multi-tier layouts must use `.row-fill` / `.row-auto` or split across slides to respect the 860px budget.
