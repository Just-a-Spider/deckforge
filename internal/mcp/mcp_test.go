package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMCPServerTools(t *testing.T) {
	tempConfigDir, err := os.MkdirTemp("", "deckforge-mcp-test-config-*")
	if err != nil {
		t.Fatalf("Failed to create temp config dir: %v", err)
	}
	defer os.RemoveAll(tempConfigDir)

	os.Setenv("DECKFORGE_CONFIG_DIR", tempConfigDir)
	defer os.Unsetenv("DECKFORGE_CONFIG_DIR")

	tempWorkspace, err := os.MkdirTemp("", "deckforge-mcp-test-ws-*")
	if err != nil {
		t.Fatalf("Failed to create temp ws: %v", err)
	}
	defer os.RemoveAll(tempWorkspace)

	// 1. Test deckforge_get_settings
	getSettingsRes := ExecuteMCPTool(tempWorkspace, "deckforge_get_settings", nil)
	var settingsData map[string]interface{}
	if err := json.Unmarshal([]byte(getSettingsRes), &settingsData); err != nil {
		t.Fatalf("Failed to unmarshal settings response: %v\nOutput: %s", err, getSettingsRes)
	}
	if _, ok := settingsData["config"]; !ok {
		t.Fatalf("Expected config key in settings")
	}

	// 2. Test deckforge_update_settings
	updateRes := ExecuteMCPTool(tempWorkspace, "deckforge_update_settings", map[string]interface{}{
		"key":   "preferredEditor",
		"value": "code",
	})
	if !strings.Contains(updateRes, "Successfully updated") {
		t.Fatalf("Expected success message, got: %s", updateRes)
	}

	// 3. Test deckforge_list_components
	listCompRes := ExecuteMCPTool(tempWorkspace, "deckforge_list_components", nil)
	var compList []map[string]interface{}
	if err := json.Unmarshal([]byte(listCompRes), &compList); err != nil {
		t.Fatalf("Failed to unmarshal component list: %v", err)
	}
	if len(compList) < 7 {
		t.Fatalf("Expected at least 7 components, got %d", len(compList))
	}

	// 4. Test deckforge_create_component
	createCompRes := ExecuteMCPTool(tempWorkspace, "deckforge_create_component", map[string]interface{}{
		"name":     "mcp-stat-box",
		"category": "metrics",
		"styles":   true,
	})
	if !strings.Contains(createCompRes, "df-mcp-stat-box") {
		t.Fatalf("Expected created component to contain df-mcp-stat-box, got: %s", createCompRes)
	}

	// 5. Verify custom component now shows up in deckforge_list_components
	listCompRes2 := ExecuteMCPTool(tempWorkspace, "deckforge_list_components", nil)
	if !strings.Contains(listCompRes2, "df-mcp-stat-box") {
		t.Fatalf("Expected df-mcp-stat-box in components list, got: %s", listCompRes2)
	}

	// 6. Test deckforge_list_themes
	listThemesRes := ExecuteMCPTool(tempWorkspace, "deckforge_list_themes", nil)
	if !strings.Contains(listThemesRes, "academic-crimson") {
		t.Fatalf("Expected academic-crimson in themes list, got: %s", listThemesRes)
	}

	// 7. Test deckforge_create_deck
	deckDir := filepath.Join(tempWorkspace, "my_mcp_deck")
	createDeckRes := ExecuteMCPTool(tempWorkspace, "deckforge_create_deck", map[string]interface{}{
		"name":   "my_mcp_deck",
		"theme":  "academic-crimson",
		"slides": float64(3),
		"path":   deckDir,
	})
	if !strings.Contains(createDeckRes, "Successfully created deck") {
		t.Fatalf("Expected deck creation success, got: %s", createDeckRes)
	}

	// 8. Test deckforge_get_deck_info
	deckInfoRes := ExecuteMCPTool(deckDir, "deckforge_get_deck_info", nil)
	if !strings.Contains(deckInfoRes, "my_mcp_deck") {
		t.Fatalf("Expected deck info to contain my_mcp_deck, got: %s", deckInfoRes)
	}

	// 9. Test deckforge_list_slides
	listSlidesRes := ExecuteMCPTool(deckDir, "deckforge_list_slides", nil)
	if !strings.Contains(listSlidesRes, "01_") {
		t.Fatalf("Expected slides list to contain 01_, got: %s", listSlidesRes)
	}

	// 10. Test deckforge_get_slide & deckforge_update_slide
	getSlideRes := ExecuteMCPTool(deckDir, "deckforge_get_slide", map[string]interface{}{
		"index": float64(1),
	})
	if !strings.Contains(getSlideRes, "<section") {
		t.Fatalf("Expected slide 1 to contain section tag, got: %s", getSlideRes)
	}

	updateSlideRes := ExecuteMCPTool(deckDir, "deckforge_update_slide", map[string]interface{}{
		"index": float64(1),
		"html":  "<section><h1>Updated Slide MCP</h1></section>",
	})
	if !strings.Contains(updateSlideRes, "Updated slide 1") {
		t.Fatalf("Expected slide update success, got: %s", updateSlideRes)
	}

	// 11. Test deckforge_insert_component
	insertRes := ExecuteMCPTool(deckDir, "deckforge_insert_component", map[string]interface{}{
		"index":    float64(1),
		"selector": "df-metric-card",
		"props": map[string]interface{}{
			"title": "ARR Growth",
			"value": "$10M",
		},
	})
	if !strings.Contains(insertRes, "Inserted df-metric-card") {
		t.Fatalf("Expected insert component success, got: %s", insertRes)
	}

	// 12. Test deckforge_audit_tokens
	auditRes := ExecuteMCPTool(deckDir, "deckforge_audit_tokens", nil)
	if !strings.Contains(auditRes, "passes_aa") {
		t.Fatalf("Expected contrast audit report, got: %s", auditRes)
	}

	// 13. Test GetMCPTools registry returns all tools
	tools := GetMCPTools()
	if len(tools) < 17 {
		t.Fatalf("Expected at least 17 tools in registry, got %d", len(tools))
	}
}
