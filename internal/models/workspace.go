package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

// ConfigDirPath returns the base directory for global deckforge configuration
func ConfigDirPath() string {
	if custom := os.Getenv("DECKFORGE_CONFIG_DIR"); custom != "" {
		return custom
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "deckforge")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".deckforge"
	}
	return filepath.Join(home, ".config", "deckforge")
}

// ConfigFilePath returns the path to config.json
func ConfigFilePath() string {
	return filepath.Join(ConfigDirPath(), "config.json")
}

// GlobalComponentsDir returns the path to ~/.config/deckforge/components
func GlobalComponentsDir() string {
	return filepath.Join(ConfigDirPath(), "components")
}

// GlobalThemesDir returns the path to ~/.config/deckforge/themes
func GlobalThemesDir() string {
	return filepath.Join(ConfigDirPath(), "themes")
}

// WorkspaceDotDir returns the path to <root>/.deckforge
func WorkspaceDotDir(root string) string {
	if root == "" {
		root = "."
	}
	return filepath.Join(root, ".deckforge")
}

// WorkspaceConfigFilePath returns the path to <root>/.deckforge/config.json
func WorkspaceConfigFilePath(root string) string {
	return filepath.Join(WorkspaceDotDir(root), "config.json")
}

// WorkspaceComponentsDir returns the path to <root>/.deckforge/components
func WorkspaceComponentsDir(root string) string {
	return filepath.Join(WorkspaceDotDir(root), "components")
}

// WorkspaceThemesDir returns the path to <root>/.deckforge/themes
func WorkspaceThemesDir(root string) string {
	return filepath.Join(WorkspaceDotDir(root), "themes")
}

// WorkspaceDecksDir returns the path to <root>/.deckforge/decks
func WorkspaceDecksDir(root string) string {
	return filepath.Join(WorkspaceDotDir(root), "decks")
}

// LoadWorkspaceConfig loads config from disk or returns sensible defaults, overlaying local .deckforge/config.json if present
func LoadWorkspaceConfig() *WorkspaceConfig {
	return LoadWorkspaceConfigForRoot("")
}

// LoadWorkspaceConfigForRoot loads config with optional workspace root override
func LoadWorkspaceConfigForRoot(root string) *WorkspaceConfig {
	cfg := &WorkspaceConfig{
		RecentRoots:       []string{},
		DefaultTheme:      "academic-crimson",
		DefaultSlideCount: 5,
		PreferredEditor:   "auto",
		ServerPort:        8080,
		AutoWatch:         true,
	}

	// 1. Load global config
	path := ConfigFilePath()
	data, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(data, cfg)
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	targetRoot := root
	if targetRoot == "" {
		targetRoot = cwd
	}
	cfg.ActiveRoot = targetRoot
	cfg.AddRecentRoot(targetRoot)

	// 2. Overlay workspace local config: <targetRoot>/.deckforge/config.json
	if targetRoot != "" {
		localPath := WorkspaceConfigFilePath(targetRoot)
		if localData, err := os.ReadFile(localPath); err == nil {
			var localCfg WorkspaceConfig
			if err := json.Unmarshal(localData, &localCfg); err == nil {
				if localCfg.PreferredEditor != "" {
					cfg.PreferredEditor = localCfg.PreferredEditor
				}
				if localCfg.DefaultTheme != "" {
					cfg.DefaultTheme = localCfg.DefaultTheme
				}
				if localCfg.DefaultSlideCount > 0 {
					cfg.DefaultSlideCount = localCfg.DefaultSlideCount
				}
				if localCfg.ServerPort > 0 {
					cfg.ServerPort = localCfg.ServerPort
				}
				var rawMap map[string]interface{}
				if err := json.Unmarshal(localData, &rawMap); err == nil {
					if _, hasWatch := rawMap["autoWatch"]; hasWatch {
						cfg.AutoWatch = localCfg.AutoWatch
					}
				}
			}
		}
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

	return cfg
}

// Save writes the workspace config to config.json
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

// SaveLocal writes the local configuration overrides to <root>/.deckforge/config.json
func (c *WorkspaceConfig) SaveLocal(root string) error {
	if root == "" {
		root = c.ActiveRoot
	}
	if root == "" {
		var err error
		root, err = os.Getwd()
		if err != nil {
			root = "."
		}
	}
	localPath := WorkspaceConfigFilePath(root)
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(localPath, data, 0644)
}

// Get retrieves a setting by name
func (c *WorkspaceConfig) Get(key string) (string, error) {
	switch strings.ToLower(key) {
	case "preferrededitor", "editor":
		return c.PreferredEditor, nil
	case "defaulttheme", "theme":
		return c.DefaultTheme, nil
	case "defaultslidecount", "slides", "slidecount":
		return strconv.Itoa(c.DefaultSlideCount), nil
	case "serverport", "port":
		return strconv.Itoa(c.ServerPort), nil
	case "autowatch", "watch":
		return strconv.FormatBool(c.AutoWatch), nil
	case "activeroot", "root":
		return c.ActiveRoot, nil
	default:
		return "", fmt.Errorf("unknown config key: %s", key)
	}
}

// ApplyValue validates and updates a configuration property in memory without writing to disk
func (c *WorkspaceConfig) ApplyValue(key, value string) error {
	switch strings.ToLower(key) {
	case "preferrededitor", "editor":
		if strings.TrimSpace(value) == "" {
			c.PreferredEditor = "auto"
		} else {
			c.PreferredEditor = strings.TrimSpace(value)
		}
	case "defaulttheme", "theme":
		clean := strings.TrimSpace(value)
		if clean == "" {
			return fmt.Errorf("defaultTheme cannot be empty")
		}
		c.DefaultTheme = clean
	case "defaultslidecount", "slides", "slidecount":
		n, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil || n < 1 || n > 50 {
			return fmt.Errorf("defaultSlideCount must be an integer between 1 and 50")
		}
		c.DefaultSlideCount = n
	case "serverport", "port":
		p, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil || p < 1 || p > 65535 {
			return fmt.Errorf("serverPort must be a valid port number (1-65535)")
		}
		c.ServerPort = p
	case "autowatch", "watch":
		v := strings.ToLower(strings.TrimSpace(value))
		if v == "true" || v == "1" || v == "yes" || v == "on" {
			c.AutoWatch = true
		} else if v == "false" || v == "0" || v == "no" || v == "off" {
			c.AutoWatch = false
		} else {
			return fmt.Errorf("autoWatch must be true or false")
		}
	case "activeroot", "root":
		clean := filepath.Clean(strings.TrimSpace(value))
		if clean == "" {
			return fmt.Errorf("activeRoot cannot be empty")
		}
		c.ActiveRoot = clean
		c.AddRecentRoot(clean)
	default:
		return fmt.Errorf("unknown config key: %s (supported: preferredEditor, defaultTheme, defaultSlideCount, serverPort, autoWatch, activeRoot)", key)
	}
	return nil
}

// Set updates a setting globally and persists to ~/.config/deckforge/config.json
func (c *WorkspaceConfig) Set(key, value string) error {
	if err := c.ApplyValue(key, value); err != nil {
		return err
	}
	return c.Save()
}

// SetLocal updates a setting locally and persists to <root>/.deckforge/config.json
func (c *WorkspaceConfig) SetLocal(key, value, root string) error {
	if err := c.ApplyValue(key, value); err != nil {
		return err
	}
	return c.SaveLocal(root)
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
