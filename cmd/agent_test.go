package cmd

import (
	"os"
	"testing"

	"deckforge/internal/models"
)

func TestAgentRouterAndLocalConfig(t *testing.T) {
	tempConfigDir, err := os.MkdirTemp("", "deckforge-agent-cfg-*")
	if err != nil {
		t.Fatalf("Failed to create temp config dir: %v", err)
	}
	defer os.RemoveAll(tempConfigDir)

	os.Setenv("DECKFORGE_CONFIG_DIR", tempConfigDir)
	defer os.Unsetenv("DECKFORGE_CONFIG_DIR")

	tempWorkspace, err := os.MkdirTemp("", "deckforge-agent-ws-*")
	if err != nil {
		t.Fatalf("Failed to create temp ws: %v", err)
	}
	defer os.RemoveAll(tempWorkspace)

	// Verify local config set and save
	cfg := models.LoadWorkspaceConfigForRoot(tempWorkspace)
	if err := cfg.SetLocal("preferredEditor", "nano", tempWorkspace); err != nil {
		t.Fatalf("SetLocal failed: %v", err)
	}

	localFile := models.WorkspaceConfigFilePath(tempWorkspace)
	if _, err := os.Stat(localFile); err != nil {
		t.Fatalf("Expected %s to exist: %v", localFile, err)
	}

	loaded := models.LoadWorkspaceConfigForRoot(tempWorkspace)
	if loaded.PreferredEditor != "nano" {
		t.Fatalf("Expected preferredEditor nano, got %s", loaded.PreferredEditor)
	}
}
