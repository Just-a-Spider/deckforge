package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"deckforge/internal/compiler"
	"deckforge/internal/models"
	"deckforge/internal/scaffold"
	"deckforge/internal/theme"
	"deckforge/internal/workspace"

	tea "github.com/charmbracelet/bubbletea"
)

func TestScaffoldNavigationNoPanic(t *testing.T) {
	v := NewScaffoldView("/tmp")

	// Press down 12 times to cycle past field 3 multiple times
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	for i := 0; i < 12; i++ {
		v.Update(downMsg)
	}

	// Press up 12 times
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	for i := 0; i < 12; i++ {
		v.Update(upMsg)
	}

	// If no panic occurred, test passes
	if v.FocusIndex < 0 || v.FocusIndex > 3 {
		t.Fatalf("FocusIndex out of bounds: %d", v.FocusIndex)
	}
}

func TestSlidesScrollingViewport(t *testing.T) {
	tm := theme.NewThemeManager(".")
	comp := compiler.NewCompiler(tm)
	cfg := models.LoadWorkspaceConfig()

	srvCtrl := NewServerController()
	v := NewSlidesView(comp, cfg, srvCtrl)

	deckPath := "decks/proteccion_datos_peru"
	if _, err := os.Stat(deckPath); err != nil {
		deckPath = "../../decks/proteccion_datos_peru"
	}
	deck, err := workspace.InspectDeck(deckPath)
	if err != nil {
		t.Fatalf("Failed to inspect test deck: %v", err)
	}

	v.SetDeck(deck)
	if len(v.CurrentDeck.Slides) != 22 {
		t.Fatalf("Expected 22 slides, got %d", len(v.CurrentDeck.Slides))
	}

	downMsg := tea.KeyMsg{Type: tea.KeyDown}

	// Navigate down to slide 15
	for i := 0; i < 15; i++ {
		v.Update(downMsg)
	}

	if v.Cursor != 15 {
		t.Fatalf("Expected cursor at 15, got %d", v.Cursor)
	}

	if v.ViewportOffset <= 0 {
		t.Fatalf("Expected ViewportOffset to advance, got %d", v.ViewportOffset)
	}

	viewOutput := v.View()
	if viewOutput == "" {
		t.Fatal("View output was empty")
	}
}

func TestTabNavigationAndSettings(t *testing.T) {
	app := NewAppModel(".")

	if app.ActiveTab != TabDecks {
		t.Fatalf("Expected initial tab TabDecks, got %d", app.ActiveTab)
	}

	// Test next tab with right bracket
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{']'}})
	if app.ActiveTab != TabSlides {
		t.Fatalf("Expected TabSlides after ']', got %d", app.ActiveTab)
	}

	// Test next tab with right arrow on Slides tab
	app.Update(tea.KeyMsg{Type: tea.KeyRight})
	if app.ActiveTab != TabScaffold {
		t.Fatalf("Expected TabScaffold after KeyRight, got %d", app.ActiveTab)
	}

	// Test settings tab
	app.ActiveTab = TabSettings
	prevEditor := app.Config.PreferredEditor
	app.SettingsView.Update(tea.KeyMsg{Type: tea.KeyRight})
	if app.Config.PreferredEditor == prevEditor && len(app.SettingsView.EditorOptions) > 1 {
		t.Fatalf("Expected preferred editor to change, remained %s", prevEditor)
	}
}

