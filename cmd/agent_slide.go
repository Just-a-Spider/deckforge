package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"deckforge/internal/compiler"
	"deckforge/internal/server"
	"deckforge/internal/theme"
	"deckforge/internal/workspace"
)

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
