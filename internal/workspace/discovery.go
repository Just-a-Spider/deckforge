package workspace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"deckforge/internal/models"
)

// IsDirDeck checks if a folder is a deck by duck-typing (has deck.json or slides/*.html)
func IsDirDeck(dir string) bool {
	deckJSON := filepath.Join(dir, "deck.json")
	if _, err := os.Stat(deckJSON); err == nil {
		return true
	}

	slidesDir := filepath.Join(dir, "slides")
	entries, err := os.ReadDir(slidesDir)
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".html") {
				return true
			}
		}
	}
	return false
}

// InspectDeck loads deck details and slides from a directory
func InspectDeck(dir string) (*models.Deck, error) {
	cleanDir, err := filepath.Abs(dir)
	if err != nil {
		cleanDir = filepath.Clean(dir)
	}

	deckName := filepath.Base(cleanDir)
	conf := models.DeckConfig{
		Title: strings.Title(strings.ReplaceAll(deckName, "_", " ")),
		Theme: "academic-crimson",
	}

	// Read deck.json if present
	confPath := filepath.Join(cleanDir, "deck.json")
	if data, err := os.ReadFile(confPath); err == nil {
		_ = json.Unmarshal(data, &conf)
	}

	// Discover and parse slides
	slidesDir := filepath.Join(cleanDir, "slides")
	slides := []models.SlideInfo{}

	entries, err := os.ReadDir(slidesDir)
	if err == nil {
		var slideFiles []string
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".html") {
				slideFiles = append(slideFiles, e.Name())
			}
		}
		sort.Strings(slideFiles)

		for idx, fn := range slideFiles {
			filePath := filepath.Join(slidesDir, fn)
			title := extractSlideTitle(filePath, fn)
			slug := strings.TrimSuffix(fn, ".html")

			slides = append(slides, models.SlideInfo{
				Index:    idx + 1,
				Filename: fn,
				Path:     filePath,
				Title:    title,
				Slug:     slug,
			})
		}
	}

	themeName := conf.Theme
	if themeName == "" {
		themeName = "academic-crimson"
	}

	return &models.Deck{
		Name:      deckName,
		Path:      cleanDir,
		Config:    conf,
		Slides:    slides,
		ThemeName: themeName,
	}, nil
}

// SetDeckTheme updates the theme property in deck.json atomically
func SetDeckTheme(deckPath, themeName string) error {
	cleanDir, err := filepath.Abs(deckPath)
	if err != nil {
		cleanDir = filepath.Clean(deckPath)
	}

	confPath := filepath.Join(cleanDir, "deck.json")
	var conf models.DeckConfig
	if data, err := os.ReadFile(confPath); err == nil {
		_ = json.Unmarshal(data, &conf)
	} else {
		deckName := filepath.Base(cleanDir)
		conf.Title = strings.Title(strings.ReplaceAll(deckName, "_", " "))
	}

	conf.Theme = themeName
	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(confPath, append(data, '\n'), 0644)
}

// ScanDecksInRoot discovers decks in the root or any child folders up to depth 2
func ScanDecksInRoot(root string) ([]*models.Deck, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		absRoot = root
	}

	// Case 1: The root directory itself is a deck
	if IsDirDeck(absRoot) {
		deck, err := InspectDeck(absRoot)
		if err == nil {
			return []*models.Deck{deck}, nil
		}
	}

	deckMap := make(map[string]*models.Deck)

	// Case 2: Scan decks/ subfolder if it exists
	decksDir := filepath.Join(absRoot, "decks")
	if entries, err := os.ReadDir(decksDir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				target := filepath.Join(decksDir, e.Name())
				if IsDirDeck(target) {
					if d, err := InspectDeck(target); err == nil {
						deckMap[target] = d
					}
				}
			}
		}
	}

	// Case 2b: Scan .deckforge/decks/ subfolder if it exists
	dotDecksDir := models.WorkspaceDecksDir(absRoot)
	if entries, err := os.ReadDir(dotDecksDir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				target := filepath.Join(dotDecksDir, e.Name())
				if IsDirDeck(target) {
					if d, err := InspectDeck(target); err == nil {
						deckMap[target] = d
					}
				}
			}
		}
	}

	// Case 3: Scan immediate child subdirectories of root (depth 1 only)
	entries, err := os.ReadDir(absRoot)
	if err == nil {
		for _, e := range entries {
			if e.IsDir() {
				name := e.Name()
				if strings.HasPrefix(name, ".") || name == "node_modules" || name == "dist" || name == "bin" || name == "engine" || name == "decks" {
					continue
				}
				subPath := filepath.Join(absRoot, name)
				if IsDirDeck(subPath) {
					if d, err := InspectDeck(subPath); err == nil {
						deckMap[subPath] = d
					}
				}
			}
		}
	}

	var results []*models.Deck
	for _, d := range deckMap {
		results = append(results, d)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results, nil
}

var (
	titleRegex   = regexp.MustCompile(`(?i)<div class="topbar-title">([^<]+)`)
	headingRegex = regexp.MustCompile(`(?i)<h[12][^>]*>([^<]+)</h[12]>`)
)

func extractSlideTitle(filePath, filename string) string {
	content, err := os.ReadFile(filePath)
	if err == nil {
		if m := titleRegex.FindSubmatch(content); len(m) > 1 {
			t := strings.TrimSpace(string(m[1]))
			if t != "" {
				return t
			}
		}
		if m := headingRegex.FindSubmatch(content); len(m) > 1 {
			t := strings.TrimSpace(string(m[1]))
			if t != "" {
				return t
			}
		}
	}

	// Fallback to filename: strip "01_" and replace underscores with spaces
	clean := strings.TrimSuffix(filename, ".html")
	parts := strings.SplitN(clean, "_", 2)
	if len(parts) == 2 {
		return strings.Title(strings.ReplaceAll(parts[1], "_", " "))
	}
	return strings.Title(strings.ReplaceAll(clean, "_", " "))
}
