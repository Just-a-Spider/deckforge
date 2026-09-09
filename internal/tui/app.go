package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"deckforge/internal/compiler"
	"deckforge/internal/models"
	"deckforge/internal/theme"
	"deckforge/internal/workspace"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TabIndex int

const (
	TabDecks TabIndex = iota
	TabSlides
	TabScaffold
	TabThemes
	TabSettings
	TabCount
)

type FileWatchTickMsg time.Time

func tickEvery(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return FileWatchTickMsg(t)
	})
}

type AppModel struct {
	ActiveTab       TabIndex
	Config          *models.WorkspaceConfig
	ThemeManager    *theme.ThemeManager
	Compiler        *compiler.Compiler
	ServerCtrl      *ServerController
	DecksView       *DecksView
	SlidesView      *SlidesView
	ScaffoldView    *ScaffoldView
	ThemeView       *ThemeStudioView
	SettingsView    *SettingsView
	Width           int
	Height          int
	LastDeckMtime   time.Time
	LastSlidesMtime time.Time
}

func NewAppModel(initialPath string) *AppModel {
	cfg := models.LoadWorkspaceConfig()
	if initialPath != "" {
		if abs, err := filepath.Abs(initialPath); err == nil {
			initialPath = abs
		}
		cfg.ActiveRoot = initialPath
		cfg.AddRecentRoot(initialPath)
		_ = cfg.Save()
	}

	tm := theme.NewThemeManager(cfg.ActiveRoot)
	comp := compiler.NewCompiler(tm)
	srvCtrl := NewServerController()

	decksView := NewDecksView(cfg.ActiveRoot, cfg.RecentRoots, comp, cfg, srvCtrl)
	slidesView := NewSlidesView(comp, cfg, srvCtrl)
	scaffoldView := NewScaffoldView(cfg.ActiveRoot)
	themeView := NewThemeStudioView(tm, cfg.ActiveRoot, cfg, comp)
	settingsView := NewSettingsView(cfg)

	// If decks are found, initialize slides view and theme view with the first deck
	if len(decksView.Decks) > 0 {
		slidesView.SetDeck(decksView.Decks[0])
		themeView.SetActiveDeck(decksView.Decks[0])
	}

	return &AppModel{
		ActiveTab:    TabDecks,
		Config:       cfg,
		ThemeManager: tm,
		Compiler:     comp,
		ServerCtrl:   srvCtrl,
		DecksView:    decksView,
		SlidesView:   slidesView,
		ScaffoldView: scaffoldView,
		ThemeView:    themeView,
		SettingsView: settingsView,
	}
}

func (m *AppModel) Init() tea.Cmd {
	return tickEvery(1500 * time.Millisecond)
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case FileWatchTickMsg:
		if m.SlidesView != nil && m.SlidesView.CurrentDeck != nil {
			deckPath := m.SlidesView.CurrentDeck.Path
			deckJSONPath := filepath.Join(deckPath, "deck.json")
			slidesPath := filepath.Join(deckPath, "slides")

			changed := false
			if fi, err := os.Stat(deckJSONPath); err == nil {
				if m.LastDeckMtime.IsZero() {
					m.LastDeckMtime = fi.ModTime()
				} else if fi.ModTime().After(m.LastDeckMtime) {
					m.LastDeckMtime = fi.ModTime()
					changed = true
				}
			}

			if fi, err := os.Stat(slidesPath); err == nil {
				if m.LastSlidesMtime.IsZero() {
					m.LastSlidesMtime = fi.ModTime()
				} else if fi.ModTime().After(m.LastSlidesMtime) {
					m.LastSlidesMtime = fi.ModTime()
					changed = true
				}
			}

			if changed {
				if refreshed, err := workspace.InspectDeck(deckPath); err == nil {
					m.SlidesView.SetDeck(refreshed)
					m.ThemeView.SetActiveDeck(refreshed)
					m.DecksView.RefreshDecks()
					m.SlidesView.StatusMsg = "Live sync: external file change reloaded"
					m.SlidesView.IsError = false
				}
			}
		}
		return m, tickEvery(1500 * time.Millisecond)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.SlidesView.SetDimensions(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" {
			if m.ServerCtrl != nil {
				_ = m.ServerCtrl.Stop()
			}
			return m, tea.Quit
		}

		// Universal bracket tab switching (works everywhere)
		switch msg.String() {
		case "[", "shift+left":
			m.prevTab()
			return m, nil
		case "]", "shift+right":
			m.nextTab()
			return m, nil
		}

		// When not in text input modes, allow direct tab switching and arrow navigation
		if !m.isTextInputActive() {
			switch msg.String() {
			case "q":
				if m.ServerCtrl != nil {
					_ = m.ServerCtrl.Stop()
				}
				return m, tea.Quit
			case "1":
				m.ActiveTab = TabDecks
				return m, nil
			case "2":
				m.ActiveTab = TabSlides
				return m, nil
			case "3":
				m.ActiveTab = TabScaffold
				return m, nil
			case "4":
				m.ActiveTab = TabThemes
				if m.DecksView != nil && m.DecksView.SelectedDeck() != nil {
					m.ThemeView.SetActiveDeck(m.DecksView.SelectedDeck())
				}
				return m, nil
			case "5":
				m.ActiveTab = TabSettings
				return m, nil

			case "left":
				if m.ActiveTab == TabScaffold && m.ScaffoldView.FocusIndex == scaffoldFieldTheme {
					break
				}
				m.prevTab()
				return m, nil

			case "right":
				if m.ActiveTab == TabScaffold && m.ScaffoldView.FocusIndex == scaffoldFieldTheme {
					break
				}
				m.nextTab()
				return m, nil
			}
		}
	}

	// Route update to active view
	var cmd tea.Cmd
	switch m.ActiveTab {
	case TabDecks:
		var switchToSlides bool
		cmd, switchToSlides = m.DecksView.Update(msg)
		if m.DecksView.NavigateToTab == TabScaffold {
			m.DecksView.NavigateToTab = TabDecks
			m.ActiveTab = TabScaffold
			return m, cmd
		}
		if switchToSlides {
			selected := m.DecksView.SelectedDeck()
			if selected != nil {
				m.SlidesView.SetDeck(selected)
				m.ThemeView.SetActiveDeck(selected)
				m.ActiveTab = TabSlides
			}
		}
		// If root changed, propagate
		if m.DecksView.ActiveRoot != m.Config.ActiveRoot {
			if m.ServerCtrl != nil {
				_ = m.ServerCtrl.Stop()
			}
			m.Config.ActiveRoot = m.DecksView.ActiveRoot
			m.Config.AddRecentRoot(m.DecksView.ActiveRoot)
			_ = m.Config.Save()
			m.ScaffoldView.SetActiveRoot(m.DecksView.ActiveRoot)
			m.ThemeManager.WorkspaceRoot = m.DecksView.ActiveRoot
			m.ThemeView.ActiveWorkspace = m.DecksView.ActiveRoot
			m.ThemeView.RefreshThemes()
		}

	case TabSlides:
		cmd = m.SlidesView.Update(msg)
		if m.SlidesView.NavigateToTab == TabThemes {
			m.SlidesView.NavigateToTab = TabSlides
			m.ThemeView.SetActiveDeck(m.SlidesView.CurrentDeck)
			m.ActiveTab = TabThemes
			return m, cmd
		}

	case TabScaffold:
		var created bool
		cmd, created = m.ScaffoldView.Update(msg)
		if created && m.ScaffoldView.CreatedDeck != nil {
			m.DecksView.RefreshDecks()
			m.SlidesView.SetDeck(m.ScaffoldView.CreatedDeck)
			m.ThemeView.SetActiveDeck(m.ScaffoldView.CreatedDeck)
			m.ActiveTab = TabSlides
		}

	case TabThemes:
		cmd = m.ThemeView.Update(msg)
		if m.ThemeView.ActiveDeck != nil {
			if m.SlidesView.CurrentDeck != nil && m.ThemeView.ActiveDeck.Path == m.SlidesView.CurrentDeck.Path {
				m.SlidesView.CurrentDeck.Config.Theme = m.ThemeView.ActiveDeck.Config.Theme
				m.SlidesView.CurrentDeck.ThemeName = m.ThemeView.ActiveDeck.ThemeName
			}
			m.DecksView.RefreshDecks()
		}

	case TabSettings:
		cmd = m.SettingsView.Update(msg)
	}

	return m, cmd
}

