package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"time"

	"deckforge/internal/compiler"
	"github.com/fsnotify/fsnotify"
)

// DeckServer encapsulates the HTTP server and fsnotify watcher
type DeckServer struct {
	Port        int
	DeckPath    string
	Compiler    *compiler.Compiler
	Broadcaster *SSEBroadcaster
	server      *http.Server
	listener    net.Listener
	watcher     *fsnotify.Watcher
	running     bool
}

// NewDeckServer initializes a server instance
func NewDeckServer(deckPath string, port int, comp *compiler.Compiler) *DeckServer {
	if port <= 0 {
		port = 8080
	}
	return &DeckServer{
		Port:        port,
		DeckPath:    deckPath,
		Compiler:    comp,
		Broadcaster: NewSSEBroadcaster(),
	}
}

// IsRunning returns whether server is actively serving
func (s *DeckServer) IsRunning() bool {
	return s.running
}

// Start launches the HTTP server and watcher
func (s *DeckServer) Start(watch bool) error {
	// First compile the deck
	if _, err := s.Compiler.Build(s.DeckPath, ""); err != nil {
		return fmt.Errorf("initial compile failed: %w", err)
	}

	distDir := filepath.Join(s.DeckPath, "dist")
	fs := http.FileServer(http.Dir(distDir))

	mux := http.NewServeMux()
	registerAPIRoutes(mux, s)
	mux.Handle("/", fs)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: mux,
	}

	ln, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("failed to bind port %d: %w", s.Port, err)
	}
	s.listener = ln
	s.running = true

	if watch {
		if err := s.startWatcher(); err != nil {
			fmt.Printf("Warning: watcher failed to start: %v\n", err)
		}
	}

	httpServer := s.server
	go func() {
		if err := httpServer.Serve(ln); err != nil && err != http.ErrServerClosed {
			s.running = false
		}
	}()

	return nil
}

// Stop shuts down the server and watcher
func (s *DeckServer) Stop() error {
	s.running = false
	if s.watcher != nil {
		_ = s.watcher.Close()
		s.watcher = nil
	}
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := s.server.Shutdown(ctx)
		s.server = nil
		s.listener = nil
		return err
	}
	if s.listener != nil {
		_ = s.listener.Close()
		s.listener = nil
	}
	return nil
}

func (s *DeckServer) startWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	s.watcher = watcher

	slidesDir := filepath.Join(s.DeckPath, "slides")
	_ = watcher.Add(slidesDir)
	_ = watcher.Add(s.DeckPath)

	go func() {
		var lastRecompile time.Time
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
					// Debounce 250ms
					if time.Since(lastRecompile) > 250*time.Millisecond {
						lastRecompile = time.Now()
						if _, err := s.Compiler.Build(s.DeckPath, ""); err == nil && s.Broadcaster != nil {
							s.Broadcaster.Broadcast("reload", `{"reason": "fsnotify"}`)
						}
					}
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
			}
		}
	}()

	return nil
}

// ServeSync is a blocking runner for CLI 'serve' command
func (s *DeckServer) ServeSync(watch bool) error {
	if err := s.Start(watch); err != nil {
		return err
	}
	fmt.Printf("\nDeckForge Live Server Active\n")
	fmt.Printf("  -> URL: http://localhost:%d/\n", s.Port)
	fmt.Printf("  -> Serving from: %s/dist\n", s.DeckPath)
	if watch {
		fmt.Printf("  -> Watching slides/ and deck.json for instant reload\n")
	}
	fmt.Println("Press Ctrl+C to terminate.")

	// Block until interrupted
	select {}
}
