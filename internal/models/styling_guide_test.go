package models

import "testing"

func TestGetStylingGuide(t *testing.T) {
	guide := GetStylingGuide()
	if guide.CanvasRules.StageWidth != 1920 || guide.CanvasRules.StageHeight != 1080 {
		t.Fatalf("Unexpected canvas rules: %v", guide.CanvasRules)
	}
	if len(guide.Categories) == 0 {
		t.Fatal("Expected non-empty class categories")
	}
	if len(guide.SpacingScale) != 6 {
		t.Fatalf("Expected 6 spacing steps, got %d", len(guide.SpacingScale))
	}
	if len(guide.ColorTokens) == 0 {
		t.Fatal("Expected non-empty color tokens")
	}

	// Verify chapter break is present in Containers
	foundChapterBreak := false
	for _, cat := range guide.Categories {
		if cat.Category == "Containers" {
			for _, cls := range cat.Classes {
				if cls.Class == ".slide.chapter-break" {
					foundChapterBreak = true
					break
				}
			}
		}
	}
	if !foundChapterBreak {
		t.Fatal("Expected .slide.chapter-break in Containers category")
	}
}
