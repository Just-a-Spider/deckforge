package tui

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"deckforge/internal/compiler"
	"deckforge/internal/models"
	"deckforge/internal/workspace"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DecksView struct {
	Decks             []*models.Deck
	Cursor            int
	ActiveRoot        string
	RecentRoots       []string
	Compiler          *compiler.Compiler
	Config            *models.WorkspaceConfig
	ServerCtrl        *ServerController
	RootInput         textinput.Model
	IsSwitchingRoot   bool
	IsPickingTheme    bool
	ThemeChoices      []models.Theme
	ThemeChoiceCursor int
	StatusMsg         string
	IsError           bool
	NavigateToTab     TabIndex
}

func NewDecksView(root string, recent []string, comp *compiler.Compiler, cfg *models.WorkspaceConfig, srvCtrl *ServerController) *DecksView {
	ti := textinput.New()
	ti.Placeholder = "/path/to/presentations"
	ti.CharLimit = 256
	ti.Width = 50

	v := &DecksView{
		ActiveRoot:  root,
		RecentRoots: recent,
		Compiler:    comp,
		Config:      cfg,
		ServerCtrl:  srvCtrl,
		RootInput:   ti,
	}
	v.RefreshDecks()
	return v
}

func (v *DecksView) RefreshDecks() {
	decks, err := workspace.ScanDecksInRoot(v.ActiveRoot)
	if err != nil {
		v.StatusMsg = fmt.Sprintf("Scan error: %v", err)
		v.IsError = true
		v.Decks = nil
		return
	}
	v.Decks = decks
	if v.Cursor >= len(v.Decks) && len(v.Decks) > 0 {
		v.Cursor = len(v.Decks) - 1
	}
}

func (v *DecksView) SelectedDeck() *models.Deck {
	if len(v.Decks) == 0 || v.Cursor < 0 || v.Cursor >= len(v.Decks) {
		return nil
	}
	return v.Decks[v.Cursor]
}

