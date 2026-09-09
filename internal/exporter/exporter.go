package exporter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"deckforge/internal/compiler"
	"deckforge/internal/workspace"
)

// ExportDeck exports the presentation to PDF using headless Chrome/Chromium
func ExportDeck(deckPath, format string, comp *compiler.Compiler) (string, error) {
	deck, err := workspace.InspectDeck(deckPath)
	if err != nil {
		return "", err
	}

	// 1. Ensure latest build
	htmlPath, err := comp.Build(deckPath, "")
	if err != nil {
		return "", fmt.Errorf("build failed before export: %w", err)
	}

	// 2. Discover Chrome / Chromium binary
	chromeBin := findChromeBinary()
	if chromeBin == "" {
		return "", fmt.Errorf("headless Chrome or Chromium not found on system PATH")
	}

	// 3. Setup output directory
	exportsDir := filepath.Join(deckPath, "exports")
	if err := os.MkdirAll(exportsDir, 0755); err != nil {
		return "", err
	}

	pdfOut := filepath.Join(exportsDir, fmt.Sprintf("%s.pdf", deck.Name))

	args := []string{
		"--headless",
		"--disable-gpu",
		"--no-pdf-header-footer",
		fmt.Sprintf("--print-to-pdf=%s", pdfOut),
		fmt.Sprintf("file://%s", htmlPath),
	}

	cmd := exec.Command(chromeBin, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("chrome export error: %w (output: %s)", err, string(out))
	}

	return pdfOut, nil
}

func findChromeBinary() string {
	candidates := []string{
		"google-chrome",
		"google-chrome-stable",
		"chromium",
		"chromium-browser",
		"/usr/bin/google-chrome",
		"/usr/bin/chromium",
		"/usr/bin/chromium-browser",
		"/home/andre/.local/bin/google-chrome",
	}

	for _, c := range candidates {
		if path, err := exec.LookPath(c); err == nil {
			return path
		}
		if fi, err := os.Stat(c); err == nil && !fi.IsDir() {
			return c
		}
	}
	return ""
}
