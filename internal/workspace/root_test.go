package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectWorkspaceRoot(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "deckforge_ws_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create workspace structure: tempDir/.deckforge/decks/my-deck
	deckDir := filepath.Join(tempDir, ".deckforge", "decks", "my-deck")
	if err := os.MkdirAll(deckDir, 0755); err != nil {
		t.Fatalf("Failed to create deckDir: %v", err)
	}

	// 1. Walking up from nested deckDir should find tempDir
	root := DetectWorkspaceRoot(deckDir)
	if root != tempDir {
		t.Errorf("Expected root %s, got %s", tempDir, root)
	}

	// 2. Walking up from a file inside deckDir
	slideFile := filepath.Join(deckDir, "slide-01.html")
	_ = os.WriteFile(slideFile, []byte("<h1>Slide</h1>"), 0644)
	rootFromFile := DetectWorkspaceRoot(slideFile)
	if rootFromFile != tempDir {
		t.Errorf("Expected root from file %s, got %s", tempDir, rootFromFile)
	}
}
