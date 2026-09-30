package editor

import (
	"os"
	"os/exec"
	"strings"
)

// EditorSpec defines configuration for an external editor candidate
type EditorSpec struct {
	ID          string   `json:"id"`
	DisplayName string   `json:"displayName"`
	Binaries    []string `json:"binaries"`
	WaitFlags   []string `json:"waitFlags,omitempty"`
	IsTerminal  bool     `json:"isTerminal"`
}

// KnownEditors returns the curated registry of supported modern editors
func KnownEditors() []EditorSpec {
	return []EditorSpec{
		{
			ID:          "cursor",
			DisplayName: "Cursor AI",
			Binaries:    []string{"cursor"},
			WaitFlags:   []string{"--wait"},
			IsTerminal:  false,
		},
		{
			ID:          "code",
			DisplayName: "VS Code",
			Binaries:    []string{"code", "code-insiders"},
			WaitFlags:   []string{"--wait"},
			IsTerminal:  false,
		},
		{
			ID:          "zed",
			DisplayName: "Zed",
			Binaries:    []string{"zed", "zeditor"},
			WaitFlags:   []string{"--wait"},
			IsTerminal:  false,
		},
		{
			ID:          "windsurf",
			DisplayName: "Windsurf",
			Binaries:    []string{"windsurf"},
			WaitFlags:   []string{"--wait"},
			IsTerminal:  false,
		},
		{
			ID:          "subl",
			DisplayName: "Sublime Text",
			Binaries:    []string{"subl", "sublime_text"},
			WaitFlags:   []string{"-w"},
			IsTerminal:  false,
		},
		{
			ID:          "nvim",
			DisplayName: "Neovim",
			Binaries:    []string{"nvim"},
			WaitFlags:   nil,
			IsTerminal:  true,
		},
		{
			ID:          "helix",
			DisplayName: "Helix",
			Binaries:    []string{"hx", "helix"},
			WaitFlags:   nil,
			IsTerminal:  true,
		},
		{
			ID:          "micro",
			DisplayName: "Micro",
			Binaries:    []string{"micro"},
			WaitFlags:   nil,
			IsTerminal:  true,
		},
		{
			ID:          "vim",
			DisplayName: "Vim",
			Binaries:    []string{"vim", "vi"},
			WaitFlags:   nil,
			IsTerminal:  true,
		},
		{
			ID:          "emacs",
			DisplayName: "Emacs",
			Binaries:    []string{"emacs"},
			WaitFlags:   []string{"-nw"},
			IsTerminal:  true,
		},
		{
			ID:          "nano",
			DisplayName: "Nano",
			Binaries:    []string{"nano"},
			WaitFlags:   nil,
			IsTerminal:  true,
		},
	}
}

// InstalledEditor records a detected editor on the current system
type InstalledEditor struct {
	ID          string   `json:"id"`
	DisplayName string   `json:"displayName"`
	BinaryPath  string   `json:"binaryPath"`
	WaitFlags   []string `json:"waitFlags,omitempty"`
	IsTerminal  bool     `json:"isTerminal"`
}

// DetectInstalledEditors scans PATH for all known editors
func DetectInstalledEditors() []InstalledEditor {
	var installed []InstalledEditor
	for _, spec := range KnownEditors() {
		for _, bin := range spec.Binaries {
			if path, err := exec.LookPath(bin); err == nil {
				installed = append(installed, InstalledEditor{
					ID:          spec.ID,
					DisplayName: spec.DisplayName,
					BinaryPath:  path,
					WaitFlags:   spec.WaitFlags,
					IsTerminal:  spec.IsTerminal,
				})
				break
			}
		}
	}
	return installed
}

// ResolveEditorCommand resolves the binary and arguments for an editor preference
func ResolveEditorCommand(pref string) (string, []string) {
	cleanPref := strings.TrimSpace(pref)

	if cleanPref != "" && cleanPref != "auto" {
		// 1. Check known editors by ID or binary name
		for _, spec := range KnownEditors() {
			if strings.EqualFold(cleanPref, spec.ID) || containsString(spec.Binaries, cleanPref) {
				for _, bin := range spec.Binaries {
					if path, err := exec.LookPath(bin); err == nil {
						return path, spec.WaitFlags
					}
				}
			}
		}

		// 2. Direct binary or custom command lookup
		parts := strings.Fields(cleanPref)
		if len(parts) > 0 {
			if path, err := exec.LookPath(parts[0]); err == nil {
				return path, parts[1:]
			}
		}
	}

	// 3. Environment variable fallbacks ($EDITOR, $VISUAL)
	if ed := os.Getenv("EDITOR"); ed != "" {
		parts := strings.Fields(ed)
		if path, err := exec.LookPath(parts[0]); err == nil {
			return path, parts[1:]
		}
	}
	if vis := os.Getenv("VISUAL"); vis != "" {
		parts := strings.Fields(vis)
		if path, err := exec.LookPath(parts[0]); err == nil {
			return path, parts[1:]
		}
	}

	// 4. Auto-detect first installed editor in priority order
	installed := DetectInstalledEditors()
	if len(installed) > 0 {
		return installed[0].BinaryPath, installed[0].WaitFlags
	}

	return "nano", nil
}

func containsString(slice []string, val string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, val) {
			return true
		}
	}
	return false
}
