# ADR-0012: Tiered Angular Component Engine, Encapsulated Styles, and Unified Global Settings

## Status
Accepted

## Context
In ADR-0003, DeckForge committed to an Angular-inspired declarative component specification (`selector`, `inputs`, `slots`, `classes`, `templateSnippet`). However, components remained hardcoded to 7 static built-in archetypes with zero filesystem discovery. Users and LLM agents could not scaffold, customize, or distribute reusable slide components locally or globally. Furthermore:
1. Component definitions lacked style encapsulation: in Angular, `@Component({ styles: [...] })` or `styleUrls` isolates component visual design. In DeckForge, custom component styles had to be manually spliced into theme files.
2. Configuration persistence was constrained: `~/.config/deckforge/config.json` did not respect `XDG_CONFIG_HOME` or `DECKFORGE_CONFIG_DIR`, causing automated test runs to pollute real user preferences.
3. External editor options in the TUI were hardcoded to 5 candidates (`auto`, `code`, `nvim`, `vim`, `nano`), omitting modern editors (Cursor, Zed, Windsurf, Sublime Text, Helix, Micro), and lacking auto-detection and blocking `--wait` flags.
4. CLI, Agent CLI, and MCP servers lacked configuration and component inspection/scaffolding primitives.

## Decision
1. **Multi-Tier Component Architecture (`internal/components/manager.go`)**:
   - Establish three-tier component resolution:
     - Tier 1: Embedded Built-in presets (`df-metric-card`, `df-architecture-node`, etc.)
     - Tier 2: Global user components (`~/.config/deckforge/components/` or `$DECKFORGE_CONFIG_DIR/components/`)
     - Tier 3: Workspace and deck root components (`<workspace>/components/` and `<deck>/components/`)
   - Precedence: Workspace Root > Global > Built-in (keyed by `selector`).
2. **Encapsulated Component Styling & Dual Bundle Formats**:
   - Add `Styles` field to `ComponentDefinition` models.
   - Support single-file JSON definitions (`<selector>.json`) and directory bundles (`components/<selector>/` containing `component.json` + companion `component.css` or `styles.css`).
   - Automatically compile and inject aggregated custom component CSS into the HTML stage `<style>` block directly following `components.css`.
3. **Component Scaffolding Engine (`internal/components/scaffold.go`)**:
   - Provide automated scaffolding for custom component definitions with normalized selectors (`df-*`), typed inputs, default slot projections, and scoped CSS classes.
   - Expose via CLI (`deckforge component create`), Agent CLI (`deckforge agent component-create`), and MCP (`deckforge_create_component`).
4. **Unified Global Settings & Modern Editor Detection (`internal/editor/detector.go`, `internal/models/workspace.go`)**:
   - Check `DECKFORGE_CONFIG_DIR` first, `XDG_CONFIG_HOME/deckforge` second, falling back to `~/.config/deckforge`.
   - Implement `DetectInstalledEditors()` probing PATH for Cursor, VS Code, Zed, Windsurf, Sublime, Neovim, Helix, Micro, Vim, Emacs, and Nano.
   - Inject required blocking flags (`--wait`, `-w`) for GUI editors to ensure smooth raw-terminal suspend and resume.
   - Dynamically highlight `[installed]` editors in the TUI Settings view.
   - Expose configuration management across CLI (`deckforge config get/set/list/editors/path`), Agent CLI (`deckforge agent config-get/set`), and MCP (`deckforge_get_settings`, `deckforge_update_settings`).
   - Ensure deck scaffolding (`deckforge new`, `deckforge agent deck-create`, `deckforge_create_deck`) respects configured default theme and slide count.

## Consequences
- **Positive**:
  - Full personalization: users and LLM agents can scaffold, edit, and share custom Angular-style slide components globally or per-workspace.
  - Custom component styles work out of the box with zero manual stylesheet editing.
  - Zero test pollution: test suites isolate configuration in temporary directories via `DECKFORGE_CONFIG_DIR`.
  - Seamless editor integration: Cursor, Zed, Windsurf, and VS Code work reliably with terminal suspend/resume.
  - Complete parity across TUI, CLI, Agent CLI, and MCP servers.
- **Negative**:
  - Global custom components must be uniquely named to prevent unintended selector collisions.
