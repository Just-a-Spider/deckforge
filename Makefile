# DeckForge - Pure Go Modular 1080p Presentation Platform
VERSION ?= 0.3.0
PREFIX ?= $(HOME)/.local
GOPATH ?= $(shell go env GOPATH)
LDFLAGS := -s -w -X 'deckforge/cmd.Version=$(VERSION)'

.PHONY: all build install completions skill install-all test clean audit

all: build

build:
	@echo "Building stripped DeckForge binary..."
	@mkdir -p bin
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -trimpath -o bin/deckforge .

install: build
	@echo "Installing DeckForge binary to $(PREFIX)/bin..."
	@mkdir -p $(PREFIX)/bin
	@cp bin/deckforge $(PREFIX)/bin/deckforge
	@chmod +x $(PREFIX)/bin/deckforge
	@if [ -d "$(GOPATH)/bin" ]; then \
		cp bin/deckforge $(GOPATH)/bin/deckforge; \
	fi
	@echo "Installed deckforge successfully: $$(which deckforge || echo $(PREFIX)/bin/deckforge)"

completions: build
	@echo "Installing Fish shell completions..."
	@mkdir -p $(HOME)/.config/fish/completions
	./bin/deckforge completion fish --install

skill:
	@echo "Installing global Gemini skill..."
	@mkdir -p $(HOME)/.gemini/skills/deckforge
	@cp assets/agent/SKILL.md $(HOME)/.gemini/skills/deckforge/SKILL.md
	@cp assets/agent/SKILL.md SKILL.md
	@echo "Global skill synced to $(HOME)/.gemini/skills/deckforge/SKILL.md"

install-all: build install completions skill
	@echo ""
	@echo "DeckForge installation complete."
	@echo "  Binary:      $$(which deckforge || echo $(PREFIX)/bin/deckforge)"
	@echo "  Completions: $(HOME)/.config/fish/completions/deckforge.fish"
	@echo "  Skill:       $(HOME)/.gemini/skills/deckforge/SKILL.md"
	@echo "  Version:     $$(./bin/deckforge version)"

test:
	go test -v ./...

audit:
	./bin/deckforge book audit

clean:
	rm -rf bin/
	rm -f /tmp/deckforge*
