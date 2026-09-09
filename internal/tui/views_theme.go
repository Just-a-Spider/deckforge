package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	"deckforge/internal/compiler"
	"deckforge/internal/models"
	"deckforge/internal/theme"
	"deckforge/internal/workspace"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ThemeStudioView struct {
	ThemeManager    *theme.ThemeManager
	Themes          []models.Theme
	ThemeCursor     int
	ThemeOffset     int
	SegmentCursor   int
	ActivePanel     int // 0 = Left (Themes List), 1 = Right (Segments Workshop)
	Segments        []string
	ActiveWorkspace string
	ActiveDeck      *models.Deck
	Compiler        *compiler.Compiler
	Config          *models.WorkspaceConfig
	StatusMsg       string
	IsError         bool
}

func NewThemeStudioView(tm *theme.ThemeManager, workspace string, cfg *models.WorkspaceConfig, comp *compiler.Compiler) *ThemeStudioView {
	v := &ThemeStudioView{
		ThemeManager:    tm,
		ActiveWorkspace: workspace,
		Config:          cfg,
		Compiler:        comp,
		Segments: []string{
			"1. Palette & Fonts (tokens.json)",
			"2. Typography Rules (typography.css)",
			"3. Surfaces & Cards (surfaces.css)",
			"4. Backdrop & Atmosphere (backdrop.css)",
		},
		ActivePanel: 0,
	}
	v.RefreshThemes()
	return v
}

func (v *ThemeStudioView) SetActiveDeck(deck *models.Deck) {
	v.ActiveDeck = deck
}

func (v *ThemeStudioView) RefreshThemes() {
	v.Themes = v.ThemeManager.ListThemes("")
	if v.ThemeCursor >= len(v.Themes) && len(v.Themes) > 0 {
		v.ThemeCursor = len(v.Themes) - 1
	}
}

func (v *ThemeStudioView) SelectedTheme() *models.Theme {
	if len(v.Themes) == 0 || v.ThemeCursor < 0 || v.ThemeCursor >= len(v.Themes) {
		return nil
	}
	return &v.Themes[v.ThemeCursor]
}

