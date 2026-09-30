package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"deckforge/internal/compiler"
	"deckforge/internal/components"
	"deckforge/internal/theme"
	"deckforge/internal/workspace"
)

func handleComponentTools(deckPath, name string, args map[string]interface{}) string {
	cm := components.NewComponentManager(deckPath)

	switch name {
	case "deckforge_get_catalog", "deckforge_list_components":
		comps := cm.ListComponents(deckPath)
		data, _ := json.MarshalIndent(comps, "", "  ")
		return string(data)

	case "deckforge_create_component":
		nameStr, ok := args["name"].(string)
		if !ok || nameStr == "" {
			return "Error: component name is required"
		}
		cat, _ := args["category"].(string)
		desc, _ := args["description"].(string)
		global, _ := args["global"].(bool)
		withStyles, _ := args["styles"].(bool)

		opts := components.ScaffoldOptions{
			Name:          nameStr,
			Category:      cat,
			Description:   desc,
			Global:        global,
			WithCSS:       withStyles,
			WorkspaceRoot: deckPath,
		}
		created, err := components.ScaffoldComponent(opts)
		if err != nil {
			return fmt.Sprintf("Error scaffolding component: %v", err)
		}
		data, _ := json.MarshalIndent(created, "", "  ")
		return string(data)

	case "deckforge_insert_component":
		deck, err := workspace.InspectDeck(deckPath)
		if err != nil {
			return fmt.Sprintf("Error inspecting deck: %v", err)
		}
		idxF, ok := args["index"].(float64)
		selector, ok2 := args["selector"].(string)
		if !ok || !ok2 || int(idxF) < 1 || int(idxF) > len(deck.Slides) {
			return "Error: invalid index or selector"
		}
		props, _ := args["props"].(map[string]interface{})
		if props == nil {
			props = make(map[string]interface{})
		}
		targetComp, err := cm.GetComponent(selector, deckPath)
		if err != nil || targetComp == nil {
			return fmt.Sprintf("Error: unknown component %s", selector)
		}
		rendered, err := targetComp.Render(props, nil)
		if err != nil {
			return fmt.Sprintf("Error rendering component: %v", err)
		}
		target := deck.Slides[int(idxF)-1]
		content, _ := os.ReadFile(target.Path)
		curStr := string(content)
		lastClose := strings.LastIndex(curStr, "</section>")
		var updated string
		if lastClose != -1 {
			updated = curStr[:lastClose] + "\n    " + rendered + "\n" + curStr[lastClose:]
		} else {
			updated = curStr + "\n" + rendered
		}
		_ = os.WriteFile(target.Path, []byte(updated), 0644)

		tm := theme.NewThemeManager(deckPath)
		comp := compiler.NewCompiler(tm)
		out, _ := comp.Build(deckPath, "")
		return fmt.Sprintf("Inserted %s into slide %d and rebuilt to %s", selector, int(idxF), out)

	default:
		return fmt.Sprintf("Unknown component tool %s", name)
	}
}
