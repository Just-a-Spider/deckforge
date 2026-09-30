package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanDecksInRoot(t *testing.T) {
	tempWorkspace, err := os.MkdirTemp("", "deckforge-disc-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempWorkspace)

	// Create a deck in legacy decks/
	legacyDeck := filepath.Join(tempWorkspace, "decks", "pitch_deck")
	if err := os.MkdirAll(filepath.Join(legacyDeck, "slides"), 0755); err != nil {
		t.Fatalf("Failed to create legacy deck: %v", err)
	}
	_ = os.WriteFile(filepath.Join(legacyDeck, "slides", "01_intro.html"), []byte("<h1>Pitch</h1>"), 0644)

	// Create a deck in .deckforge/decks/
	dotDeck := filepath.Join(tempWorkspace, ".deckforge", "decks", "keynote_2026")
	if err := os.MkdirAll(filepath.Join(dotDeck, "slides"), 0755); err != nil {
		t.Fatalf("Failed to create dot deck: %v", err)
	}
	_ = os.WriteFile(filepath.Join(dotDeck, "slides", "01_welcome.html"), []byte("<h1>Welcome</h1>"), 0644)

	decks, err := ScanDecksInRoot(tempWorkspace)
	if err != nil {
		t.Fatalf("ScanDecksInRoot error: %v", err)
	}

	if len(decks) != 2 {
		t.Fatalf("Expected 2 decks, got %d", len(decks))
	}

	foundLegacy := false
	foundDot := false
	for _, d := range decks {
		if d.Name == "pitch_deck" {
			foundLegacy = true
		}
		if d.Name == "keynote_2026" {
			foundDot = true
		}
	}

	if !foundLegacy {
		t.Fatalf("Legacy deck pitch_deck not found")
	}
	if !foundDot {
		t.Fatalf("Dot deck keynote_2026 not found")
	}
}