func (v *DecksView) Update(msg tea.Msg) (tea.Cmd, bool) {
	if v.IsSwitchingRoot {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				newPath := strings.TrimSpace(v.RootInput.Value())
				if newPath != "" {
					v.ActiveRoot = newPath
					v.RefreshDecks()
					v.StatusMsg = fmt.Sprintf("Root switched to: %s", newPath)
					v.IsError = false
				}
				v.IsSwitchingRoot = false
				v.RootInput.Blur()
				return nil, false
			case "esc":
				v.IsSwitchingRoot = false
				v.RootInput.Blur()
				return nil, false
			}
		}
		var cmd tea.Cmd
		v.RootInput, cmd = v.RootInput.Update(msg)
		return cmd, false
	}

	if v.IsPickingTheme {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				if len(v.ThemeChoices) > 0 {
					v.ThemeChoiceCursor--
					if v.ThemeChoiceCursor < 0 {
						v.ThemeChoiceCursor = len(v.ThemeChoices) - 1
					}
				}
				return nil, false
			case "down", "j":
				if len(v.ThemeChoices) > 0 {
					v.ThemeChoiceCursor++
					if v.ThemeChoiceCursor >= len(v.ThemeChoices) {
						v.ThemeChoiceCursor = 0
					}
				}
				return nil, false
			case "enter":
				if len(v.ThemeChoices) > 0 && v.ThemeChoiceCursor >= 0 && v.ThemeChoiceCursor < len(v.ThemeChoices) {
					d := v.SelectedDeck()
					if d != nil {
						chosen := v.ThemeChoices[v.ThemeChoiceCursor]
						err := workspace.SetDeckTheme(d.Path, chosen.Tokens.Name)
						if err != nil {
							v.StatusMsg = fmt.Sprintf("Failed to set theme: %v", err)
							v.IsError = true
						} else {
							d.Config.Theme = chosen.Tokens.Name
							d.ThemeName = chosen.Tokens.Name
							if v.Compiler != nil {
								outFile, err := v.Compiler.Build(d.Path, "")
								if err == nil {
									v.StatusMsg = fmt.Sprintf("Applied theme '%s' to '%s' (compiled: %s)", chosen.Tokens.Name, d.Name, filepath.Base(outFile))
								} else {
									v.StatusMsg = fmt.Sprintf("Applied theme '%s' to '%s'", chosen.Tokens.Name, d.Name)
								}
							} else {
								v.StatusMsg = fmt.Sprintf("Applied theme '%s' to '%s'", chosen.Tokens.Name, d.Name)
							}
							v.IsError = false
						}
					}
				}
				v.IsPickingTheme = false
				return nil, false
			case "esc", "q":
				v.IsPickingTheme = false
				return nil, false
			}
		}
		return nil, false
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if v.Cursor > 0 {
				v.Cursor--
			}
		case "down", "j":
			if v.Cursor < len(v.Decks)-1 {
				v.Cursor++
			}
		case "w":
			v.IsSwitchingRoot = true
			v.RootInput.SetValue(v.ActiveRoot)
			v.RootInput.Focus()
			return textinput.Blink, false
		case "t":
			d := v.SelectedDeck()
			if d != nil && v.Compiler != nil && v.Compiler.ThemeManager != nil {
				v.ThemeChoices = v.Compiler.ThemeManager.ListThemes(d.Path)
				if len(v.ThemeChoices) > 0 {
					v.ThemeChoiceCursor = 0
					for i, th := range v.ThemeChoices {
						if th.Tokens.Name == d.ThemeName {
							v.ThemeChoiceCursor = i
							break
						}
					}
					v.IsPickingTheme = true
					return nil, false
				}
			}
		case "r":
			v.RefreshDecks()
			v.StatusMsg = "Workspace rescanned."
			v.IsError = false
		case "b":
			d := v.SelectedDeck()
			if d != nil {
				outFile, err := v.Compiler.Build(d.Path, "")
				if err != nil {
					v.StatusMsg = fmt.Sprintf("Build failed: %v", err)
					v.IsError = true
				} else {
					v.StatusMsg = fmt.Sprintf("Compiled successfully to %s", outFile)
					v.IsError = false
				}
			}
		case "s":
			d := v.SelectedDeck()
			if d != nil && v.ServerCtrl != nil {
				port := 8080
				watch := true
				if v.Config != nil {
					if v.Config.ServerPort > 0 {
						port = v.Config.ServerPort
					}
					watch = v.Config.AutoWatch
				}
				_, msg, err := v.ServerCtrl.Toggle(d, port, watch, v.Compiler)
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
		case "o":
			d := v.SelectedDeck()
			if d != nil {
				if v.ServerCtrl != nil && v.ServerCtrl.IsServing(d.Path) {
					openBrowser(fmt.Sprintf("http://localhost:%d/", v.ServerCtrl.ActivePort()))
				} else {
					outFile, err := v.Compiler.Build(d.Path, "")
					if err == nil {
						openBrowser(fmt.Sprintf("file://%s", outFile))
					}
				}
			}
		case "e":
			d := v.SelectedDeck()
			if d != nil {
				pref := "auto"
				if v.Config != nil {
					pref = v.Config.PreferredEditor
				}
				return OpenInEditor(pref, filepath.Join(d.Path, "deck.json")), false
			}
		case "n":
			v.NavigateToTab = TabScaffold
			return nil, false
		case "enter":
			// Select deck and switch to Slides tab (signal to parent app)
			return nil, true
		}
	}

	return nil, false
}

