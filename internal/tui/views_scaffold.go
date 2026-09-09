package tui

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"deckforge/internal/models"
	"deckforge/internal/scaffold"
	"deckforge/internal/theme"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ScaffoldView struct {
	ActiveRoot   string
	FocusIndex   int
	IsEditing    bool
	Inputs       []textinput.Model
	ThemePresets []models.Theme
	ThemeIndex   int
	StatusMsg    string
	IsError      bool
	CreatedDeck  *models.Deck
}

const (
	scaffoldFieldTitle  = 0
	scaffoldFieldPath   = 1
	scaffoldFieldCount  = 2
	scaffoldFieldTheme  = 3
	scaffoldFieldSubmit = 4
	scaffoldFieldTotal  = 5
)

func NewScaffoldView(root string) *ScaffoldView {
	presets := theme.GetBuiltinPresets()

	titleInput := textinput.New()
	titleInput.Placeholder = "Deck Name (e.g. distributed_systems)"
	titleInput.CharLimit = 64
	titleInput.Width = 50

	pathInput := textinput.New()
	pathInput.Placeholder = "/custom/path/to/my_deck (or relative)"
	pathInput.CharLimit = 256
	pathInput.Width = 50

	countInput := textinput.New()
	countInput.Placeholder = "5"
	countInput.CharLimit = 3
	countInput.Width = 10
	countInput.SetValue("5")

	return &ScaffoldView{
		ActiveRoot:   root,
		FocusIndex:   0,
		IsEditing:    false,
		Inputs:       []textinput.Model{titleInput, pathInput, countInput},
		ThemePresets: presets,
		ThemeIndex:   0,
	}
}

func (v *ScaffoldView) SetActiveRoot(root string) {
	v.ActiveRoot = root
}

func (v *ScaffoldView) Update(msg tea.Msg) (tea.Cmd, bool) {
	if v.IsEditing {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				if v.FocusIndex < len(v.Inputs) {
					v.Inputs[v.FocusIndex].Blur()
				}
				v.IsEditing = false
				return nil, false

			case "enter":
				if v.FocusIndex < len(v.Inputs) {
					v.Inputs[v.FocusIndex].Blur()
				}
				v.IsEditing = false
				v.FocusIndex = (v.FocusIndex + 1) % scaffoldFieldTotal
				return nil, false

			case "tab":
				if v.FocusIndex < len(v.Inputs) {
					v.Inputs[v.FocusIndex].Blur()
				}
				v.FocusIndex = (v.FocusIndex + 1) % scaffoldFieldTotal
				if v.FocusIndex < len(v.Inputs) {
					v.Inputs[v.FocusIndex].Focus()
				} else {
					v.IsEditing = false
				}
				return nil, false
			}
		}

		if v.FocusIndex < len(v.Inputs) {
			var cmd tea.Cmd
			v.Inputs[v.FocusIndex], cmd = v.Inputs[v.FocusIndex].Update(msg)
			return cmd, false
		}
		return nil, false
	}

	// Navigation Mode
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			v.FocusIndex--
			if v.FocusIndex < 0 {
				v.FocusIndex = scaffoldFieldTotal - 1
			}
			return nil, false

		case "down", "j", "tab":
			v.FocusIndex = (v.FocusIndex + 1) % scaffoldFieldTotal
			return nil, false

		case "shift+tab":
			v.FocusIndex--
			if v.FocusIndex < 0 {
				v.FocusIndex = scaffoldFieldTotal - 1
			}
			return nil, false

		case "left", "h":
			if v.FocusIndex == scaffoldFieldTheme {
				v.ThemeIndex--
				if v.ThemeIndex < 0 {
					v.ThemeIndex = len(v.ThemePresets) - 1
				}
				return nil, false
			}

		case "right", "l":
			if v.FocusIndex == scaffoldFieldTheme {
				v.ThemeIndex = (v.ThemeIndex + 1) % len(v.ThemePresets)
				return nil, false
			}

		case "enter", "i", " ":
			if v.FocusIndex < len(v.Inputs) {
				// Enter edit mode
				v.IsEditing = true
				v.Inputs[v.FocusIndex].Focus()
				return textinput.Blink, false
			} else if v.FocusIndex == scaffoldFieldTheme {
				// Cycle theme
				v.ThemeIndex = (v.ThemeIndex + 1) % len(v.ThemePresets)
				return nil, false
			} else if v.FocusIndex == scaffoldFieldSubmit {
				// Create presentation
				return v.submit()
			}
		}
	}

	return nil, false
}

