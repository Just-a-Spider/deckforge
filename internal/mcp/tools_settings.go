package mcp

import (
	"encoding/json"
	"fmt"

	"deckforge/internal/editor"
	"deckforge/internal/models"
)

func handleSettingsTools(deckPath, name string, args map[string]interface{}) string {
	switch name {
	case "deckforge_get_settings":
		cfg := models.LoadWorkspaceConfig()
		installed := editor.DetectInstalledEditors()
		res := map[string]interface{}{
			"config":            cfg,
			"config_path":       models.ConfigFilePath(),
			"installed_editors": installed,
		}
		data, _ := json.MarshalIndent(res, "", "  ")
		return string(data)

	case "deckforge_update_settings":
		key, _ := args["key"].(string)
		val, _ := args["value"].(string)
		if key == "" {
			return "Error: key is required"
		}
		cfg := models.LoadWorkspaceConfig()
		if err := cfg.Set(key, val); err != nil {
			return fmt.Sprintf("Error setting config: %v", err)
		}
		return fmt.Sprintf("Successfully updated %s = %s in %s", key, val, models.ConfigFilePath())

	default:
		return fmt.Sprintf("Unknown settings tool %s", name)
	}
}
