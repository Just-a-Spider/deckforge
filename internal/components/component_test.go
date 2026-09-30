package components

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestComponentManagerBuiltin(t *testing.T) {
	cm := NewComponentManager("")
	comps := cm.ListComponents("")

	if len(comps) < 7 {
		t.Fatalf("Expected at least 7 built-in components, got %d", len(comps))
	}

	foundMetric := false
	for _, c := range comps {
		if c.Selector == "df-metric-card" {
			foundMetric = true
			if c.Scope != "builtin" {
				t.Fatalf("Expected builtin scope, got %s", c.Scope)
			}
		}
	}
	if !foundMetric {
		t.Fatalf("df-metric-card not found in catalog")
	}
}

func TestComponentScaffoldSingleFileAndDiscovery(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "deckforge-comp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	opts := ScaffoldOptions{
		Name:          "stat-gauge",
		Category:      "metrics",
		Description:   "A custom test gauge",
		WithCSS:       false,
		WorkspaceRoot: tempDir,
	}

	scaffolded, err := ScaffoldComponent(opts)
	if err != nil {
		t.Fatalf("Scaffold error: %v", err)
	}
	if scaffolded.Selector != "df-stat-gauge" {
		t.Fatalf("Expected selector df-stat-gauge, got %s", scaffolded.Selector)
	}

	cm := NewComponentManager(tempDir)
	comps := cm.ListComponents("")

	found := false
	for _, c := range comps {
		if c.Selector == "df-stat-gauge" {
			found = true
			if c.Scope != "workspace" {
				t.Fatalf("Expected scope workspace, got %s", c.Scope)
			}
			if c.Category != "metrics" {
				t.Fatalf("Expected category metrics, got %s", c.Category)
			}
		}
	}
	if !found {
		t.Fatalf("Custom component df-stat-gauge not discovered by ComponentManager")
	}

	css := cm.CompileComponentCSS("")
	if !strings.Contains(css, "df-stat-gauge") {
		t.Fatalf("Expected df-stat-gauge in compiled CSS, got:\n%s", css)
	}
}

func TestComponentScaffoldBundleWithCSS(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "deckforge-comp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	opts := ScaffoldOptions{
		Name:          "hero-badge-pill",
		Category:      "cards",
		WithCSS:       true,
		WorkspaceRoot: tempDir,
	}

	scaffolded, err := ScaffoldComponent(opts)
	if err != nil {
		t.Fatalf("Scaffold error: %v", err)
	}
	if scaffolded.Selector != "df-hero-badge-pill" {
		t.Fatalf("Expected selector df-hero-badge-pill, got %s", scaffolded.Selector)
	}

	// Verify bundle files exist on disk in .deckforge/components
	bundleDir := filepath.Join(tempDir, ".deckforge", "components", "df-hero-badge-pill")
	if _, err := os.Stat(filepath.Join(bundleDir, "component.json")); err != nil {
		t.Fatalf("component.json missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(bundleDir, "component.css")); err != nil {
		t.Fatalf("component.css missing: %v", err)
	}

	cm := NewComponentManager(tempDir)
	retrieved, err := cm.GetComponent("df-hero-badge-pill", "")
	if err != nil {
		t.Fatalf("GetComponent error: %v", err)
	}
	if retrieved.Scope != "workspace" {
		t.Fatalf("Expected scope workspace, got %s", retrieved.Scope)
	}
	if !strings.Contains(retrieved.Styles, "df-hero-badge-pill") {
		t.Fatalf("Expected Styles to contain df-hero-badge-pill")
	}
}

func TestComponentPrecedence(t *testing.T) {
	tempWorkspace, err := os.MkdirTemp("", "deckforge-comp-ws-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempWorkspace)

	tempGlobal, err := os.MkdirTemp("", "deckforge-comp-global-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempGlobal)

	// Create custom df-metric-card in workspace (overriding builtin)
	opts := ScaffoldOptions{
		Name:          "metric-card",
		Category:      "metrics",
		Description:   "Workspace override for metric card",
		WithCSS:       false,
		WorkspaceRoot: tempWorkspace,
	}
	_, err = ScaffoldComponent(opts)
	if err != nil {
		t.Fatalf("Scaffold error: %v", err)
	}

	cm := &ComponentManager{
		WorkspaceRoot: tempWorkspace,
		GlobalDir:     tempGlobal,
	}

	comp, err := cm.GetComponent("df-metric-card", "")
	if err != nil {
		t.Fatalf("Failed to get component: %v", err)
	}
	if comp.Scope != "workspace" {
		t.Fatalf("Expected workspace override for df-metric-card, got scope: %s", comp.Scope)
	}
	if comp.Description != "Workspace override for metric card" {
		t.Fatalf("Expected workspace description, got %s", comp.Description)
	}
}

func TestComponentLegacyWorkspaceDiscovery(t *testing.T) {
	tempWorkspace, err := os.MkdirTemp("", "deckforge-legacy-comp-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempWorkspace)

	legacyDir := filepath.Join(tempWorkspace, "components")
	if err := os.MkdirAll(legacyDir, 0755); err != nil {
		t.Fatalf("Failed to create legacy components dir: %v", err)
	}

	legacyJSON := `{"selector": "df-legacy-badge", "name": "Legacy Badge", "category": "badges"}`
	if err := os.WriteFile(filepath.Join(legacyDir, "df-legacy-badge.json"), []byte(legacyJSON), 0644); err != nil {
		t.Fatalf("Failed to write legacy component: %v", err)
	}

	cm := NewComponentManager(tempWorkspace)
	comp, err := cm.GetComponent("df-legacy-badge", "")
	if err != nil {
		t.Fatalf("Failed to discover legacy component: %v", err)
	}
	if comp.Scope != "workspace" {
		t.Fatalf("Expected workspace scope, got %s", comp.Scope)
	}
}

