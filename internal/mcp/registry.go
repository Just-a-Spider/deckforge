package mcp

// GetMCPTools returns the full catalog of DeckForge tools available to MCP clients
func GetMCPTools() []map[string]interface{} {
	return []map[string]interface{}{
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
		{
			"name":        "deckforge_list_components",
			"description": "Lists all available presentation components merged across builtin, global, and workspace scopes",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "deckforge_create_component",
			"description": "Scaffolds a new custom component definition in workspace or global scope",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name":        map[string]interface{}{"type": "string", "description": "Component name (e.g. kpi-gauge)"},
					"category":    map[string]interface{}{"type": "string", "description": "Category (metrics, cards, layout, typography, code, media)"},
					"global":      map[string]interface{}{"type": "boolean", "description": "If true, save to global ~/.config/deckforge/components/"},
					"styles":      map[string]interface{}{"type": "boolean", "description": "If true, create companion CSS file for encapsulated styling"},
					"description": map[string]interface{}{"type": "string", "description": "Human-readable description"},
				},
				"required": []string{"name"},
			},
		},
		{
			"name":        "deckforge_get_settings",
			"description": "Returns current global DeckForge configuration, active workspace root, and installed external editors",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "deckforge_update_settings",
			"description": "Updates a global DeckForge configuration property (preferredEditor, defaultTheme, defaultSlideCount, serverPort, autoWatch)",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"key":   map[string]interface{}{"type": "string", "description": "Configuration key name"},
					"value": map[string]interface{}{"type": "string", "description": "New configuration value"},
				},
				"required": []string{"key", "value"},
			},
		},
	}
}