func (v *ScaffoldView) submit() (tea.Cmd, bool) {
	name := strings.TrimSpace(v.Inputs[scaffoldFieldTitle].Value())
	if name == "" {
		v.StatusMsg = "Deck name cannot be empty. Press Enter on field 1 to set name."
		v.IsError = true
		v.FocusIndex = scaffoldFieldTitle
		return nil, false
	}

	targetPath := strings.TrimSpace(v.Inputs[scaffoldFieldPath].Value())
	if targetPath == "" {
		targetPath = filepath.Join(v.ActiveRoot, "decks", name)
	} else if !filepath.IsAbs(targetPath) {
		targetPath = filepath.Join(v.ActiveRoot, targetPath)
	}

	countStr := strings.TrimSpace(v.Inputs[scaffoldFieldCount].Value())
	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		count = 5
	}

	selectedTheme := v.ThemePresets[v.ThemeIndex].Tokens.Name

	deck, err := scaffold.ScaffoldDeck(targetPath, name, selectedTheme, count)
	if err != nil {
		v.StatusMsg = fmt.Sprintf("Scaffold error: %v", err)
		v.IsError = true
		return nil, false
	}

	v.CreatedDeck = deck
	v.StatusMsg = fmt.Sprintf("Scaffolded %d slides in %s", len(deck.Slides), targetPath)
	v.IsError = false

	// Reset inputs
	v.Inputs[scaffoldFieldTitle].Reset()
	v.Inputs[scaffoldFieldPath].Reset()
	v.Inputs[scaffoldFieldCount].SetValue("5")
	v.FocusIndex = 0
	v.IsEditing = false

	return nil, true
}

func (v *ScaffoldView) View() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("%s\n",
		lipgloss.NewStyle().Bold(true).Foreground(crimson).Render("SCAFFOLD NEW PRESENTATION WIZARD"),
	))

	// Mode Banner
	if v.IsEditing {
		s.WriteString(lipgloss.NewStyle().Foreground(lime).Bold(true).Render("▶ INPUT MODE: Type your value. Press Enter or Esc when done.\n\n"))
	} else {
		s.WriteString(lipgloss.NewStyle().Foreground(faint).Italic(true).Render("▶ NAVIGATION MODE: j/k or Up/Down selects field. Enter edits. ←/→ switches tabs.\n\n"))
	}

	renderField := func(idx int, label, valView string) {
		selected := v.FocusIndex == idx
		prefix := "  "
		labelStyle := lipgloss.NewStyle().Foreground(light)

		if selected {
			prefix = "▸ "
			if v.IsEditing {
				labelStyle = lipgloss.NewStyle().Bold(true).Foreground(lime)
			} else {
				labelStyle = lipgloss.NewStyle().Bold(true).Foreground(white).Background(crimsonDark).Padding(0, 1)
			}
		}

		s.WriteString(fmt.Sprintf("%s%-34s %s\n", prefix, labelStyle.Render(label), valView))
	}

	// 1. Name
	val1 := v.Inputs[scaffoldFieldTitle].View()
	if !v.IsEditing || v.FocusIndex != 0 {
		val := v.Inputs[scaffoldFieldTitle].Value()
		if val == "" {
			val = lipgloss.NewStyle().Foreground(muted).Render("(Press Enter to type deck name)")
		}
		val1 = val
	}
	renderField(0, "1. Deck Identifier / Name:", val1)

	// 2. Path
	val2 := v.Inputs[scaffoldFieldPath].View()
	if !v.IsEditing || v.FocusIndex != 1 {
		val := v.Inputs[scaffoldFieldPath].Value()
		if val == "" {
			name := v.Inputs[scaffoldFieldTitle].Value()
			if name == "" {
				name = "<name>"
			}
			val = lipgloss.NewStyle().Foreground(muted).Render(fmt.Sprintf("./decks/%s (default)", name))
		}
		val2 = val
	}
	renderField(1, "2. Destination Directory:", val2)

	// 3. Count
	val3 := v.Inputs[scaffoldFieldCount].View()
	if !v.IsEditing || v.FocusIndex != 2 {
		val3 = fmt.Sprintf("%s slides (Press Enter to edit count)", v.Inputs[scaffoldFieldCount].Value())
	}
	renderField(2, "3. Initial Slide Count (1-50):", val3)

	// 4. Theme
	activePreset := v.ThemePresets[v.ThemeIndex]
	themeSelector := fmt.Sprintf("%s %s %s  %s",
		lipgloss.NewStyle().Bold(true).Foreground(crimson).Render("◀"),
		lipgloss.NewStyle().Bold(true).Foreground(white).Background(navy).Padding(0, 1).Render(activePreset.Tokens.DisplayName),
		lipgloss.NewStyle().Bold(true).Foreground(crimson).Render("▶"),
		lipgloss.NewStyle().Foreground(faint).Render(activePreset.Tokens.Vibe),
	)
	renderField(3, "4. Base Theme (←/→ to cycle):", themeSelector)

	// 5. Submit Button
	submitStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(slate).
		Padding(0, 2)
	if v.FocusIndex == scaffoldFieldSubmit {
		submitStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Background(lime).
			Padding(0, 2)
	}

	submitBtn := submitStyle.Render("CREATE PRESENTATION")
	if v.FocusIndex == scaffoldFieldSubmit {
		s.WriteString(fmt.Sprintf("\n▸ %s  %s\n", submitBtn, KeyHelpStyle.Render("[Press Enter to Scaffold]")))
	} else {
		s.WriteString(fmt.Sprintf("\n  %s\n", submitBtn))
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
