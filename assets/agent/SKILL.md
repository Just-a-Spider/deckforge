---
name: deckforge
description: Autonomous and interactive management of DeckForge modular 1080p slide presentations
---

# DeckForge Skill

Use this skill whenever working in repositories containing DeckForge presentations (`deck.json`, `slides/*.html`, or running `deckforge`).

## 1. Non-Negotiable Geometry Constraints
- All slides render on a fixed **1920x1080 canvas** (`Math.min(W/1920, H/1080)`).
- Never use dynamic viewport height/width (`vh`, `vw`).
- Never let layouts wrap uncontrollably. Always use `.grid-2col`, `.grid-3col`, or `.grid-split-hero`.

## 2. Discrete Slide Structure
- Slides are discrete HTML files in `slides/NN_slug.html`.
- Outer wrapper is `<section class="slide [active]" data-slide="N">`.
- Inside wrapper: `<div class="slide-frame">`.
- Header: `<div class="slide-topbar"><div class="topbar-title">...</div></div>`.
- Content: `<div class="slide-content-area">`.

## 3. DeckForge Semantic Classes
- Layouts: `.grid-2col`, `.grid-3col`, `.grid-4col`, `.grid-split-hero`, `.grid-metrics`.
- Panels & Surfaces: `.glass-panel`, `.lime-edge`, `.accent-edge`.
- Badges: `.card-tag`, `.topbar-tag`, `.lime-badge`, `.purple-badge`.
- Typography: `.panel-headline`, `.micro-kicker`, `.bullet-list`, `.lime-hl`, `.accent-hl`.

## 4. Agent CLI Subcommands
Always prefer the structured JSON CLI commands over manual raw text edits when possible:
- `deckforge agent catalog` (retrieve complete component schemas & templates)
- `deckforge agent styling-guide` (retrieve all CSS classes, layout primitives, and tokens)
- `deckforge agent deck-create <name> [--theme <preset>] [--slides <count>] [--path <dir>]`
- `deckforge agent theme-list [deck-path]`
- `deckforge agent theme-set <deck-path> <theme-name>`
- `deckforge agent theme-create <name> [--base <preset>] [--global] [--scss]`
- `deckforge agent slide-get <deck-path> <index>`
- `deckforge agent slide-set <deck-path> <index> --content <html_or_file>`
- `deckforge agent slide-reorder <deck-path> --order <1,3,2,...>`
- `deckforge agent tokens-audit <deck-path>`
- `deckforge agent tokens-fix <deck-path> --apply`
- `deckforge agent component-insert <deck-path> <index> <selector> [--props <json>]`
- `deckforge book log "<achievement>"`

## 5. Model Context Protocol (MCP) Tools
When connected via stdio (`deckforge mcp <deck-path>`):
- `deckforge_get_catalog`: Introspect all component schemas, inputs, and slots.
- `deckforge_get_styling_guide`: Introspect all styling rules, CSS classes, and tokens.
- `deckforge_create_deck`: Scaffold new presentation.
- `deckforge_list_themes`: Enumerate available themes with palette metadata.
- `deckforge_set_theme`: Switch theme in `deck.json` and recompile.
- `deckforge_create_theme`: Scaffold new theme skeleton or preset clone with CSS or SCSS.
- `deckforge_list_slides`: Enumerate all slides.
- `deckforge_get_slide`: Read slide HTML by index.
- `deckforge_update_slide`: Write clean slide HTML and recompile.
- `deckforge_reorder_slides`: Reorder slides on disk.
- `deckforge_audit_tokens`: Audit WCAG 2.1 contrast compliance.
- `deckforge_apply_token_fix`: Auto-fix theme contrast tokens.
- `deckforge_insert_component`: Insert declarative component.


