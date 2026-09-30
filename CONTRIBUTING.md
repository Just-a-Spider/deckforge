# Contributing to DeckForge

Thank you for your interest in contributing to DeckForge! We welcome bug reports, feature requests, documentation improvements, theme additions, and code contributions.

---

## Development Setup

### Prerequisites

- **Go**: 1.23+
- **Chrome / Chromium**: For headless PDF and PNG exports.
- **Node.js (Optional)**: If using Sass/SCSS external transpilers or the optional Node PDF exporter.
- **Git**

### Build from Source

```bash
git clone https://github.com/Just-a-Spider/deckforge.git
cd deckforge
make build
```

Run tests to ensure everything is green:

```bash
make test
```

---

## Submitting Changes

1. **Fork the repository** on GitHub.
2. **Create a branch** for your feature or bug fix:
   ```bash
   git checkout -b feat/your-feature-name
   # or
   git checkout -b fix/issue-description
   ```
3. **Write tests** for any new features or bug fixes.
4. **Follow code conventions**:
   - Keep Go code idiomatic (`gofmt`, standard naming).
   - Maintain the Living Development Book in `book/` for significant architectural changes (record an ADR in `book/03_adrs/` if making architectural shifts).
5. **Ensure all tests pass**:
   ```bash
   go test -v ./...
   ```
6. **Commit with conventional commits**:
   - `feat: ...`
   - `fix: ...`
   - `docs: ...`
   - `test: ...`
   - `refactor: ...`
7. **Open a Pull Request** against the `main` branch.

---

## Reporting Issues

- Check existing issues before opening a new one.
- Provide minimal reproduction steps, OS environment, and Go version (`go version`).
