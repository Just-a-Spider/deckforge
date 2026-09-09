package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"deckforge/internal/compiler"
	"deckforge/internal/models"
	"deckforge/internal/scaffold"
	"deckforge/internal/theme"
)

func TestSanitizeSlideHTML(t *testing.T) {
	raw := `<section class="slide active df-selected-node df-hovered" data-slide="1" contenteditable="true" df-selected="true" data-studio-id="elem-123">
    <h1 class="topbar-title df-editable-active" contenteditable="true">Title</h1>
</section>`

	cleaned := SanitizeSlideHTML(raw)

	if strings.Contains(cleaned, "contenteditable") {
		t.Errorf("cleaned HTML still contains contenteditable: %s", cleaned)
	}
	if strings.Contains(cleaned, "df-selected") {
		t.Errorf("cleaned HTML still contains df-selected: %s", cleaned)
	}
	if strings.Contains(cleaned, "data-studio-id") {
		t.Errorf("cleaned HTML still contains data-studio-id: %s", cleaned)
	}
	if strings.Contains(cleaned, "df-selected-node") {
		t.Errorf("cleaned HTML still contains df-selected-node: %s", cleaned)
	}
	if strings.Contains(cleaned, "df-hovered") {
		t.Errorf("cleaned HTML still contains df-hovered: %s", cleaned)
	}
	if strings.Contains(cleaned, "df-editable-active") {
		t.Errorf("cleaned HTML still contains df-editable-active: %s", cleaned)
	}
}

func TestSlideUpdateAPI(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "deckforge_api_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	deck, err := scaffold.ScaffoldDeck(tempDir, "api_test_deck", "academic-crimson", 3)
	if err != nil {
		t.Fatalf("failed to scaffold deck: %v", err)
	}

	tm := theme.NewThemeManager(tempDir)
	comp := compiler.NewCompiler(tm)
	srv := NewDeckServer(tempDir, 8080, comp)

	mux := http.NewServeMux()
	registerAPIRoutes(mux, srv)

	// Test GET /api/deck
	req := httptest.NewRequest(http.MethodGet, "/api/deck", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /api/deck returned %d", w.Code)
	}

	// Test POST /api/slides/update
	updatedSlideHTML := `<section class="slide" data-slide="1" contenteditable="true"><h1>Updated Slide Title</h1></section>`
	updatePayload := SlideUpdateRequest{
		Index: 1,
		HTML:  updatedSlideHTML,
	}
	payloadBytes, _ := json.Marshal(updatePayload)
	reqUpdate := httptest.NewRequest(http.MethodPost, "/api/slides/update", bytes.NewReader(payloadBytes))
	wUpdate := httptest.NewRecorder()
	mux.ServeHTTP(wUpdate, reqUpdate)

	if wUpdate.Code != http.StatusOK {
		t.Fatalf("POST /api/slides/update returned %d: %s", wUpdate.Code, wUpdate.Body.String())
	}

	// Verify file on disk
	firstSlidePath := deck.Slides[0].Path
	diskContent, err := os.ReadFile(firstSlidePath)
	if err != nil {
		t.Fatalf("failed to read first slide file: %v", err)
	}

	if !strings.Contains(string(diskContent), "Updated Slide Title") {
		t.Errorf("expected updated title on disk, got: %s", string(diskContent))
	}
	if strings.Contains(string(diskContent), "contenteditable") {
		t.Errorf("saved slide should not contain contenteditable attribute: %s", string(diskContent))
	}
}

