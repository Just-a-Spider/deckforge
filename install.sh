#!/usr/bin/env bash
# DeckForge One-Command Automated Installer
set -e

echo "=== DeckForge Automated Installer ==="

if ! command -v go >/dev/null 2>&1; then
    echo "Error: Go compiler is required to build DeckForge."
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

VERSION="0.3.0"
INSTALL_BIN="${HOME}/.local/bin"
mkdir -p "$INSTALL_BIN"

echo "-> Compiling stripped binary (v${VERSION})..."
go build -ldflags="-s -w -X 'deckforge/cmd.Version=${VERSION}'" -trimpath -o bin/deckforge .

echo "-> Installing binary to ${INSTALL_BIN}/deckforge..."
cp bin/deckforge "${INSTALL_BIN}/deckforge"
chmod +x "${INSTALL_BIN}/deckforge"

GOPATH_BIN="$(go env GOPATH)/bin"
if [ -d "$GOPATH_BIN" ]; then
    cp bin/deckforge "${GOPATH_BIN}/deckforge"
fi

# Fish shell completions
if [ -d "${HOME}/.config/fish" ]; then
    echo "-> Installing Fish autocompletions..."
    mkdir -p "${HOME}/.config/fish/completions"
    ./bin/deckforge completion fish > "${HOME}/.config/fish/completions/deckforge.fish"
fi

# Global Gemini skill
if [ -d "${HOME}/.gemini" ]; then
    echo "-> Registering global Gemini skill..."
    mkdir -p "${HOME}/.gemini/skills/deckforge"
    cp assets/agent/SKILL.md "${HOME}/.gemini/skills/deckforge/SKILL.md"
fi

echo ""
echo "DeckForge successfully installed!"
"${INSTALL_BIN}/deckforge" version