func (v *ThemeStudioView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case EditorFinishedMsg:
		v.RefreshThemes()
		v.StatusMsg = "Editor closed. Theme segment refreshed."
		v.IsError = false
		return nil

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			v.ActivePanel = (v.ActivePanel + 1) % 2
			return nil

		case "h":
			v.ActivePanel = 0
			return nil

		case "l":
			v.ActivePanel = 1
			return nil

		case "esc":
			v.ActivePanel = 0
			return nil

		case "up", "k":
			if v.ActivePanel == 0 {
				if v.ThemeCursor > 0 {
					v.ThemeCursor--
					if v.ThemeCursor < v.ThemeOffset {
						v.ThemeOffset = v.ThemeCursor
					}
				}
			} else {
				if v.SegmentCursor > 0 {
					v.SegmentCursor--
				}
			}

		case "down", "j":
			if v.ActivePanel == 0 {
				if v.ThemeCursor < len(v.Themes)-1 {
					v.ThemeCursor++
					maxVisible := 8
					if v.ThemeCursor >= v.ThemeOffset+maxVisible {
						v.ThemeOffset = v.ThemeCursor - maxVisible + 1
					}
				}
			} else {
				if v.SegmentCursor < len(v.Segments)-1 {
					v.SegmentCursor++
				}
			}

		case "f":
			th := v.SelectedTheme()
			if th != nil {
				reports := models.EvaluateThemeContrast(th.Tokens.Palette)
				fixed := 0
				for _, rep := range reports {
					if rep.RecommendedAction != nil {
						rec := rep.RecommendedAction
						switch rec.TargetToken {
						case "textPrimary":
							th.Tokens.Palette.TextPrimary = rec.SuggestedHex
							fixed++
						case "textSecondary":
							th.Tokens.Palette.TextSecondary = rec.SuggestedHex
							fixed++
						case "accentPrimary":
							th.Tokens.Palette.AccentPrimary = rec.SuggestedHex
							fixed++
						}
					}
				}
				if fixed > 0 {
					themeDir := th.Dir
					if themeDir == "" {
						targetDir := filepath.Join(v.ActiveWorkspace, "themes", th.Tokens.Name)
						if v.ActiveWorkspace == "" {
							targetDir = filepath.Join("themes", th.Tokens.Name)
						}
						themeDir = targetDir
						th.Dir = themeDir
					}
					_ = theme.SaveSegmentedTheme(themeDir, *th)
					v.RefreshThemes()
					v.StatusMsg = fmt.Sprintf("Applied %d WCAG contrast fixes to %s", fixed, th.Tokens.Name)
					v.IsError = false
				} else {
					v.StatusMsg = "All token pairs already meet WCAG AA contrast."
					v.IsError = false
				}
			}

		case "c":
			// Clone highlighted preset into active workspace themes/
			th := v.SelectedTheme()
			if th != nil {
				cloneName := fmt.Sprintf("%s-custom", th.Tokens.Name)
				targetDir := filepath.Join(v.ActiveWorkspace, "themes", cloneName)
				if v.ActiveWorkspace == "" {
					targetDir = filepath.Join("themes", cloneName)
				}
				cloned, err := v.ThemeManager.ClonePreset(th.Tokens.Name, cloneName, targetDir)
				if err != nil {
					v.StatusMsg = fmt.Sprintf("Clone failed: %v", err)
					v.IsError = true
				} else {
					v.RefreshThemes()
					v.StatusMsg = fmt.Sprintf("Cloned segmented theme to: %s", cloned.Dir)
					v.IsError = false
					// Switch cursor to cloned theme
					for i, t := range v.Themes {
						if t.Tokens.Name == cloneName {
							v.ThemeCursor = i
							break
						}
					}
				}
			}

		case "a":
			th := v.SelectedTheme()
			if th != nil && v.ActiveDeck != nil {
				err := workspace.SetDeckTheme(v.ActiveDeck.Path, th.Tokens.Name)
				if err != nil {
					v.StatusMsg = fmt.Sprintf("Failed to set theme: %v", err)
					v.IsError = true
				} else {
					v.ActiveDeck.Config.Theme = th.Tokens.Name
					v.ActiveDeck.ThemeName = th.Tokens.Name
					if v.Compiler != nil {
						outFile, err := v.Compiler.Build(v.ActiveDeck.Path, "")
						if err == nil {
							v.StatusMsg = fmt.Sprintf("Applied theme '%s' to '%s' (compiled: %s)", th.Tokens.Name, v.ActiveDeck.Name, filepath.Base(outFile))
						} else {
							v.StatusMsg = fmt.Sprintf("Applied theme '%s' to '%s'", th.Tokens.Name, v.ActiveDeck.Name)
						}
					} else {
						v.StatusMsg = fmt.Sprintf("Applied theme '%s' to '%s'", th.Tokens.Name, v.ActiveDeck.Name)
					}
					v.IsError = false
				}
			} else if th != nil && v.ActiveDeck == nil {
				v.StatusMsg = fmt.Sprintf("No active presentation selected. Selected theme: %s", th.Tokens.Name)
				v.IsError = false
			}
			return nil

		case "e", "enter":
			if v.ActivePanel == 0 && msg.String() == "enter" {
				th := v.SelectedTheme()
				if th != nil && v.ActiveDeck != nil {
					err := workspace.SetDeckTheme(v.ActiveDeck.Path, th.Tokens.Name)
					if err != nil {
						v.StatusMsg = fmt.Sprintf("Failed to set theme: %v", err)
						v.IsError = true
					} else {
						v.ActiveDeck.Config.Theme = th.Tokens.Name
						v.ActiveDeck.ThemeName = th.Tokens.Name
						if v.Compiler != nil {
							outFile, err := v.Compiler.Build(v.ActiveDeck.Path, "")
							if err == nil {
								v.StatusMsg = fmt.Sprintf("Applied theme '%s' to '%s' (compiled: %s)", th.Tokens.Name, v.ActiveDeck.Name, filepath.Base(outFile))
							} else {
								v.StatusMsg = fmt.Sprintf("Applied theme '%s' to '%s'", th.Tokens.Name, v.ActiveDeck.Name)
							}
						} else {
							v.StatusMsg = fmt.Sprintf("Applied theme '%s' to '%s'", th.Tokens.Name, v.ActiveDeck.Name)
						}
						v.IsError = false
					}
				} else if th != nil && v.ActiveDeck == nil {
					v.StatusMsg = fmt.Sprintf("No active presentation selected. Selected theme: %s", th.Tokens.Name)
					v.IsError = false
				}
				return nil
			}

			// Open selected segment in external editor
			th := v.SelectedTheme()
			if th != nil {
				themeDir := th.Dir
				if themeDir == "" {
					targetDir := filepath.Join(v.ActiveWorkspace, "themes", th.Tokens.Name)
					if v.ActiveWorkspace == "" {
						targetDir = filepath.Join("themes", th.Tokens.Name)
					}
					themeDir = targetDir
					_ = theme.SaveSegmentedTheme(themeDir, *th)
					th.Dir = themeDir
				}

				var targetFile string
				switch v.SegmentCursor {
				case 0:
					targetFile = filepath.Join(themeDir, "tokens.json")
				case 1:
					targetFile = filepath.Join(themeDir, "typography.css")
				case 2:
					targetFile = filepath.Join(themeDir, "surfaces.css")
				case 3:
					targetFile = filepath.Join(themeDir, "backdrop.css")
				}

				pref := "auto"
				if v.Config != nil {
					pref = v.Config.PreferredEditor
				}
				return OpenInEditor(pref, targetFile)
			}
		}
	}

	return nil
}

