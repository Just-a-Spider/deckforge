package models

// DeckConfig represents the JSON metadata in deck.json
type DeckConfig struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle,omitempty"`
	Theme    string `json:"theme"`
	Author   string `json:"author,omitempty"`
}

// SlideInfo represents a single discovered slide file
type SlideInfo struct {
	Index    int    `json:"index"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Title    string `json:"title"`
	Slug     string `json:"slug"`
}

// Deck represents an active or discovered slide deck
type Deck struct {
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Config    DeckConfig  `json:"config"`
	Slides    []SlideInfo `json:"slides"`
	ThemeName string      `json:"theme_name"`
}