func TestScaffoldTabNavigationModes(t *testing.T) {
	app := NewAppModel(".")
	app.ActiveTab = TabScaffold

	// When entering TabScaffold, should be in navigation mode
	if app.isTextInputActive() {
		t.Fatal("Expected isTextInputActive to be false initially on TabScaffold")
	}

	// Arrow keys should navigate tabs from TabScaffold (when not on theme selector)
	app.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if app.ActiveTab != TabSlides {
		t.Fatalf("Expected TabSlides after KeyLeft from Scaffold, got %d", app.ActiveTab)
	}

	app.Update(tea.KeyMsg{Type: tea.KeyRight})
	if app.ActiveTab != TabScaffold {
		t.Fatalf("Expected TabScaffold after KeyRight, got %d", app.ActiveTab)
	}

	// Number keys switch tabs directly in navigation mode
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	if app.ActiveTab != TabDecks {
		t.Fatalf("Expected TabDecks after pressing '1', got %d", app.ActiveTab)
	}

	// Return to TabScaffold
	app.ActiveTab = TabScaffold
	app.ScaffoldView.FocusIndex = 0

	// Press Enter to enter Edit Mode
	app.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !app.ScaffoldView.IsEditing {
		t.Fatal("Expected ScaffoldView.IsEditing to be true after Enter")
	}
	if !app.isTextInputActive() {
		t.Fatal("Expected isTextInputActive to be true during Edit Mode")
	}

	// Pressing '2' while editing should type into input, NOT switch tabs
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	if app.ActiveTab != TabScaffold {
		t.Fatalf("Tab switched while in Edit Mode! Expected TabScaffold, got %d", app.ActiveTab)
	}
	if app.ScaffoldView.Inputs[0].Value() != "2" {
		t.Fatalf("Expected input value '2', got '%s'", app.ScaffoldView.Inputs[0].Value())
	}

	// Press Esc to exit Edit Mode
	app.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if app.ScaffoldView.IsEditing {
		t.Fatal("Expected ScaffoldView.IsEditing to be false after Esc")
	}
	if app.isTextInputActive() {
		t.Fatal("Expected isTextInputActive to be false after Esc")
	}

	// Number keys now switch tabs again
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}})
	if app.ActiveTab != TabSettings {
		t.Fatalf("Expected TabSettings after pressing '5', got %d", app.ActiveTab)
	}
}

func TestSlidesResponsiveCompactMode(t *testing.T) {
	tm := theme.NewThemeManager(".")
	comp := compiler.NewCompiler(tm)
	cfg := models.LoadWorkspaceConfig()

	srvCtrl := NewServerController()
	v := NewSlidesView(comp, cfg, srvCtrl)

	deckPath := "decks/proteccion_datos_peru"
	if _, err := os.Stat(deckPath); err != nil {
		deckPath = "../../decks/proteccion_datos_peru"
	}
	deck, err := workspace.InspectDeck(deckPath)
	if err != nil {
		t.Fatalf("Failed to inspect test deck: %v", err)
	}
	v.SetDeck(deck)

	// Small terminal height (e.g. 24 rows standard)
	v.SetDimensions(80, 24)
	if !v.IsCompact {
		t.Fatal("Expected IsCompact to be true for height 24")
	}
	if v.MaxVisibleCards <= 0 {
		t.Fatalf("Expected positive MaxVisibleCards, got %d", v.MaxVisibleCards)
	}

	compactView := v.View()
	if compactView == "" {
		t.Fatal("Compact view output was empty")
	}

	// Large terminal height (e.g. 50 rows)
	v.SetDimensions(120, 50)
	if v.IsCompact {
		t.Fatal("Expected IsCompact to be false for height 50")
	}
	if v.MaxVisibleCards <= 0 {
		t.Fatalf("Expected positive MaxVisibleCards, got %d", v.MaxVisibleCards)
	}

	normalView := v.View()
	if normalView == "" {
		t.Fatal("Normal view output was empty")
	}
}

func TestTab4ThemeNavigationAndArrowSwitching(t *testing.T) {
	app := NewAppModel(".")
	app.ActiveTab = TabThemes

	// Arrow keys left/right should switch tabs away from Tab 4
	app.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if app.ActiveTab != TabScaffold {
		t.Fatalf("Expected TabScaffold after KeyLeft on TabThemes, got %d", app.ActiveTab)
	}

	app.Update(tea.KeyMsg{Type: tea.KeyRight})
	if app.ActiveTab != TabThemes {
		t.Fatalf("Expected TabThemes after KeyRight on TabScaffold, got %d", app.ActiveTab)
	}

	app.Update(tea.KeyMsg{Type: tea.KeyRight})
	if app.ActiveTab != TabSettings {
		t.Fatalf("Expected TabSettings after KeyRight on TabThemes, got %d", app.ActiveTab)
	}

	// Return to TabThemes
	app.ActiveTab = TabThemes
	prevCursor := app.ThemeView.ThemeCursor

	// In TabThemes, down/up should navigate themes in panel 0
	app.Update(tea.KeyMsg{Type: tea.KeyDown})
	if app.ThemeView.ThemeCursor == prevCursor && len(app.ThemeView.Themes) > 1 {
		t.Fatalf("Expected ThemeCursor to advance after KeyDown, remained %d", prevCursor)
	}

	app.Update(tea.KeyMsg{Type: tea.KeyUp})
	if app.ThemeView.ThemeCursor != prevCursor {
		t.Fatalf("Expected ThemeCursor to return to %d after KeyUp, got %d", prevCursor, app.ThemeView.ThemeCursor)
	}
}

