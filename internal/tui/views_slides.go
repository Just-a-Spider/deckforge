package tui

import (
	"fmt"
	"os"
	"strings"

	"deckforge/internal/compiler"
	"deckforge/internal/models"
	"deckforge/internal/scaffold"
	"deckforge/internal/workspace"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SlidesView struct {
	CurrentDeck     *models.Deck
	Cursor          int
	ViewportOffset  int
	MaxVisibleCards int
	Width           int
	Height          int
	IsCompact       bool
	Compiler        *compiler.Compiler
	Config          *models.WorkspaceConfig
	ServerCtrl      *ServerController
	IsAdding        bool
	TitleInput      textinput.Model
	LayoutIndex     int
	Layouts         []string
	StatusMsg       string
	IsError         bool
	NavigateToTab   TabIndex
}

func NewSlidesView(comp *compiler.Compiler, cfg *models.WorkspaceConfig, srvCtrl *ServerController) *SlidesView {
	ti := textinput.New()
	ti.Placeholder = "New Slide Title (e.g. Distributed Consensus)"
	ti.CharLimit = 80
	ti.Width = 45

	return &SlidesView{
		Compiler:        comp,
		Config:          cfg,
		ServerCtrl:      srvCtrl,
		TitleInput:      ti,
		Layouts:         []string{"2col", "3col", "hero", "code"},
		LayoutIndex:     0,
		MaxVisibleCards: 10,
		ViewportOffset:  0,
		Height:          24,
	}
}

func (v *SlidesView) SetDimensions(width, height int) {
	v.Width = width
	v.Height = height
	v.IsCompact = height < 30

	// Total chrome overhead (app header, tabs, footer, deck info, scroll info)
	overhead := 10
	if v.IsAdding {
		overhead += 5
	}

	available := height - overhead
	if available < 4 {
		available = 4
	}

	// 1 selected card takes 2 lines, remaining visible cards take 1 line
	v.MaxVisibleCards = available - 1
	if v.MaxVisibleCards < 3 {
		v.MaxVisibleCards = 3
	}

	// Re-clamp viewport strictly
	if v.Cursor < v.ViewportOffset {
		v.ViewportOffset = v.Cursor
	} else if v.Cursor >= v.ViewportOffset+v.MaxVisibleCards {
		v.ViewportOffset = v.Cursor - v.MaxVisibleCards + 1
	}
	if v.ViewportOffset < 0 {
		v.ViewportOffset = 0
	}
}

func (v *SlidesView) SetDeck(d *models.Deck) {
	v.CurrentDeck = d
	v.Cursor = 0
	v.ViewportOffset = 0
	v.ReloadSlides()
}

func (v *SlidesView) ReloadSlides() {
	if v.CurrentDeck == nil {
		return
	}
	refreshed, err := workspace.InspectDeck(v.CurrentDeck.Path)
	if err == nil {
		v.CurrentDeck = refreshed
		if v.Cursor >= len(v.CurrentDeck.Slides) && len(v.CurrentDeck.Slides) > 0 {
			v.Cursor = len(v.CurrentDeck.Slides) - 1
		}
	}
}

func (v *SlidesView) SelectedSlide() *models.SlideInfo {
	if v.CurrentDeck == nil || len(v.CurrentDeck.Slides) == 0 || v.Cursor < 0 || v.Cursor >= len(v.CurrentDeck.Slides) {
		return nil
	}
	return &v.CurrentDeck.Slides[v.Cursor]
}

func (v *SlidesView) Update(msg tea.Msg) tea.Cmd {
	if v.IsAdding {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "tab":
				v.LayoutIndex = (v.LayoutIndex + 1) % len(v.Layouts)
				return nil
			case "enter":
				title := strings.TrimSpace(v.TitleInput.Value())
				if title != "" && v.CurrentDeck != nil {
					layout := v.Layouts[v.LayoutIndex]
					slide, err := scaffold.AddSlide(v.CurrentDeck.Path, title, layout)
					if err != nil {
						v.StatusMsg = fmt.Sprintf("Failed to add slide: %v", err)
						v.IsError = true
					} else {
						v.ReloadSlides()
						v.StatusMsg = fmt.Sprintf("Added Slide %02d: %s", slide.Index, slide.Title)
						v.IsError = false
						v.Cursor = len(v.CurrentDeck.Slides) - 1
					}
				}
				v.IsAdding = false
				v.TitleInput.Blur()
				return nil
			case "esc":
				v.IsAdding = false
				v.TitleInput.Blur()
				return nil
			}
		}
		var cmd tea.Cmd
		v.TitleInput, cmd = v.TitleInput.Update(msg)
		return cmd
	}

	switch msg := msg.(type) {
	case EditorFinishedMsg:
		v.ReloadSlides()
		v.StatusMsg = "Editor closed. Deck slides reloaded."
		v.IsError = false
		// Auto recompile
		if v.CurrentDeck != nil {
			_, _ = v.Compiler.Build(v.CurrentDeck.Path, "")
		}
		return nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if v.Cursor > 0 {
				v.Cursor--
				if v.Cursor < v.ViewportOffset {
					v.ViewportOffset = v.Cursor
				}
			}
		case "down", "j":
			if v.CurrentDeck != nil && v.Cursor < len(v.CurrentDeck.Slides)-1 {
				v.Cursor++
				if v.Cursor >= v.ViewportOffset+v.MaxVisibleCards {
					v.ViewportOffset = v.Cursor - v.MaxVisibleCards + 1
				}
			}
		case "K", "shift+up":
			if v.CurrentDeck != nil && v.Cursor > 0 && len(v.CurrentDeck.Slides) > 1 {
				total := len(v.CurrentDeck.Slides)
				newOrder := make([]int, total)
				for i := 0; i < total; i++ {
					newOrder[i] = i + 1
				}
				newOrder[v.Cursor-1], newOrder[v.Cursor] = newOrder[v.Cursor], newOrder[v.Cursor-1]
				if err := workspace.ReorderSlides(v.CurrentDeck.Path, newOrder); err != nil {
					v.StatusMsg = fmt.Sprintf("Reorder error: %v", err)
					v.IsError = true
				} else {
					v.ReloadSlides()
					v.Cursor--
					if v.Cursor < v.ViewportOffset {
						v.ViewportOffset = v.Cursor
					}
					v.StatusMsg = fmt.Sprintf("Moved slide up to %02d", v.Cursor+1)
					v.IsError = false
					_, _ = v.Compiler.Build(v.CurrentDeck.Path, "")
				}
			}
		case "J", "shift+down":
			if v.CurrentDeck != nil && v.Cursor < len(v.CurrentDeck.Slides)-1 && len(v.CurrentDeck.Slides) > 1 {
				total := len(v.CurrentDeck.Slides)
				newOrder := make([]int, total)
				for i := 0; i < total; i++ {
					newOrder[i] = i + 1
				}
				newOrder[v.Cursor], newOrder[v.Cursor+1] = newOrder[v.Cursor+1], newOrder[v.Cursor]
				if err := workspace.ReorderSlides(v.CurrentDeck.Path, newOrder); err != nil {
					v.StatusMsg = fmt.Sprintf("Reorder error: %v", err)
					v.IsError = true
				} else {
					v.ReloadSlides()
					v.Cursor++
					if v.Cursor >= v.ViewportOffset+v.MaxVisibleCards {
						v.ViewportOffset = v.Cursor - v.MaxVisibleCards + 1
					}
					v.StatusMsg = fmt.Sprintf("Moved slide down to %02d", v.Cursor+1)
					v.IsError = false
					_, _ = v.Compiler.Build(v.CurrentDeck.Path, "")
				}
			}
		case "pgup":
			v.Cursor -= v.MaxVisibleCards
			if v.Cursor < 0 {
				v.Cursor = 0
			}
			if v.Cursor < v.ViewportOffset {
				v.ViewportOffset = v.Cursor
			}
		case "pgdown":
			if v.CurrentDeck != nil {
				v.Cursor += v.MaxVisibleCards
				if v.Cursor >= len(v.CurrentDeck.Slides) {
					v.Cursor = len(v.CurrentDeck.Slides) - 1
				}
				if v.Cursor >= v.ViewportOffset+v.MaxVisibleCards {
					v.ViewportOffset = v.Cursor - v.MaxVisibleCards + 1
				}
			}
		case "e":
			s := v.SelectedSlide()
			if s != nil {
				pref := "auto"
				if v.Config != nil {
					pref = v.Config.PreferredEditor
				}
				return OpenInEditor(pref, s.Path)
			}
		case "E", "a":
			if v.CurrentDeck != nil && len(v.CurrentDeck.Slides) > 0 {
				var paths []string
				for _, sl := range v.CurrentDeck.Slides {
					paths = append(paths, sl.Path)
				}
				pref := "auto"
				if v.Config != nil {
					pref = v.Config.PreferredEditor
				}
				return OpenInEditor(pref, paths...)
			}
		case "n":
			if v.CurrentDeck != nil {
				v.IsAdding = true
				v.TitleInput.Reset()
				v.TitleInput.Focus()
				return textinput.Blink
			}
		case "d":
			s := v.SelectedSlide()
			if s != nil {
				_ = os.Remove(s.Path)
				v.ReloadSlides()
				v.StatusMsg = fmt.Sprintf("Deleted slide %s", s.Filename)
				v.IsError = false
			}
		case "r":
			v.ReloadSlides()
			v.StatusMsg = "Slides refreshed."
			v.IsError = false
		case "b":
			if v.CurrentDeck != nil {
				out, err := v.Compiler.Build(v.CurrentDeck.Path, "")
				if err != nil {
					v.StatusMsg = fmt.Sprintf("Build failed: %v", err)
					v.IsError = true
				} else {
					v.StatusMsg = fmt.Sprintf("Compiled to %s", out)
					v.IsError = false
				}
			}
		case "s":
			if v.CurrentDeck != nil && v.ServerCtrl != nil {
				port := 8080
				watch := true
				if v.Config != nil {
					if v.Config.ServerPort > 0 {
						port = v.Config.ServerPort
					}
					watch = v.Config.AutoWatch
				}
				_, msg, err := v.ServerCtrl.Toggle(v.CurrentDeck, port, watch, v.Compiler)
				v.StatusMsg = msg
				v.IsError = err != nil
			}
		case "x":
			if v.ServerCtrl != nil && v.ServerCtrl.IsRunning() {
				_ = v.ServerCtrl.Stop()
				v.StatusMsg = "Server stopped."
				v.IsError = false
			} else {
				v.StatusMsg = "No server currently running."
				v.IsError = false
			}
		case "t":
			v.NavigateToTab = TabThemes
			return nil
		case "o":
			if v.CurrentDeck != nil {
				if v.ServerCtrl != nil && v.ServerCtrl.IsServing(v.CurrentDeck.Path) {
					openBrowser(fmt.Sprintf("http://localhost:%d/", v.ServerCtrl.ActivePort()))
				} else {
					outFile, err := v.Compiler.Build(v.CurrentDeck.Path, "")
					if err == nil {
						openBrowser(fmt.Sprintf("file://%s", outFile))
					}
				}
			}
		}
	}

	return nil
}

