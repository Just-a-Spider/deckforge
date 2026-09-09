package theme

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFullSpectrumThemeTokens(t *testing.T) {
	preset, ok := GetPresetByName("swiss-minimal")
	if !ok {
		t.Fatalf("expected swiss-minimal preset to be found")
	}

	css := preset.CompileFullCSS()

	// Verify newly introduced full-spectrum token variables in :root
	expectedTokens := []string{
		"--radius-card: 0px;",
		"--radius-sm: 0px;",
		"--border-width: 2px;",
		"--shadow-card: 6px 6px 0 #000000;",
		"--backdrop-blur: 0px;",
		"--card-padding: 24px;",
		"--slide-transition: none;",
		"--animation-speed: 0.15s;",
	}

	for _, token := range expectedTokens {
		if !strings.Contains(css, token) {
			t.Errorf("expected compiled CSS to contain token %q, but was missing", token)
		}
	}
}

func TestModularSegmentedThemeLoading(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "df-modular-theme-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	themeDir := filepath.Join(tmpDir, "custom-modular")
	_ = os.MkdirAll(filepath.Join(themeDir, "styles"), 0755)

	tokensJSON := `{
  "name": "custom-modular",
  "displayName": "Custom Modular",
  "palette": {
    "bgCanvas": "#111111",
    "bgSurface": "#222222",
    "textPrimary": "#ffffff",
    "textSecondary": "#cccccc",
    "accentPrimary": "#00ff88",
    "borderColor": "#00ff88"
  },
  "fonts": {
    "display": "sans-serif",
    "body": "sans-serif",
    "mono": "monospace"
  },
  "geometry": {
    "radiusCard": "16px"
  }
}`
	_ = os.WriteFile(filepath.Join(themeDir, "tokens.json"), []byte(tokensJSON), 0644)
	_ = os.WriteFile(filepath.Join(themeDir, "typography.css"), []byte(".custom-heading { font-size: 40px; }"), 0644)
	_ = os.WriteFile(filepath.Join(themeDir, "components.css"), []byte(".df-metric-card { border-style: dashed; }"), 0644)
	_ = os.WriteFile(filepath.Join(themeDir, "animations.css"), []byte("@keyframes customPulse { 0% { opacity: 0; } }"), 0644)
	_ = os.WriteFile(filepath.Join(themeDir, "styles", "cards.css"), []byte(".special-card { transform: scale(1.02); }"), 0644)

	loaded, err := LoadSegmentedTheme(themeDir)
	if err != nil {
		t.Fatalf("LoadSegmentedTheme failed: %v", err)
	}

	if loaded.Tokens.Geometry.RadiusCard != "16px" {
		t.Errorf("expected RadiusCard to be 16px, got %s", loaded.Tokens.Geometry.RadiusCard)
	}
	if !strings.Contains(loaded.ComponentsCSS, ".df-metric-card") {
		t.Errorf("expected ComponentsCSS to be loaded")
	}
	if !strings.Contains(loaded.AnimationsCSS, "customPulse") {
		t.Errorf("expected AnimationsCSS to be loaded")
	}
	if !strings.Contains(loaded.ExtraCSS["styles/cards"], ".special-card") {
		t.Errorf("expected ExtraCSS to contain styles/cards")
	}

	compiled := loaded.CompileFullCSS()
	if !strings.Contains(compiled, "/* Components Segment */") {
		t.Errorf("expected compiled CSS to contain Components Segment")
	}
	if !strings.Contains(compiled, "/* Animations Segment */") {
		t.Errorf("expected compiled CSS to contain Animations Segment")
	}
	if !strings.Contains(compiled, "/* Custom Segment: styles/cards */") {
		t.Errorf("expected compiled CSS to contain Custom Segment for styles/cards")
	}
}

func TestSCSSCompilationOrFallback(t *testing.T) {
	rawSCSS := ".card { color: red; & .child { color: blue; } }"
	out, err := CompileSCSS(rawSCSS)
	if out == "" {
		t.Fatalf("expected non-empty output from CompileSCSS")
	}
	// Even if fallback, output should contain valid syntax
	if !strings.Contains(out, ".card") {
		t.Errorf("expected output to contain .card")
	}
	_ = err // fallback error is expected if sass isn't on host
}

func TestThemeCreationTooling(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "df-theme-create-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tm := &ThemeManager{
		WorkspaceRoot: tmpDir,
		GlobalDir:     filepath.Join(tmpDir, "global"),
	}

	// 1. Create workspace theme with SCSS
	created, err := tm.CreateTheme("my-custom", "cyber-dark", false, true)
	if err != nil {
		t.Fatalf("CreateTheme failed: %v", err)
	}
	if created.Scope != "workspace" {
		t.Errorf("expected scope workspace, got %s", created.Scope)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "themes", "my-custom", "surfaces.scss")); err != nil {
		t.Errorf("expected surfaces.scss to exist: %v", err)
	}

	// 2. Create global theme with base preset
	globalTheme, err := tm.CreateTheme("global-frost", "academic-crimson", true, true)
	if err != nil {
		t.Fatalf("CreateTheme global failed: %v", err)
	}
	if globalTheme.Scope != "global" {
		t.Errorf("expected scope global, got %s", globalTheme.Scope)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "global", "global-frost", "tokens.json")); err != nil {
		t.Errorf("expected tokens.json to exist in global dir: %v", err)
	}
}
