package mcp

import (
	"encoding/json"
	"fmt"

	"deckforge/internal/compiler"
	"deckforge/internal/theme"
	"deckforge/internal/workspace"
)

func handleThemeTools(deckPath, name string, args map[string]interface{}) string {
	tm := theme.NewThemeManager(deckPath)
	comp := compiler.NewCompiler(tm)

	switch name {
	case "deckforge_list_themes":
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
		data, _ := json.MarshalIndent(list, "", "  ")
		return string(data)

	case "deckforge_set_theme":
		themeStr, ok := args["theme"].(string)
		if !ok || themeStr == "" {
			return "Error: theme name is required"
		}
		if err := workspace.SetDeckTheme(deckPath, themeStr); err != nil {
			return fmt.Sprintf("Error setting theme: %v", err)
		}
		outFile, err := comp.Build(deckPath, themeStr)
		if err != nil {
			return fmt.Sprintf("Error rebuilding deck: %v", err)
		}
		return fmt.Sprintf("Successfully updated theme to '%s' and recompiled to %s", themeStr, outFile)

	case "deckforge_create_theme":
		nameStr, ok := args["name"].(string)
		if !ok || nameStr == "" {
			return "Error: theme name is required"
		}
		baseStr, _ := args["base"].(string)
		global, _ := args["global"].(bool)
		scss, _ := args["scss"].(bool)

		created, err := tm.CreateTheme(nameStr, baseStr, global, scss)
		if err != nil {
			return fmt.Sprintf("Error creating theme: %v", err)
		}
		return fmt.Sprintf("Successfully created theme '%s' (scope: %s) at %s", created.Tokens.Name, created.Scope, created.Dir)

	default:
		return fmt.Sprintf("Unknown theme tool %s", name)
	}
}
