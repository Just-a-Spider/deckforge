package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var dataSlideRegex = regexp.MustCompile(`data-slide=["']\d+["']`)
var prefixRegex = regexp.MustCompile(`^(\d+)[_-](.*)$`)

// ReorderSlides atomically renames slides in deckPath based on newOrder (1-based indices).
// e.g. newOrder = [3, 1, 2] means the old 3rd slide becomes 1st, 1st becomes 2nd, 2nd becomes 3rd.
func ReorderSlides(deckPath string, newOrder []int) error {
	deck, err := InspectDeck(deckPath)
	if err != nil {
		return fmt.Errorf("failed to inspect deck: %w", err)
	}

	total := len(deck.Slides)
	if total == 0 {
		return fmt.Errorf("deck has no slides")
	}

	if len(newOrder) != total {
		return fmt.Errorf("newOrder length (%d) does not match total slides (%d)", len(newOrder), total)
	}

	// Validate permutation
	seen := make(map[int]bool)
	for _, idx := range newOrder {
		if idx < 1 || idx > total {
			return fmt.Errorf("invalid slide index in newOrder: %d (expected 1-%d)", idx, total)
		}
		if seen[idx] {
			return fmt.Errorf("duplicate slide index in newOrder: %d", idx)
		}
		seen[idx] = true
	}

	slidesDir := filepath.Join(deckPath, "slides")

	// Step 1: Stage files with temporary names to avoid collisions
	type stagedFile struct {
		tempPath  string
		finalName string
		newIndex  int
	}

	staged := make([]stagedFile, total)

	for newPos, oldIdx := range newOrder {
		newNum := newPos + 1
		oldSlide := deck.Slides[oldIdx-1]

		// Extract base slug without leading number
		baseSlug := oldSlide.Slug
		if m := prefixRegex.FindStringSubmatch(oldSlide.Filename); len(m) > 2 {
			baseSlug = strings.TrimSuffix(m[2], ".html")
		}

		finalName := fmt.Sprintf("%02d_%s.html", newNum, baseSlug)
		tempName := fmt.Sprintf(".stage_%d_%s", newNum, finalName)

		oldPath := oldSlide.Path
		tempPath := filepath.Join(slidesDir, tempName)

		if err := os.Rename(oldPath, tempPath); err != nil {
			return fmt.Errorf("staging rename failed for %s -> %s: %w", oldPath, tempPath, err)
		}

		staged[newPos] = stagedFile{
			tempPath:  tempPath,
			finalName: finalName,
			newIndex:  newNum,
		}
	}

	// Step 2: Update content data-slide attribute and rename to final name
	for _, sf := range staged {
		finalPath := filepath.Join(slidesDir, sf.finalName)

		content, err := os.ReadFile(sf.tempPath)
		if err == nil {
			updatedContent := dataSlideRegex.ReplaceAllString(string(content), fmt.Sprintf(`data-slide="%d"`, sf.newIndex))
			_ = os.WriteFile(sf.tempPath, []byte(updatedContent), 0644)
		}

		if err := os.Rename(sf.tempPath, finalPath); err != nil {
			return fmt.Errorf("final rename failed for %s -> %s: %w", sf.tempPath, finalPath, err)
		}
	}

	return nil
}