func (v *ThemeStudioView) View() string {
	var s strings.Builder

	deckHeader := ""
	if v.ActiveDeck != nil {
		deckHeader = fmt.Sprintf("  • Active Deck: %s", lipgloss.NewStyle().Foreground(lime).Bold(true).Render(v.ActiveDeck.Config.Title))
	}
	s.WriteString(fmt.Sprintf("%s  %s%s\n\n",
		lipgloss.NewStyle().Bold(true).Foreground(crimson).Render("SEGMENTED THEME WORKSHOP"),
		lipgloss.NewStyle().Foreground(faint).Render("• Modular CSS Token Refinement"),
		deckHeader,
	))

	th := v.SelectedTheme()
	if th == nil {
		s.WriteString("No themes loaded.\n")
		return s.String()
	}

	// 1. Left Column: Theme List
	leftCol := strings.Builder{}
	leftTitle := "Available Themes"
	if v.ActivePanel == 0 {
		leftTitle = "▶ " + leftTitle + " (Active)"
	}
	leftCol.WriteString(lipgloss.NewStyle().Bold(true).Foreground(white).Render(leftTitle) + "\n")

	maxVisible := 8
	total := len(v.Themes)
	start := v.ThemeOffset
	if start < 0 {
		start = 0
	}
	end := start + maxVisible
	if end > total {
		end = total
	}

	for i := start; i < end; i++ {
		t := v.Themes[i]
		selected := i == v.ThemeCursor
		prefix := "  "
		name := t.Tokens.DisplayName
		if len(name) > 16 {
			name = name[:16]
		}

		customBadge := ""
		if t.IsCustom {
			customBadge = " [C]"
		}
		activeBadge := ""
		if v.ActiveDeck != nil && (v.ActiveDeck.Config.Theme == t.Tokens.Name || v.ActiveDeck.ThemeName == t.Tokens.Name) {
			activeBadge = " " + lipgloss.NewStyle().Foreground(lime).Bold(true).Render("[ON]")
		}

		if selected {
			prefix = "▸ "
			bg := crimsonDark
			if v.ActivePanel != 0 {
				bg = slate
			}
			style := lipgloss.NewStyle().Bold(true).Foreground(white).Background(bg).Padding(0, 1)
			leftCol.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, style.Render(fmt.Sprintf("%-14s", name)), customBadge, activeBadge))
		} else {
			style := lipgloss.NewStyle().Foreground(light)
			leftCol.WriteString(fmt.Sprintf("%s%-16s%s%s\n", prefix, style.Render(name), customBadge, activeBadge))
		}
	}

	if total > maxVisible {
		leftCol.WriteString(KeyHelpStyle.Render(fmt.Sprintf("(%d/%d themes)\n", v.ThemeCursor+1, total)))
	}

	leftPanelStyle := lipgloss.NewStyle().
		Width(26).
		MarginRight(1)
	if v.ActivePanel == 0 {
		leftPanelStyle = leftPanelStyle.Border(lipgloss.RoundedBorder()).BorderForeground(crimson).Padding(0, 1)
	} else {
		leftPanelStyle = leftPanelStyle.Border(lipgloss.RoundedBorder()).BorderForeground(muted).Padding(0, 1)
	}

	// 2. Right Column: Inspector & Segment Actions
	rightCol := strings.Builder{}
	rightTitle := "Theme Inspector & Segments"
	if v.ActivePanel == 1 {
		rightTitle = "▶ " + rightTitle + " (Active)"
	}
	rightCol.WriteString(lipgloss.NewStyle().Bold(true).Foreground(white).Render(rightTitle) + "\n")

	// Theme Header
	rightCol.WriteString(fmt.Sprintf("%s  %s • %s\n",
		ItemTitle.Render(th.Tokens.DisplayName),
		BadgeStyle.Render(th.Tokens.Name),
		lipgloss.NewStyle().Italic(true).Foreground(faint).Render(th.Tokens.Vibe),
	))

	// Live Miniature Slide Stage Preview
	rightCol.WriteString(v.renderStagePreview(th) + "\n")

	// Palette Swatches
	p := th.Tokens.Palette
	rightCol.WriteString(fmt.Sprintf("Palette: %s %s %s %s\n",
		renderSwatch(p.BgCanvas, "Canvas"),
		renderSwatch(p.BgSurface, "Card"),
		renderSwatch(p.TextPrimary, "Text"),
		renderSwatch(p.AccentPrimary, "Accent"),
	))

	// Live WCAG Contrast Diagnostic
	reports := models.EvaluateThemeContrast(p)
	var passesAll = true
	var failMsg = ""
	for _, r := range reports {
		if !r.PassesAA {
			passesAll = false
			failMsg = fmt.Sprintf("%s vs %s (%.1f:1)", r.FgToken, r.BgToken, r.Ratio)
			break
		}
	}
	if passesAll {
		rightCol.WriteString(fmt.Sprintf("Contrast: %s (WCAG AA/AAA Compliant)\n",
			lipgloss.NewStyle().Bold(true).Foreground(lime).Render("PASS")))
	} else {
		rightCol.WriteString(fmt.Sprintf("Contrast: %s %s • Press [f] to Auto-Tune\n",
			lipgloss.NewStyle().Bold(true).Foreground(crimson).Render("FAIL"),
			lipgloss.NewStyle().Foreground(faint).Render(failMsg)))
	}

	// Modular Segments
	for i, seg := range v.Segments {
		selected := i == v.SegmentCursor && v.ActivePanel == 1
		prefix := "  "
		style := lipgloss.NewStyle().Foreground(faint)
		if selected {
			prefix = "▸ "
			style = lipgloss.NewStyle().Bold(true).Foreground(lime)
		}
		rightCol.WriteString(fmt.Sprintf("%s%s\n", prefix, style.Render(seg)))
	}

	rightPanelStyle := lipgloss.NewStyle().
		Width(60)
	if v.ActivePanel == 1 {
		rightPanelStyle = rightPanelStyle.Border(lipgloss.RoundedBorder()).BorderForeground(lime).Padding(0, 1)
	} else {
		rightPanelStyle = rightPanelStyle.Border(lipgloss.RoundedBorder()).BorderForeground(muted).Padding(0, 1)
	}

	// Join Panels
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
		leftPanelStyle.Render(leftCol.String()),
		rightPanelStyle.Render(rightCol.String()),
	))

	// Keybindings
	s.WriteString("\n" + KeyHelpStyle.Render("[Tab] Switch Panel  [j/k] Navigate  [a/Enter] Apply to Deck  [f] Auto-Fix Contrast  [e] Edit Segment  [c] Clone"))

	if v.StatusMsg != "" {
		style := SuccessStyle
		if v.IsError {
			style = ErrorStyle
		}
		s.WriteString(fmt.Sprintf("\n%s\n", style.Render(v.StatusMsg)))
	}

	return s.String()
}