func (v *SlidesView) View() string {
	var s strings.Builder

	if v.CurrentDeck == nil {
		s.WriteString(fmt.Sprintf("%s\n%s\n",
			lipgloss.NewStyle().Foreground(muted).Render("No active presentation selected."),
			KeyHelpStyle.Render("Go to [1] Decks tab and press Enter on a presentation to inspect its slides."),
		))
		return s.String()
	}

	// Active Deck Header
	serverBadge := ""
	if v.ServerCtrl != nil && v.ServerCtrl.IsServing(v.CurrentDeck.Path) {
		serverBadge = " " + ActiveBadgeStyle.Render(fmt.Sprintf("LIVE: http://localhost:%d", v.ServerCtrl.ActivePort()))
	}
	themeBadge := ""
	if v.CurrentDeck.Config.Theme != "" {
		themeBadge = " " + BadgeStyle.Render("THEME: "+v.CurrentDeck.Config.Theme)
	}
	s.WriteString(fmt.Sprintf("%s %s  %s%s%s\n",
		lipgloss.NewStyle().Bold(true).Foreground(crimson).Render("PRESENTATION:"),
		ItemTitle.Render(v.CurrentDeck.Config.Title),
		BadgeStyle.Render(fmt.Sprintf("%d Slides", len(v.CurrentDeck.Slides))),
		themeBadge,
		serverBadge,
	))
	s.WriteString(fmt.Sprintf("%s\n\n", RootPathStyle.Render(v.CurrentDeck.Path)))

	if v.IsAdding {
		s.WriteString(fmt.Sprintf("%s\n%s\nLayout (Tab to cycle): %s\n%s\n\n",
			lipgloss.NewStyle().Bold(true).Foreground(lime).Render("Add New Slide:"),
			v.TitleInput.View(),
			ActiveBadgeStyle.Render(v.Layouts[v.LayoutIndex]),
			KeyHelpStyle.Render("Enter to create, Esc to cancel"),
		))
	}

	// Slide List
	if len(v.CurrentDeck.Slides) == 0 {
		s.WriteString(lipgloss.NewStyle().Foreground(muted).Render("This deck has no slides. Press 'n' to add one.\n"))
	} else {
		total := len(v.CurrentDeck.Slides)
		start := v.ViewportOffset
		if start < 0 {
			start = 0
		}
		if start >= total {
			start = 0
		}
		end := start + v.MaxVisibleCards
		if end > total {
			end = total
		}

		scrollInfo := fmt.Sprintf("Showing slides %d - %d of %d  [Scroll: ▲ / ▼ or j / k or PgUp/PgDn]", start+1, end, total)
		s.WriteString(lipgloss.NewStyle().Foreground(faint).Italic(true).Render(scrollInfo))
		s.WriteString("\n\n")

		if start > 0 {
			earlierText := fmt.Sprintf("   ▲ ... %d earlier slides above ...", start)
			s.WriteString(lipgloss.NewStyle().Foreground(crimson).Render(earlierText))
			s.WriteString("\n")
		}

		for i := start; i < end; i++ {
			sl := v.CurrentDeck.Slides[i]
			selected := i == v.Cursor

			if selected {
				prefix := "▸ "
				numBadge := lipgloss.NewStyle().Bold(true).Foreground(white).Background(crimsonDark).Padding(0, 1)
				titleStyle := lipgloss.NewStyle().Bold(true).Foreground(white)
				fnameStyle := lipgloss.NewStyle().Foreground(lime)

				row := fmt.Sprintf("%s%s %-40s %s\n    %s",
					prefix,
					numBadge.Render(fmt.Sprintf("%02d", sl.Index)),
					titleStyle.Render(sl.Title),
					fnameStyle.Render(sl.Filename),
					KeyHelpStyle.Render("↳ [e] Edit  [E] ALL  [t] Theme  [s] Server  [o] Browser  [b] Build  [n] Add  [d] Delete"),
				)
				s.WriteString(row + "\n")
			} else {
				prefix := "  "
				numBadge := BadgeStyle
				titleStyle := ItemTitle
				fnameStyle := ItemDesc

				row := fmt.Sprintf("%s%s %-40s %s",
					prefix,
					numBadge.Render(fmt.Sprintf("%02d", sl.Index)),
					titleStyle.Render(sl.Title),
					fnameStyle.Render(sl.Filename),
				)
				s.WriteString(row + "\n")
			}
		}

		if end < total {
			remainingText := fmt.Sprintf("   ▼ ... %d more slides below ...", total-end)
			s.WriteString(lipgloss.NewStyle().Foreground(crimson).Render(remainingText))
			s.WriteString("\n")
		}
	}

	if v.StatusMsg != "" {
		style := SuccessStyle
		if v.IsError {
			style = ErrorStyle
		}
		s.WriteString(fmt.Sprintf("\n%s\n", style.Render(v.StatusMsg)))
	}

	return s.String()
}
