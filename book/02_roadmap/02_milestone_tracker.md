# Milestone Tracker

## Status Overview
- Current Active Phase: **Complete (All 7 Phases Verified & Deployed)**
- Overall Progress: **100%**

---

## Phase 0: Development Book Foundations
- [x] Create `book/` directory tree with all 6 sections
- [x] Document Vision, Principles, and Triad Architecture
- [x] Record ADR-0001 through ADR-0005
- [x] Document technical specifications (Components, Tokens, MCP)
- [x] Implement `cmd/book.go` (`deckforge book status/log/audit`)
- [x] Verify `deckforge book audit` passes

## Phase 1: Go Server API & Bi-directional Persistence
- [x] Implement `internal/workspace/reorder.go` (atomic slide renaming)
- [x] Implement `internal/server/api.go` (slide CRUD, reorder, token endpoints, SSE)
- [x] Sanitize runtime attributes on disk write (strip contenteditable, df-*)
- [x] Unit tests for reordering and HTML sanitization

## Phase 2: DeckForge Studio Canvas Engine
- [x] Implement `engine/core/studio.css` (selection outlines, floating toolbars)
- [x] Implement `engine/core/studio.js` (selector, element moving, inline editor)
- [x] Update `engine/core/controller.js` to hook into Studio Mode (`S` key)
- [x] Connect bi-directional auto-save over `/api/slides/update`

## Phase 3: Intelligent Tokens & WCAG Contrast Engine
- [x] Implement `internal/models/tokens_math.go` (relative luminance, WCAG 2.1 ratio)
- [x] Implement automated Recommended Action generator (binary search lightness delta)
- [x] Implement `engine/core/tokens_studio.js` (live palette editor with 1-click fix)
- [x] Unit tests for contrast math and auto-fix suggestions

## Phase 4: Angular-Inspired Component System
- [x] Implement `internal/models/component.go` (declarative schema, inputs, slots)
- [x] Implement `engine/core/catalog.json` (7 core archetypes)
- [x] Implement Studio component inserter with slot projection

## Phase 5: First-Class LLM & Agent Layer
- [x] Implement `cmd/agent.go` (`deckforge agent slide-get/set/reorder/tokens-audit`)
- [x] Implement `internal/server/mcp.go` (native Model Context Protocol tools)
- [x] Exportable `assets/agent/SKILL.md`

## Phase 6: TUI Companion Upgrades & Polish
- [x] Implement slide reordering in TUI (`Shift+Up` / `Shift+Down`)
- [x] Real-time WCAG contrast diagnostics in Theme Studio view
- [x] Full end-to-end verification and release logging

## Phase 7: Session Persistence, Dynamic Theme Switching & Deck Scaffolding
- [x] Add `SetDeckTheme` to workspace layer and update `deck.json` cleanly
- [x] Add `POST /api/deck/theme`, `GET /api/themes`, and `POST /api/deck/create` to server REST API
- [x] Implement `sessionStorage` state persistence in `studio.js` (`df_studio_active`, `df_inspector_open`)
- [x] Add Theme Selector dropdown in `tokens_studio.js`
- [x] Add CLI commands `deckforge new` / `create` and `deckforge theme set`
- [x] Add Agent CLI commands `deck-create`, `theme-set`, `theme-list`
- [x] Add MCP tools `deckforge_create_deck`, `deckforge_set_theme`, `deckforge_list_themes`
- [x] Record ADR-0006 and update API specifications

## Phase 8: Self-Describing Agent Introspection, Standardized Themes & Live TUI Sync
- [x] Implement `internal/models/styling_guide.go` with geometry rules, classes, and tokens
- [x] Implement CLI subcommands `deckforge agent catalog` and `deckforge agent styling-guide`
- [x] Implement MCP tools `deckforge_get_catalog` and `deckforge_get_styling_guide`
- [x] Implement REST endpoints `GET /api/catalog` and `GET /api/styling-guide`
- [x] Standardize theme resolution on `/themes/<name>/` across workspace and global config
- [x] Implement `deckforge theme seed <name> [--global]` command for zero-clutter on-demand seeding
- [x] Implement passive `mtime` live synchronization ticker in TUI companion
- [x] Add atomic spacing scale (`.p-1` to `.p-6`, `.gap-1` to `.gap-6`) and `.slide.chapter-break`
- [x] Record ADR-0007 and update Living Development Book

## Phase 9: CWD-First Workspace Resolution and In-TUI Theme Operations
- [x] Update `cmd/root.go` and `internal/tui/app.go` to default active root to `os.Getwd()` absolute path
- [x] Fix TUI `e` key to target `deck.json` directly instead of crashing on directory paths
- [x] Implement interactive Theme Picker Modal overlay on `t` key in `internal/tui/views_decks.go`
- [x] Auto-recompile deck and update in-memory state upon selecting a theme in TUI
- [x] Synchronize selected deck to `ThemeStudioView.ActiveDeck` across tab switches
- [x] Replace all `engine/themes` paths with `themes/` in `internal/tui/views_theme.go`
- [x] Purge legacy `engine/themes/` directory from repository
- [x] Record ADR-0008 and verify Living Development Book audit health

