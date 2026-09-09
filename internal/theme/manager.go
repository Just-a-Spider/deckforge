package theme

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"deckforge/internal/models"
)

// ThemeManager coordinates discovery, loading, editing, and compiling of segmented themes
type ThemeManager struct {
	WorkspaceRoot string
	GlobalDir     string
}

// NewThemeManager initializes manager with workspace and global paths
func NewThemeManager(workspaceRoot string) *ThemeManager {
	home, _ := os.UserHomeDir()
	globalDir := filepath.Join(home, ".config", "deckforge", "themes")
	return &ThemeManager{
		WorkspaceRoot: workspaceRoot,
		GlobalDir:     globalDir,
	}
}

// ListThemes returns all available themes (workspace, global, and embedded)
func (m *ThemeManager) ListThemes(deckPath string) []models.Theme {
	themeMap := make(map[string]models.Theme)

	// 1. Add all 14 built-in presets first (Scope: "builtin")
	for _, p := range GetBuiltinPresets() {
		p.Scope = "builtin"
		themeMap[p.Tokens.Name] = p
	}

	// 2. Scan global directory (~/.config/deckforge/themes/<name>/) (Scope: "global")
	if m.GlobalDir != "" {
		scanThemesInDir(m.GlobalDir, "global", themeMap)
	}

	// 3. Scan workspace themes (<workspace>/themes/<name>/) (Scope: "workspace")
	if m.WorkspaceRoot != "" {
		scanThemesInDir(filepath.Join(m.WorkspaceRoot, "themes"), "workspace", themeMap)
	}

	var result []models.Theme
	for _, t := range themeMap {
		result = append(result, t)
	}
	return result
}

// SeedTheme exports a built-in preset template to workspace ./themes/ or user config ~/.config/deckforge/themes/
func (m *ThemeManager) SeedTheme(presetName string, global bool) (models.Theme, error) {
	return m.SeedThemeAs(presetName, presetName, global, false)
}

// SeedThemeAs exports a built-in preset with custom name, global flag, and optional SCSS format
func (m *ThemeManager) SeedThemeAs(presetName, targetName string, global bool, scss bool) (models.Theme, error) {
	if targetName == "" {
		targetName = presetName
	}
	var targetDir string
	if global {
		targetDir = filepath.Join(m.GlobalDir, targetName)
	} else {
		targetDir = filepath.Join(m.WorkspaceRoot, "themes", targetName)
	}
	cloned, err := m.ClonePresetSCSS(presetName, targetName, targetDir, scss)
	if err != nil {
		return models.Theme{}, err
	}
	if global {
		cloned.Scope = "global"
	} else {
		cloned.Scope = "workspace"
	}
	return cloned, nil
}