func (v *DecksView) View() string {
	var s strings.Builder

	// Workspace Header
	s.WriteString(fmt.Sprintf("%s %s\n",
		lipgloss.NewStyle().Bold(true).Foreground(crimson).Render("ROOT:"),
		RootPathStyle.Render(v.ActiveRoot),
	))

	if v.IsSwitchingRoot {
		s.WriteString(fmt.Sprintf("\n%s\n%s\n%s\n\n",
			lipgloss.NewStyle().Bold(true).Foreground(orange).Render("Switch Workspace Root Path:"),
			v.RootInput.View(),
			KeyHelpStyle.Render("Press Enter to scan path, Esc to cancel"),
		))
	}

	if v.IsPickingTheme {
		d := v.SelectedDeck()
		deckName := ""
		if d != nil {
			deckName = d.Name
		}
		totalThemes := len(v.ThemeChoices)
		maxVisible := 7
		startIdx := 0
		if v.ThemeChoiceCursor >= maxVisible {
			startIdx = v.ThemeChoiceCursor - maxVisible + 1
		}
		endIdx := startIdx + maxVisible
		if endIdx > totalThemes {
			endIdx = totalThemes
			startIdx = endIdx - maxVisible
			if startIdx < 0 {
				startIdx = 0
			}
		}

		var modal strings.Builder
		modal.WriteString(fmt.Sprintf("%s\n\n",
			lipgloss.NewStyle().Bold(true).Foreground(crimson).Render(fmt.Sprintf("Select Theme for: %s  (%d/%d)", deckName, v.ThemeChoiceCursor+1, totalThemes)),
		))

		for i := startIdx; i < endIdx; i++ {
			th := v.ThemeChoices[i]
			selected := i == v.ThemeChoiceCursor
			prefix := "  "
			nameStyle := lipgloss.NewStyle().Foreground(white)
			if selected {
				prefix = "▸ "
				nameStyle = lipgloss.NewStyle().Bold(true).Foreground(lime)
			}
			scopeBadge := BadgeStyle.Render(strings.ToUpper(th.Scope))
			desc := th.Tokens.Vibe
			if desc == "" {
				desc = th.Tokens.Description
			}

			// Pre-pad theme name before applying ANSI codes so columns stay strictly aligned
			paddedName := fmt.Sprintf("%-20s", th.Tokens.Name)
			modal.WriteString(fmt.Sprintf("%s%s %s %s\n",
				prefix,
				nameStyle.Render(paddedName),
				scopeBadge,
				lipgloss.NewStyle().Foreground(muted).Render(desc),
			))
		}
		modal.WriteString(fmt.Sprintf("\n%s\n", KeyHelpStyle.Render("[↑/k/↓/j] Navigate  [Enter] Apply Theme & Recompile  [Esc/q] Cancel")))

		modalBox := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(crimson).
			Padding(1, 2).
			Render(modal.String())

		s.WriteString("\n")
		s.WriteString(modalBox)
		s.WriteString("\n")

		if v.StatusMsg != "" {
			style := SuccessStyle
			if v.IsError {
				style = ErrorStyle
			}
			s.WriteString(fmt.Sprintf("\n%s\n", style.Render(v.StatusMsg)))
		}

		return s.String()
	}

	// Decks List
	if len(v.Decks) == 0 {
		s.WriteString(fmt.Sprintf("\n%s\n%s\n",
			lipgloss.NewStyle().Foreground(muted).Render("No slide decks discovered in this root."),
			KeyHelpStyle.Render("Press 'w' to switch root, or switch to [3] Scaffold tab to create one."),
		))
	} else {
		s.WriteString(fmt.Sprintf("\nDiscovered %d Presentation Decks:\n\n", len(v.Decks)))
		for i, d := range v.Decks {
			selected := i == v.Cursor

			prefix := "  "
			itemStyle := CardStyle
			if selected {
				prefix = "▸ "
				itemStyle = ActiveCardStyle
			}

			serverBadge := ""
			if v.ServerCtrl != nil && v.ServerCtrl.IsServing(d.Path) {
				serverBadge = ActiveBadgeStyle.Render(fmt.Sprintf("LIVE: http://localhost:%d", v.ServerCtrl.ActivePort()))
			}

			content := fmt.Sprintf("%s%s %s %s\n   %s | %d Slides | Theme: %s\n   %s",
				prefix,
				ItemTitle.Render(d.Config.Title),
				BadgeStyle.Render(d.Name),
				serverBadge,
				ItemDesc.Render(d.Path),
				len(d.Slides),
				lipgloss.NewStyle().Foreground(blue).Render(d.ThemeName),
				KeyHelpStyle.Render("[Enter] Slides  [t] Theme  [n] New Deck  [b] Build  [s] Server  [x] Stop  [o] Browser  [e] Edit deck.json"),
			)

			s.WriteString(itemStyle.Render(content))
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

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}

	_ = exec.Command(cmd, args...).Start()
}
