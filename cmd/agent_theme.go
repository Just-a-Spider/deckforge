package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"deckforge/internal/compiler"
	"deckforge/internal/theme"
	"deckforge/internal/workspace"
)

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
