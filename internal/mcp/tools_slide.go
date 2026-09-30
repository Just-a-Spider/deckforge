package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"deckforge/internal/compiler"
	"deckforge/internal/theme"
	"deckforge/internal/workspace"
)

var (
	contentEditableRegex = regexp.MustCompile(`\s*contenteditable="[^"]*"`)
	dfSelectedAttrRegex  = regexp.MustCompile(`\s*data-df-selected="[^"]*"`)
	dataStudioRegex      = regexp.MustCompile(`\s*data-studio-[a-z0-9_-]+="[^"]*"`)
	runtimeClassRegex    = regexp.MustCompile(`\s*(df-selected|df-drag-over|df-active-drop)`)
	emptyClassRegex      = regexp.MustCompile(`\s*class="(?:\s*)"`)
)

func sanitizeSlideHTML(raw string) string {
	cleaned := contentEditableRegex.ReplaceAllString(raw, "")
	cleaned = dfSelectedAttrRegex.ReplaceAllString(cleaned, "")
	cleaned = dataStudioRegex.ReplaceAllString(cleaned, "")

	cleaned = runtimeClassRegex.ReplaceAllString(cleaned, "")
	cleaned = emptyClassRegex.ReplaceAllString(cleaned, "")

	lines := strings.Split(cleaned, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t\r")
	}
	res := strings.Join(lines, "\n")
	return strings.TrimSpace(res) + "\n"
}

func handleSlideTools(deckPath, name string, args map[string]interface{}) string {
	deck, err := workspace.InspectDeck(deckPath)
	if err != nil {
		return fmt.Sprintf("Error inspecting deck: %v", err)
	}

	tm := theme.NewThemeManager(deckPath)
	comp := compiler.NewCompiler(tm)

	switch name {
	case "deckforge_list_slides":
		data, _ := json.MarshalIndent(deck.Slides, "", "  ")
		return string(data)

	case "deckforge_get_slide":
		idxF, ok := args["index"].(float64)
		if !ok || int(idxF) < 1 || int(idxF) > len(deck.Slides) {
			return "Error: slide index out of bounds"
		}
		target := deck.Slides[int(idxF)-1]
		content, err := os.ReadFile(target.Path)
		if err != nil {
			return fmt.Sprintf("Error reading slide: %v", err)
		}
		return string(content)

	case "deckforge_update_slide":
		idxF, ok := args["index"].(float64)
		htmlStr, ok2 := args["html"].(string)
		if !ok || !ok2 || int(idxF) < 1 || int(idxF) > len(deck.Slides) {
			return "Error: invalid index or html"
		}
		target := deck.Slides[int(idxF)-1]
		cleaned := sanitizeSlideHTML(htmlStr)
		if err := os.WriteFile(target.Path, []byte(cleaned), 0644); err != nil {
			return fmt.Sprintf("Error writing slide: %v", err)
		}
		out, _ := comp.Build(deckPath, "")
		return fmt.Sprintf("Updated slide %d and rebuilt to %s", int(idxF), out)

	case "deckforge_reorder_slides":
		orderRaw, ok := args["order"].([]interface{})
		if !ok {
			return "Error: order array required"
		}
		var newOrder []int
		for _, o := range orderRaw {
			if v, ok := o.(float64); ok {
				newOrder = append(newOrder, int(v))
			}
		}
		if err := workspace.ReorderSlides(deckPath, newOrder); err != nil {
			return fmt.Sprintf("Reorder error: %v", err)
		}
		out, _ := comp.Build(deckPath, "")
		return fmt.Sprintf("Reordered deck slides and rebuilt to %s", out)

	default:
		return fmt.Sprintf("Unknown slide tool %s", name)
	}
}
