package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"deckforge/internal/compiler"
	"deckforge/internal/models"
	"deckforge/internal/scaffold"
	"deckforge/internal/theme"
	"deckforge/internal/workspace"
)

type jsonRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type jsonRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// RunMCPServer runs the MCP JSON-RPC 2.0 stdio server for deckPath
func RunMCPServer(deckPath string) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		line = []byte(strings.TrimSpace(string(line)))
		if len(line) == 0 {
			continue
		}

		var req jsonRPCRequest
		if err := json.Unmarshal(line, &req); err != nil {
			continue
		}

		resp := handleMCPRequest(deckPath, req)
		if resp != nil {
			data, _ := json.Marshal(resp)
			fmt.Printf("%s\n", string(data))
		}
	}
}

func handleMCPRequest(deckPath string, req jsonRPCRequest) *jsonRPCResponse {
	switch req.Method {
	case "initialize":
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"serverInfo": map[string]string{
					"name":    "deckforge-mcp",
					"version": "1.2.0",
				},
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
			},
		}

	case "notifications/initialized":
		return nil

	case "tools/list":
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"tools": []map[string]interface{}{
					{
						"name":        "deckforge_get_catalog",
						"description": "Returns complete catalog of Angular-inspired component definitions with input schemas, slots, and HTML templates",
						"inputSchema": map[string]interface{}{
							"type":       "object",
							"properties": map[string]interface{}{},
						},
					},
					{
						"name":        "deckforge_get_styling_guide",
						"description": "Returns comprehensive styling rules, semantic CSS class dictionary, spacing scale, and token variables",
						"inputSchema": map[string]interface{}{
							"type":       "object",
							"properties": map[string]interface{}{},
						},
					},
					{
						"name":        "deckforge_create_deck",
						"description": "Scaffolds a new presentation deck with chosen theme and slide count",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"name":   map[string]interface{}{"type": "string", "description": "Deck name (e.g. cloud_architecture)"},
								"theme":  map[string]interface{}{"type": "string", "description": "Preset theme name (default: cyber-dark)"},
								"slides": map[string]interface{}{"type": "integer", "description": "Number of initial slides (default: 5)"},
								"path":   map[string]interface{}{"type": "string", "description": "Target directory path"},
							},
							"required": []string{"name"},
						},
					},
					{
						"name":        "deckforge_list_themes",
						"description": "Lists all available themes with displayName, vibe, and custom status",
						"inputSchema": map[string]interface{}{
							"type":       "object",
							"properties": map[string]interface{}{},
						},
					},
					{
						"name":        "deckforge_set_theme",
						"description": "Updates the deck's theme in deck.json and recompiles the presentation",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"theme": map[string]interface{}{"type": "string", "description": "Theme name to apply"},
							},
							"required": []string{"theme"},
						},
					},
					{
						"name":        "deckforge_create_theme",
						"description": "Scaffolds a new theme skeleton or clones a preset in workspace or global scope with CSS or SCSS modular files",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"name":   map[string]interface{}{"type": "string", "description": "Theme name slug (e.g. frosted-crimson)"},
								"base":   map[string]interface{}{"type": "string", "description": "Optional base preset to clone tokens and styles from (e.g. academic-crimson)"},
								"global": map[string]interface{}{"type": "boolean", "description": "If true, installs into ~/.config/deckforge/themes/ for global availability across all decks"},
								"scss":   map[string]interface{}{"type": "boolean", "description": "If true, scaffolds modular .scss stylesheets instead of .css"},
							},
							"required": []string{"name"},
						},
					},
					{
						"name":        "deckforge_get_deck_info",
						"description": "Returns deck metadata including title, subtitle, theme, slide count, and slide list",
						"inputSchema": map[string]interface{}{
							"type":       "object",
							"properties": map[string]interface{}{},
						},
					},
					{
						"name":        "deckforge_set_deck_meta",
						"description": "Updates presentation metadata (title, subtitle, theme) in deck.json and recompiles presentation",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"title":    map[string]interface{}{"type": "string", "description": "Presentation title"},
								"subtitle": map[string]interface{}{"type": "string", "description": "Presentation subtitle"},
								"theme":    map[string]interface{}{"type": "string", "description": "Theme name to apply"},
							},
						},
					},
					{
						"name":        "deckforge_list_slides",
						"description": "Lists all sequential slides in active presentation deck",
						"inputSchema": map[string]interface{}{
							"type":       "object",
							"properties": map[string]interface{}{},
						},
					},
					{
						"name":        "deckforge_get_slide",
						"description": "Returns raw HTML and metadata for slide at index",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"index": map[string]interface{}{"type": "integer", "description": "1-based slide index"},
							},
							"required": []string{"index"},
						},
					},
					{
						"name":        "deckforge_update_slide",
						"description": "Overwrites and cleans slide HTML, auto-recompiling presentation",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"index": map[string]interface{}{"type": "integer", "description": "1-based slide index"},
								"html":  map[string]interface{}{"type": "string", "description": "Clean slide HTML"},
							},
							"required": []string{"index", "html"},
						},
					},
					{
						"name":        "deckforge_reorder_slides",
						"description": "Reorders slides on disk according to 1-based index array",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"order": map[string]interface{}{"type": "array", "items": map[string]string{"type": "integer"}},
							},
							"required": []string{"order"},
						},
					},
					{
						"name":        "deckforge_audit_tokens",
						"description": "Audits active presentation theme for WCAG 2.1 contrast compliance",
						"inputSchema": map[string]interface{}{
							"type":       "object",
							"properties": map[string]interface{}{},
						},
					},
					{
						"name":        "deckforge_apply_token_fix",
						"description": "Auto-fixes failing contrast tokens in theme tokens.json",
						"inputSchema": map[string]interface{}{
							"type":       "object",
							"properties": map[string]interface{}{},
						},
					},
					{
						"name":        "deckforge_insert_component",
						"description": "Renders and inserts a modular component into slide",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"index":    map[string]interface{}{"type": "integer", "description": "1-based slide index"},
								"selector": map[string]interface{}{"type": "string", "description": "Component selector (e.g. df-metric-card)"},
								"props":    map[string]interface{}{"type": "object", "description": "Input properties"},
							},
							"required": []string{"index", "selector"},
						},
					},
				},
			},
		}

	case "tools/call":
		var callParams struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments"`
		}
		_ = json.Unmarshal(req.Params, &callParams)

		toolResult := executeMCPTool(deckPath, callParams.Name, callParams.Arguments)
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"content": []map[string]string{
					{
						"type": "text",
						"text": toolResult,
					},
				},
			},
		}

	default:
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: map[string]interface{}{
				"code":    -32601,
				"message": "Method not found",
			},
		}
	}
}

