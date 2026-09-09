# ADR-0007: Self-Describing Agent Introspection, Standardized Themes and TUI Live Sync

## Status
Accepted

## Context
When running DeckForge in production environments (e.g., compiled standalone binary in `/usr/local/bin`), repository source files (`engine/core/style.css`, `catalog.json`) do not exist on disk because they are embedded via `go:embed`. External autonomous agents (Claude Code, Gemini CLI, Antigravity) connecting via MCP or CLI had no way to discover available CSS classes, layout containers, or component schemas without manual user spoon-feeding. Additionally, theme locations suffered from ad-hoc path fragmentation, and the TUI companion was passive—failing to update when files were modified by external agents or editors.

## Decision
1. **Runtime Self-Describing Introspection**:
   - Provide `internal/models/styling_guide.go` returning canvas rules, semantic CSS class dictionary, spacing scale, and token variables directly from memory.
   - Expose `deckforge agent catalog` and `deckforge agent styling-guide` subcommands in CLI.
   - Expose `deckforge_get_catalog` and `deckforge_get_styling_guide` tools in native MCP server.
   - Expose `GET /api/catalog` and `GET /api/styling-guide` in REST server.
2. **Standardized Theme Directory Layout & On-Demand Seeding**:
   - Standardize strictly on `/themes/<name>/` across workspace (`./themes/`) and global user config (`~/.config/deckforge/themes/`), falling back to built-in embedded presets.
   - Maintain zero initial disk clutter by default; provide `deckforge theme seed <name> [--global]` to populate templates on demand.
3. **Passive TUI Live Synchronization**:
   - Implement a lightweight 1.5s background ticker (`FileWatchTickMsg`) in Bubble Tea companion.
   - Monitor modification timestamps (`mtime`) on `deck.json` and `slides/`.
   - Automatically re-inspect deck and refresh slide view when external changes occur without requiring manual `r` keypresses.
4. **Spacing & Chapter Primitives**:
   - Standardize atomic spacing scale (`--space-1` to `--space-6`) and utility classes (`.p-1` to `.p-6`, `.gap-1` to `.gap-6`).
   - Standardize `.slide.chapter-break` layout with `.chapter-kicker`, `.chapter-num`, `.chapter-title`, `.chapter-divider`, and `.chapter-summary`.

## Consequences
- **Positive**: Autonomous agents can self-discover all styling rules and components with zero repo access and zero spoon-feeding; themes follow clean directory semantics without proliferation; TUI reflects agent and editor changes in real time.
- **Negative**: Adds continuous 1.5s periodic stat calls in TUI runtime.
