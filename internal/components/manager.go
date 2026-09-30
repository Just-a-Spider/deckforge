package components

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"deckforge/internal/models"
)

// ComponentManager coordinates discovery, loading, and compilation of multi-tier components
type ComponentManager struct {
	WorkspaceRoot string
	GlobalDir     string
}

// NewComponentManager initializes a component manager for a workspace
func NewComponentManager(workspaceRoot string) *ComponentManager {
	return &ComponentManager{
		WorkspaceRoot: workspaceRoot,
		GlobalDir:     models.GlobalComponentsDir(),
	}
}

// ListComponents returns all available components merged across tiers
func (m *ComponentManager) ListComponents(deckPath string) []models.ComponentDefinition {
	compMap := make(map[string]models.ComponentDefinition)

	// 1. Built-in presets (tier 1)
	for _, c := range models.GetBuiltinComponents() {
		c.Scope = "builtin"
		compMap[c.Selector] = c
	}

	// 2. Global scope: ~/.config/deckforge/components/ (tier 2)
	if m.GlobalDir != "" {
		scanComponentsDir(m.GlobalDir, "global", compMap)
	}

	// 3. Workspace root: legacy <workspace>/components/ and preferred <workspace>/.deckforge/components/ (tier 3)
	if m.WorkspaceRoot != "" {
		scanComponentsDir(filepath.Join(m.WorkspaceRoot, "components"), "workspace", compMap)
		scanComponentsDir(models.WorkspaceComponentsDir(m.WorkspaceRoot), "workspace", compMap)
	}

	// 4. Scan deck local components and walk up from deckPath to discover enclosing components
	if deckPath != "" {
		absDeck, err := filepath.Abs(deckPath)
		if err == nil {
			var ancestors []string
			curr := absDeck
			for {
				ancestors = append([]string{curr}, ancestors...)
				parent := filepath.Dir(curr)
				if parent == curr || parent == "." || parent == "/" {
					break
				}
				curr = parent
			}
			for _, p := range ancestors {
				scanComponentsDir(filepath.Join(p, "components"), "workspace", compMap)
				scanComponentsDir(models.WorkspaceComponentsDir(p), "workspace", compMap)
			}
		} else {
			scanComponentsDir(filepath.Join(deckPath, "components"), "workspace", compMap)
			scanComponentsDir(models.WorkspaceComponentsDir(deckPath), "workspace", compMap)
		}
	}

	// 5. Scan current working directory if not already covered
	if cwd, err := os.Getwd(); err == nil {
		scanComponentsDir(filepath.Join(cwd, "components"), "workspace", compMap)
		scanComponentsDir(models.WorkspaceComponentsDir(cwd), "workspace", compMap)
	}

	var list []models.ComponentDefinition
	for _, c := range compMap {
		list = append(list, c)
	}

	// Sort deterministically: built-in first, then by selector
	sort.Slice(list, func(i, j int) bool {
		if list[i].Scope != list[j].Scope {
			scopeWeight := func(s string) int {
				switch s {
				case "workspace":
					return 1
				case "global":
					return 2
				default:
					return 3
				}
			}
			return scopeWeight(list[i].Scope) < scopeWeight(list[j].Scope)
		}
		return list[i].Selector < list[j].Selector
	})

	return list
}

// GetComponent finds a component by its custom tag selector (e.g. df-metric-card)
func (m *ComponentManager) GetComponent(selector string, deckPath string) (*models.ComponentDefinition, error) {
	for _, c := range m.ListComponents(deckPath) {
		if c.Selector == selector {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("component '%s' not found", selector)
}

// CompileComponentCSS aggregates all custom styles declared by components
func (m *ComponentManager) CompileComponentCSS(deckPath string) string {
	components := m.ListComponents(deckPath)
	var sb strings.Builder

	for _, c := range components {
		cleanCSS := strings.TrimSpace(c.Styles)
		if cleanCSS != "" {
			sb.WriteString(fmt.Sprintf("/* --- Custom Component: %s (%s) --- */\n", c.Selector, c.Scope))
			sb.WriteString(cleanCSS)
			sb.WriteString("\n\n")
		}
	}

	return strings.TrimSpace(sb.String())
}

func scanComponentsDir(dir string, scope string, compMap map[string]models.ComponentDefinition) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			// Directory bundle: components/<selector>/
			loadComponentBundle(path, scope, compMap)
		} else if strings.HasSuffix(entry.Name(), ".json") {
			// Single file: components/<selector>.json
			loadComponentFile(path, scope, compMap)
		}
	}
}

func loadComponentFile(filePath string, scope string, compMap map[string]models.ComponentDefinition) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	var comp models.ComponentDefinition
	if err := json.Unmarshal(data, &comp); err != nil {
		return
	}

	if comp.Selector == "" {
		base := strings.TrimSuffix(filepath.Base(filePath), ".json")
		if !strings.HasPrefix(base, "df-") {
			base = "df-" + base
		}
		comp.Selector = base
	}
	if comp.Name == "" {
		comp.Name = comp.Selector
	}
	comp.Scope = scope
	comp.SourcePath = filePath

	compMap[comp.Selector] = comp
}

func loadComponentBundle(bundleDir string, scope string, compMap map[string]models.ComponentDefinition) {
	// Look for component.json or <bundle-name>.json
	jsonCandidates := []string{
		filepath.Join(bundleDir, "component.json"),
		filepath.Join(bundleDir, filepath.Base(bundleDir)+".json"),
	}

	var jsonFile string
	for _, cand := range jsonCandidates {
		if fi, err := os.Stat(cand); err == nil && !fi.IsDir() {
			jsonFile = cand
			break
		}
	}

	if jsonFile == "" {
		return
	}

	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return
	}

	var comp models.ComponentDefinition
	if err := json.Unmarshal(data, &comp); err != nil {
		return
	}

	if comp.Selector == "" {
		base := filepath.Base(bundleDir)
		if !strings.HasPrefix(base, "df-") {
			base = "df-" + base
		}
		comp.Selector = base
	}
	if comp.Name == "" {
		comp.Name = comp.Selector
	}

	// Look for companion CSS: component.css or styles.css
	cssCandidates := []string{
		filepath.Join(bundleDir, "component.css"),
		filepath.Join(bundleDir, "styles.css"),
		filepath.Join(bundleDir, filepath.Base(bundleDir)+".css"),
	}

	for _, cssCand := range cssCandidates {
		if cssData, err := os.ReadFile(cssCand); err == nil {
			if comp.Styles != "" {
				comp.Styles = comp.Styles + "\n" + string(cssData)
			} else {
				comp.Styles = string(cssData)
			}
			break
		}
	}

	comp.Scope = scope
	comp.SourcePath = bundleDir
	compMap[comp.Selector] = comp
}
