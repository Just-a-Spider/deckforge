# The Triad Architecture

## Core System Overview

DeckForge operates as a synchronized Triad:
1. **The Terminal Companion (TUI)**: Bubble Tea terminal application for high-speed slide reordering, directory switching, scaffolding, and daemon management.
2. **The Visual Studio (Browser Canvas)**: Deterministic 1920x1080 stage with click-to-select DOM hierarchy, WYSIWYG inline text editing, CSS class inspector, component palette, and live token tweaker.
3. **The Agentic Protocol Layer (MCP & CLI)**: Model Context Protocol endpoint and headless JSON CLI subcommands enabling autonomous LLMs to query, construct, edit, and audit decks.

All three surfaces converge on a single source of truth: **The Disk Storage**.

```
                        +----------------------------+
                        |     Disk File System       |
                        | (slides/*.html, tokens)    |
                        +--------------+-------------+
                                       ^
                                       |
                +----------------------+----------------------+
                |                      |                      |
                v                      v                      v
        +---------------+      +---------------+      +---------------+
        |  Terminal TUI |      | Web Studio    |      | Agent / MCP   |
        |  (Bubble Tea) |      | (1080p Canvas)|      | (Claude/Gemini|
        +---------------+      +---------------+      +---------------+
```

## Data Flow and Persistence Contract

1. **Slide Edits**:
   - In Browser Studio: User edits text or toggles CSS classes. After a 500ms debounce, Studio issues `POST /api/slides/update` to the Go backend.
   - The Go backend strips temporary runtime DOM attributes (such as `contenteditable="true"` or `df-selected="true"`) and overwrites `slides/NN_*.html`.
   - The fsnotify watcher or SSE broadcaster notifies any connected viewers of the update.
2. **Slide Reordering**:
   - In TUI or Studio: User moves slide 3 to position 1.
   - An atomic filesystem rename transaction (`03_foo.html` -> `.tmp_...` -> `01_foo.html`) runs.
   - Both TUI and Web Studio immediately refresh their sequence indices.
3. **Design Tokens**:
   - Live color updates in Studio update `:root` variables in real time.
   - Clicking "Save" or "Apply Recommended Action" sends `POST /api/tokens/update`.
   - The Go engine updates `themes/<name>/tokens.json` and triggers re-compilation.
