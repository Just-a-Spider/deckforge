package editor

import (
	"testing"
)

func TestDetectInstalledEditors(t *testing.T) {
	installed := DetectInstalledEditors()
	// All development machines should have at least one editor (e.g. nano, vim, vi, code)
	if len(installed) == 0 {
		t.Logf("No editors detected via PATH")
	} else {
		t.Logf("Detected %d installed editors", len(installed))
		for _, ed := range installed {
			t.Logf("  - %s (%s): %s (wait: %v)", ed.ID, ed.DisplayName, ed.BinaryPath, ed.WaitFlags)
		}
	}
}

func TestResolveEditorCommand(t *testing.T) {
	// Nano or vim should resolve
	bin, args := ResolveEditorCommand("nano")
	if bin == "" {
		t.Fatalf("Expected nano binary path, got empty")
	}
	t.Logf("nano resolved to: %s with args: %v", bin, args)

	// "auto" should resolve to a valid candidate
	autoBin, autoArgs := ResolveEditorCommand("auto")
	if autoBin == "" {
		t.Fatalf("Expected auto to resolve to an editor")
	}
	t.Logf("auto resolved to: %s with args: %v", autoBin, autoArgs)
}
