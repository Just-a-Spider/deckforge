package theme

import (
	"deckforge/internal/models"
)

// GetBuiltinPresets returns all 14 curated themes derived from frontend-slides skill
func GetBuiltinPresets() []models.Theme {
	return []models.Theme{
		academicCrimsonPreset(),
		cyberDarkPreset(),
		swissMinimalPreset(),
		boldSignalPreset(),
		electricStudioPreset(),
		creativeVoltagePreset(),
		darkBotanicalPreset(),
		notebookTabsPreset(),
		pastelGeometryPreset(),
		splitPastelPreset(),
		vintageEditorialPreset(),
		neonCyberPreset(),
		terminalGreenPreset(),
		paperAndInkPreset(),
	}
}

// GetPresetByName finds a preset by its identifier
func GetPresetByName(name string) (models.Theme, bool) {
	for _, p := range GetBuiltinPresets() {
		if p.Tokens.Name == name {
			return p, true
		}
	}
	return models.Theme{}, false
}
