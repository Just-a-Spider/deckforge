# Master Goals

## Strategic Objectives

### Goal 1: Living Development Book (Codex)
Establish a living documentation system (`book/`) tracking vision, decisions, specifications, milestones, and achievements. Provide CLI commands (`deckforge book status/log/audit`) for continuous maintenance.

### Goal 2: Bi-directional Disk Synchronization
Eliminate the `localStorage` gap. Ensure any edit made in the browser visual canvas is sanitized and saved directly back into individual `slides/NN_*.html` files on disk. Support live hot-reloading across TUI and browser via SSE.

### Goal 3: Interactive DeckForge Studio (Browser 1080p Canvas)
Transform the slide preview into an interactive workstation:
- Click-to-select DOM elements with visual bounding outlines and breadcrumb hierarchy.
- Reorder elements (move up/down/left/right in grid/flex layouts).
- Inline WYSIWYG text formatting toolbar.
- Class inspector drawer to toggle and test utility classes.

### Goal 4: Angular-Inspired Declarative Component System
Define modular slide components using structured schemas (selector, inputs, slots, templates, encapsulated classes). In Phase 1-2, render clean HTML with utility classes; prepare the schema contract for an internal Go-based component compiler.

### Goal 5: Intelligent Tokens & WCAG Contrast Engine
Implement WCAG 2.1 relative luminance and contrast ratio algorithms. Provide automated "Recommended Actions" with 1-click luminosity adjustments to guarantee accessible presentations.

### Goal 6: First-Class AI Agent Layer (MCP & JSON CLI)
Expose headless CLI JSON subcommands and a native Model Context Protocol (MCP) server so autonomous agents (Claude Code, Gemini, Antigravity, OpenClaw) can inspect, build, reorder, and audit slide decks.

### Goal 7: Terminal Companion Ergonomics
Enable keyboard-first slide reordering (`Shift+Up` / `Shift+Down` / `J`/`K`), live contrast diagnostics in theme view, and seamless external `$EDITOR` suspend/resume.

### Goal 8: Studio Session Persistence, Dynamic Theme Switching & Deck Scaffolding
Preserve active Studio canvas and inspector drawer states across browser reloads via `sessionStorage`. Provide declarative theme switching across REST, Agent CLI, MCP, and Studio UI to eliminate manual `deck.json` edits. Provide first-class deck creation tooling (`deckforge new`, `deckforge agent deck-create`, `deckforge_create_deck`, `POST /api/deck/create`).

### Goal 9: Self-Describing Agent Introspection, Standardized Theme Scoping & Live TUI Sync
Eliminate agent context spoon-feeding by providing runtime introspection endpoints (`deckforge_get_catalog`, `deckforge_get_styling_guide`, `deckforge agent catalog/styling-guide`). Standardize theme directory resolution on clean `/themes/<name>/` across workspace and user global configs. Enable passive real-time background external file change synchronization in the TUI companion.

### Goal 10: CWD-First Workspace Resolution and In-TUI Theme Management
Eliminate CLI path ambiguity by automatically scoping the active workspace root to current working directory (`cwd`) when launched without explicit arguments. Enable direct, in-place theme switching on any deck inside the TUI via interactive modal (`t`), auto-recompiling upon selection. Fix external `$EDITOR` integration on `e` key to target `deck.json` directly instead of crashing on directory paths, and cleanly purge legacy `engine/themes/` artifacts.

### Goal 11: Typography Cascade Precedence, Universal Master Index & Vertical Budget Primitives
Resolve browser font fallback failures by injecting theme Google Web Fonts `<link>` tags directly into `<head>` and ordering `<style>` theme CSS first to obey W3C `@import` precedence rules. Synchronize presentation controller slide counter to dynamic `#hudCounter`. Generalize Master Index into a universal English drawer with dynamic title extraction and dedicated CSS styling. Provide programmatic agent deck metadata inspection/mutation (`deck-info`, `deck-meta`, `deckforge_get_deck_info`, `deckforge_set_deck_meta`) to eliminate manual file hacking, enforce an explicit 860px vertical content budget, and provide diverse card archetypes (`.panel-bordered`, `.panel-top-accent`, `.panel-elevated`, `.panel-outline`).

### Goal 12: Full-Spectrum Theme Tokens, Modular Stylesheets & SCSS Pipeline
Decouple immutable application engine boundaries (1920x1080 stage scaling, viewport coordinates, slide lifecycle, 860px budget) from complete presentation design authority. Expand theme system to manage all visual tokens (`tokens.json` with geometry/radii, elevation/shadows, materials/blur, density/spacing, motion). Support granular modular stylesheets (`typography`, `surfaces`, `backdrop`, `components`, `animations`, `styles/*.css`), optional SCSS preprocessor transpilation (`sass` / `dart-sass` / `npx sass`), and fragment Go codebases into modular, cluster-oriented packages.

### Goal 13: Subprocess Caching, Engine Acceleration, and One-Command Installation Automation
Eliminate preprocessor performance bottlenecks by caching compiled SCSS via SHA-256 in memory (`sync.Map`) and persistently on disk (`~/.cache/deckforge/scss/`), achieving a 600x compilation speedup (3.6s down to 6ms). Memoize static core engine assets in the HTML compiler. Upgrade file watcher with a 100ms quiet-period timer debouncer and active theme directory monitoring. Build lean, stripped binaries (`-ldflags="-s -w"`, 31% reduction from 16MB to 11MB) with semantic version introspection. Automate workstation deployment with native Fish shell completions (`deckforge completion fish`), a declarative `Makefile` (`make install-all`), and standalone `install.sh`.