func TestServerControllerToggleSameDeck(t *testing.T) {
	tm := theme.NewThemeManager(".")
	comp := compiler.NewCompiler(tm)
	ctrl := NewServerController()

	deckPath := "decks/starter_deck"
	if _, err := os.Stat(deckPath); err != nil {
		deckPath = "../../decks/starter_deck"
	}
	deck, err := workspace.InspectDeck(deckPath)
	if err != nil {
		t.Fatalf("Failed to inspect test deck: %v", err)
	}

	port := 18081
	running, msg, err := ctrl.Toggle(deck, port, false, comp)
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer ctrl.Stop()

	if !running || !ctrl.IsRunning() {
		t.Fatalf("Expected server to be running, msg: %s", msg)
	}
	if !ctrl.IsServing(deck.Path) {
		t.Fatalf("Expected server to be serving %s", deck.Path)
	}
	if ctrl.ActivePort() != port {
		t.Fatalf("Expected active port %d, got %d", port, ctrl.ActivePort())
	}

	// Toggle second time -> should turn OFF
	running, msg, err = ctrl.Toggle(deck, port, false, comp)
	if err != nil {
		t.Fatalf("Error toggling off: %v", err)
	}
	if running || ctrl.IsRunning() {
		t.Fatalf("Expected server to be stopped, msg: %s", msg)
	}
	if ctrl.IsServing(deck.Path) {
		t.Fatal("Expected server not to be serving after toggle off")
	}
}

func TestServerControllerSwitchDecks(t *testing.T) {
	tm := theme.NewThemeManager(".")
	comp := compiler.NewCompiler(tm)
	ctrl := NewServerController()

	deck1Path := "decks/starter_deck"
	if _, err := os.Stat(deck1Path); err != nil {
		deck1Path = "../../decks/starter_deck"
	}
	deck1, err := workspace.InspectDeck(deck1Path)
	if err != nil {
		t.Fatalf("Failed to inspect deck1: %v", err)
	}

	deck2Path := "decks/proteccion_datos_peru"
	if _, err := os.Stat(deck2Path); err != nil {
		deck2Path = "../../decks/proteccion_datos_peru"
	}
	deck2, err := workspace.InspectDeck(deck2Path)
	if err != nil {
		t.Fatalf("Failed to inspect deck2: %v", err)
	}

	port := 18082
	// Start Deck 1
	running, _, err := ctrl.Toggle(deck1, port, false, comp)
	if err != nil {
		t.Fatalf("Failed to start deck1: %v", err)
	}
	defer ctrl.Stop()

	if !running || !ctrl.IsServing(deck1.Path) {
		t.Fatal("Expected deck1 to be serving")
	}

	// Toggle on Deck 2 -> should switch active deck
	running, msg, err := ctrl.Toggle(deck2, port, false, comp)
	if err != nil {
		t.Fatalf("Failed to switch to deck2: %v", err)
	}
	if !running {
		t.Fatalf("Expected server to be running on deck2, msg: %s", msg)
	}
	if ctrl.IsServing(deck1.Path) {
		t.Fatal("Expected deck1 to NO longer be serving")
	}
	if !ctrl.IsServing(deck2.Path) {
		t.Fatal("Expected deck2 to be serving")
	}

	// Explicit Stop
	if err := ctrl.Stop(); err != nil {
		t.Fatalf("Stop error: %v", err)
	}
	if ctrl.IsRunning() {
		t.Fatal("Expected server to be stopped")
	}
}

func TestDecksViewServerToggle(t *testing.T) {
	tm := theme.NewThemeManager(".")
	comp := compiler.NewCompiler(tm)
	cfg := models.LoadWorkspaceConfig()
	cfg.ServerPort = 18083
	cfg.AutoWatch = false
	ctrl := NewServerController()

	root := "."
	if _, err := os.Stat("decks"); err != nil {
		root = "../.."
	}
	v := NewDecksView(root, cfg.RecentRoots, comp, cfg, ctrl)
	defer ctrl.Stop()

	if len(v.Decks) == 0 {
		t.Fatal("No decks discovered in test root")
	}

	// Turn ON with 's'
	v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	if !ctrl.IsRunning() {
		t.Fatalf("Expected server to start on 's', status: %s", v.StatusMsg)
	}
	selected := v.SelectedDeck()
	if !ctrl.IsServing(selected.Path) {
		t.Fatalf("Expected server to be serving %s", selected.Path)
	}

	// Turn OFF with 's'
	v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	if ctrl.IsRunning() {
		t.Fatalf("Expected server to stop on second 's', status: %s", v.StatusMsg)
	}

	// Turn ON with 's', then stop with 'x'
	v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	if !ctrl.IsRunning() {
		t.Fatalf("Expected server to start on 's', status: %s", v.StatusMsg)
	}
	v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	if ctrl.IsRunning() {
		t.Fatalf("Expected server to stop on 'x', status: %s", v.StatusMsg)
	}
}