func executeMCPTool(deckPath, name string, args map[string]interface{}) string {
	tm := theme.NewThemeManager(deckPath)
	comp := compiler.NewCompiler(tm)

	if name == "deckforge_create_deck" {
		nameStr, ok := args["name"].(string)
		if !ok || nameStr == "" {
			return "Error: deck name is required"
		}
		themeStr, _ := args["theme"].(string)
		if themeStr == "" {
			themeStr = "cyber-dark"
		}
		slidesCount := 5
		if sF, ok := args["slides"].(float64); ok && int(sF) > 0 {
			slidesCount = int(sF)
		}
		targetPath, _ := args["path"].(string)
		if targetPath == "" {
			targetPath = filepath.Join("decks", nameStr)
		}
		deck, err := scaffold.ScaffoldDeck(targetPath, nameStr, themeStr, slidesCount)
		if err != nil {
			return fmt.Sprintf("Error creating deck: %v", err)
		}
		newTM := theme.NewThemeManager(targetPath)
		newComp := compiler.NewCompiler(newTM)
		outFile, _ := newComp.Build(targetPath, themeStr)
		return fmt.Sprintf("Successfully created deck '%s' with %d slides at %s (Compiled: %s)", deck.Name, len(deck.Slides), deck.Path, outFile)
	}

	if name == "deckforge_list_themes" {
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
	}

	if name == "deckforge_get_catalog" {
		components := models.GetBuiltinComponents()
		data, _ := json.MarshalIndent(components, "", "  ")
		return string(data)
	}

	if name == "deckforge_get_styling_guide" {
		guide := models.GetStylingGuide()
		data, _ := json.MarshalIndent(guide, "", "  ")
		return string(data)
	}

	deck, err := workspace.InspectDeck(deckPath)
	if err != nil {
		return fmt.Sprintf("Error inspecting deck: %v", err)
	}

	switch name {
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

	case "deckforge_get_deck_info":
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
		cleaned := SanitizeSlideHTML(htmlStr)
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

	case "deckforge_insert_component":
		idxF, ok := args["index"].(float64)
		selector, ok2 := args["selector"].(string)
		if !ok || !ok2 || int(idxF) < 1 || int(idxF) > len(deck.Slides) {
			return "Error: invalid index or selector"
		}
		props, _ := args["props"].(map[string]interface{})
		if props == nil {
			props = make(map[string]interface{})
		}
		var targetComp *models.ComponentDefinition
		for _, c := range models.GetBuiltinComponents() {
			if c.Selector == selector {
				targetComp = &c
				break
			}
		}
		if targetComp == nil {
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
		out, _ := comp.Build(deckPath, "")
		return fmt.Sprintf("Inserted %s into slide %d and rebuilt to %s", selector, int(idxF), out)

	default:
		return fmt.Sprintf("Unknown tool %s", name)
	}
}
