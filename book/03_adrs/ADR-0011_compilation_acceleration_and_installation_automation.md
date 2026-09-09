# ADR-0011: Compilation Acceleration, Subprocess Caching, and Installation Automation

## Status
Accepted

## Context
Following the implementation of modular SCSS stylesheets in Phase 11, compilation latency of SCSS-enabled themes was throttled by repeated invocation of external CLI compilers (`npx -y sass` / `sass`). Because Node/NPX startup overhead is between 200ms and 800ms per file, transpiling a multi-segment theme (`typography.scss`, `surfaces.scss`, `backdrop.scss`, `components.scss`) spawned 4 external subprocesses, producing a 3.6-second compile latency on every slide change or rebuild.

Furthermore:
1. Core embedded CSS and JS assets were repeatedly read and checked against the disk via filesystem stat syscalls on every build.
2. The live server watcher used timestamp-based debouncing that could miss trailing companion file writes during rapid saves, and lacked file watching on active theme directories (`themes/<name>/` and `~/.config/deckforge/themes/`).
3. Binary distribution lacked automated compilation flags (`-ldflags="-s -w"`), resulting in an unnecessarily bloated 16MB executable.
4. Developers and users had no single-command installation pipeline, no Fish shell completion integration, and no automated global skill deployment.

## Decision
1. **Multi-Tier SCSS Subprocess Caching (`internal/theme/scss.go`)**:
   - Compute `SHA-256` content hash of raw SCSS content prior to transpilation.
   - Tier 1: Check in-memory `sync.Map` cache for instant zero-latency retrieval in long-running processes (`deckforge serve`, `deckforge tui`, MCP stdio server).
   - Tier 2: Check persistent disk cache at `~/.cache/deckforge/scss/<hash>.css` across independent CLI invocations (`deckforge build`, `deckforge agent`).
   - Memoize compiler binary discovery via `sync.Once` (`detectSassCompiler`), avoiding repeated `exec.LookPath` calls.
2. **Compiler Static Asset Memoization (`internal/compiler/compiler.go`)**:
   - Cache embedded core stylesheets (`stage.css`, `hud.css`, `components.css`, `studio.css`) and scripts in a thread-safe memory buffer after initial warm-up, bypassing redundant filesystem stat calls.
3. **Timer-Reset Quiet-Window Watcher (`internal/server/server.go`)**:
   - Implement `time.AfterFunc` (100ms quiet period) that resets upon every incoming filesystem event, guaranteeing that batches of file writes are captured in a single atomic compile.
   - Automatically register active workspace and global theme directories in `fsnotify.Watcher`.
4. **Binary Footprint Stripping & Semantic Versioning**:
   - Strip DWARF symbol tables with `-ldflags="-s -w -X 'deckforge/cmd.Version=0.3.0'"`, reducing binary size from 16MB to 11MB (~31% reduction).
   - Expose `deckforge version` command with operating system, architecture, and Go runtime details.
5. **Fish Shell Autocompletion & Automated Installation**:
   - Author native Fish completion generator (`cmd/completion.go` and `completions/deckforge.fish`) supporting dynamic theme and presentation autocompletion.
   - Provide `deckforge completion fish [--install]`.
   - Provide root [`Makefile`](../../Makefile) (`make install-all`) and standalone [`install.sh`](../../install.sh) automating binary deployment to `~/.local/bin/`, Fish completions to `~/.config/fish/completions/`, and global Gemini skill synchronization to `~/.gemini/skills/deckforge/`.

## Consequences
- **Positive**:
  - Rebuild latency for SCSS themes drops from 3,630ms to 6ms (a 600x speedup; 99.8% latency reduction).
  - Memory-efficient 11MB single binary with zero external runtime dependencies.
  - Complete Fish shell integration with tab completion for commands, subcommands, flags, and themes.
  - One-command setup (`make install-all` or `./install.sh`) for any Linux/macOS developer workstation.
- **Negative**:
  - Disk cache stored in `~/.cache/deckforge/scss/` requires periodic cleanup if thousands of distinct SCSS permutations are generated.