## Phase 10: Typography Cascade, Universal Master Index, Vertical Budget & Agent Metadata
- [x] Inject `<link rel="stylesheet">` Google Fonts directly in `<head>` and place theme CSS first in `<style>`
- [x] Unify presentation controller HUD slide counter lookup (`#deckPageNum` / `#hudCounter`)
- [x] Standardize Master Index on universal English, dynamic title extraction, and rich CSS styling
- [x] Add vertical layout budget primitives (`.slide-content-area` with 860px max-height, `.row-fill`, `.row-auto`)
- [x] Add diverse card archetypes (`.panel-bordered`, `.panel-top-accent`, `.panel-elevated`, `.panel-outline`)
- [x] Implement Agent CLI subcommands `deckforge agent deck-info` and `deckforge agent deck-meta`
- [x] Implement Native MCP tools `deckforge_get_deck_info` and `deckforge_set_deck_meta`
- [x] Fix TUI Theme Selector modal ANSI padding glitches and add 7-item windowed scrolling
- [x] Record ADR-0009 and verify Living Development Book audit health

## Phase 11: Full-Spectrum Theme Tokens, Modular Stylesheets & SCSS Pipeline
- [x] Create `internal/models/theme_tokens.go` with geometry, elevation, material, spacing, and motion token schemas
- [x] Refactor `engine/core/components.css` and `stage.css` to use token variables with fallbacks
- [x] Implement `internal/theme/scss.go` with automatic SCSS compilation and modern CSS fallback
- [x] Implement `internal/theme/loader.go` with open directory scanning for `components.css`, `animations.css`, and `styles/*.css`
- [x] Fragment 14 presets into modular aesthetic clusters (`presets_cyber.go`, `presets_editorial.go`, `presets_modern.go`, `presets_pastel.go`)
- [x] Enrich all 14 presets with full-spectrum design tokens
- [x] Update `internal/models/styling_guide.go` to expose complete design tokens to LLM agents
- [x] Record ADR-0010 and update technical specifications
- [x] Verify complete test suite and Living Development Book health

## Phase 12: Subprocess Caching, Engine Acceleration & Installation Automation
- [x] Implement SHA-256 in-memory (`sync.Map`) and persistent on-disk (`~/.cache/deckforge/scss/`) cache in `internal/theme/scss.go`
- [x] Memoize Sass binary detection (`detectSassCompiler`) via `sync.Once`
- [x] Memoize core embedded static CSS and JS assets in `internal/compiler/compiler.go`
- [x] Upgrade file watcher in `internal/server/server.go` with 100ms timer-reset debouncer and active theme directory monitoring
- [x] Build stripped 11MB binaries with `-ldflags="-s -w -X 'deckforge/cmd.Version=0.3.0'"`
- [x] Add `deckforge version` command with runtime and platform details
- [x] Implement native Fish shell autocompletions (`cmd/completion.go`, `completions/deckforge.fish`, `deckforge completion fish --install`)
- [x] Create root `Makefile` with `make install-all` (build, install, completions, skill)
- [x] Create standalone `install.sh` automated installation script
- [x] Record ADR-0011 and verify Living Development Book audit health

## Phase 13: Tiered Angular Component Engine & Unified Global Settings
- [x] Implement `internal/editor/detector.go` with auto-detection for Cursor, VS Code, Zed, Windsurf, Sublime, Neovim, Helix, Micro, Vim, Emacs, Nano and blocking `--wait` flags
- [x] Update `internal/models/workspace.go` with `DECKFORGE_CONFIG_DIR`, `XDG_CONFIG_HOME`, key-value `Get`/`Set`, and config isolation for unit tests
- [x] Integrate installed editor indicators and dynamic cycling in TUI Settings view (`internal/tui/views_settings.go`)
- [x] Implement `internal/components/manager.go` with 3-tier component resolution (Builtin, Global, Workspace) and precedence overrides
- [x] Support dual custom component formats: single JSON files and modular directory bundles with encapsulated `component.css`
- [x] Implement `internal/components/scaffold.go` for declarative component scaffolding
- [x] Inject custom component CSS into HTML compilation stage in `internal/compiler/compiler.go`
- [x] Implement CLI commands `deckforge config [list|get|set|editors|path]` and `deckforge component [list|create|inspect]`
- [x] Implement Agent CLI commands `component-list`, `component-create`, `config-get`, `config-set`, and dynamic catalog/component-insert
- [x] Implement Native MCP tools `deckforge_list_components`, `deckforge_create_component`, `deckforge_get_settings`, `deckforge_update_settings`, and dynamic catalog/insert
- [x] Respect configured default theme and slide count across `deckforge init` and `deckforge agent deck-create`

