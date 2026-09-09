package tui

import (
	"fmt"
	"path/filepath"
	"sync"

	"deckforge/internal/compiler"
	"deckforge/internal/models"
	"deckforge/internal/server"
)

// ServerController coordinates the live HTTP server lifecycle across TUI views
type ServerController struct {
	mu     sync.Mutex
	Server *server.DeckServer
}

// NewServerController instantiates an empty server controller
func NewServerController() *ServerController {
	return &ServerController{}
}

// IsRunning reports whether a deck server is actively listening
func (c *ServerController) IsRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Server != nil && c.Server.IsRunning()
}

// IsServing reports whether the server is actively serving the specified deck path
func (c *ServerController) IsServing(deckPath string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Server == nil || !c.Server.IsRunning() {
		return false
	}
	return filepath.Clean(c.Server.DeckPath) == filepath.Clean(deckPath)
}

// ActiveDeckPath returns the path of the deck currently being served
func (c *ServerController) ActiveDeckPath() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Server != nil && c.Server.IsRunning() {
		return c.Server.DeckPath
	}
	return ""
}

// ActivePort returns the port of the active server or fallback 8080
func (c *ServerController) ActivePort() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Server != nil && c.Server.Port > 0 {
		return c.Server.Port
	}
	return 8080
}

// Stop terminates any currently active server
func (c *ServerController) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Server != nil {
		err := c.Server.Stop()
		c.Server = nil
		return err
	}
	return nil
}

// Start launches a new server for the target deck path, terminating any existing instance
func (c *ServerController) Start(deckPath string, port int, watch bool, comp *compiler.Compiler) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Server != nil {
		_ = c.Server.Stop()
		c.Server = nil
	}

	srv := server.NewDeckServer(deckPath, port, comp)
	if err := srv.Start(watch); err != nil {
		return err
	}
	c.Server = srv
	return nil
}

// Toggle toggles the server state for the given deck:
// - If already serving this deck -> stops server (turns OFF)
// - If serving another deck or stopped -> starts server for this deck (turns ON)
func (c *ServerController) Toggle(deck *models.Deck, port int, watch bool, comp *compiler.Compiler) (bool, string, error) {
	if deck == nil {
		return false, "No deck selected.", nil
	}

	if port <= 0 {
		port = 8080
	}

	c.mu.Lock()
	isCurrent := c.Server != nil && c.Server.IsRunning() && filepath.Clean(c.Server.DeckPath) == filepath.Clean(deck.Path)
	if c.Server != nil {
		_ = c.Server.Stop()
		c.Server = nil
	}
	c.mu.Unlock()

	if isCurrent {
		return false, fmt.Sprintf("Server stopped for %s.", deck.Name), nil
	}

	srv := server.NewDeckServer(deck.Path, port, comp)
	if err := srv.Start(watch); err != nil {
		return false, fmt.Sprintf("Server error: %v", err), err
	}

	c.mu.Lock()
	c.Server = srv
	c.mu.Unlock()

	return true, fmt.Sprintf("Server running on http://localhost:%d/ (Serving: %s)", port, deck.Name), nil
}