func TestThemesAndDeckCreateAPI(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "deckforge_theme_api_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	_, err = scaffold.ScaffoldDeck(tempDir, "theme_test_deck", "cyber-dark", 2)
	if err != nil {
		t.Fatalf("failed to scaffold deck: %v", err)
	}

	tm := theme.NewThemeManager(tempDir)
	comp := compiler.NewCompiler(tm)
	srv := NewDeckServer(tempDir, 8080, comp)

	mux := http.NewServeMux()
	registerAPIRoutes(mux, srv)

	// 1. Test GET /api/themes
	reqThemes := httptest.NewRequest(http.MethodGet, "/api/themes", nil)
	wThemes := httptest.NewRecorder()
	mux.ServeHTTP(wThemes, reqThemes)

	if wThemes.Code != http.StatusOK {
		t.Fatalf("GET /api/themes returned %d", wThemes.Code)
	}
	var themeList []map[string]interface{}
	if err := json.Unmarshal(wThemes.Body.Bytes(), &themeList); err != nil {
		t.Fatalf("failed to parse /api/themes response: %v", err)
	}
	if len(themeList) == 0 {
		t.Errorf("expected at least 1 theme in list, got 0")
	}

	// 2. Test POST /api/deck/theme
	themePayload := []byte(`{"theme": "academic-crimson"}`)
	reqSetTheme := httptest.NewRequest(http.MethodPost, "/api/deck/theme", bytes.NewReader(themePayload))
	wSetTheme := httptest.NewRecorder()
	mux.ServeHTTP(wSetTheme, reqSetTheme)

	if wSetTheme.Code != http.StatusOK {
		t.Fatalf("POST /api/deck/theme returned %d: %s", wSetTheme.Code, wSetTheme.Body.String())
	}

	// Verify deck.json updated
	deckData, err := os.ReadFile(tempDir + "/deck.json")
	if err != nil {
		t.Fatalf("failed to read deck.json: %v", err)
	}
	if !strings.Contains(string(deckData), "academic-crimson") {
		t.Errorf("deck.json was not updated with academic-crimson: %s", string(deckData))
	}

	// 3. Test POST /api/deck/create
	newDeckDir := tempDir + "_subdeck"
	defer os.RemoveAll(newDeckDir)
	createPayload := []byte(`{"name": "sub_deck", "theme": "cyber-dark", "slides": 3, "path": "` + newDeckDir + `"}`)
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/deck/create", bytes.NewReader(createPayload))
	wCreate := httptest.NewRecorder()
	mux.ServeHTTP(wCreate, reqCreate)

	if wCreate.Code != http.StatusOK {
		t.Fatalf("POST /api/deck/create returned %d: %s", wCreate.Code, wCreate.Body.String())
	}

	// Verify new deck on disk
	if _, err := os.Stat(newDeckDir + "/deck.json"); os.IsNotExist(err) {
		t.Errorf("created deck deck.json does not exist at %s", newDeckDir)
	}
}

func TestAPICatalogAndStylingGuide(t *testing.T) {
	tempDir := t.TempDir()
	_, err := scaffold.ScaffoldDeck(tempDir, "api_introspect_deck", "cyber-dark", 1)
	if err != nil {
		t.Fatalf("failed to scaffold deck: %v", err)
	}
	tm := theme.NewThemeManager(tempDir)
	comp := compiler.NewCompiler(tm)
	srv := NewDeckServer(tempDir, 8080, comp)
	mux := http.NewServeMux()
	registerAPIRoutes(mux, srv)

	// Test GET /api/catalog
	wCat := httptest.NewRecorder()
	mux.ServeHTTP(wCat, httptest.NewRequest(http.MethodGet, "/api/catalog", nil))
	if wCat.Code != http.StatusOK {
		t.Fatalf("GET /api/catalog returned %d", wCat.Code)
	}
	var components []models.ComponentDefinition
	if err := json.Unmarshal(wCat.Body.Bytes(), &components); err != nil {
		t.Fatalf("failed to unmarshal catalog: %v", err)
	}
	if len(components) != 7 {
		t.Fatalf("expected 7 components, got %d", len(components))
	}

	// Test GET /api/styling-guide
	wGuide := httptest.NewRecorder()
	mux.ServeHTTP(wGuide, httptest.NewRequest(http.MethodGet, "/api/styling-guide", nil))
	if wGuide.Code != http.StatusOK {
		t.Fatalf("GET /api/styling-guide returned %d", wGuide.Code)
	}
	var guide models.StylingGuide
	if err := json.Unmarshal(wGuide.Body.Bytes(), &guide); err != nil {
		t.Fatalf("failed to unmarshal styling guide: %v", err)
	}
	if guide.CanvasRules.StageWidth != 1920 {
		t.Fatalf("expected stage width 1920, got %d", guide.CanvasRules.StageWidth)
	}
}


