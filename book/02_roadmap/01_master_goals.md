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

### Goal 14: Tiered Angular Component Engine & Unified Global Settings
Establish a multi-tier component discovery and scaffolding engine supporting 3 scopes: Built-in presets, Global user components (`~/.config/deckforge/components/`), and Workspace root components (`<workspace>/components/`). Support encapsulated component styles compiled directly into the 1080p canvas stage, with dual bundle formats (single JSON and directory bundle with CSS). Implement unified settings with modern editor auto-detection (Cursor, VS Code, Zed, Windsurf, Sublime, Neovim, Helix, Micro, Vim, Emacs, Nano) and automatic blocking `--wait` flags, with full parity across CLI, Agent CLI, and Native MCP.

### Goal 15: Local .deckforge Workspace Directory Architecture & Agent/MCP Decoupling
Establish `.deckforge/` as the canonical project-level folder for non-global presentations, custom components, themes, and configuration overrides (`.deckforge/config.json`, `.deckforge/components/`, `.deckforge/themes/`, `.deckforge/decks/`). Support a two-tier configuration hierarchy (global defaults overlaid by local workspace overrides). Decouple the monolithic agent and MCP files into modular, domain-driven packages (`internal/mcp/` and `cmd/agent_*.go`), eliminating coupling between the HTTP preview server and stdio JSON-RPC protocol handling while retaining complete backward compatibility.

### Goal 16: Scoped Directory Listing, Self-Contained Agent Context & 1:1 Landscape PDF Engine
Scope presentation listing cleanly to the current project directory without foreign repository bleed or unsolicited theme dumps. Ensure complete user- and system-agnostic execution by removing all hardcoded paths, prioritizing `cwd` over stale persisted workspace roots, and providing a zero-engine-dependency protocol for autonomous agents. Re-engineer the PDF export pipeline to guarantee 1:1 1920x1080 landscape rendering with full theme canvas backgrounds, correct page counts, and zero artifact leakage.

### Goal 17: Industrial CDP PDF Export Fidelity & 1:1 Theme Parity
Eliminate color, theme, and font discrepancies between interactive browser view and exported PDFs by migrating from raw CLI process spawning to a dedicated Chrome DevTools Protocol (CDP) Go library (`github.com/chromedp/chromedp`), complemented by a zero-npm-dependency companion Node script (`scripts/export-pdf.mjs`). Synchronously await web font loading (`document.fonts.ready`) to eliminate system font fallbacks. Preserve full radial ambient glows, cyber grids, mesh textures, and OLED canvas backdrops across all slides via a dedicated per-slide backdrop architecture with exact print color adjustment.

### Goal 18: Universal Workspace Anchoring & Embedded Portable Runners
Eliminate workspace path ambiguity, nested deck theme mismatches, and external script dependencies across all commands. Automatically detect project roots from nested deck paths via upward filesystem traversal, anchoring theme and component managers to the enclosing workspace. Implement recursive ancestor scanning in theme and component managers with explicit fallback warnings to prevent silent style regressions. Embed companion runner scripts directly into the Go binary (`//go:embed`), enabling zero-setup portable export across all platforms and external consumer workspaces.