func TestSlidesViewServerToggle(t *testing.T) {
	tm := theme.NewThemeManager(".")
	comp := compiler.NewCompiler(tm)
	cfg := models.LoadWorkspaceConfig()
	cfg.ServerPort = 18084
	cfg.AutoWatch = false
	ctrl := NewServerController()

	v := NewSlidesView(comp, cfg, ctrl)
	defer ctrl.Stop()

	deckPath := "decks/starter_deck"
	if _, err := os.Stat(deckPath); err != nil {
		deckPath = "../../decks/starter_deck"
	}
	deck, err := workspace.InspectDeck(deckPath)
	if err != nil {
		t.Fatalf("Failed to inspect test deck: %v", err)
	}
	v.SetDeck(deck)

	// Initially not running
	view := v.View()
	if strings.Contains(view, "LIVE: http://localhost:") {
		t.Fatal("Expected no LIVE badge initially")
	}

	// Turn ON with 's'
	v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	if !ctrl.IsRunning() {
		t.Fatalf("Expected server to start on 's' in SlidesView, status: %s", v.StatusMsg)
	}
	view = v.View()
	if !strings.Contains(view, "LIVE: http://localhost:18084") {
		t.Fatalf("Expected LIVE badge in SlidesView header, got: %s", view)
	}

	// Turn OFF with 's'
	v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	if ctrl.IsRunning() {
		t.Fatalf("Expected server to stop on second 's', status: %s", v.StatusMsg)
	}
	view = v.View()
	if strings.Contains(view, "LIVE: http://localhost:") {
		t.Fatal("Expected LIVE badge to be gone after turn off")
	}
}

func TestTUIThemeApplicationAndNavigation(t *testing.T) {
	app := NewAppModel(".")

	// Test 'n' key in DecksView routes to TabScaffold
	app.ActiveTab = TabDecks
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	if app.ActiveTab != TabScaffold {
		t.Fatalf("Expected ActiveTab to switch to TabScaffold on 'n', got %d", app.ActiveTab)
	}

	// Test 't' key in SlidesView routes to TabThemes
	app.ActiveTab = TabSlides
	app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	if app.ActiveTab != TabThemes {
		t.Fatalf("Expected ActiveTab to switch to TabThemes on 't', got %d", app.ActiveTab)
	}

	// Test applying theme to deck via ThemeStudioView
	tmpDir := t.TempDir()
	deckDir := filepath.Join(tmpDir, "tui_test_deck")
	tempDeck, err := scaffold.ScaffoldDeck(deckDir, "tui_test_deck", "cyber-dark", 2)
	if err != nil {
		t.Fatalf("Failed to scaffold temp deck: %v", err)
	}

	app.ThemeView.SetActiveDeck(tempDeck)
	// Select first theme (e.g. academic-crimson)
	app.ThemeView.ThemeCursor = 0
	targetTheme := app.ThemeView.Themes[0].Tokens.Name

	// Press 'a' to apply
	app.ThemeView.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

	refreshedDeck, err := workspace.InspectDeck(deckDir)
	if err != nil {
		t.Fatalf("Failed to re-inspect deck: %v", err)
	}
	if refreshedDeck.Config.Theme != targetTheme {
		t.Fatalf("Expected deck theme to be '%s', got '%s'", targetTheme, refreshedDeck.Config.Theme)
	}

	viewStr := app.ThemeView.View()
	if !strings.Contains(viewStr, "[ON]") {
		t.Fatalf("Expected [ON] badge in ThemeView, got: %s", viewStr)
	}
}

