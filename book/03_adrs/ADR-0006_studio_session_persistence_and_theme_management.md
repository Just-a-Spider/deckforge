# ADR-0006: Studio Session State Persistence and Declarative Theme Management

## Status
Accepted

## Context
When editing presentations in DeckForge Studio or modifying slide layouts, Server-Sent Events (SSE) recompile and refresh the browser viewport. Previously, this caused the browser to lose its active Studio state (`this.isActive = false`), disabling visual outlines, inspector drawers, and requiring manual reactivation via keypress `S`. Additionally, switching a presentation's theme required manually modifying `deck.json`, and autonomous agents had no primitives for scaffolding new presentations.

## Decision
1. **Client-Side Session Continuity**:
   - Utilize browser `sessionStorage` (`df_studio_active`, `df_inspector_open`, `df_last_slide`) to record workspace state on every toggle and SSE event.
   - Automatically restore Studio mode and open inspector drawers immediately upon DOM initialization.
2. **Declarative Theme Switching**:
   - Implement `workspace.SetDeckTheme` in Go to atomically update `deck.json`.
   - Expose `POST /api/deck/theme` in the REST server and live theme selector in Tokens Studio UI.
   - Expose `deckforge theme set` in CLI, `deckforge agent theme-set` / `theme-list` in Agent CLI, and `deckforge_set_theme` in MCP.
3. **Programmatic Deck Creation**:
   - Provide `deckforge new` / `create` CLI aliases.
   - Expose `deckforge agent deck-create`, `POST /api/deck/create`, and `deckforge_create_deck` MCP tool for autonomous multi-slide scaffolding.

## Consequences
- **Positive**: Zero workflow interruption during live edits and recompiles; seamless theme experimentation directly in UI or CLI without manual JSON editing; full autonomous agent lifecycle from deck creation to styling.
- **Negative**: Adds reliance on browser `sessionStorage` for canvas runtime state.
