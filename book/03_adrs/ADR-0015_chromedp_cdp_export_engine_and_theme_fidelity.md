# ADR-0015: Chromedp CDP Export Engine, Companion Node Script, and 1:1 Theme Color Parity

## Status
Accepted

## Context
Previous iterations of DeckForge exported presentations by executing raw CLI commands (`google-chrome --headless=new --print-to-pdf`). This created three systemic flaws:
1. **Theme Backdrop Stripping**: In `stage.css`, `@media print` hidden `.stage-backdrop`, erasing radial ambient lighting, mesh gradients, cyber grids, and accent glow across all slides. Slides collapsed to flat solid backgrounds.
2. **Font Race Conditions**: The Chromium CLI print pipeline did not wait for web fonts (`document.fonts.ready`), triggering system fallback fonts (e.g. DejaVu / Liberation / Helvetica) instead of Poppins, Plus Jakarta Sans, and JetBrains Mono.
3. **No DevTools Protocol Control**: Spawning raw CLI commands lacked event-driven page lifecycle awareness, could not synchronize layout calculations, and was vulnerable to desktop browser wrapper interference.

Developers requested a dedicated library or a companion mini Node script for reliable, color-accurate PDF generation.

## Decision
1. **Primary Built-in Engine (Go `chromedp` Library)**:
   - Integrated `github.com/chromedp/chromedp` directly into `internal/exporter/`.
   - Communicates with Chromium over WebSocket Chrome DevTools Protocol (CDP).
   - Synchronously awaits `document.fonts.ready` promise and layout stabilization.
   - Issues `Page.printToPDF` with `printBackground: true`, `preferCSSPageSize: true`, and 0 margins.
   - Retains zero runtime external dependencies: single static binary.
2. **Extensible Companion Node Script (`scripts/export-pdf.mjs`)**:
   - Implemented a zero-npm-dependency Node CDP runner utilizing native Node 22+ `fetch` and `WebSocket`.
   - Accessible via CLI flag `deckforge export --engine=node` or standalone `node scripts/export-pdf.mjs <deck-path>`.
3. **Per-Slide Backdrop Architecture**:
   - Updated `internal/compiler/compiler.go` to inject `<div class="slide-backdrop"></div>` inside every slide `<section>`.
   - Updated `CompileFullCSS` in `internal/models/theme.go` so all theme `BackdropCSS` rules target `.stage-backdrop, .slide-backdrop`.
   - Updated `engine/core/stage.css` so each printed slide renders the complete backdrop gradient, grid, and canvas background.
4. **CLI Flags**:
   - Added `--engine` (`chromedp` [default], `node`, `cli`) and `--timeout` to `deckforge export`.

## Consequences
- **Positive**:
  - Full theme color parity: dark themes retain emerald/purple glow and OLED black, light themes retain crimson glow and fine graph grids.
  - Pixel-perfect typography: Google Fonts are guaranteed to finish downloading before PDF rasterization.
  - Fast execution: average 1.5–2s export time per 12-slide presentation.
  - Complete flexibility: single Go binary out of the box, with Node script option for JavaScript-centric environments.
- **Negative**:
  - Added `github.com/chromedp/chromedp` to `go.mod`.
