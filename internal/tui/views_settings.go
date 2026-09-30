package tui

import (
	"fmt"
	"strings"

	"deckforge/internal/editor"
	"deckforge/internal/models"
	"deckforge/internal/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SettingsField int

const (
	SettingEditor SettingsField = iota
	SettingTheme
	SettingSlideCount
	SettingPort
	SettingWatch
	SettingCount
)

type SettingsView struct {
	Config        *models.WorkspaceConfig
	FocusIndex    SettingsField
	EditorOptions []string
	EditorIndex   int
	ThemePresets  []models.Theme
	ThemeIndex    int
	Ports         []int
	PortIndex     int
	StatusMsg     string
}

func NewSettingsView(cfg *models.WorkspaceConfig) *SettingsView {
	installedMap := make(map[string]bool)
	for _, inst := range editor.DetectInstalledEditors() {
		installedMap[inst.ID] = true
	}

	editors := []string{"auto"}
	// First add installed known editors
	for _, k := range editor.KnownEditors() {
		if installedMap[k.ID] {
			editors = append(editors, k.ID)
		}
	}
	// Then add remaining known editors
	for _, k := range editor.KnownEditors() {
		if !installedMap[k.ID] {
			editors = append(editors, k.ID)
		}
	}

	// Preserve current preference if custom
	found := false
	for _, ed := range editors {
		if ed == cfg.PreferredEditor {
			found = true
			break
		}
	}
	if !found && cfg.PreferredEditor != "" {
		editors = append(editors, cfg.PreferredEditor)
	}

	presets := theme.GetBuiltinPresets()
	ports := []int{8080, 3000, 8000, 5000, 9000}

	edIdx := 0
	for i, ed := range editors {
		if ed == cfg.PreferredEditor {
			edIdx = i
			break
		}
	}

	thIdx := 0
	for i, th := range presets {
		if th.Tokens.Name == cfg.DefaultTheme {
			thIdx = i
			break
		}
	}

	portIdx := 0
	for i, p := range ports {
		if p == cfg.ServerPort {
			portIdx = i
			break
		}
	}

	return &SettingsView{
		Config:        cfg,
		FocusIndex:    SettingEditor,
		EditorOptions: editors,
		EditorIndex:   edIdx,
		ThemePresets:  presets,
		ThemeIndex:    thIdx,
		Ports:         ports,
		PortIndex:     portIdx,
	}
}

func (v *SettingsView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if v.FocusIndex > 0 {
				v.FocusIndex--
			} else {
				v.FocusIndex = SettingCount - 1
			}
		case "down", "j", "tab":
			v.FocusIndex = (v.FocusIndex + 1) % SettingCount

		case "left", "-", "h":
			v.adjustSetting(-1)
		case "right", "+", "=", "l", "enter", "space":
			v.adjustSetting(1)
		}
	}
	return nil
}

func (v *SettingsView) adjustSetting(delta int) {
	switch v.FocusIndex {
	case SettingEditor:
		v.EditorIndex += delta
		if v.EditorIndex < 0 {
			v.EditorIndex = len(v.EditorOptions) - 1
		} else if v.EditorIndex >= len(v.EditorOptions) {
			v.EditorIndex = 0
		}
		v.Config.PreferredEditor = v.EditorOptions[v.EditorIndex]
		_ = v.Config.Save()
		v.StatusMsg = fmt.Sprintf("Preferred editor saved: %s", v.Config.PreferredEditor)

	case SettingTheme:
		v.ThemeIndex += delta
		if v.ThemeIndex < 0 {
			v.ThemeIndex = len(v.ThemePresets) - 1
		} else if v.ThemeIndex >= len(v.ThemePresets) {
			v.ThemeIndex = 0
		}
		v.Config.DefaultTheme = v.ThemePresets[v.ThemeIndex].Tokens.Name
		_ = v.Config.Save()
		v.StatusMsg = fmt.Sprintf("Default theme saved: %s", v.Config.DefaultTheme)

	case SettingSlideCount:
		v.Config.DefaultSlideCount += delta
		if v.Config.DefaultSlideCount < 1 {
			v.Config.DefaultSlideCount = 1
		} else if v.Config.DefaultSlideCount > 50 {
			v.Config.DefaultSlideCount = 50
		}
		_ = v.Config.Save()
		v.StatusMsg = fmt.Sprintf("Default slide count saved: %d", v.Config.DefaultSlideCount)

	case SettingPort:
		v.PortIndex += delta
		if v.PortIndex < 0 {
			v.PortIndex = len(v.Ports) - 1
		} else if v.PortIndex >= len(v.Ports) {
			v.PortIndex = 0
		}
		v.Config.ServerPort = v.Ports[v.PortIndex]
		_ = v.Config.Save()
		v.StatusMsg = fmt.Sprintf("Default server port saved: %d", v.Config.ServerPort)

	case SettingWatch:
		v.Config.AutoWatch = !v.Config.AutoWatch
		_ = v.Config.Save()
		state := "Enabled"
		if !v.Config.AutoWatch {
			state = "Disabled"
		}
		v.StatusMsg = fmt.Sprintf("Live watcher: %s", state)
	}
}