// CreateTheme creates a new theme skeleton or cloned preset in workspace or global directory
func (m *ThemeManager) CreateTheme(name string, basePreset string, global bool, scss bool) (models.Theme, error) {
	if name == "" {
		return models.Theme{}, fmt.Errorf("theme name cannot be empty")
	}

	if basePreset != "" {
		return m.SeedThemeAs(basePreset, name, global, scss)
	}

	var targetDir string
	scope := "workspace"
	if global {
		targetDir = filepath.Join(m.GlobalDir, name)
		scope = "global"
	} else {
		targetDir = filepath.Join(m.WorkspaceRoot, "themes", name)
	}

	theme := models.Theme{
		Tokens: models.ThemeTokens{
			Name:        name,
			DisplayName: strings.Title(strings.ReplaceAll(name, "-", " ")),
			Description: fmt.Sprintf("Custom theme %s", name),
			Vibe:        "Custom modular aesthetic",
			Palette: models.ThemePalette{
				BgCanvas:          "#0f172a",
				BgSurface:         "#1e293b",
				BgSurfaceElevated: "#334155",
				TextPrimary:       "#f8fafc",
				TextSecondary:     "#94a3b8",
				TextMuted:         "#64748b",
				AccentPrimary:     "#3b82f6",
				AccentSecondary:   "#60a5fa",
				AccentGlow:        "rgba(59, 130, 246, 0.2)",
				BorderColor:       "#334155",
				BorderFaint:       "rgba(255, 255, 255, 0.1)",
			},
			Fonts: models.ThemeFonts{
				Display: "'Poppins', system-ui, sans-serif",
				Body:    "'Plus Jakarta Sans', system-ui, sans-serif",
				Mono:    "'JetBrains Mono', monospace",
			},
			Geometry: models.ThemeGeometry{
				RadiusCard:  "12px",
				RadiusSm:    "4px",
				RadiusLg:    "16px",
				BorderWidth: "1.5px",
				BorderStyle: "solid",
			},
			Shadows: models.ThemeShadows{
				Card:     "0 10px 30px rgba(0, 0, 0, 0.3)",
				Elevated: "0 20px 40px rgba(0, 0, 0, 0.4)",
				Glow:     "0 0 20px rgba(59, 130, 246, 0.25)",
			},
			Surfaces: models.ThemeSurfaces{
				BackdropBlur:   "16px",
				SurfaceOpacity: "0.85",
			},
			Spacing: models.ThemeSpacing{
				CardPadding: "28px",
			},
			Motion: models.ThemeMotion{
				SlideTransition: "ease-out",
				AnimationSpeed:  "0.2s",
			},
		},
		SurfacesCSS:   ".glass-panel {\n  background: var(--bg-surface);\n  border: var(--border-width) var(--border-style) var(--border-color);\n  border-radius: var(--radius-card);\n}\n",
		BackdropCSS:   ".stage-backdrop {\n  background-color: var(--bg-canvas);\n}\n",
		ComponentsCSS: ".lime-badge {\n  background: rgba(59, 130, 246, 0.15);\n  color: var(--accent-primary);\n}\n",
		TypographyCSS: ".topbar-title { font-family: var(--font-display); }\n",
		Dir:           targetDir,
		IsCustom:      true,
		Scope:         scope,
	}

	var err error
	if scss {
		err = SaveSegmentedThemeSCSS(targetDir, theme)
	} else {
		err = SaveSegmentedTheme(targetDir, theme)
	}
	if err != nil {
		return models.Theme{}, err
	}
	return theme, nil
}

// ResolveTheme finds a theme by name with hierarchical resolution
func (m *ThemeManager) ResolveTheme(name string, deckPath string) (models.Theme, error) {
	all := m.ListThemes(deckPath)
	for _, t := range all {
		if strings.EqualFold(t.Tokens.Name, name) {
			return t, nil
		}
	}
	// Fallback to academic-crimson preset
	fallback, found := GetPresetByName("academic-crimson")
	if found {
		return fallback, nil
	}
	return models.Theme{}, fmt.Errorf("theme '%s' not found and fallback failed", name)
}

// ClonePreset creates a new segmented theme on disk based on a built-in preset
func (m *ThemeManager) ClonePreset(presetName, newName, targetDir string) (models.Theme, error) {
	return m.ClonePresetSCSS(presetName, newName, targetDir, false)
}

// ClonePresetSCSS creates a new segmented theme on disk based on a preset, using CSS or SCSS
func (m *ThemeManager) ClonePresetSCSS(presetName, newName, targetDir string, scss bool) (models.Theme, error) {
	preset, ok := GetPresetByName(presetName)
	if !ok {
		return models.Theme{}, fmt.Errorf("preset '%s' not found", presetName)
	}

	preset.Tokens.Name = newName
	preset.Tokens.DisplayName = strings.Title(strings.ReplaceAll(newName, "-", " "))
	preset.Tokens.Description = fmt.Sprintf("Customized variant of %s", presetName)
	preset.IsCustom = true
	preset.Dir = targetDir

	var err error
	if scss {
		err = SaveSegmentedThemeSCSS(targetDir, preset)
	} else {
		err = SaveSegmentedTheme(targetDir, preset)
	}
	if err != nil {
		return models.Theme{}, err
	}

	return preset, nil
}