func (v *ThemeStudioView) renderStagePreview(th *models.Theme) string {
	bg := th.Tokens.Palette.BgCanvas
	cardBg := th.Tokens.Palette.BgSurface
	txt := th.Tokens.Palette.TextPrimary
	acc := th.Tokens.Palette.AccentPrimary

	canvasStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(bg)).
		Foreground(lipgloss.Color(txt)).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(acc)).
		Padding(0, 1).
		Width(56)

	cardStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(cardBg)).
		Foreground(lipgloss.Color(txt)).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(acc)).
		Padding(0, 1).
		Width(23)

	topbar := fmt.Sprintf("SLIDE TITLE                         %s",
		lipgloss.NewStyle().Foreground(lipgloss.Color(acc)).Bold(true).Render("1080P STAGE"),
	)
	cardA := cardStyle.Render("Heading & Body")
	cardB := cardStyle.Render("Card Surfaces")
	cardsRow := lipgloss.JoinHorizontal(lipgloss.Top, cardA, " ", cardB)

	return canvasStyle.Render(fmt.Sprintf("%s\n%s", topbar, cardsRow))
}

func renderSwatch(hexCode, label string) string {
	if hexCode == "" {
		hexCode = "#000000"
	}
	cleanHex := hexCode
	if !strings.HasPrefix(cleanHex, "#") {
		cleanHex = "#" + cleanHex
	}

	box := lipgloss.NewStyle().
		Background(lipgloss.Color(cleanHex)).
		Foreground(lipgloss.Color("#000000")).
		Padding(0, 1).
		Render(" ")

	return fmt.Sprintf("%s %s", box, cleanHex)
}
