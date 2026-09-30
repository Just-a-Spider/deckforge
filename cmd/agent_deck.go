package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"deckforge/internal/compiler"
	"deckforge/internal/models"
	"deckforge/internal/scaffold"
	"deckforge/internal/theme"
	"deckforge/internal/workspace"
)

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

func runAgentDeckCreate(args []string) {
	if len(args) < 1 {
		fmt.Println(`{"error": "usage: deckforge agent deck-create <name> [--theme <preset>] [--slides <count>] [--path <dir>]"}`)
		os.Exit(1)
	}
	name := args[0]
	cfg := models.LoadWorkspaceConfig()
	themeName := "cyber-dark"
	slideCount := 5
	if cfg != nil {
		if cfg.DefaultTheme != "" {
			themeName = cfg.DefaultTheme
		}
		if cfg.DefaultSlideCount > 0 {
			slideCount = cfg.DefaultSlideCount
		}
	}
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
		dotDir := models.WorkspaceDotDir(cwd)
		if fi, err := os.Stat(dotDir); err == nil && fi.IsDir() {
			targetPath = filepath.Join(dotDir, "decks", name)
		} else {
			targetPath = filepath.Join(cwd, "decks", name)
		}
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
