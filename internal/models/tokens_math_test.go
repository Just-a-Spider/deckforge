package models

import (
	"math"
	"testing"
)

func TestContrastRatio(t *testing.T) {
	// Pure black on pure white must be 21:1
	ratio, err := ContrastRatio("#000000", "#ffffff")
	if err != nil {
		t.Fatalf("ContrastRatio failed: %v", err)
	}
	if math.Abs(ratio-21.0) > 0.1 {
		t.Errorf("expected ~21.0, got %.2f", ratio)
	}

	// Same color must be 1.0:1
	ratioSame, err := ContrastRatio("#333333", "#333333")
	if err != nil {
		t.Fatalf("ContrastRatio same failed: %v", err)
	}
	if math.Abs(ratioSame-1.0) > 0.01 {
		t.Errorf("expected 1.0, got %.2f", ratioSame)
	}
}

func TestGenerateRecommendedAction(t *testing.T) {
	// Low contrast: dark gray text (#444444) on dark background (#0f172a)
	rec := GenerateRecommendedAction("#444444", "#0f172a", "textPrimary")
	if rec == nil {
		t.Fatalf("expected recommended action for low contrast pair, got nil")
	}

	if rec.RatioBefore >= 4.5 {
		t.Errorf("expected RatioBefore < 4.5, got %.2f", rec.RatioBefore)
	}
	if rec.RatioAfter < 4.5 {
		t.Errorf("expected RatioAfter >= 4.5, got %.2f", rec.RatioAfter)
	}

	// Verify the suggested hex actually satisfies contrast against the bg
	suggestedRatio, err := ContrastRatio(rec.SuggestedHex, "#0f172a")
	if err != nil {
		t.Fatalf("ContrastRatio failed on suggested hex: %v", err)
	}
	if suggestedRatio < 4.5 {
		t.Errorf("suggested hex %s failed WCAG AA: ratio=%.2f", rec.SuggestedHex, suggestedRatio)
	}
}

func TestEvaluateThemeContrast(t *testing.T) {
	palette := ThemePalette{
		BgCanvas:      "#090d16",
		BgSurface:     "#111827",
		TextPrimary:   "#f8fafc",
		TextSecondary: "#94a3b8",
		AccentPrimary: "#10b981",
	}

	reports := EvaluateThemeContrast(palette)
	if len(reports) != 4 {
		t.Fatalf("expected 4 reports, got %d", len(reports))
	}

	// TextPrimary (#f8fafc) on BgCanvas (#090d16) should pass AA and AAA
	if !reports[0].PassesAA {
		t.Errorf("expected textPrimary on bgCanvas to pass AA, ratio: %.2f", reports[0].Ratio)
	}
}
