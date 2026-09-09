# Vision and Principles

## Manifesto: Deterministic Presentation Engineering

Presentations in traditional tools (PowerPoint, Keynote, Canva, Google Slides) suffer from:
1. **Opaque binary or proprietary cloud formats**: cannot be version-controlled, diffed, or reviewed via pull requests.
2. **Fragile layout reflows**: flexbox wrapping and dynamic DOM cascades cause text and cards to break across different projector aspect ratios.
3. **Disconnected tooling**: editing text in a web UI cannot be automated or driven by terminal workflows or external editors ($EDITOR).
4. **AI-illiterate architecture**: LLM agents struggle to inspect, reason about, or deterministically manipulate arbitrary UI canvas layouts.

DeckForge exists to solve this with four non-negotiable principles:

### Principle 1: Deterministic 1080p Canvas
Every slide stage is strictly pinned to a 1920x1080 canvas scaled uniformly via `Math.min(W/1920, H/1080)`. What the developer sees on an ultra-wide monitor is mathematically identical to what is projected on a conference room screen. Zero flex reflow surprises.

### Principle 2: Pure Discrete HTML as Single Source of Truth
Every slide lives as a discrete file (`slides/01_title.html`, `slides/02_architecture.html`). No opaque binary format. No monolithic database. A slide is pure semantic HTML with utility classes. Developers can use Neovim, VSCode, or the DeckForge Studio interchangeably without format conversion artifacts.

### Principle 3: Zero-Dependency Portability
DeckForge builds to a standalone 100% pure Go binary. Slides compile into self-contained zero-dependency HTML files. An exported presentation runs anywhere an HTML5 browser exists, completely offline, forever.

### Principle 4: First-Class Human & Agent Ergonomics
DeckForge is designed simultaneously for two types of operators:
- **Humans**: who desire fast terminal keyboard navigation (TUI) and an interactive 1080p visual canvas (Studio).
- **AI Agents**: (Claude Code, Gemini, Antigravity, OpenClaw) who need clean JSON CLI subcommands, an MCP server, and explicit semantic token/component dictionaries to create and edit decks autonomously.
