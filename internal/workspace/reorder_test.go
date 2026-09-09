package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestReorderSlides(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "deckforge_reorder_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	slidesDir := filepath.Join(tempDir, "slides")
	if err := os.MkdirAll(slidesDir, 0755); err != nil {
		t.Fatalf("failed to create slides dir: %v", err)
	}

	// Create 3 slides
	files := []string{"01_intro.html", "02_architecture.html", "03_summary.html"}
	for i, f := range files {
		content := fmt.Sprintf(`<section class="slide" data-slide="%d"><h1>Slide %d</h1></section>`, i+1, i+1)
		if err := os.WriteFile(filepath.Join(slidesDir, f), []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test slide: %v", err)
		}
	}

	// Reorder [3, 1, 2] -> old 3 becomes 1, old 1 becomes 2, old 2 becomes 3
	newOrder := []int{3, 1, 2}
	if err := ReorderSlides(tempDir, newOrder); err != nil {
		t.Fatalf("ReorderSlides failed: %v", err)
	}

	deck, err := InspectDeck(tempDir)
	if err != nil {
		t.Fatalf("InspectDeck failed: %v", err)
	}

	if len(deck.Slides) != 3 {
		t.Fatalf("expected 3 slides, got %d", len(deck.Slides))
	}

	expectedNames := []string{"01_summary.html", "02_intro.html", "03_architecture.html"}
	for i, expected := range expectedNames {
		if deck.Slides[i].Filename != expected {
			t.Errorf("slide %d: expected filename %s, got %s", i+1, expected, deck.Slides[i].Filename)
		}
		// Verify data-slide inside content was updated
		content, err := os.ReadFile(deck.Slides[i].Path)
		if err != nil {
			t.Errorf("failed to read slide %s: %v", deck.Slides[i].Filename, err)
		}
		expectedAttr := fmt.Sprintf(`data-slide="%d"`, i+1)
		if !stringsContains(string(content), expectedAttr) {
			t.Errorf("slide %s does not contain %s", deck.Slides[i].Filename, expectedAttr)
		}
	}
}

func stringsContains(s, substr string) bool {
	return filepath.Base(s) != "" && (len(s) >= len(substr)) && (s == substr || len(s) > 0 && searchSubstr(s, substr))
}

func searchSubstr(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestSetDeckTheme(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "deckforge_theme_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	if err := SetDeckTheme(tempDir, "cyber-dark"); err != nil {
		t.Fatalf("SetDeckTheme failed: %v", err)
	}

	deck, err := InspectDeck(tempDir)
	if err != nil {
		t.Fatalf("InspectDeck failed: %v", err)
	}

	if deck.ThemeName != "cyber-dark" {
		t.Errorf("expected theme cyber-dark, got %s", deck.ThemeName)
	}

	// Update to another theme
	if err := SetDeckTheme(tempDir, "academic-crimson"); err != nil {
		t.Fatalf("SetDeckTheme second call failed: %v", err)
	}

	deck2, err := InspectDeck(tempDir)
	if err != nil {
		t.Fatalf("InspectDeck failed: %v", err)
	}

	if deck2.ThemeName != "academic-crimson" {
		t.Errorf("expected theme academic-crimson, got %s", deck2.ThemeName)
	}
}
