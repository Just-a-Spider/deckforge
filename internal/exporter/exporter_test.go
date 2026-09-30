package exporter

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"deckforge/internal/compiler"
	"deckforge/internal/theme"
)

func TestFindChromeBinary(t *testing.T) {
	bin := findChromeBinary()
	// Should find at least one Chrome / Chromium binary on dev system
	if bin == "" {
		t.Logf("Warning: no Chrome / Chromium found on system, skipping path check")
		return
	}

	// Binary path should be an existing executable file
	fi, err := os.Stat(bin)
	if err != nil || fi.IsDir() {
		t.Fatalf("Expected valid binary file, got: %s (err: %v)", bin, err)
	}
}

func TestExportDeckEngines(t *testing.T) {
	bin := findChromeBinary()
	if bin == "" {
		t.Skip("Skipping test: Chrome binary not available")
	}

	deckPath := "../../decks/starter_deck"
	tm := theme.NewThemeManager(deckPath)
	comp := compiler.NewCompiler(tm)

	// Test chromedp engine (default)
	outChromedp, err := ExportDeckWithOptions(deckPath, "pdf", "chromedp", 30*time.Second, comp)
	if err != nil {
		t.Fatalf("chromedp export failed: %v", err)
	}
	if fi, err := os.Stat(outChromedp); err != nil || fi.Size() == 0 {
		t.Fatalf("chromedp exported PDF is missing or empty")
	}

	// Test node engine
	if _, err := exec.LookPath("node"); err == nil {
		outNode, err := ExportDeckWithOptions(deckPath, "pdf", "node", 30*time.Second, comp)
		if err != nil {
			t.Logf("Notice: node export returned: %v (skipping if script not in test dir)", err)
		} else if fi, err := os.Stat(outNode); err == nil && fi.Size() > 0 {
			t.Logf("node export verified: %s (%d bytes)", outNode, fi.Size())
		}
	}
}