func TestTUILiveChangeDetection(t *testing.T) {
	tmpDir := t.TempDir()
	deckDir := filepath.Join(tmpDir, "live_sync_deck")
	deck, err := scaffold.ScaffoldDeck(deckDir, "live_sync_deck", "cyber-dark", 2)
	if err != nil {
		t.Fatalf("failed to scaffold temp deck: %v", err)
	}

	app := NewAppModel(tmpDir)
	app.SlidesView.SetDeck(deck)

	// First tick initializes baseline mtime
	app.Update(FileWatchTickMsg(time.Now()))

	// Simulate external modification to deck.json
	deckJSONPath := filepath.Join(deckDir, "deck.json")
	newConf := `{"title": "Updated Title", "theme": "paper-ink"}`
	if err := os.WriteFile(deckJSONPath, []byte(newConf), 0644); err != nil {
		t.Fatalf("failed to write updated deck.json: %v", err)
	}
	futureTime := time.Now().Add(2 * time.Second)
	_ = os.Chtimes(deckJSONPath, futureTime, futureTime)

	// Second tick detects modification and auto-reloads
	app.Update(FileWatchTickMsg(time.Now()))

	if app.SlidesView.CurrentDeck.Config.Theme != "paper-ink" {
		t.Fatalf("expected live sync to update theme to 'paper-ink', got '%s'", app.SlidesView.CurrentDeck.Config.Theme)
	}
	if !strings.Contains(app.SlidesView.StatusMsg, "Live sync") {
		t.Fatalf("expected Live sync status message, got: %s", app.SlidesView.StatusMsg)
	}
}

func TestDecksViewThemePicker(t *testing.T) {
	tmpDir := t.TempDir()
	deckDir := filepath.Join(tmpDir, "picker_deck")
	_, err := scaffold.ScaffoldDeck(deckDir, "picker_deck", "academic-crimson", 2)
	if err != nil {
		t.Fatalf("failed to scaffold deck: %v", err)
	}

	tm := theme.NewThemeManager(tmpDir)
	comp := compiler.NewCompiler(tm)
	cfg := models.LoadWorkspaceConfig()
	cfg.ActiveRoot = tmpDir
	srvCtrl := NewServerController()

	v := NewDecksView(tmpDir, nil, comp, cfg, srvCtrl)
	if len(v.Decks) == 0 {
		t.Fatalf("expected at least 1 deck discovered")
	}

	// 1. Press 't' to open theme picker modal
	v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	if !v.IsPickingTheme {
		t.Fatalf("expected IsPickingTheme to be true after pressing 't'")
	}
	if len(v.ThemeChoices) == 0 {
		t.Fatalf("expected ThemeChoices to be populated")
	}

	// 2. Navigate down with 'j'
	initialCursor := v.ThemeChoiceCursor
	v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	expectedCursor := (initialCursor + 1) % len(v.ThemeChoices)
	if v.ThemeChoiceCursor != expectedCursor {
		t.Fatalf("expected cursor %d, got %d", expectedCursor, v.ThemeChoiceCursor)
	}

	chosenTheme := v.ThemeChoices[v.ThemeChoiceCursor].Tokens.Name

	// 3. Press enter to apply theme
	v.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if v.IsPickingTheme {
		t.Fatalf("expected IsPickingTheme to be false after pressing Enter")
	}

	// Check deck in memory
	d := v.SelectedDeck()
	if d.ThemeName != chosenTheme {
		t.Fatalf("expected deck ThemeName to be '%s', got '%s'", chosenTheme, d.ThemeName)
	}

	// Check deck on disk
	inspected, err := workspace.InspectDeck(d.Path)
	if err != nil {
		t.Fatalf("failed to inspect deck: %v", err)
	}
	if inspected.Config.Theme != chosenTheme {
		t.Fatalf("expected deck.json theme '%s', got '%s'", chosenTheme, inspected.Config.Theme)
	}

	// Check View renders modal when active
	v.IsPickingTheme = true
	viewOut := v.View()
	if !strings.Contains(viewOut, "Select Theme for:") {
		t.Fatalf("expected Theme modal header in View output, got: %s", viewOut)
	}
}

func TestCWDRootResolution(t *testing.T) {
	tmpDir := t.TempDir()
	app := NewAppModel(tmpDir)

	absExpected, _ := filepath.Abs(tmpDir)
	if app.Config.ActiveRoot != absExpected {
		t.Fatalf("expected ActiveRoot to be '%s', got '%s'", absExpected, app.Config.ActiveRoot)
	}
	if app.DecksView.ActiveRoot != absExpected {
		t.Fatalf("expected DecksView.ActiveRoot to be '%s', got '%s'", absExpected, app.DecksView.ActiveRoot)
	}
}



