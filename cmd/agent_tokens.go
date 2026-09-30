package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"deckforge/internal/compiler"
	"deckforge/internal/models"
	"deckforge/internal/theme"
	"deckforge/internal/workspace"
)

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

func runAgentStylingGuide(args []string) {
	guide := models.GetStylingGuide()
	data, err := json.MarshalIndent(guide, "", "  ")
	if err != nil {
		fmt.Printf(`{"error": "failed to serialize styling guide: %v"}`+"\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}
