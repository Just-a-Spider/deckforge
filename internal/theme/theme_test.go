package theme

import (
	"os"
	"path/filepath"
	"testing"
)

func TestThemeScopingAndSeeding(t *testing.T) {
	tmpDir := t.TempDir()
	workspaceRoot := filepath.Join(tmpDir, "workspace")
	globalDir := filepath.Join(tmpDir, "global_config")
	_ = os.MkdirAll(workspaceRoot, 0755)
	_ = os.MkdirAll(globalDir, 0755)

	tm := &ThemeManager{
		WorkspaceRoot: workspaceRoot,
		GlobalDir:     globalDir,
	}

	// 1. Initially should have built-in presets (all with scope: builtin)
	themes := tm.ListThemes("")
	if len(themes) < 14 {
		t.Fatalf("Expected at least 14 presets, got %d", len(themes))
	}
	for _, th := range themes {
		if th.Scope != "builtin" {
			t.Fatalf("Expected initial theme scope to be 'builtin', got '%s'", th.Scope)
		}
	}

	// 2. Seed a theme to workspace
	seededWork, err := tm.SeedTheme("cyber-dark", false)
	if err != nil {
		t.Fatalf("Failed to seed theme to workspace: %v", err)
	}
	if seededWork.Scope != "workspace" {
		t.Fatalf("Expected scope 'workspace', got '%s'", seededWork.Scope)
	}
	if _, err := os.Stat(filepath.Join(seededWork.Dir, "tokens.json")); err != nil {
		t.Fatalf("Expected tokens.json in workspace themes dir: %v", err)
	}

	// 3. Seed a theme to global
	seededGlob, err := tm.SeedTheme("bold-signal", true)
	if err != nil {
		t.Fatalf("Failed to seed theme to global: %v", err)
	}
	if seededGlob.Scope != "global" {
		t.Fatalf("Expected scope 'global', got '%s'", seededGlob.Scope)
	}
	if _, err := os.Stat(filepath.Join(globalDir, "bold-signal", "tokens.json")); err != nil {
		t.Fatalf("Expected tokens.json in global themes dir: %v", err)
	}

	// 4. Verify ListThemes reflects updated scopes
	updatedThemes := tm.ListThemes("")
	foundWorkspace := false
	foundGlobal := false
	for _, th := range updatedThemes {
		if th.Tokens.Name == "cyber-dark" && th.Scope == "workspace" {
			foundWorkspace = true
		}
		if th.Tokens.Name == "bold-signal" && th.Scope == "global" {
			foundGlobal = true
		}
	}
	if !foundWorkspace {
		t.Fatal("Expected cyber-dark to be resolved with workspace scope")
	}
	if !foundGlobal {
		t.Fatal("Expected bold-signal to be resolved with global scope")
	}
}
