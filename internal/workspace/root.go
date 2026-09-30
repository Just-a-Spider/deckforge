package workspace

import (
	"os"
	"path/filepath"
)

// DetectWorkspaceRoot determines the canonical workspace root from a target path or cwd.
// It walks up the directory tree from targetPath to discover an enclosing .deckforge/
// directory, deckforge.json configuration, or decks/themes structure.
func DetectWorkspaceRoot(targetPath string) string {
	var startDir string
	if targetPath == "" || targetPath == "." {
		if cwd, err := os.Getwd(); err == nil {
			startDir = cwd
		}
	} else {
		if abs, err := filepath.Abs(targetPath); err == nil {
			if fi, err := os.Stat(abs); err == nil && !fi.IsDir() {
				startDir = filepath.Dir(abs)
			} else {
				startDir = abs
			}
		}
	}

	if startDir == "" {
		cwd, _ := os.Getwd()
		return cwd
	}

	curr := startDir
	for {
		// 1. Check if curr itself is the .deckforge directory
		if filepath.Base(curr) == ".deckforge" {
			return filepath.Dir(curr)
		}

		// 2. Check for child .deckforge directory
		dotDir := filepath.Join(curr, ".deckforge")
		if fi, err := os.Stat(dotDir); err == nil && fi.IsDir() {
			return curr
		}

		// 2. Check for deckforge.json config
		confFile := filepath.Join(curr, "deckforge.json")
		if fi, err := os.Stat(confFile); err == nil && !fi.IsDir() {
			return curr
		}

		// 3. Check for co-presence of decks/ and themes/
		decksDir := filepath.Join(curr, "decks")
		themesDir := filepath.Join(curr, "themes")
		hasDecks, _ := os.Stat(decksDir)
		hasThemes, _ := os.Stat(themesDir)
		if hasDecks != nil && hasThemes != nil && hasDecks.IsDir() && hasThemes.IsDir() {
			return curr
		}

		parent := filepath.Dir(curr)
		if parent == curr || parent == "." || parent == "/" {
			break
		}
		curr = parent
	}

	// 4. If not found walking up, check cwd
	if cwd, err := os.Getwd(); err == nil {
		dotDir := filepath.Join(cwd, ".deckforge")
		if fi, err := os.Stat(dotDir); err == nil && fi.IsDir() {
			return cwd
		}
		return cwd
	}

	return startDir
}
