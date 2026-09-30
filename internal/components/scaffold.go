package components

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"deckforge/internal/models"
)

var slugSanitizer = regexp.MustCompile(`[^a-z0-9-]+`)

// ScaffoldOptions specifies parameters for generating a new component skeleton
type ScaffoldOptions struct {
	Name          string `json:"name"`
	Category      string `json:"category"`
	Description   string `json:"description"`
	Global        bool   `json:"global"`
	WithCSS       bool   `json:"withCss"`
	TargetDir     string `json:"targetDir,omitempty"`
	WorkspaceRoot string `json:"workspaceRoot,omitempty"`
}

// ScaffoldComponent creates a new Angular-inspired component schema on disk
func ScaffoldComponent(opts ScaffoldOptions) (*models.ComponentDefinition, error) {
	if strings.TrimSpace(opts.Name) == "" {
		return nil, fmt.Errorf("component name cannot be empty")
	}

	cleanSlug := strings.ToLower(strings.TrimSpace(opts.Name))
	cleanSlug = strings.ReplaceAll(cleanSlug, "_", "-")
	cleanSlug = slugSanitizer.ReplaceAllString(cleanSlug, "")
	cleanSlug = strings.Trim(cleanSlug, "-")

	var selector string
	if strings.HasPrefix(cleanSlug, "df-") {
		selector = cleanSlug
	} else {
		selector = "df-" + cleanSlug
	}

	rawName := strings.TrimPrefix(selector, "df-")
	words := strings.Split(rawName, "-")
	var titleParts []string
	for _, w := range words {
		if len(w) > 0 {
			titleParts = append(titleParts, strings.ToUpper(w[:1])+w[1:])
		}
	}
	humanName := strings.Join(titleParts, " ") + " Component"

	category := strings.ToLower(strings.TrimSpace(opts.Category))
	if category == "" {
		category = "cards"
	}

	desc := opts.Description
	if desc == "" {
		desc = fmt.Sprintf("Custom %s component with responsive card styling and token awareness", humanName)
	}

	// Resolve destination directory
	var destDir string
	scope := "workspace"
	if opts.TargetDir != "" {
		destDir = opts.TargetDir
	} else if opts.Global {
		destDir = models.GlobalComponentsDir()
		scope = "global"
	} else {
		root := opts.WorkspaceRoot
		if root == "" {
			root = "."
		}
		destDir = models.WorkspaceComponentsDir(root)
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create components directory: %w", err)
	}

	comp := &models.ComponentDefinition{
		Selector:    selector,
		Name:        humanName,
		Category:    category,
		Description: desc,
		Inputs: []models.ComponentInput{
			{Name: "title", Type: "string", Default: humanName, Description: "Component headline text"},
			{Name: "value", Type: "string", Default: "100%", Description: "Primary metric or content value"},
			{Name: "badge", Type: "string", Default: "ACTIVE", Description: "Optional header status badge"},
			{Name: "accentEdge", Type: "boolean", Default: true, Description: "Accent line indicator"},
		},
		Slots:   []string{"notes"},
		Classes: []string{"glass-panel"},
		Scope:   scope,
		TemplateSnippet: `<div class="glass-panel{{if .accentEdge}} lime-edge{{end}}" style="padding: 24px;">
    {{if .badge}}<div class="card-tag">{{.badge}}</div>{{end}}
    <div class="panel-headline" style="margin-bottom: 8px;">{{.title}}</div>
    <div style="font-family: var(--font-display); font-size: 38px; font-weight: 800; color: var(--accent-primary); line-height: 1.1; margin-bottom: 8px;">{{.value}}</div>
    {{if .slot_notes}}<div style="font-size: 13px; color: var(--text-secondary); margin-top: 8px;">{{.slot_notes}}</div>{{end}}
</div>`,
	}

	sampleCSS := fmt.Sprintf(`/* Encapsulated Styles for %s */
.%s {
    border-radius: var(--radius-card, 10px);
    transition: transform var(--animation-speed, 0.2s) ease;
}
`, selector, selector)

	if opts.WithCSS {
		// Create bundle folder: <destDir>/<selector>/
		bundlePath := filepath.Join(destDir, selector)
		if err := os.MkdirAll(bundlePath, 0755); err != nil {
			return nil, fmt.Errorf("failed to create bundle dir: %w", err)
		}

		comp.Styles = sampleCSS
		comp.SourcePath = bundlePath

		jsonData, err := json.MarshalIndent(comp, "", "  ")
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(filepath.Join(bundlePath, "component.json"), jsonData, 0644); err != nil {
			return nil, err
		}
		if err := os.WriteFile(filepath.Join(bundlePath, "component.css"), []byte(sampleCSS), 0644); err != nil {
			return nil, err
		}
	} else {
		// Single JSON file: <destDir>/<selector>.json
		comp.Styles = sampleCSS
		filePath := filepath.Join(destDir, selector+".json")
		comp.SourcePath = filePath

		jsonData, err := json.MarshalIndent(comp, "", "  ")
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
			return nil, err
		}
	}

	return comp, nil
}
