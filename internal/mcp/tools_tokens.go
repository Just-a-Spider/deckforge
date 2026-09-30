package mcp

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"deckforge/internal/compiler"
	"deckforge/internal/models"
	"deckforge/internal/theme"
	"deckforge/internal/workspace"
)

func handleTokenTools(deckPath, name string, args map[string]interface{}) string {
	if name == "deckforge_get_styling_guide" {
		guide := models.GetStylingGuide()
		data, _ := json.MarshalIndent(guide, "", "  ")
		return string(data)
	}

	deck, err := workspace.InspectDeck(deckPath)
	if err != nil {
		return fmt.Sprintf("Error inspecting deck: %v", err)
	}

	tm := theme.NewThemeManager(deckPath)
	comp := compiler.NewCompiler(tm)

	switch name {
	case "deckforge_audit_tokens":
		resolved, err := tm.ResolveTheme(deck.ThemeName, deckPath)
		if err != nil {
			return fmt.Sprintf("Error resolving theme: %v", err)
		}
		reports := models.EvaluateThemeContrast(resolved.Tokens.Palette)
		data, _ := json.MarshalIndent(reports, "", "  ")
		return string(data)

	case "deckforge_apply_token_fix":
		resolved, err := tm.ResolveTheme(deck.ThemeName, deckPath)
		if err != nil {
			return fmt.Sprintf("Error resolving theme: %v", err)
		}
		reports := models.EvaluateThemeContrast(resolved.Tokens.Palette)
		fixed := 0
		for _, r := range reports {
			if r.RecommendedAction != nil {
				rec := r.RecommendedAction
				switch rec.TargetToken {
				case "textPrimary":
					resolved.Tokens.Palette.TextPrimary = rec.SuggestedHex
					fixed++
				case "textSecondary":
					resolved.Tokens.Palette.TextSecondary = rec.SuggestedHex
					fixed++
				case "accentPrimary":
					resolved.Tokens.Palette.AccentPrimary = rec.SuggestedHex
					fixed++
				}
			}
		}
		if fixed > 0 {
			dir := resolved.Dir
			if dir == "" {
				dir = filepath.Join(deckPath, "themes", resolved.Tokens.Name)
			}
			_ = theme.SaveSegmentedTheme(dir, resolved)
			_, _ = comp.Build(deckPath, "")
		}
		return fmt.Sprintf("Applied %d contrast token adjustments to theme %s", fixed, resolved.Tokens.Name)

	default:
		return fmt.Sprintf("Unknown token tool %s", name)
	}
}
