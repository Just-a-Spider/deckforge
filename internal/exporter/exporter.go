package exporter

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"deckforge/internal/compiler"
	"deckforge/internal/workspace"
)

// ExportDeck exports the presentation to PDF using the default chromedp engine
func ExportDeck(deckPath, format string, comp *compiler.Compiler) (string, error) {
	return ExportDeckWithOptions(deckPath, format, "chromedp", 30*time.Second, comp)
}

// ExportDeckWithOptions exports the presentation to PDF with explicit engine and timeout controls
func ExportDeckWithOptions(deckPath, format, engine string, timeout time.Duration, comp *compiler.Compiler) (string, error) {
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
	if chromeBin == "" && engine != "node" {
		return "", fmt.Errorf("headless Chrome or Chromium not found on system (checked PATH, $CHROME_BIN, standard OS locations)")
	}

	// 3. Setup output directory
	exportsDir := filepath.Join(deckPath, "exports")
	if err := os.MkdirAll(exportsDir, 0755); err != nil {
		return "", err
	}

	pdfOut := filepath.Join(exportsDir, fmt.Sprintf("%s.pdf", deck.Name))
	absHtmlPath, err := filepath.Abs(htmlPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute html path: %w", err)
	}

	absPdfOut, err := filepath.Abs(pdfOut)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute pdf path: %w", err)
	}

	// Remove old PDF if present to ensure fresh generation
	_ = os.Remove(absPdfOut)

	var exportErr error
	switch strings.ToLower(engine) {
	case "node":
		exportErr = exportWithNode(absHtmlPath, absPdfOut)
	case "cli":
		exportErr = exportWithCLI(chromeBin, absHtmlPath, absPdfOut)
	default: // "chromedp" or auto
		exportErr = exportWithChromedp(chromeBin, absHtmlPath, absPdfOut, timeout)
		if exportErr != nil {
			// Graceful fallback to CLI if chromedp encounters an unexpected environment failure
			if cliErr := exportWithCLI(chromeBin, absHtmlPath, absPdfOut); cliErr == nil {
				exportErr = nil
			}
		}
	}

	if exportErr != nil {
		return "", exportErr
	}

	if fi, err := os.Stat(absPdfOut); err != nil || fi.Size() == 0 {
		return "", fmt.Errorf("pdf export failed: output file %s was not created", absPdfOut)
	}

	return pdfOut, nil
}

func exportWithCLI(chromeBin, absHtmlPath, absPdfOut string) error {
	baseArgs := []string{
		"--disable-gpu",
		"--no-pdf-header-footer",
		"--run-all-compositor-stages-before-draw",
		"--virtual-time-budget=3000",
		fmt.Sprintf("--print-to-pdf=%s", absPdfOut),
		fmt.Sprintf("file://%s", absHtmlPath),
	}

	cmd := exec.Command(chromeBin, append([]string{"--headless=new"}, baseArgs...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		legacyCmd := exec.Command(chromeBin, append([]string{"--headless"}, baseArgs...)...)
		legacyOut, legacyErr := legacyCmd.CombinedOutput()
		if legacyErr != nil {
			return fmt.Errorf("chrome export error: %w (output: %s; legacy output: %s)", err, string(out), string(legacyOut))
		}
	}
	return nil
}

//go:embed export-pdf.mjs
var embeddedNodeRunner string

func exportWithNode(absHtmlPath, absPdfOut string) error {
	nodeBin, err := exec.LookPath("node")
	if err != nil {
		return fmt.Errorf("node runtime not found in PATH for node export engine")
	}

	cwd, _ := os.Getwd()
	candidates := []string{
		filepath.Join(cwd, "scripts", "export-pdf.mjs"),
		filepath.Join(cwd, "..", "scripts", "export-pdf.mjs"),
		filepath.Join(cwd, "..", "..", "scripts", "export-pdf.mjs"),
	}

	var scriptPath string
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && !fi.IsDir() {
			scriptPath = c
			break
		}
	}

	if scriptPath == "" {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			cacheDir = os.TempDir()
		}
		dfCache := filepath.Join(cacheDir, "deckforge")
		if err := os.MkdirAll(dfCache, 0755); err != nil {
			return fmt.Errorf("failed to create cache dir for node runner: %w", err)
		}
		cachedScript := filepath.Join(dfCache, "export-pdf.mjs")
		if err := os.WriteFile(cachedScript, []byte(embeddedNodeRunner), 0644); err != nil {
			return fmt.Errorf("failed to extract embedded node runner: %w", err)
		}
		scriptPath = cachedScript
	}

	cmd := exec.Command(nodeBin, scriptPath, absHtmlPath, absPdfOut)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("node export runner failed: %w (output: %s)", err, string(out))
	}

	return nil
}

func findChromeBinary() string {
	// 1. Check Chrome-specific environment variables
	for _, env := range []string{"CHROME_BIN", "CHROMIUM_BIN"} {
		if val := os.Getenv(env); val != "" {
			if path, err := exec.LookPath(val); err == nil {
				return path
			}
			if fi, err := os.Stat(val); err == nil && !fi.IsDir() {
				return val
			}
		}
	}

	// 2. Common command names on system PATH
	commands := []string{
		"google-chrome",
		"google-chrome-stable",
		"chromium",
		"chromium-browser",
		"brave-browser",
		"microsoft-edge",
		"microsoft-edge-stable",
		"edge",
	}
	for _, cmd := range commands {
		if path, err := exec.LookPath(cmd); err == nil {
			return path
		}
	}

	// 3. User home directory candidates (dynamic, user-agnostic)
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		userCandidates := []string{
			filepath.Join(home, ".local", "bin", "google-chrome"),
			filepath.Join(home, ".local", "bin", "chromium"),
			filepath.Join(home, ".local", "bin", "chromium-browser"),
			filepath.Join(home, "bin", "google-chrome"),
			filepath.Join(home, "bin", "chromium"),
		}
		for _, u := range userCandidates {
			if fi, err := os.Stat(u); err == nil && !fi.IsDir() {
				return u
			}
		}
	}

	// 4. Platform-specific standard installation locations
	switch runtime.GOOS {
	case "darwin":
		macCandidates := []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
			"/Applications/Brave Browser.app/Contents/MacOS/Brave Browser",
			"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
		}
		for _, m := range macCandidates {
			if fi, err := os.Stat(m); err == nil && !fi.IsDir() {
				return m
			}
		}
	case "windows":
		winRoots := []string{
			os.Getenv("ProgramFiles"),
			os.Getenv("ProgramFiles(x86)"),
			os.Getenv("LocalAppData"),
		}
		winSubpaths := []string{
			`Google\Chrome\Application\chrome.exe`,
			`Chromium\Application\chrome.exe`,
			`BraveSoftware\Brave-Browser\Application\brave.exe`,
			`Microsoft\Edge\Application\msedge.exe`,
		}
		for _, root := range winRoots {
			if root == "" {
				continue
			}
			for _, sub := range winSubpaths {
				target := filepath.Join(root, sub)
				if fi, err := os.Stat(target); err == nil && !fi.IsDir() {
					return target
				}
			}
		}
	default: // Linux and other Unix-like OS
		linuxCandidates := []string{
			"/usr/bin/google-chrome",
			"/usr/bin/google-chrome-stable",
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
			"/snap/bin/chromium",
			"/var/lib/flatpak/exports/bin/org.chromium.Chromium",
		}
		for _, l := range linuxCandidates {
			if fi, err := os.Stat(l); err == nil && !fi.IsDir() {
				return l
			}
		}
	}

	return ""
}
