package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"deckforge/internal/compiler"
	"deckforge/internal/models"
	"deckforge/internal/scaffold"
	"deckforge/internal/theme"
	"deckforge/internal/workspace"
)

var (
	contentEditableRegex = regexp.MustCompile(`\s*contenteditable(=["'][^"']*["'])?`)
	dfSelectedAttrRegex  = regexp.MustCompile(`\s*df-selected(=["'][^"']*["'])?`)
	dataStudioRegex      = regexp.MustCompile(`\s*data-studio-[a-z0-9_-]+(=["'][^"']*["'])?`)
	runtimeClassRegex    = regexp.MustCompile(`\b(df-selected-node|df-hovered|df-editable-active)\b`)
	emptyClassRegex      = regexp.MustCompile(`\s*class=["']\s*["']`)
	multiSpaceRegex      = regexp.MustCompile(`[ \t]+`)
)

// SanitizeSlideHTML strips temporary Studio/editor attributes and classes
func SanitizeSlideHTML(raw string) string {
	cleaned := contentEditableRegex.ReplaceAllString(raw, "")
	cleaned = dfSelectedAttrRegex.ReplaceAllString(cleaned, "")
	cleaned = dataStudioRegex.ReplaceAllString(cleaned, "")

	// Clean class attributes
	cleaned = runtimeClassRegex.ReplaceAllString(cleaned, "")
	cleaned = emptyClassRegex.ReplaceAllString(cleaned, "")

	// Clean up consecutive spaces inside tags
	lines := strings.Split(cleaned, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t\r")
	}
	res := strings.Join(lines, "\n")
	res = strings.TrimSpace(res) + "\n"
	return res
}

// SSEBroadcaster manages active Server-Sent Events subscriber connections
type SSEBroadcaster struct {
	mu      sync.Mutex
	clients map[chan string]bool
}

func NewSSEBroadcaster() *SSEBroadcaster {
	return &SSEBroadcaster{
		clients: make(map[chan string]bool),
	}
}

func (b *SSEBroadcaster) Subscribe() chan string {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan string, 16)
	b.clients[ch] = true
	return ch
}

func (b *SSEBroadcaster) Unsubscribe(ch chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.clients[ch] {
		delete(b.clients, ch)
		close(ch)
	}
}

func (b *SSEBroadcaster) Broadcast(event, data string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	msg := fmt.Sprintf("event: %s\ndata: %s\n\n", event, data)
	for ch := range b.clients {
		select {
		case ch <- msg:
		default:
		}
	}
}

// SlideUpdateRequest payload
type SlideUpdateRequest struct {
	Index int    `json:"index"`
	HTML  string `json:"html"`
}

// SlideReorderRequest payload
type SlideReorderRequest struct {
	NewOrder []int `json:"newOrder"`
}

// SlideCreateRequest payload
type SlideCreateRequest struct {
	Title  string `json:"title"`
	Layout string `json:"layout"`
}

// TokenUpdateRequest payload
type TokenUpdateRequest struct {
	Palette models.ThemePalette `json:"palette"`
}

