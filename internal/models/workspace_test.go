package models

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	tempConfigDir, err := os.MkdirTemp("", "deckforge-models-test-config-*")
	if err == nil {
		os.Setenv("DECKFORGE_CONFIG_DIR", tempConfigDir)
		defer os.RemoveAll(tempConfigDir)
	}
	os.Exit(m.Run())
}

func TestWorkspaceConfigEnvOverride(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "deckforge-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	os.Setenv("DECKFORGE_CONFIG_DIR", tempDir)
	defer os.Unsetenv("DECKFORGE_CONFIG_DIR")

	if ConfigDirPath() != tempDir {
		t.Fatalf("Expected ConfigDirPath %s, got %s", tempDir, ConfigDirPath())
	}
	expectedFile := filepath.Join(tempDir, "config.json")
	if ConfigFilePath() != expectedFile {
		t.Fatalf("Expected ConfigFilePath %s, got %s", expectedFile, ConfigFilePath())
	}

	cfg := LoadWorkspaceConfig()
	cfg.PreferredEditor = "cursor"
	cfg.DefaultTheme = "cyber-dark"
	cfg.DefaultSlideCount = 8
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Reload and verify
	reloaded := LoadWorkspaceConfig()
	if reloaded.PreferredEditor != "cursor" {
		t.Fatalf("Expected preferredEditor cursor, got %s", reloaded.PreferredEditor)
	}
	if reloaded.DefaultTheme != "cyber-dark" {
		t.Fatalf("Expected defaultTheme cyber-dark, got %s", reloaded.DefaultTheme)
	}
	if reloaded.DefaultSlideCount != 8 {
		t.Fatalf("Expected defaultSlideCount 8, got %d", reloaded.DefaultSlideCount)
	}
}

func TestWorkspaceConfigGetSet(t *testing.T) {
	cfg := &WorkspaceConfig{
		PreferredEditor:   "code",
		DefaultTheme:      "academic-crimson",
		DefaultSlideCount: 5,
		ServerPort:        8080,
		AutoWatch:         true,
	}

	val, err := cfg.Get("preferredEditor")
	if err != nil || val != "code" {
		t.Fatalf("Expected code, got %s, err: %v", val, err)
	}

	val, err = cfg.Get("theme")
	if err != nil || val != "academic-crimson" {
		t.Fatalf("Expected academic-crimson, got %s, err: %v", val, err)
	}

	if err := cfg.Set("preferredEditor", "zed"); err != nil {
		t.Fatalf("Failed to set preferredEditor: %v", err)
	}
	if cfg.PreferredEditor != "zed" {
		t.Fatalf("Expected zed, got %s", cfg.PreferredEditor)
	}

	if err := cfg.Set("slides", "12"); err != nil {
		t.Fatalf("Failed to set slides: %v", err)
	}
	if cfg.DefaultSlideCount != 12 {
		t.Fatalf("Expected 12 slides, got %d", cfg.DefaultSlideCount)
	}

	// Test invalid keys / values
	if err := cfg.Set("slides", "99"); err == nil {
		t.Fatalf("Expected error for slides > 50")
	}
	if err := cfg.Set("unknownKey", "val"); err == nil {
		t.Fatalf("Expected error for unknown key")
	}
}
