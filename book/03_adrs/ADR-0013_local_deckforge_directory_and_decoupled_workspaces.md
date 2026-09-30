# ADR-0013: Local .deckforge Workspace Architecture and Decoupled Agent/MCP Subsystems

## Status
Accepted

## Context
As DeckForge expanded to support multi-tier component libraries, modular SCSS/CSS themes, global user settings, and autonomous agent/MCP workflows, two architectural pressures emerged:
1. **Workspace Clutter and Lack of Project Isolation**:
   Non-global presentations, custom components, custom themes, and project-specific configuration overrides were previously scattered in the working directory root (`./themes`, `./components`, `./decks`). This caused repository clutter in consumer repositories and prevented projects from declaring repository-local configuration overrides without mutating global user preferences (`~/.config/deckforge/config.json`).
2. **Subsystem Monoliths**:
   The CLI agent interface (`cmd/agent.go`) had grown to 804 lines and the MCP JSON-RPC server (`internal/server/mcp.go`) had grown to 679 lines. Furthermore, MCP was inappropriately housed inside `internal/server/`, tightly coupling stdio JSON-RPC protocol handling with the live HTTP preview and SSE server.

## Decision
1. **Local `.deckforge/` Directory Architecture**:
   - Establish `.deckforge/` as the canonical project-level directory in workspace roots:
     - Configuration: `<root>/.deckforge/config.json`
     - Components: `<root>/.deckforge/components/`
     - Themes: `<root>/.deckforge/themes/`
     - Decks: `<root>/.deckforge/decks/`
   - Two-tier configuration hierarchy: Global `~/.config/deckforge/config.json` provides base defaults, cleanly overlaid by workspace `<root>/.deckforge/config.json`.
   - Backward-compatible precedence:
     - Themes: Workspace `.deckforge/themes` > Legacy workspace `./themes` > Global `~/.config/deckforge/themes` > Built-in presets.
     - Components: Workspace `.deckforge/components` > Legacy workspace `./components` > Global `~/.config/deckforge/components` > Built-in catalog.
     - Decks: Discovery scans `<root>/.deckforge/decks/` alongside `<root>/decks/` and root folders.
   - CLI enhancements:
     - `deckforge config set <key> <value> [--local]`
     - `deckforge config get <key> [--local]`
     - `deckforge config path [--local]`
     - `deckforge new` automatically targets `.deckforge/decks/<name>` if `.deckforge` exists.
2. **Decoupled MCP Architecture (`internal/mcp/`)**:
   - Extracted all MCP code from `internal/server/` into a dedicated package `internal/mcp/`:
     - `server.go`: JSON-RPC 2.0 stdio protocol loop and dispatcher.
     - `registry.go`: Complete tool metadata and input schema definitions (`GetMCPTools()`).
     - `tools_deck.go`: Presentation scaffolding, inspection, and metadata updates.
     - `tools_slide.go`: Slide listing, extraction, sanitization, updates, and reordering.
     - `tools_theme.go`: Theme listing, assignment, and scaffolding.
     - `tools_component.go`: Component catalog, discovery, scaffolding, and slide insertion.
     - `tools_tokens.go`: Contrast audits, automatic color correction, and styling guide.
     - `tools_settings.go`: Global/local configuration reading and updates.
     - `mcp_test.go`: Isolated unit test suite for all MCP tools.
3. **Decoupled Agent Subsystem (`cmd/agent*.go`)**:
   - Split monolithic `cmd/agent.go` into focused domain modules:
     - `cmd/agent.go`: High-level command router, help, and skill exporter.
     - `cmd/agent_deck.go`: `deck-info`, `deck-meta`, `deck-create`.
     - `cmd/agent_slide.go`: `slide-get`, `slide-set`, `slide-reorder`.
     - `cmd/agent_theme.go`: `theme-list`, `theme-set`, `theme-create`.
     - `cmd/agent_component.go`: `component-list`, `component-create`, `component-insert`, `catalog`.
     - `cmd/agent_tokens.go`: `tokens-audit`, `tokens-fix`, `styling-guide`.
     - `cmd/agent_config.go`: `config-get`, `config-set` (with `--local` support).

## Consequences
- **Positive**:
  - Clean consumer repositories: all DeckForge artifacts can live inside `.deckforge/` without polluting root directories.
  - Granular project overrides: teams can commit `.deckforge/config.json`, `.deckforge/themes/`, and `.deckforge/components/` into version control.
  - Zero breaking changes: legacy `./themes/`, `./components/`, and `./decks/` remain fully supported via fallback resolution.
  - Clean architectural separation: HTTP preview server (`internal/server/`) and stdio JSON-RPC (`internal/mcp/`) are completely decoupled.
  - High maintainability: agent and MCP implementations are broken into small, testable, single-responsibility files.
- **Negative**:
  - Tool authors and users must be aware of the two-tier configuration hierarchy when debugging overrides.
