# ADR-0002: Clean HTML Slide Storage as Single Source of Truth

## Status
Accepted

## Context
Many presentation frameworks store slides in proprietary JSON schemas, monolithic Markdown files with delimiters, or binary databases. When visual editors manipulate these formats, they often inject framework-specific metadata, IDs, and AST wrappers that make manual editing in `$EDITOR` hostile and cause noisy Git diffs.

## Decision
Store each slide as an individual, clean HTML file on disk (`slides/01_title.html`, `slides/02_architecture.html`).
- The file contains pure HTML element snippets (`<section class="slide ...">...</section>`) using semantic DeckForge utility classes.
- When the visual Studio saves changes via `POST /api/slides/update`, the Go backend parses and sanitizes the HTML, stripping temporary editor attributes (`contenteditable`, `df-selected`, `data-studio-id`) before writing to disk.

## Consequences
- **Positive**: 100% human-readable. Clean Git diffs. Developers can edit with Neovim, VSCode, or Studio interchangeably.
- **Positive**: Zero vendor lock-in.
- **Negative**: Requires robust server-side sanitization to prevent studio runtime state from polluting source files.