func (m *AppModel) prevTab() {
	if m.ActiveTab > 0 {
		m.ActiveTab--
	} else {
		m.ActiveTab = TabCount - 1
	}
	if m.ActiveTab == TabThemes && m.DecksView != nil && m.DecksView.SelectedDeck() != nil {
		m.ThemeView.SetActiveDeck(m.DecksView.SelectedDeck())
	}
}

func (m *AppModel) nextTab() {
	m.ActiveTab = (m.ActiveTab + 1) % TabCount
	if m.ActiveTab == TabThemes && m.DecksView != nil && m.DecksView.SelectedDeck() != nil {
		m.ThemeView.SetActiveDeck(m.DecksView.SelectedDeck())
	}
}

func (m *AppModel) isTextInputActive() bool {
	if m.ActiveTab == TabDecks && (m.DecksView.IsSwitchingRoot || m.DecksView.IsPickingTheme) {
		return true
	}
	if m.ActiveTab == TabSlides && m.SlidesView.IsAdding {
		return true
	}
	if m.ActiveTab == TabScaffold && m.ScaffoldView.IsEditing {
		return true
	}
	return false
}

func (m *AppModel) View() string {
	var s strings.Builder

	// Header Bar
	headerTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(crimsonDark).
		Padding(0, 1).
		Render("DECKFORGE")

	headerSub := lipgloss.NewStyle().
		Foreground(faint).
		Render(" 1080p Presentation Platform")

	s.WriteString(fmt.Sprintf("%s%s\n", headerTitle, headerSub))

	// Navigation Tabs (5 Tabs)
	tabNames := []string{"[1] DECKS", "[2] SLIDES", "[3] SCAFFOLD", "[4] THEMES", "[5] SETTINGS"}
	var renderedTabs []string
	for i, name := range tabNames {
		if TabIndex(i) == m.ActiveTab {
			renderedTabs = append(renderedTabs, ActiveTab.Render(name))
		} else {
			renderedTabs = append(renderedTabs, InactiveTab.Render(name))
		}
	}
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...))
	s.WriteString("\n\n")

	// View Content
	switch m.ActiveTab {
	case TabDecks:
		s.WriteString(m.DecksView.View())
	case TabSlides:
		s.WriteString(m.SlidesView.View())
	case TabScaffold:
		s.WriteString(m.ScaffoldView.View())
	case TabThemes:
		s.WriteString(m.ThemeView.View())
	case TabSettings:
		s.WriteString(m.SettingsView.View())
	}

	// Status & Keybinding Bar
	s.WriteString("\n")
	helpBar := KeyHelpStyle.Render("[1-5 or ←/→ or [/]] Switch Tabs  [q] Quit  [w] Switch Root  [e] Edit in $EDITOR")
	s.WriteString(StatusBarStyle.Render(helpBar))

	return DocStyle.Render(s.String())
}