func (v *SettingsView) View() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("%s\n\n",
		lipgloss.NewStyle().Bold(true).Foreground(crimson).Render("GLOBAL PREFERENCES & SETTINGS"),
	))

	renderRow := func(field SettingsField, title, value, hint string) {
		selected := v.FocusIndex == field
		prefix := "  "
		lblStyle := lipgloss.NewStyle().Foreground(light)
		valStyle := ActiveBadgeStyle

		if selected {
			prefix = "▸ "
			lblStyle = lipgloss.NewStyle().Bold(true).Foreground(white).Background(crimsonDark).Padding(0, 1)
			valStyle = lipgloss.NewStyle().Bold(true).Foreground(white).Background(lime).Padding(0, 1)
		} else {
			valStyle = lipgloss.NewStyle().Foreground(slate).Background(light).Padding(0, 1)
		}

		s.WriteString(fmt.Sprintf("%s%-32s %s  %s\n",
			prefix,
			lblStyle.Render(title),
			valStyle.Render(value),
			lipgloss.NewStyle().Foreground(faint).Render(hint),
		))
	}

	edKey := v.EditorOptions[v.EditorIndex]
	edDisplay := edKey
	if edKey == "auto" {
		edDisplay = "auto ($EDITOR)"
	} else {
		// Check if installed
		for _, inst := range editor.DetectInstalledEditors() {
			if inst.ID == edKey {
				edDisplay = fmt.Sprintf("%s [installed]", edKey)
				break
			}
		}
	}
	renderRow(SettingEditor, "1. External Editor ($EDITOR):", edDisplay, "cursor, code, zed, windsurf, nvim, micro")

	thName := v.ThemePresets[v.ThemeIndex].Tokens.DisplayName
	renderRow(SettingTheme, "2. Default Presentation Theme:", thName, "Preset for new decks")

	slideCountStr := fmt.Sprintf("%d Slides", v.Config.DefaultSlideCount)
	renderRow(SettingSlideCount, "3. Default Slide Count:", slideCountStr, "Initial generated slides")

	portStr := fmt.Sprintf("Port :%d", v.Config.ServerPort)
	renderRow(SettingPort, "4. Presentation Server Port:", portStr, "Local HTTP port")

	watchStr := "Enabled (Hot-Reload)"
	if !v.Config.AutoWatch {
		watchStr = "Disabled"
	}
	renderRow(SettingWatch, "5. Live File Watcher:", watchStr, "Recompile on slide edits")

	s.WriteString(fmt.Sprintf("\n%s\n", KeyHelpStyle.Render("[Enter / Space / +/-] Change value  [j/k] Navigate fields  [←/→] Switch tabs")))

	// Config file location
	s.WriteString(fmt.Sprintf("%s %s\n",
		lipgloss.NewStyle().Foreground(muted).Render("Config:"),
		RootPathStyle.Render(models.ConfigFilePath()),
	))

	if v.StatusMsg != "" {
		s.WriteString(fmt.Sprintf("\n%s\n", SuccessStyle.Render(v.StatusMsg)))
	}

	return s.String()
}
