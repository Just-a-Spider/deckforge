package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"deckforge/internal/compiler"
	"deckforge/internal/models"
	"deckforge/internal/scaffold"
	"deckforge/internal/server"
	"deckforge/internal/theme"
	"deckforge/internal/workspace"
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
  deckforge agent catalog
  deckforge agent styling-guide
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
  deckforge agent component-insert <deck-path> <index> <selector> [--props <json>]
  deckforge agent export-skill [target-dir]`)
}

func runAgentSlideGet(args []string) {
	if len(args) < 2 {
		fmt.Println(`{"error": "usage: deckforge agent slide-get <deck-path> <index>"}`)
		os.Exit(1)
	}
	deckPath := args[0]
	idx, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println(`{"error": "index must be an integer"}`)
		os.Exit(1)
	}

	deck, err := workspace.InspectDeck(deckPath)
	if err != nil || idx < 1 || idx > len(deck.Slides) {
		fmt.Println(`{"error": "slide index out of bounds"}`)
		os.Exit(1)
	}

	target := deck.Slides[idx-1]
	content, err := os.ReadFile(target.Path)
	if err != nil {
		fmt.Printf(`{"error": "failed to read slide file: %s"}`+"\n", err)
		os.Exit(1)
	}

	res := map[string]interface{}{
		"index":    target.Index,
		"filename": target.Filename,
		"path":     target.Path,
		"title":    target.Title,
		"slug":     target.Slug,
		"html":     string(content),
	}
	_ = json.NewEncoder(os.Stdout).Encode(res)
}

func runAgentSlideSet(args []string) {
	if len(args) < 4 || args[2] != "--content" {
		fmt.Println(`{"error": "usage: deckforge agent slide-set <deck-path> <index> --content <html_or_file>"}`)
		os.Exit(1)
	}

	deckPath := args[0]
	idx, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println(`{"error": "index must be an integer"}`)
		os.Exit(1)
	}

	rawContent := args[3]
	if fi, err := os.Stat(rawContent); err == nil && !fi.IsDir() {
		data, err := os.ReadFile(rawContent)
		if err == nil {
			rawContent = string(data)
		}
	}

	deck, err := workspace.InspectDeck(deckPath)
	if err != nil || idx < 1 || idx > len(deck.Slides) {
		fmt.Println(`{"error": "slide index out of bounds"}`)
		os.Exit(1)
	}

	target := deck.Slides[idx-1]
	cleaned := server.SanitizeSlideHTML(rawContent)

	if err := os.WriteFile(target.Path, []byte(cleaned), 0644); err != nil {
		fmt.Printf(`{"error": "failed to write slide: %s"}`+"\n", err)
		os.Exit(1)
	}

	tm := theme.NewThemeManager(deckPath)
	comp := compiler.NewCompiler(tm)
	out, err := comp.Build(deckPath, "")
	if err != nil {
		fmt.Printf(`{"error": "recompile failed: %s"}`+"\n", err)
		os.Exit(1)
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
		"status":      "ok",
		"index":       idx,
		"path":        target.Path,
		"compiled_to": out,
	})
}

func runAgentSlideReorder(args []string) {
	if len(args) < 3 || args[1] != "--order" {
		fmt.Println(`{"error": "usage: deckforge agent slide-reorder <deck-path> --order <1,3,2,...>"}`)
		os.Exit(1)
	}

	deckPath := args[0]
	parts := strings.Split(args[2], ",")
	var newOrder []int
	for _, p := range parts {
		num, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			fmt.Println(`{"error": "invalid order index"}`)
			os.Exit(1)
		}
		newOrder = append(newOrder, num)
	}

	if err := workspace.ReorderSlides(deckPath, newOrder); err != nil {
		fmt.Printf(`{"error": "reorder failed: %s"}`+"\n", err)
		os.Exit(1)
	}

	tm := theme.NewThemeManager(deckPath)
	comp := compiler.NewCompiler(tm)
	out, _ := comp.Build(deckPath, "")

	_ = json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
		"status":      "ok",
		"new_order":   newOrder,
		"compiled_to": out,
	})
}

func runAgentTokensAudit(args []string) {
	if len(args) < 1 {
		fmt.Println(`{"error": "usage: deckforge agent tokens-audit <deck-path>"}`)
		os.Exit(1)
	}
	deckPath := args[0]
	deck, err := workspace.InspectDeck(deckPath)
	if err != nil {
		fmt.Printf(`{"error": "inspect failed: %s"}`+"\n", err)
		os.Exit(1)
	}

	tm := theme.NewThemeManager(deckPath)
	resolved, err := tm.ResolveTheme(deck.ThemeName, deckPath)
	if err != nil {
		fmt.Printf(`{"error": "theme resolve failed: %s"}`+"\n", err)
		os.Exit(1)
	}

	reports := models.EvaluateThemeContrast(resolved.Tokens.Palette)
	_ = json.NewEncoder(os.Stdout).Encode(reports)
}

func runAgentTokensFix(args []string) {
	if len(args) < 1 {
		fmt.Println(`{"error": "usage: deckforge agent tokens-fix <deck-path> [--apply]"}`)
		os.Exit(1)
	}
	deckPath := args[0]
	apply := len(args) > 1 && args[1] == "--apply"

	deck, err := workspace.InspectDeck(deckPath)
	if err != nil {
		fmt.Printf(`{"error": "inspect failed: %s"}`+"\n", err)
		os.Exit(1)
	}

	tm := theme.NewThemeManager(deckPath)
	resolved, err := tm.ResolveTheme(deck.ThemeName, deckPath)
	if err != nil {
		fmt.Printf(`{"error": "theme resolve failed: %s"}`+"\n", err)
		os.Exit(1)
	}

	reports := models.EvaluateThemeContrast(resolved.Tokens.Palette)
	fixedCount := 0

	if apply {
		for _, rep := range reports {
			if rep.RecommendedAction != nil {
				rec := rep.RecommendedAction
				switch rec.TargetToken {
				case "textPrimary":
					resolved.Tokens.Palette.TextPrimary = rec.SuggestedHex
					fixedCount++
				case "textSecondary":
					resolved.Tokens.Palette.TextSecondary = rec.SuggestedHex
					fixedCount++
				case "accentPrimary":
					resolved.Tokens.Palette.AccentPrimary = rec.SuggestedHex
					fixedCount++
				}
			}
		}

		if fixedCount > 0 {
			themeDir := resolved.Dir
			if themeDir == "" {
				themeDir = filepath.Join(deckPath, "themes", resolved.Tokens.Name)
			}
			_ = theme.SaveSegmentedTheme(themeDir, resolved)
			comp := compiler.NewCompiler(tm)
			_, _ = comp.Build(deckPath, "")
		}
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
		"applied":     apply,
		"fixed_count": fixedCount,
		"reports":     reports,
	})
}

func runAgentComponentInsert(args []string) {
	if len(args) < 3 {
		fmt.Println(`{"error": "usage: deckforge agent component-insert <deck-path> <index> <selector> [--props <json>]"}`)
		os.Exit(1)
	}

	deckPath := args[0]
	idx, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println(`{"error": "index must be an integer"}`)
		os.Exit(1)
	}
	selector := args[2]

	inputs := make(map[string]interface{})
	if len(args) > 4 && args[3] == "--props" {
		_ = json.Unmarshal([]byte(args[4]), &inputs)
	}

	var targetComp *models.ComponentDefinition
	for _, c := range models.GetBuiltinComponents() {
		if c.Selector == selector {
			targetComp = &c
			break
		}
	}
	if targetComp == nil {
		fmt.Println(`{"error": "unknown component selector"}`)
		os.Exit(1)
	}

	html, err := targetComp.Render(inputs, nil)
	if err != nil {
		fmt.Printf(`{"error": "render failed: %s"}`+"\n", err)
		os.Exit(1)
	}

	deck, err := workspace.InspectDeck(deckPath)
	if err != nil || idx < 1 || idx > len(deck.Slides) {
		fmt.Println(`{"error": "slide index out of bounds"}`)
		os.Exit(1)
	}

	target := deck.Slides[idx-1]
	curContent, err := os.ReadFile(target.Path)
	if err != nil {
		fmt.Printf(`{"error": "failed to read slide: %s"}`+"\n", err)
		os.Exit(1)
	}

	// Insert before closing </section>
	curStr := string(curContent)
	lastClose := strings.LastIndex(curStr, "</section>")
	var updated string
	if lastClose != -1 {
		updated = curStr[:lastClose] + "\n    " + html + "\n" + curStr[lastClose:]
	} else {
		updated = curStr + "\n" + html
	}

	_ = os.WriteFile(target.Path, []byte(updated), 0644)
	tm := theme.NewThemeManager(deckPath)
	comp := compiler.NewCompiler(tm)
	out, _ := comp.Build(deckPath, "")

	_ = json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
		"status":      "ok",
		"selector":    selector,
		"index":       idx,
		"compiled_to": out,
	})
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

## Core Directives
1. **Deterministic 1080p Stage**: Always respect 1920x1080 master stage geometry. Do NOT write responsive flex wrapping or arbitrary dynamic viewport units.
2. **Discrete Slide Chunks**: Every slide lives in 'slides/NN_slug.html'. Edit slides cleanly without adding runtime attributes (no contenteditable, no df-*).
3. **Semantic Classes**:
   - Layouts: '.grid-2col', '.grid-3col', '.grid-4col', '.grid-split-hero'
   - Panels: '.glass-panel', '.lime-edge', '.accent-edge'
   - Badges: '.card-tag', '.topbar-tag', '.lime-badge', '.purple-badge'
   - Typography: '.topbar-title', '.panel-headline', '.micro-kicker', '.bullet-list', '.lime-hl', '.accent-hl'
4. **Agent Commands**:
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

func runAgentDeckCreate(args []string) {
	if len(args) < 1 {
		fmt.Println(`{"error": "usage: deckforge agent deck-create <name> [--theme <preset>] [--slides <count>] [--path <dir>]"}`)
		os.Exit(1)
	}
	name := args[0]
	themeName := "cyber-dark"
	slideCount := 5
	targetPath := ""

	for i := 1; i < len(args); i++ {
		if args[i] == "--theme" && i+1 < len(args) {
			themeName = args[i+1]
			i++
		} else if args[i] == "--slides" && i+1 < len(args) {
			if count, err := strconv.Atoi(args[i+1]); err == nil {
				slideCount = count
			}
			i++
		} else if args[i] == "--path" && i+1 < len(args) {
			targetPath = args[i+1]
			i++
		}
	}

	if targetPath == "" {
		cwd, _ := os.Getwd()
		targetPath = filepath.Join(cwd, "decks", name)
	}

	deck, err := scaffold.ScaffoldDeck(targetPath, name, themeName, slideCount)
	if err != nil {
		fmt.Printf(`{"error": "failed to create deck: %s"}`+"\n", err)
		os.Exit(1)
	}

	tm := theme.NewThemeManager(targetPath)
	comp := compiler.NewCompiler(tm)
	outFile, _ := comp.Build(targetPath, themeName)

	res := map[string]interface{}{
		"status":      "ok",
		"name":        deck.Name,
		"path":        deck.Path,
		"theme":       deck.ThemeName,
		"slides":      len(deck.Slides),
		"compiled_to": outFile,
	}
	_ = json.NewEncoder(os.Stdout).Encode(res)
}

func runAgentThemeList(args []string) {
	deckPath := "."
	if len(args) > 0 {
		deckPath = args[0]
	}
	tm := theme.NewThemeManager(deckPath)
	themes := tm.ListThemes(deckPath)

	type ThemeSummary struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		Description string `json:"description,omitempty"`
		Vibe        string `json:"vibe"`
		IsCustom    bool   `json:"isCustom"`
		Scope       string `json:"scope,omitempty"`
	}
	var list []ThemeSummary
	for _, t := range themes {
		list = append(list, ThemeSummary{
			Name:        t.Tokens.Name,
			DisplayName: t.Tokens.DisplayName,
			Description: t.Tokens.Description,
			Vibe:        t.Tokens.Vibe,
			IsCustom:    t.IsCustom,
			Scope:       t.Scope,
		})
	}
	_ = json.NewEncoder(os.Stdout).Encode(list)
}

func runAgentThemeSet(args []string) {
	if len(args) < 2 {
		fmt.Println(`{"error": "usage: deckforge agent theme-set <deck-path> <theme-name>"}`)
		os.Exit(1)
	}
	deckPath := args[0]
	themeName := args[1]

	if err := workspace.SetDeckTheme(deckPath, themeName); err != nil {
		fmt.Printf(`{"error": "failed to set theme: %s"}`+"\n", err)
		os.Exit(1)
	}

	tm := theme.NewThemeManager(deckPath)
	comp := compiler.NewCompiler(tm)
	outFile, err := comp.Build(deckPath, themeName)
	if err != nil {
		fmt.Printf(`{"error": "recompile failed: %s"}`+"\n", err)
		os.Exit(1)
	}

	res := map[string]interface{}{
		"status":      "ok",
		"deck_path":   deckPath,
		"theme":       themeName,
		"compiled_to": outFile,
	}
	_ = json.NewEncoder(os.Stdout).Encode(res)
}

func runAgentThemeCreate(args []string) {
	if len(args) < 1 {
		fmt.Println(`{"error": "usage: deckforge agent theme-create <name> [--base <preset>] [--global] [--scss]"}`)
		os.Exit(1)
	}
	name := args[0]
	base := ""
	global := false
	scss := false
	for i := 1; i < len(args); i++ {
		if (args[i] == "--base" || args[i] == "-b") && i+1 < len(args) {
			base = args[i+1]
			i++
		} else if args[i] == "--global" || args[i] == "-g" {
			global = true
		} else if args[i] == "--scss" {
			scss = true
		}
	}

	tm := theme.NewThemeManager(".")
	created, err := tm.CreateTheme(name, base, global, scss)
	if err != nil {
		fmt.Printf(`{"error": "failed to create theme: %s"}`+"\n", err)
		os.Exit(1)
	}

	ext := ".css"
	if scss {
		ext = ".scss"
	}
	res := map[string]interface{}{
		"status":       "ok",
		"name":         created.Tokens.Name,
		"display_name": created.Tokens.DisplayName,
		"scope":        created.Scope,
		"dir":          created.Dir,
		"files": []string{
			"tokens.json",
			"surfaces" + ext,
			"backdrop" + ext,
			"components" + ext,
			"typography" + ext,
		},
	}
	_ = json.NewEncoder(os.Stdout).Encode(res)
}

func runAgentCatalog(args []string) {
	components := models.GetBuiltinComponents()
	data, err := json.MarshalIndent(components, "", "  ")
	if err != nil {
		fmt.Printf(`{"error": "failed to serialize catalog: %v"}`+"\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func runAgentStylingGuide(args []string) {
	guide := models.GetStylingGuide()
	data, err := json.MarshalIndent(guide, "", "  ")
	if err != nil {
		fmt.Printf(`{"error": "failed to serialize styling guide: %v"}`+"\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func runAgentDeckInfo(args []string) {
	if len(args) < 1 {
		fmt.Println(`{"error": "usage: deckforge agent deck-info <deck-path>"}`)
		os.Exit(1)
	}
	deckPath := args[0]
	deck, err := workspace.InspectDeck(deckPath)
	if err != nil {
		fmt.Printf(`{"error": "failed to inspect deck: %s"}`+"\n", err)
		os.Exit(1)
	}
	type DeckSummary struct {
		Name       string             `json:"name"`
		Path       string             `json:"path"`
		Title      string             `json:"title"`
		Subtitle   string             `json:"subtitle"`
		Theme      string             `json:"theme"`
		SlideCount int                `json:"slide_count"`
		Slides     []models.SlideInfo `json:"slides"`
	}
	summary := DeckSummary{
		Name:       deck.Name,
		Path:       deck.Path,
		Title:      deck.Config.Title,
		Subtitle:   deck.Config.Subtitle,
		Theme:      deck.Config.Theme,
		SlideCount: len(deck.Slides),
		Slides:     deck.Slides,
	}
	data, _ := json.MarshalIndent(summary, "", "  ")
	fmt.Println(string(data))
}

func runAgentDeckMeta(args []string) {
	if len(args) < 1 {
		fmt.Println(`{"error": "usage: deckforge agent deck-meta <deck-path> [--title <title>] [--subtitle <subtitle>] [--theme <theme>]"}`)
		os.Exit(1)
	}
	deckPath := args[0]
	deck, err := workspace.InspectDeck(deckPath)
	if err != nil {
		fmt.Printf(`{"error": "failed to inspect deck: %s"}`+"\n", err)
		os.Exit(1)
	}

	norm := normalizeArgs(args[1:])
	fs := flag.NewFlagSet("deck-meta", flag.ExitOnError)
	title := fs.String("title", "", "Update presentation title")
	subtitle := fs.String("subtitle", "", "Update presentation subtitle")
	themeName := fs.String("theme", "", "Update presentation theme")
	_ = fs.Parse(norm)

	if *title != "" {
		deck.Config.Title = *title
	}
	if *subtitle != "" {
		deck.Config.Subtitle = *subtitle
	}
	if *themeName != "" {
		deck.Config.Theme = *themeName
		deck.ThemeName = *themeName
	}

	configPath := filepath.Join(deckPath, "deck.json")
	data, err := json.MarshalIndent(deck.Config, "", "  ")
	if err != nil {
		fmt.Printf(`{"error": "failed to serialize deck config: %s"}`+"\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		fmt.Printf(`{"error": "failed to write deck.json: %s"}`+"\n", err)
		os.Exit(1)
	}

	fmt.Printf(`{"status": "success", "title": "%s", "subtitle": "%s", "theme": "%s"}`+"\n",
		deck.Config.Title, deck.Config.Subtitle, deck.Config.Theme)
}


