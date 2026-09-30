package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"deckforge/internal/compiler"
	"deckforge/internal/components"
	"deckforge/internal/theme"
	"deckforge/internal/workspace"
)

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

	cm := components.NewComponentManager(deckPath)
	targetComp, err := cm.GetComponent(selector, deckPath)
	if err != nil || targetComp == nil {
		fmt.Printf(`{"error": "unknown component selector: %s"}`+"\n", selector)
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

func runAgentComponentList(args []string) {
	deckPath := "."
	if len(args) > 0 {
		deckPath = args[0]
	}
	cm := components.NewComponentManager(deckPath)
	comps := cm.ListComponents(deckPath)

	type CompSummary struct {
		Selector     string `json:"selector"`
		Name         string `json:"name"`
		Category     string `json:"category"`
		Scope        string `json:"scope"`
		Description  string `json:"description,omitempty"`
		HasCustomCSS bool   `json:"has_custom_css"`
	}
	var list []CompSummary
	for _, c := range comps {
		list = append(list, CompSummary{
			Selector:     c.Selector,
			Name:         c.Name,
			Category:     c.Category,
			Scope:        c.Scope,
			Description:  c.Description,
			HasCustomCSS: c.Styles != "",
		})
	}
	_ = json.NewEncoder(os.Stdout).Encode(list)
}

func runAgentComponentCreate(args []string) {
	if len(args) < 1 {
		fmt.Println(`{"error": "usage: deckforge agent component-create <name> [--category <cat>] [--global] [--styles]"}`)
		os.Exit(1)
	}
	name := args[0]
	category := "cards"
	global := false
	withCSS := false

	for i := 1; i < len(args); i++ {
		if (args[i] == "--category" || args[i] == "-c") && i+1 < len(args) {
			category = args[i+1]
			i++
		} else if args[i] == "--global" || args[i] == "-g" {
			global = true
		} else if args[i] == "--styles" || args[i] == "--css" {
			withCSS = true
		}
	}

	opts := components.ScaffoldOptions{
		Name:          name,
		Category:      category,
		Global:        global,
		WithCSS:       withCSS,
		WorkspaceRoot: ".",
	}

	created, err := components.ScaffoldComponent(opts)
	if err != nil {
		fmt.Printf(`{"error": "failed to create component: %s"}`+"\n", err)
		os.Exit(1)
	}

	res := map[string]interface{}{
		"status":      "ok",
		"selector":    created.Selector,
		"name":        created.Name,
		"category":    created.Category,
		"scope":       created.Scope,
		"source_path": created.SourcePath,
	}
	_ = json.NewEncoder(os.Stdout).Encode(res)
}

func runAgentCatalog(args []string) {
	deckPath := "."
	if len(args) > 0 {
		deckPath = args[0]
	}
	cm := components.NewComponentManager(deckPath)
	comps := cm.ListComponents(deckPath)
	data, err := json.MarshalIndent(comps, "", "  ")
	if err != nil {
		fmt.Printf(`{"error": "failed to serialize catalog: %v"}`+"\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}