func registerAPIRoutes(mux *http.ServeMux, s *DeckServer) {
	mux.HandleFunc("/api/deck", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		deck, err := workspace.InspectDeck(s.DeckPath)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(deck)
	})

	mux.HandleFunc("/api/catalog", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		components := models.GetBuiltinComponents()
		_ = json.NewEncoder(w).Encode(components)
	})

	mux.HandleFunc("/api/styling-guide", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		guide := models.GetStylingGuide()
		_ = json.NewEncoder(w).Encode(guide)
	})

	mux.HandleFunc("/api/themes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var themes []models.Theme
		if s.Compiler != nil && s.Compiler.ThemeManager != nil {
			themes = s.Compiler.ThemeManager.ListThemes(s.DeckPath)
		} else {
			tm := theme.NewThemeManager(s.DeckPath)
			themes = tm.ListThemes(s.DeckPath)
		}

		type ThemeSummary struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
			Description string `json:"description"`
			Vibe        string `json:"vibe"`
			IsCustom    bool   `json:"isCustom"`
			Scope       string `json:"scope"`
		}

		var summaries []ThemeSummary
		for _, t := range themes {
			summaries = append(summaries, ThemeSummary{
				Name:        t.Tokens.Name,
				DisplayName: t.Tokens.DisplayName,
				Description: t.Tokens.Description,
				Vibe:        t.Tokens.Vibe,
				IsCustom:    t.IsCustom,
				Scope:       t.Scope,
			})
		}
		_ = json.NewEncoder(w).Encode(summaries)
	})

	mux.HandleFunc("/api/deck/theme", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		type ThemeReq struct {
			Theme string `json:"theme"`
		}
		var req ThemeReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Theme == "" {
			http.Error(w, `{"error": "invalid theme name"}`, http.StatusBadRequest)
			return
		}

		if err := workspace.SetDeckTheme(s.DeckPath, req.Theme); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		_, _ = s.Compiler.Build(s.DeckPath, req.Theme)
		s.Broadcaster.Broadcast("reload", fmt.Sprintf(`{"reason": "theme_changed", "theme": "%s"}`, req.Theme))

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "ok", "theme": "%s"}`, req.Theme)
	})

	mux.HandleFunc("/api/deck/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		type CreateReq struct {
			Name   string `json:"name"`
			Theme  string `json:"theme"`
			Slides int    `json:"slides"`
			Path   string `json:"path"`
		}
		var req CreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "invalid json"}`, http.StatusBadRequest)
			return
		}
		if req.Name == "" {
			req.Name = "new_presentation"
		}
		if req.Theme == "" {
			req.Theme = "cyber-dark"
		}
		if req.Slides <= 0 {
			req.Slides = 5
		}

		targetPath := req.Path
		if targetPath == "" {
			parent := filepath.Dir(s.DeckPath)
			if filepath.Base(parent) == "decks" {
				targetPath = filepath.Join(parent, req.Name)
			} else {
				targetPath = filepath.Join("decks", req.Name)
			}
		}

		deck, err := scaffold.ScaffoldDeck(targetPath, req.Name, req.Theme, req.Slides)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		tm := theme.NewThemeManager(targetPath)
		comp := compiler.NewCompiler(tm)
		_, _ = comp.Build(targetPath, req.Theme)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(deck)
	})

	mux.HandleFunc("/api/slides/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		var req SlideUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "invalid json: %s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		deck, err := workspace.InspectDeck(s.DeckPath)
		if err != nil || req.Index < 1 || req.Index > len(deck.Slides) {
			http.Error(w, `{"error": "slide index out of bounds"}`, http.StatusBadRequest)
			return
		}

		targetSlide := deck.Slides[req.Index-1]
		cleaned := SanitizeSlideHTML(req.HTML)

		if err := os.WriteFile(targetSlide.Path, []byte(cleaned), 0644); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "failed to write slide: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		// Rebuild
		_, _ = s.Compiler.Build(s.DeckPath, "")
		s.Broadcaster.Broadcast("slide_updated", fmt.Sprintf(`{"index": %d}`, req.Index))

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "ok", "index": %d}`, req.Index)
	})

	mux.HandleFunc("/api/slides/reorder", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		var req SlideReorderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "invalid json: %s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		if err := workspace.ReorderSlides(s.DeckPath, req.NewOrder); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "reorder failed: %s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		_, _ = s.Compiler.Build(s.DeckPath, "")
		s.Broadcaster.Broadcast("reload", `{"reason": "reorder"}`)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status": "ok"}`)
	})

	mux.HandleFunc("/api/slides/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		var req SlideCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "invalid json"}`, http.StatusBadRequest)
			return
		}
		if req.Title == "" {
			req.Title = "New Slide"
		}
		if req.Layout == "" {
			req.Layout = "2col"
		}

		slide, err := scaffold.AddSlide(s.DeckPath, req.Title, req.Layout)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "failed to add slide: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		_, _ = s.Compiler.Build(s.DeckPath, "")
		s.Broadcaster.Broadcast("reload", `{"reason": "create"}`)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(slide)
	})

	mux.HandleFunc("/api/slides/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodDelete {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		type DeleteReq struct {
			Index int `json:"index"`
		}
		var req DeleteReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "invalid json"}`, http.StatusBadRequest)
			return
		}

		deck, err := workspace.InspectDeck(s.DeckPath)
		if err != nil || req.Index < 1 || req.Index > len(deck.Slides) {
			http.Error(w, `{"error": "slide index out of bounds"}`, http.StatusBadRequest)
			return
		}

		_ = os.Remove(deck.Slides[req.Index-1].Path)

		// Re-read and re-sequence
		deckAfter, _ := workspace.InspectDeck(s.DeckPath)
		if len(deckAfter.Slides) > 0 {
			var newOrder []int
			for i := 1; i <= len(deckAfter.Slides); i++ {
				newOrder = append(newOrder, i)
			}
			_ = workspace.ReorderSlides(s.DeckPath, newOrder)
		}

		_, _ = s.Compiler.Build(s.DeckPath, "")
		s.Broadcaster.Broadcast("reload", `{"reason": "delete"}`)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status": "ok"}`)
	})

	mux.HandleFunc("/api/tokens/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		var req TokenUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "invalid json"}`, http.StatusBadRequest)
			return
		}

		deck, err := workspace.InspectDeck(s.DeckPath)
		if err != nil {
			http.Error(w, `{"error": "failed to inspect deck"}`, http.StatusInternalServerError)
			return
		}

		resolved, err := s.Compiler.ThemeManager.ResolveTheme(deck.ThemeName, s.DeckPath)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "resolve theme failed: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		// Update palette
		resolved.Tokens.Palette = req.Palette
		themeDir := resolved.Dir
		if themeDir == "" {
			themeDir = filepath.Join(s.DeckPath, "themes", resolved.Tokens.Name)
		}
		if err := theme.SaveSegmentedTheme(themeDir, resolved); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "save theme failed: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		_, _ = s.Compiler.Build(s.DeckPath, "")
		s.Broadcaster.Broadcast("reload", `{"reason": "tokens"}`)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status": "ok"}`)
	})

	mux.HandleFunc("/api/tokens/audit", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		deck, err := workspace.InspectDeck(s.DeckPath)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		resolved, err := s.Compiler.ThemeManager.ResolveTheme(deck.ThemeName, s.DeckPath)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		reports := models.EvaluateThemeContrast(resolved.Tokens.Palette)
		_ = json.NewEncoder(w).Encode(reports)
	})

	mux.HandleFunc("/api/tokens/autofix", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		deck, err := workspace.InspectDeck(s.DeckPath)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		resolved, err := s.Compiler.ThemeManager.ResolveTheme(deck.ThemeName, s.DeckPath)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		reports := models.EvaluateThemeContrast(resolved.Tokens.Palette)
		fixedCount := 0

		for _, rep := range reports {
			if rep.RecommendedAction != nil {
				rec := rep.RecommendedAction
				switch rec.TargetToken {
				case "textPrimary":
					resolved.Tokens.Palette.TextPrimary = rec.SuggestedHex
					fixedCount++
				case "textSecondary":
					resolved.Tokens.Palette.TextSecondary = rec.SuggestedHex
					fixedCount++
				case "accentPrimary":
					resolved.Tokens.Palette.AccentPrimary = rec.SuggestedHex
					fixedCount++
				}
			}
		}

		if fixedCount > 0 {
			themeDir := resolved.Dir
			if themeDir == "" {
				themeDir = filepath.Join(s.DeckPath, "themes", resolved.Tokens.Name)
			}
			_ = theme.SaveSegmentedTheme(themeDir, resolved)
			_, _ = s.Compiler.Build(s.DeckPath, "")
			s.Broadcaster.Broadcast("reload", `{"reason": "tokens_autofix"}`)
		}

		updatedReports := models.EvaluateThemeContrast(resolved.Tokens.Palette)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"fixed_count": fixedCount,
			"reports":     updatedReports,
			"palette":     resolved.Tokens.Palette,
		})
	})

	mux.HandleFunc("/api/components", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		comps := models.GetBuiltinComponents()
		_ = json.NewEncoder(w).Encode(comps)
	})

	mux.HandleFunc("/api/components/render", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		type RenderReq struct {
			Selector string                 `json:"selector"`
			Inputs   map[string]interface{} `json:"inputs"`
			Slots    map[string]string      `json:"slots"`
		}
		var req RenderReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "invalid json"}`, http.StatusBadRequest)
			return
		}

		var targetComp *models.ComponentDefinition
		for _, c := range models.GetBuiltinComponents() {
			if c.Selector == req.Selector {
				targetComp = &c
				break
			}
		}
		if targetComp == nil {
			http.Error(w, `{"error": "unknown component selector"}`, http.StatusBadRequest)
			return
		}

		html, err := targetComp.Render(req.Inputs, req.Slots)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"html": html})
	})

	mux.HandleFunc("/api/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		ch := s.Broadcaster.Subscribe()
		defer s.Broadcaster.Unsubscribe(ch)

		// Initial connection event
		fmt.Fprintf(w, "event: connected\ndata: {\"deck\": \"%s\"}\n\n", filepath.Base(s.DeckPath))
		flusher.Flush()

		notify := r.Context().Done()
		for {
			select {
			case <-notify:
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				fmt.Fprint(w, msg)
				flusher.Flush()
			}
		}
	})
}
