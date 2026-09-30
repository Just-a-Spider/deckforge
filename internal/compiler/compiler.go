package compiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"deckforge/internal/assets"
	"deckforge/internal/components"
	"deckforge/internal/theme"
	"deckforge/internal/workspace"
)

var embeddedAssetCache sync.Map

// Compiler compiles any deck directory into standalone 1080p HTML
type Compiler struct {
	ThemeManager *theme.ThemeManager
}

// NewCompiler creates a new compiler instance
func NewCompiler(tm *theme.ThemeManager) *Compiler {
	return &Compiler{
		ThemeManager: tm,
	}
}

// Build compiles the deck at deckPath into self-contained HTML files
func (c *Compiler) Build(deckPath, overrideTheme string) (string, error) {
	deck, err := workspace.InspectDeck(deckPath)
	if err != nil {
		return "", fmt.Errorf("failed to inspect deck at %s: %w", deckPath, err)
	}

	if len(deck.Slides) == 0 {
		return "", fmt.Errorf("no slides found in %s/slides", deckPath)
	}

	themeName := overrideTheme
	if themeName == "" {
		themeName = deck.ThemeName
	}
	if themeName == "" {
		themeName = "academic-crimson"
	}

	resolvedTheme, err := c.ThemeManager.ResolveTheme(themeName, deckPath)
	if err != nil {
		return "", fmt.Errorf("theme resolution error: %w", err)
	}

	themeCSS := resolvedTheme.CompileFullCSS()

	// Load core CSS
	stageCSS := c.loadAsset("core/stage.css")
	hudCSS := c.loadAsset("core/hud.css")
	compCSS := c.loadAsset("core/components.css")
	studioCSS := c.loadAsset("core/studio.css")

	// Inject custom component CSS if any custom components exist
	cm := components.NewComponentManager(deckPath)
	if customCompCSS := cm.CompileComponentCSS(deckPath); customCompCSS != "" {
		compCSS = compCSS + "\n\n/* --- Multi-Tier Custom Components --- */\n" + customCompCSS
	}

	// Load core JS
	captureJS := c.loadAsset("core/capture.js")
	studioJS := c.loadAsset("core/studio.js")
	tokensStudioJS := c.loadAsset("core/tokens_studio.js")
	controllerJS := c.loadAsset("core/controller.js")
	masterIndexJS := c.loadAsset("core/master_index.js")

	// Collect slides HTML
	var slideContents []string
	for _, s := range deck.Slides {
		data, err := os.ReadFile(s.Path)
		if err != nil {
			return "", fmt.Errorf("error reading slide %s: %w", s.Path, err)
		}
		raw := strings.TrimSpace(string(data))
		if strings.Contains(raw, "<section") && !strings.Contains(raw, "slide-backdrop") {
			if idx := strings.Index(raw, ">"); idx != -1 {
				raw = raw[:idx+1] + "\n        <div class=\"slide-backdrop\"></div>" + raw[idx+1:]
			}
		}
		slideContents = append(slideContents, raw)
	}
	allSlidesHTML := strings.Join(slideContents, "\n\n")

	title := deck.Config.Title
	if title == "" {
		title = deck.Name
	}

	fontLink := ""
	if resolvedTheme.Tokens.Fonts.ImportURL != "" {
		fontLink = fmt.Sprintf("\n    <link rel=\"stylesheet\" href=\"%s\">", resolvedTheme.Tokens.Fonts.ImportURL)
	}

	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s • DeckForge</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>%s
    <style>
/* --- Theme Segments (Fonts, Variables, Surfaces) --- */
%s

/* --- Core Stage Geometry --- */
%s

/* --- Core HUD & Modals --- */
%s

/* --- Core Components & Grids --- */
%s

/* --- DeckForge Studio Canvas Engine --- */
%s
    </style>
</head>
<body>
    <div class="deck-viewport">
        <!-- 1920x1080 MASTER CANVAS (PRISTINE SLIDE STAGE) -->
        <div class="deck-stage hairlines" id="deckStage">
            <div class="stage-backdrop"></div>

%s
        </div>

        <!-- UNIFIED FLOATING HUD DOCK -->
        <div class="deck-hud visible" id="deckHud">
            <div class="hud-group">
                <button class="hud-btn" id="hudPrevBtn" title="Previous Slide (← / PageUp)">
                    <svg viewBox="0 0 24 24"><polyline points="15 18 9 12 15 6"/></svg>
                </button>
                <span class="hud-counter" id="hudCounter">01 / %02d</span>
                <button class="hud-btn" id="hudNextBtn" title="Next Slide (→ / Space / PageDown)">
                    <svg viewBox="0 0 24 24"><polyline points="9 18 15 12 9 6"/></svg>
                </button>
            </div>
            <div class="hud-divider"></div>
            <div class="hud-group">
                <button class="hud-btn index-toggle-btn" id="indexToggleBtn" title="Master Index Drawer (M)">
                    <svg viewBox="0 0 24 24"><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
                    <span>INDEX</span>
                    <kbd>M</kbd>
                </button>
                <button class="hud-btn studio-toggle" id="hudStudioBtn" title="Toggle Studio Canvas (S)">
                    <svg viewBox="0 0 24 24"><rect x="3" y="3" width="18" height="18" rx="2"/><path d="M9 3v18"/><path d="M15 9h6"/></svg>
                    <span>STUDIO</span>
                    <kbd>S</kbd>
                </button>
                <button class="hud-btn edit-toggle" id="editToggle" title="Toggle Inline Edit (E)">
                    <svg viewBox="0 0 24 24"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
                    <span>EDIT</span>
                    <kbd>E</kbd>
                </button>
                <button class="hud-btn" id="hudCaptureBtn" title="Capture HD Slide (P)">
                    <svg viewBox="0 0 24 24"><path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/><circle cx="12" cy="13" r="4"/></svg>
                    <span>CAPTURE</span>
                    <kbd>P</kbd>
                </button>
                <button class="hud-btn" id="hudFullscreenBtn" title="Fullscreen (F)">
                    <svg viewBox="0 0 24 24"><path d="M8 3H5a2 2 0 0 0-2 2v3m18 0V5a2 2 0 0 0-2-2h-3m0 18h3a2 2 0 0 0 2-2v-3M3 16v3a2 2 0 0 0 2 2h3"/></svg>
                    <kbd>F</kbd>
                </button>
                <button class="hud-btn" id="hudHelpBtn" title="Keyboard Shortcuts (?)">
                    <svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="10"/><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
                    <kbd>?</kbd>
                </button>
            </div>
        </div>

        <!-- KEYBOARD SHORTCUTS HELP MODAL -->
        <div class="keyboard-help-modal" id="keyboardHelpModal">
            <div class="keyboard-help-card">
                <div class="keyboard-help-header">
                    <div class="keyboard-help-title">
                        <svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="10"/><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
                        <span>System Navigation Shortcuts</span>
                    </div>
                    <button class="keyboard-help-close" id="keyboardHelpCloseBtn" title="Close (Esc / ?)">
                        <svg viewBox="0 0 24 24"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                    </button>
                </div>
                <div class="keyboard-help-grid">
                    <div class="shortcut-row">
                        <div class="shortcut-keys"><kbd>→</kbd> <kbd>Space</kbd> <kbd>PgDn</kbd></div>
                        <div class="shortcut-desc">Advance to next slide</div>
                    </div>
                    <div class="shortcut-row">
                        <div class="shortcut-keys"><kbd>←</kbd> <kbd>PgUp</kbd></div>
                        <div class="shortcut-desc">Return to previous slide</div>
                    </div>
                    <div class="shortcut-row">
                        <div class="shortcut-keys"><kbd>M</kbd></div>
                        <div class="shortcut-desc">Toggle Master Index drawer</div>
                    </div>
                    <div class="shortcut-row">
                        <div class="shortcut-keys"><kbd>S</kbd></div>
                        <div class="shortcut-desc">Toggle Studio Canvas mode</div>
                    </div>
                    <div class="shortcut-row">
                        <div class="shortcut-keys"><kbd>E</kbd></div>
                        <div class="shortcut-desc">Toggle live inline editing mode</div>
                    </div>
                    <div class="shortcut-row">
                        <div class="shortcut-keys"><kbd>P</kbd></div>
                        <div class="shortcut-desc">Capture slide in 4K UHD PNG</div>
                    </div>
                    <div class="shortcut-row">
                        <div class="shortcut-keys"><kbd>F</kbd></div>
                        <div class="shortcut-desc">Toggle fullscreen</div>
                    </div>
                    <div class="shortcut-row">
                        <div class="shortcut-keys"><kbd>?</kbd> or <kbd>H</kbd></div>
                        <div class="shortcut-desc">Show / hide this shortcuts guide</div>
                    </div>
                    <div class="shortcut-row">
                        <div class="shortcut-keys"><kbd>Esc</kbd></div>
                        <div class="shortcut-desc">Close modals, drawers and exit edit mode</div>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <script>
%s

%s

%s

%s

%s

        window.addEventListener('DOMContentLoaded', () => {
            const presentation = new SlidePresentation();
            const masterIndex = new MasterIndex(presentation);
        });
    </script>
</body>
</html>
`,
		title,
		fontLink,
		themeCSS,
		stageCSS,
		hudCSS,
		compCSS,
		studioCSS,
		allSlidesHTML,
		len(deck.Slides),
		captureJS,
		studioJS,
		tokensStudioJS,
		controllerJS,
		masterIndexJS,
	)

	// Save to deck dist directory
	deckDist := filepath.Join(deck.Path, "dist")
	if err := os.MkdirAll(deckDist, 0755); err != nil {
		return "", err
	}

	outFile := filepath.Join(deckDist, fmt.Sprintf("%s.html", deck.Name))
	indexFile := filepath.Join(deckDist, "index.html")

	if err := os.WriteFile(outFile, []byte(htmlContent), 0644); err != nil {
		return "", err
	}
	if err := os.WriteFile(indexFile, []byte(htmlContent), 0644); err != nil {
		return "", err
	}

	return outFile, nil
}

func (c *Compiler) loadAsset(relPath string) string {
	// 1. If WorkspaceRoot is set, check disk first
	if c.ThemeManager != nil && c.ThemeManager.WorkspaceRoot != "" {
		diskPath := filepath.Join(c.ThemeManager.WorkspaceRoot, "engine", relPath)
		if content, err := os.ReadFile(diskPath); err == nil {
			return string(content)
		}
	}

	// 2. Check embedded asset cache
	if cached, ok := embeddedAssetCache.Load(relPath); ok {
		return cached.(string)
	}

	// 3. Fall back to embedded asset and memoize
	embedded, err := assets.LoadCoreAsset(relPath)
	if err == nil {
		embeddedAssetCache.Store(relPath, embedded)
		return embedded
	}

	return fmt.Sprintf("/* Missing asset: %s */", relPath)
}
