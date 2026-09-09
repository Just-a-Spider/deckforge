# AI Agent Playbook

## Instructions for Autonomous & Interactive Coding Agents
*(Claude Code, Gemini CLI, Antigravity, OpenClaw, OpenCode)*

When working with a DeckForge project, you MUST adhere to the following rules:

### 1. Canvas Geometry Rules
- All slides render on a fixed **1920x1080 canvas**.
- NEVER use dynamic viewport units (`100vw`, `100vh`) inside slide HTML.
- ALWAYS use fixed pixel values (`34px`, `70px`, `18px`) or percentages relative to the 1920x1080 stage.
- Do NOT allow cards or text blocks to wrap unpredictably. Use `.grid-2col`, `.grid-3col`, or `.grid-split-hero`.

### 2. Semantic CSS Class Dictionary
- Containers: `.slide`, `.slide-frame`, `.slide-content-area`, `.slide.chapter-break`.
- Layouts: `.grid-2col`, `.grid-3col`, `.grid-4col`, `.grid-split-hero`, `.grid-metrics`.
- Panels & Cards: `.glass-panel`, `.lime-edge`, `.accent-edge`.
- Typography: `.slide-topbar`, `.topbar-title`, `.topbar-tag`, `.micro-kicker`, `.panel-headline`, `.card-tag`, `.bullet-list`, `.chapter-kicker`, `.chapter-num`, `.chapter-title`, `.chapter-summary`.
- Highlights: `.lime-hl`, `.accent-hl`, `.lime-badge`, `.purple-badge`, `.badge-danger`.
- Spacing: `.p-1` to `.p-6`, `.px-1` to `.px-6`, `.py-1` to `.py-6`, `.gap-1` to `.gap-6`, `.m-auto`, `.mt-1` to `.mt-6`, `.mb-1` to `.mb-6`.

### 3. Agent Tooling Workflow (Zero-Guesswork Introspection)
- **DO NOT grep internal engine files**. All rules, components, and classes are self-describing at runtime:
  - Run `deckforge agent styling-guide` or use MCP `deckforge_get_styling_guide` for complete CSS class and token context.
  - Run `deckforge agent catalog` or use MCP `deckforge_get_catalog` for reusable Angular-style component schemas.
- To scaffold a new presentation: run `deckforge agent deck-create <name> [--theme <theme>] [--slides <count>]` or use MCP tool `deckforge_create_deck`.
- To inspect available themes: run `deckforge agent theme-list <path>` or use MCP tool `deckforge_list_themes`.
- To switch active presentation theme: run `deckforge agent theme-set <path> <theme>` or use MCP tool `deckforge_set_theme`.
- To scaffold a new custom theme (workspace or global): run `deckforge agent theme-create <name> [--base <preset>] [--global] [--scss]` or use MCP tool `deckforge_create_theme`.
- To inspect deck state: run `deckforge agent slide-get <path> <index>` or use MCP tool `deckforge_get_slide`.
- To update slide text: run `deckforge agent slide-set <path> <index> --content <html_or_file>` or use MCP tool `deckforge_update_slide`.
- To audit accessibility: run `deckforge agent tokens-audit <path>`.
- To log milestone achievements: run `deckforge book log "<description>"`.

### 4. Studio Mode & TUI Live Continuity
- Studio mode and Inspector state persist across slide navigation and live reloads via `sessionStorage` (`df_studio_active`, `df_inspector_open`, `df_last_slide`).
- TUI companion passively detects external file modifications to `deck.json` and `slides/` every 1.5 seconds, auto-refreshing presentation state without user keypresses.
- Agents making automated edits via API, MCP, or file writes trigger background live updates across both browser and terminal interfaces.


