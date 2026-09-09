package models

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// WorkspaceConfig records recent presentation roots and user preferences
type WorkspaceConfig struct {
	ActiveRoot        string   `json:"activeRoot"`
	RecentRoots       []string `json:"recentRoots"`
	DefaultTheme      string   `json:"defaultTheme"`
	DefaultSlideCount int      `json:"defaultSlideCount"`
	PreferredEditor   string   `json:"preferredEditor,omitempty"`
	ServerPort        int      `json:"serverPort"`
	AutoWatch         bool     `json:"autoWatch"`
}

// ConfigFilePath returns the path to ~/.config/deckforge/config.json
func ConfigFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "deckforge_config.json"
	}
	return filepath.Join(home, ".config", "deckforge", "config.json")
}

// LoadWorkspaceConfig loads config from disk or returns sensible defaults
func LoadWorkspaceConfig() *WorkspaceConfig {
	cfg := &WorkspaceConfig{
		RecentRoots:       []string{},
		DefaultTheme:      "academic-crimson",
		DefaultSlideCount: 5,
		PreferredEditor:   "auto",
		ServerPort:        8080,
		AutoWatch:         true,
	}

	path := ConfigFilePath()
	data, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(data, cfg)
	}

	if cfg.ServerPort <= 0 {
		cfg.ServerPort = 8080
	}
	if cfg.DefaultSlideCount <= 0 {
		cfg.DefaultSlideCount = 5
	}
	if cfg.PreferredEditor == "" {
		cfg.PreferredEditor = "auto"
	}

	cwd, err := os.Getwd()
	if err == nil {
		if cfg.ActiveRoot == "" {
			cfg.ActiveRoot = cwd
		}
		cfg.AddRecentRoot(cwd)
	}

	return cfg
}

// Save writes the workspace config to ~/.config/deckforge/config.json
func (c *WorkspaceConfig) Save() error {
	path := ConfigFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// AddRecentRoot adds a root to RecentRoots with deduplication
func (c *WorkspaceConfig) AddRecentRoot(root string) {
	if root == "" {
		return
	}
	clean := filepath.Clean(root)
	newRoots := []string{clean}
	for _, r := range c.RecentRoots {
		if filepath.Clean(r) != clean {
			newRoots = append(newRoots, r)
		}
		if len(newRoots) >= 10 {
			break
		}
	}
	c.RecentRoots = newRoots
}
