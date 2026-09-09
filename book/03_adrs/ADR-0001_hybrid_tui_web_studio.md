# ADR-0001: Hybrid TUI and Web Studio Architecture

## Status
Accepted

## Context
DeckForge began as a pure Bubble Tea TUI in Go. While the terminal excels at file discovery, quick scaffolding, server orchestration, and invoking external editors (`$EDITOR`), terminal cells (fixed grid, monochrome or 256/24-bit character cells) cannot support:
1. Pixel-precise 1920x1080 canvas preview.
2. Direct drag-and-drop or visual element selection.
3. Accurate web font weight rendering, kerning, and CSS transitions.

Conversely, a pure web application loses the high-efficiency terminal-driven developer workflow, headless automation, and pure single-binary deployment.

## Decision
Adopt a **Hybrid Architecture**:
- Keep the Bubble Tea TUI for rapid developer workflow: deck browsing, root switching, slide reordering, server daemon control, and compiler execution.
- Upgrade the embedded web preview into **DeckForge Studio**: an interactive 1080p canvas with element selection, WYSIWYG text editing, CSS class inspector, and token tweaks.
- Connect the two surfaces via the local Go HTTP/REST/SSE server so all changes synchronize bi-directionally on disk.

## Consequences
- **Positive**: Best of both worlds. Fast terminal operations and precise visual manipulation without compromise.
- **Negative**: Requires maintaining two frontends (Bubble Tea in Go and Canvas Studio in vanilla JS/CSS).
- **Mitigation**: Vanilla JS with zero npm/bundler dependencies embedded directly into the Go binary (`embed.FS`).
