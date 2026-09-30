package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func runAgent(args []string) {
	if len(args) < 1 {
		printAgentHelp()
		return
	}

	sub := args[0]
	rest := args[1:]

	switch sub {
	case "deck-info":
		runAgentDeckInfo(rest)
	case "deck-meta":
		runAgentDeckMeta(rest)
	case "deck-create":
		runAgentDeckCreate(rest)
	case "theme-list":
		runAgentThemeList(rest)
	case "theme-set":
		runAgentThemeSet(rest)
	case "theme-create":
		runAgentThemeCreate(rest)
	case "slide-get":
		runAgentSlideGet(rest)
	case "slide-set":
		runAgentSlideSet(rest)
	case "slide-reorder":
		runAgentSlideReorder(rest)
	case "tokens-audit":
		runAgentTokensAudit(rest)
	case "tokens-fix":
		runAgentTokensFix(rest)
	case "component-insert":
		runAgentComponentInsert(rest)
	case "component-list":
		runAgentComponentList(rest)
	case "component-create":
		runAgentComponentCreate(rest)
	case "config-get":
		runAgentConfigGet(rest)
	case "config-set":
		runAgentConfigSet(rest)
	case "catalog":
		runAgentCatalog(rest)
	case "styling-guide":
		runAgentStylingGuide(rest)
	case "export-skill":
		runAgentExportSkill(rest)
	default:
		fmt.Printf(`{"error": "unknown agent subcommand: %s"}`+"\n", sub)
		os.Exit(1)
	}
}

func printAgentHelp() {
	fmt.Println(`DeckForge Agent Interface (AI-Native CLI)

USAGE:
  deckforge agent catalog [deck-path]
  deckforge agent styling-guide
  deckforge agent component-list [deck-path]
  deckforge agent component-create <name> [--category <cat>] [--global] [--styles]
  deckforge agent component-insert <deck-path> <index> <selector> [--props <json>]
  deckforge agent config-get <key> [deck-path]
  deckforge agent config-set <key> <value> [--local]
  deckforge agent deck-info <deck-path>
  deckforge agent deck-meta <deck-path> [--title <title>] [--subtitle <subtitle>] [--theme <theme>]
  deckforge agent deck-create <name> [--theme <preset>] [--slides <n>] [--path <dir>]
  deckforge agent theme-list [deck-path]
  deckforge agent theme-set <deck-path> <theme-name>
  deckforge agent theme-create <name> [--base <preset>] [--global] [--scss]
  deckforge agent slide-get <deck-path> <index>
  deckforge agent slide-set <deck-path> <index> --content <html_or_file>
  deckforge agent slide-reorder <deck-path> --order <1,3,2,...>
  deckforge agent tokens-audit <deck-path>
  deckforge agent tokens-fix <deck-path> [--apply]
  deckforge agent export-skill [target-dir]`)
}

func runAgentExportSkill(args []string) {
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}
	_ = os.MkdirAll(targetDir, 0755)

	skillPath := filepath.Join(targetDir, "SKILL.md")
	content := `# DeckForge Skill

Instructions and conventions for autonomous coding agents (Claude Code, Gemini, Antigravity, OpenClaw, OpenCode) interacting with DeckForge presentations.

## 0. Zero Engine Dependency (Self-Contained Protocol)
- **NEVER search for, clone, or read files in the deckforge engine repository**.
- DeckForge is an installed binary and completely self-contained.
- Query components dynamically: 'deckforge agent catalog' (or MCP 'deckforge_get_catalog').
- Query CSS classes and design tokens: 'deckforge agent styling-guide' (or MCP 'deckforge_get_styling_guide').
- Query theme presets: 'deckforge agent theme-list' (or MCP 'deckforge_list_themes').
- Presentations belong in '<workspace>/.deckforge/decks/<name>/' or './decks/<name>/'.

## 1. Core Directives
1. **Deterministic 1080p Stage**: Always respect 1920x1080 master stage geometry. Do NOT write responsive flex wrapping or arbitrary dynamic viewport units.
2. **Discrete Slide Chunks**: Every slide lives in 'slides/NN_slug.html'. Edit slides cleanly without adding runtime attributes (no contenteditable, no df-*).
3. **Semantic Classes**:
   - Layouts: '.grid-2col', '.grid-3col', '.grid-4col', '.grid-split-hero'
   - Panels: '.glass-panel', '.lime-edge', '.accent-edge'
   - Badges: '.card-tag', '.topbar-tag', '.lime-badge', '.purple-badge'
   - Typography: '.topbar-title', '.panel-headline', '.micro-kicker', '.bullet-list', '.lime-hl', '.accent-hl'
4. **Agent Commands**:
   - 'deckforge agent catalog [deck-path]'
   - 'deckforge agent styling-guide'
   - 'deckforge agent deck-create <name> [--theme <preset>] [--slides <count>] [--path <dir>]'
   - 'deckforge agent theme-list [deck-path]'
   - 'deckforge agent theme-set <deck-path> <theme-name>'
   - 'deckforge agent theme-create <name> [--base <preset>] [--global] [--scss]'
   - 'deckforge agent slide-get <deck-path> <index>'
   - 'deckforge agent slide-set <deck-path> <index> --content <html_or_file>'
   - 'deckforge agent slide-reorder <deck-path> --order 1,3,2'
   - 'deckforge agent tokens-audit <deck-path>'
   - 'deckforge agent tokens-fix <deck-path> --apply'
   - 'deckforge agent component-insert <deck-path> <index> <selector>'
   - 'deckforge agent component-list [deck-path]'
   - 'deckforge agent component-create <name> [--category <cat>] [--global] [--styles]'
   - 'deckforge agent config-get <key> [deck-path]'
   - 'deckforge agent config-set <key> <val> [--local]'
`
	if err := os.WriteFile(skillPath, []byte(content), 0644); err != nil {
		fmt.Printf(`{"error": "failed to write skill: %s"}`+"\n", err)
		os.Exit(1)
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]string{
		"status": "ok",
		"path":   skillPath,
	})
}