## Phase 14: Local .deckforge Workspace Directory Architecture & Agent/MCP Decoupling
- [x] Establish `.deckforge/` project-level directory standard (`.deckforge/config.json`, `.deckforge/components/`, `.deckforge/themes/`, `.deckforge/decks/`)
- [x] Implement two-tier configuration hierarchy (global baseline with local `.deckforge/config.json` overlay and `SetLocal`/`SaveLocal`)
- [x] Update Theme Manager to scan and scaffold into `.deckforge/themes/` with fallback to `./themes/`
- [x] Update Component Manager to scan and scaffold into `.deckforge/components/` with fallback to `./components/`
- [x] Update Deck Discovery to scan `.deckforge/decks/` alongside `./decks/` and root folders
- [x] Add `--local` flag to `deckforge config [set|get|path]` and `deckforge agent config-set`
- [x] Decouple MCP JSON-RPC server from `internal/server/` into standalone `internal/mcp/` package with domain-specific tool modules
- [x] Decouple `cmd/agent.go` into domain files (`agent_deck.go`, `agent_slide.go`, `agent_theme.go`, `agent_component.go`, `agent_tokens.go`, `agent_config.go`)
- [x] Record ADR-0013 and verify Living Development Book audit health

## Phase 15: Scoped Directory Listing, Self-Contained Agent Context & 1:1 Landscape PDF Engine
- [x] Scope `deckforge list` exclusively to presentations in current directory / workspace (`.deckforge/decks/`, `decks/`, or root if deck)
- [x] Remove unsolicited global theme dump from `deckforge list` (available via `deckforge theme list` or `deckforge list --themes`)
- [x] Restrict `ScanDecksInRoot` to immediate child directories (depth 1 only), preventing cross-repo traversal
- [x] Guarantee `targetRoot` in `models.LoadWorkspaceConfigForRoot` defaults to `cwd` to prevent foreign workspace hijacking
- [x] Add Zero Engine Dependency directive in `SKILL.md` and `runAgentExportSkill` to isolate agent context to consumer repos
- [x] Purge `deckforge book log` from end-user agent skills
- [x] Remove hardcoded user paths in `internal/exporter/exporter.go` and implement dynamic cross-platform discovery
- [x] Implement modern headless Chrome flags (`--headless=new`, `--run-all-compositor-stages-before-draw`, `--virtual-time-budget=3000`)
- [x] Add `@page { size: 1920px 1080px; margin: 0; }` and `print-color-adjust: exact` in `stage.css`
- [x] Preserve theme canvas background and suppress non-slide DOM elements in print CSS
- [x] Make `book/SUMMARY.md` and `cmd/book.go audit` 100% portable with relative links
- [x] Record ADR-0014 and verify Living Development Book audit health

## Phase 16: Chromedp CDP Export Engine, Companion Node Script, and 1:1 Theme Color Parity
- [x] Integrate `github.com/chromedp/chromedp` library into `internal/exporter/` for DevTools Protocol export
- [x] Synchronously await `document.fonts.ready` before PDF generation, eliminating fallback fonts
- [x] Create standalone zero-dependency Node CDP runner in `scripts/export-pdf.mjs`
- [x] Add `--engine` (`chromedp`, `node`, `cli`) and `--timeout` flags to `deckforge export`
- [x] Inject `.slide-backdrop` inside each slide in `internal/compiler/compiler.go`
- [x] Map `.stage-backdrop` to `.stage-backdrop, .slide-backdrop` in `internal/models/theme.go`
- [x] Enable `.slide-backdrop` in `engine/core/stage.css` `@media print` with proper z-indexing
- [x] Verify visual color parity, radial gradients, cyber grids, and typography across light and dark decks
- [x] Record ADR-0015 and verify DevBook audit health

## Phase 17: Universal Workspace Anchoring, Upward Theme Traversal & Embedded Portable Runners
- [x] Implement `DetectWorkspaceRoot(targetPath)` in `internal/workspace/root.go` with upward recursive tree discovery
- [x] Add unit tests in `internal/workspace/root_test.go` covering nested decks, slide files, and direct `.deckforge` references
- [x] Implement recursive upward ancestor tree scanning and `cwd` scanning in `ThemeManager.ListThemes(deckPath)`
- [x] Implement diagnostic warnings on fallback in `ThemeManager.ResolveTheme` to eliminate silent style regressions
- [x] Implement recursive upward ancestor tree scanning and `cwd` scanning in `ComponentManager.ListComponents(deckPath)`
- [x] Anchor `runExport`, `runBuild`, `runServe`, `runList`, `runTheme`, and `runComponent` to `workspace.DetectWorkspaceRoot`
- [x] Embed `scripts/export-pdf.mjs` into `internal/exporter` binary package via `//go:embed`
- [x] Fallback to cache directory extraction in `exportWithNode` so `deckforge export --engine node` works standalone anywhere
- [x] Verify global installation (`make install`) in external workspace (`~/Desktop/Projects/SysMon`) with 100% theme color accuracy
- [x] Record ADR-0016 and verify DevBook audit health

