package theme

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"deckforge/internal/models"
)

// readStyleFile attempts to read a .css or .scss file and transpiles .scss if needed
func readStyleFile(dir, baseName string) string {
	cssPath := filepath.Join(dir, baseName+".css")
	if data, err := os.ReadFile(cssPath); err == nil {
		return string(data)
	}

	scssPath := filepath.Join(dir, baseName+".scss")
	if data, err := os.ReadFile(scssPath); err == nil {
		compiled, _ := CompileSCSS(string(data))
		return compiled
	}

	return ""
}

// LoadSegmentedTheme loads the structured tokens and all modular style segments
func LoadSegmentedTheme(dir string) (models.Theme, error) {
	tokensFile := filepath.Join(dir, "tokens.json")
	data, err := os.ReadFile(tokensFile)
	if err != nil {
		return models.Theme{}, err
	}

	var tokens models.ThemeTokens
	if err := json.Unmarshal(data, &tokens); err != nil {
		return models.Theme{}, err
	}

	// Canonical segments
	typoCSS := readStyleFile(dir, "typography")
	surfacesCSS := readStyleFile(dir, "surfaces")
	backdropCSS := readStyleFile(dir, "backdrop")
	componentsCSS := readStyleFile(dir, "components")
	animationsCSS := readStyleFile(dir, "animations")

	extraCSS := make(map[string]string)
	canonicalNames := map[string]bool{
		"typography": true,
		"surfaces":   true,
		"backdrop":   true,
		"components": true,
		"animations": true,
	}

	// Helper to scan a directory for extra styles
	scanExtra := func(scanDir string, prefix string) {
		entries, err := os.ReadDir(scanDir)
		if err != nil {
			return
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			ext := filepath.Ext(name)
			if ext != ".css" && ext != ".scss" {
				continue
			}
			base := strings.TrimSuffix(name, ext)
			if scanDir == dir && canonicalNames[base] {
				continue
			}
			content, err := os.ReadFile(filepath.Join(scanDir, name))
			if err != nil {
				continue
			}
			str := string(content)
			if ext == ".scss" {
				if compiled, err := CompileSCSS(str); err == nil {
					str = compiled
				}
			}
			key := base
			if prefix != "" {
				key = prefix + "/" + base
			}
			extraCSS[key] = str
		}
	}

	// Scan theme root for other css/scss files
	scanExtra(dir, "")

	// Scan optional styles/ subfolder for modular stylesheets
	stylesSubDir := filepath.Join(dir, "styles")
	if info, err := os.Stat(stylesSubDir); err == nil && info.IsDir() {
		scanExtra(stylesSubDir, "styles")
	}

	return models.Theme{
		Tokens:        tokens,
		TypographyCSS: typoCSS,
		SurfacesCSS:   surfacesCSS,
		BackdropCSS:   backdropCSS,
		ComponentsCSS: componentsCSS,
		AnimationsCSS: animationsCSS,
		ExtraCSS:      extraCSS,
		Dir:           dir,
		IsCustom:      true,
	}, nil
}

// SaveSegmentedTheme persists all structured tokens and modular segments into target directory
func SaveSegmentedTheme(dir string, theme models.Theme) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tokensData, err := json.MarshalIndent(theme.Tokens, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(dir, "tokens.json"), tokensData, 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "typography.css"), []byte(theme.TypographyCSS), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "surfaces.css"), []byte(theme.SurfacesCSS), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "backdrop.css"), []byte(theme.BackdropCSS), 0644); err != nil {
		return err
	}

	if strings.TrimSpace(theme.ComponentsCSS) != "" {
		if err := os.WriteFile(filepath.Join(dir, "components.css"), []byte(theme.ComponentsCSS), 0644); err != nil {
			return err
		}
	}
	if strings.TrimSpace(theme.AnimationsCSS) != "" {
		if err := os.WriteFile(filepath.Join(dir, "animations.css"), []byte(theme.AnimationsCSS), 0644); err != nil {
			return err
		}
	}

	if len(theme.ExtraCSS) > 0 {
		stylesDir := filepath.Join(dir, "styles")
		_ = os.MkdirAll(stylesDir, 0755)
		for k, content := range theme.ExtraCSS {
			cleanName := filepath.Base(k) + ".css"
			_ = os.WriteFile(filepath.Join(stylesDir, cleanName), []byte(content), 0644)
		}
	}

	return nil
}

// SaveSegmentedThemeSCSS persists all structured tokens and modular segments as .scss into target directory
func SaveSegmentedThemeSCSS(dir string, theme models.Theme) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tokensData, err := json.MarshalIndent(theme.Tokens, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(dir, "tokens.json"), tokensData, 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "typography.scss"), []byte(theme.TypographyCSS), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "surfaces.scss"), []byte(theme.SurfacesCSS), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "backdrop.scss"), []byte(theme.BackdropCSS), 0644); err != nil {
		return err
	}

	compContent := theme.ComponentsCSS
	if strings.TrimSpace(compContent) == "" {
		compContent = "// Modular SCSS components overrides\n"
	}
	if err := os.WriteFile(filepath.Join(dir, "components.scss"), []byte(compContent), 0644); err != nil {
		return err
	}

	if strings.TrimSpace(theme.AnimationsCSS) != "" {
		if err := os.WriteFile(filepath.Join(dir, "animations.scss"), []byte(theme.AnimationsCSS), 0644); err != nil {
			return err
		}
	}

	if len(theme.ExtraCSS) > 0 {
		stylesDir := filepath.Join(dir, "styles")
		_ = os.MkdirAll(stylesDir, 0755)
		for k, content := range theme.ExtraCSS {
			cleanName := filepath.Base(k) + ".scss"
			_ = os.WriteFile(filepath.Join(stylesDir, cleanName), []byte(content), 0644)
		}
	}

	return nil
}

// scanThemesInDir inspects a directory for segmented themes or single-file .css definitions
func scanThemesInDir(dir, scope string, themeMap map[string]models.Theme) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			tokensFile := filepath.Join(path, "tokens.json")
			if _, err := os.Stat(tokensFile); err == nil {
				theme, err := LoadSegmentedTheme(path)
				if err == nil {
					theme.IsCustom = true
					theme.Scope = scope
					themeMap[theme.Tokens.Name] = theme
				}
			}
		} else if strings.HasSuffix(entry.Name(), ".css") {
			name := strings.TrimSuffix(entry.Name(), ".css")
			if _, exists := themeMap[name]; !exists {
				cssContent, err := os.ReadFile(path)
				if err == nil {
					themeMap[name] = models.Theme{
						Tokens: models.ThemeTokens{
							Name:        name,
							DisplayName: strings.Title(strings.ReplaceAll(name, "-", " ")),
							Palette: models.ThemePalette{
								BgCanvas:      "#0f172a",
								TextPrimary:   "#ffffff",
								AccentPrimary: "#3b82f6",
							},
							Fonts: models.ThemeFonts{
								Display: "sans-serif",
								Body:    "sans-serif",
								Mono:    "monospace",
							},
						},
						SurfacesCSS: string(cssContent),
						Dir:         dir,
						IsCustom:    true,
						Scope:       scope,
					}
				}
			}
		}
	}
}
