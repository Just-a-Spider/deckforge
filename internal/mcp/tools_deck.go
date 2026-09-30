package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"deckforge/internal/compiler"
	"deckforge/internal/models"
	"deckforge/internal/scaffold"
	"deckforge/internal/theme"
	"deckforge/internal/workspace"
)

func handleDeckTools(deckPath, name string, args map[string]interface{}) string {
	tm := theme.NewThemeManager(deckPath)
	comp := compiler.NewCompiler(tm)

	switch name {
	case "deckforge_create_deck":
		nameStr, ok := args["name"].(string)
		if !ok || nameStr == "" {
			return "Error: deck name is required"
		}
		cfg := models.LoadWorkspaceConfig()
		defaultTheme := "cyber-dark"
		defaultSlides := 5
		if cfg != nil {
			if cfg.DefaultTheme != "" {
				defaultTheme = cfg.DefaultTheme
			}
			if cfg.DefaultSlideCount > 0 {
				defaultSlides = cfg.DefaultSlideCount
			}
		}

		themeStr, _ := args["theme"].(string)
		if themeStr == "" {
			themeStr = defaultTheme
		}
		slidesCount := defaultSlides
		if sF, ok := args["slides"].(float64); ok && int(sF) > 0 {
			slidesCount = int(sF)
		}
		targetPath, _ := args["path"].(string)
		if targetPath == "" {
			root := deckPath
			if root == "" {
				root, _ = os.Getwd()
			}
			dotDir := models.WorkspaceDotDir(root)
			if fi, err := os.Stat(dotDir); err == nil && fi.IsDir() {
				targetPath = filepath.Join(dotDir, "decks", nameStr)
			} else {
				targetPath = filepath.Join("decks", nameStr)
			}
		}
		deck, err := scaffold.ScaffoldDeck(targetPath, nameStr, themeStr, slidesCount)
		if err != nil {
			return fmt.Sprintf("Error creating deck: %v", err)
		}
		newTM := theme.NewThemeManager(targetPath)
		newComp := compiler.NewCompiler(newTM)
		outFile, _ := newComp.Build(targetPath, themeStr)
		return fmt.Sprintf("Successfully created deck '%s' with %d slides at %s (Compiled: %s)", deck.Name, len(deck.Slides), deck.Path, outFile)

	case "deckforge_get_deck_info":
		deck, err := workspace.InspectDeck(deckPath)
		if err != nil {
			return fmt.Sprintf("Error inspecting deck: %v", err)
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
		return string(data)

	case "deckforge_set_deck_meta":
		deck, err := workspace.InspectDeck(deckPath)
		if err != nil {
			return fmt.Sprintf("Error inspecting deck: %v", err)
		}
		changed := false
		if t, ok := args["title"].(string); ok && t != "" {
			deck.Config.Title = t
			changed = true
		}
		if s, ok := args["subtitle"].(string); ok && s != "" {
			deck.Config.Subtitle = s
			changed = true
		}
		if th, ok := args["theme"].(string); ok && th != "" {
			deck.Config.Theme = th
			deck.ThemeName = th
			changed = true
		}
		if !changed {
			return "No metadata changes specified"
		}
		configPath := filepath.Join(deckPath, "deck.json")
		data, err := json.MarshalIndent(deck.Config, "", "  ")
		if err != nil {
			return fmt.Sprintf("Error serializing config: %v", err)
		}
		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return fmt.Sprintf("Error writing deck.json: %v", err)
		}
		outFile, _ := comp.Build(deckPath, deck.ThemeName)
		return fmt.Sprintf("Updated deck metadata (title: '%s', subtitle: '%s', theme: '%s') and recompiled to %s",
			deck.Config.Title, deck.Config.Subtitle, deck.Config.Theme, outFile)

	default:
		return fmt.Sprintf("Unknown deck tool %s", name)
	}
}
