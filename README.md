# DeckForge

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](go.mod)
[![CI](https://github.com/Just-a-Spider/deckforge/actions/workflows/ci.yml/badge.svg)](https://github.com/Just-a-Spider/deckforge/actions/workflows/ci.yml)

High-performance, modular 1080p slide presentation platform written in 100% pure Go. Designed for terminal-driven workflows, zero-dependency HTML distribution, segmented theme management, and seamless external editor integration (`$EDITOR`).

---

## Installation

### One-Line Automated Install

```bash
git clone https://github.com/Just-a-Spider/deckforge.git
cd deckforge
./install.sh
```

### Build from Source

```bash
git clone https://github.com/Just-a-Spider/deckforge.git
cd deckforge
make install
```

This compiles a stripped, static binary and installs it to `~/.local/bin/deckforge` (and `$GOPATH/bin/deckforge` if configured). Optional shell autocompletions (`fish`) and AI agent skills can be installed via `make install-all`.

---

## Core Capabilities

1. **Deterministic 1920x1080 Canvas**: Fixed master stage geometry scaled uniformly using `Math.min(W/1920, H/1080)`. No flexbox wrapping or responsive layout breaking across different screens.
2. **100% Pure Go Binary**: Single portable static executable. Core assets (`stage.css`, `hud.css`, `components.css`, `controller.js`, `capture.js`, `master_index.js`, and 14 theme presets) are embedded directly into the binary via `embed.FS`.
3. **Multi-Root Project Discovery**: Run `deckforge` from any directory. Automatically duck-types presentations (`deck.json` or `slides/*.html`), scans subfolders up to depth 2, and allows switching workspace roots dynamically inside the TUI with `w`.
4. **Segmented Theme Architecture**: Dismantles monolithic stylesheets into 4 isolated, easily editable segments:
   - `tokens.json`: Semantic color palette and font definitions
   - `typography.css`: Specific line-heights, weights, and letter-spacings
   - `surfaces.css`: Card borders, glassmorphism filters, and box shadows
   - `backdrop.css`: Atmospheric ambient radial gradients, grids, and glyphs
5. **Configurable Scaffolding**: Create presentations anywhere with a chosen number of slides (1 to 50). Automatically populates diverse layout archetypes (Cover, Agenda, 2-Column, 3-Column, Split Hero, and Code Deep-Dive).
6. **External Editor Hooks (`$EDITOR`)**:
   - `e`: Launches `$EDITOR` (or `$VISUAL`, `code`, `nano`, `vim`) on the highlighted slide file or theme segment.
   - `E` / `a`: Launches `$EDITOR` opening all slides simultaneously.
   - Smoothly suspends and resumes terminal raw mode without disrupting the TUI.
7. **Built-in Live Server & File Watcher**: Integrated Go HTTP server with real-time `fsnotify` recompiler. Automatically updates output on file save.

---

## Interactive TUI Overview

Launch the TUI in any directory:

```bash
./bin/deckforge
# or specify a target path:
./bin/deckforge /path/to/presentations
```

### Workspaces

- **`[1] DECKS`**:
  - Browse discovered presentations in active root.
  - Press `w` to switch root path or pick from recent workspaces.
  - Press `b` to compile presentation into standalone HTML.
  - Press `s` to toggle local Go HTTP server for highlighted presentation.
  - Press `x` to stop running presentation server.
  - Press `o` to launch presentation in default web browser.
  - Press `e` to open presentation directory in external editor.
  - Press `Enter` to inspect slides in Slide Manager.

- **`[2] SLIDES`**:
  - View all sequential slides for active deck.
  - Press `s` to toggle live server for current presentation.
  - Press `x` to stop running presentation server.
  - Press `o` to launch presentation in default web browser.
  - Press `e` to edit highlighted slide in `$EDITOR`.
  - Press `E` or `a` to open all slides simultaneously in `$EDITOR`.
  - Press `n` to add a new slide with chosen layout (`2col`, `3col`, `hero`, `code`).
  - Press `d` to delete highlighted slide.
  - Press `b` to rebuild presentation.

- **`[3] SCAFFOLD`**:
  - Step-by-step wizard to create a new deck.
  - Choose presentation name, destination directory, slide count (1 - 50), and base preset.
  - Press `Enter` to generate files and jump directly to Slide Manager.

- **`[4] THEMES`**:
  - Two-panel workshop: Left theme directory, Right inspector & segment workshop.
  - Interactive terminal slide stage preview rendered with theme colors in ANSI 24-bit RGB.
  - Press `Tab` or `h`/`l` to switch focus between theme list and segment workshop.
  - Press `e` or `Enter` on any segment (`tokens.json`, `typography.css`, `surfaces.css`, `backdrop.css`) to edit only that segment in `$EDITOR`.
  - Press `c` to clone any preset to a custom workspace theme.

- **`[5] SETTINGS`**:
  - Configure preferred external editor (`$EDITOR`, `code --wait`, `nvim`, `vim`, `nano`).
  - Configure default presentation theme and initial slide count.
  - Configure embedded web server port and live `fsnotify` file watcher toggle.
  - Automatically persisted to `~/.config/deckforge/config.json`.

---

## Command Line Interface (CLI)

```bash
# Launch interactive TUI
deckforge [path]

# Scaffolding
deckforge init <name> [--path <target_dir>] [--theme <preset>] [--slides <count>]

# Headless compilation to standalone HTML
deckforge build [path] [--theme <override>]

# Live web server with auto-recompile
deckforge serve [path] [--port 8080] [--watch]

# Headless PDF export (via Chromium / Chrome)
deckforge export [path] [--pdf]

# Slide management
deckforge slide add <path> <title> [--layout 2col|3col|hero|code]

# Theme management
deckforge theme list
deckforge theme clone <preset> <new-name> [target-dir]

# Discovery
deckforge list [--path <root>]
```

---

## 14 Curated Skill Presets

All presets from the `frontend-slides` skill are embedded and available everywhere:

| Preset Name | Typography Pairing | Color Essence | Aesthetic Vibe |
|---|---|---|---|
| `academic-crimson` | Poppins + Plus Jakarta Sans | Crimson `#820024`, Navy `#0f172a`, Light `#f8f9fc` | Institutional, authoritative, Canva ledger |
| `cyber-dark` | Poppins + JetBrains Mono | OLED `#090d16`, Neon `#10b981`, Purple `#818cf8` | Futuristic, DevSecOps, high-contrast |
| `swiss-minimal` | Archivo + Nunito | Pure White, Pure Black, Red `#ff3300` | Stark Bauhaus, typographic discipline |
| `bold-signal` | Archivo Black + Space Grotesk | Dark `#1a1a1a`, Vibrant Orange `#ff5722` | Confident, bold, modern, high-impact |
| `electric-studio` | Manrope + Manrope | Dark `#0a0a0a`, Royal Blue `#4361ee` | Clean studio, high-contrast vertical split |
| `creative-voltage` | Syne + Space Mono | Cobalt `#1a1a2e`, Neon Lime `#d4ff00` | Energetic, retro-modern, creative flair |
| `dark-botanical` | Cormorant + IBM Plex Sans | Charcoal `#0f0f0f`, Warm Gold `#d4a574` | Elegant, sophisticated, premium editorial |
| `notebook-tabs` | Bodoni Moda + DM Sans | Charcoal `#2d2d2d`, Cream Paper `#f8f6f1` | Tactile editorial, pastel section tabs |
| `pastel-geometry` | Plus Jakarta Sans | Soft Sky `#c8d9e6`, Pill Badges `#7c6aad` | Friendly, approachable, modern cards |
| `split-pastel` | Outfit + Outfit | Peach `#f5e6dc`, Lavender `#e4dff0` | Playful vertical split, two-tone cards |
| `vintage-editorial` | Fraunces + Work Sans | Warm Cream `#f5f3ee`, Ink `#1a1a1a` | Witty, literary, bordered publication |
| `neon-cyber` | Poppins + Plus Jakarta Sans | Navy `#0a0f1c`, Cyan `#00ffcc`, Magenta `#ff00aa`| High-tech, neon nightscape, particle glow |
| `terminal-green` | JetBrains Mono | GitHub Dark `#0d1117`, Matrix Green `#39d353`| Developer, hacker, phosphor console |
| `paper-ink` | Cormorant Garamond + Source Serif 4 | Cream `#faf9f7`, Charcoal `#1a1a1a`, Crimson `#c41e3a` | Literary, thoughtful, classical publication |

---

## Keyboard Shortcuts During Presentation

When viewing compiled slides in a browser:

- `→` / `Space` / `PageDown`: Advance to next slide
- `←` / `PageUp`: Return to previous slide
- `M`: Open / close Master Index drawer
- `E`: Toggle live inline editing mode (with localStorage persistence)
- `P`: Capture current slide in 4K UHD PNG (3840x2160)
- `F`: Toggle fullscreen
- `?` / `H`: Show keyboard navigation guide
- `Esc`: Close drawers, modals, and exit edit mode

---

## AI Agent & MCP Integration

DeckForge includes a native [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server and headless agent CLI for autonomous slide generation, theme updates, and WCAG contrast validation.

### Configure MCP Server in Claude Desktop / Cursor

Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "deckforge": {
      "command": "deckforge",
      "args": ["agent", "mcp"]
    }
  }
}
```

### Agent CLI Subcommands

```bash
# JSON inspection of decks
deckforge agent deck-list
deckforge agent deck-inspect <path>

# Programmatic slide modifications
deckforge agent slide-add <deck_path> "Slide Title" --layout 2col
deckforge agent slide-edit <slide_file> --content "..."

# Token diagnostics & WCAG compliance
deckforge agent tokens-diagnose <deck_path>
deckforge agent tokens-set <deck_path> --bg "#0f172a" --text "#f8f9fc"

# Component insertion
deckforge agent component-insert <slide_file> <archetype>
```

---

## Attribution & Origins

DeckForge was created based on the visual design principles and preset archetypes from the official `frontend-slides` skill. It expands those curated aesthetics into a standalone, pure Go presentation engine featuring:

- Interactive Bubble Tea TUI
- Headless Chromium CDP PDF export engine
- Embedded Model Context Protocol (MCP) server & AI CLI interface
- Modular CSS/SCSS token architecture
- Living Development Book specification tracker

---

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) and the [Code of Conduct](CODE_OF_CONDUCT.md) before submitting pull requests or reporting issues.

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

