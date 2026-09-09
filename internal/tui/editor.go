package tui

import (
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// EditorFinishedMsg is sent after the external editor process terminates
type EditorFinishedMsg struct {
	Err error
}

// FindPreferredEditor resolves editor executable and initial flags based on preference
func FindPreferredEditor(pref string) (string, []string) {
	if pref != "" && pref != "auto" {
		switch pref {
		case "code", "code --wait":
			if path, err := exec.LookPath("code"); err == nil {
				return path, []string{"--wait"}
			}
		default:
			if path, err := exec.LookPath(pref); err == nil {
				return path, nil
			}
		}
	}

	if ed := os.Getenv("EDITOR"); ed != "" {
		return ed, nil
	}
	if vis := os.Getenv("VISUAL"); vis != "" {
		return vis, nil
	}

	for _, cand := range []string{"nvim", "vim", "nano", "vi"} {
		if path, err := exec.LookPath(cand); err == nil {
			return path, nil
		}
	}

	if codePath, err := exec.LookPath("code"); err == nil {
		return codePath, []string{"--wait"}
	}

	return "nano", nil
}

// OpenInEditor spawns the external editor suspended in the foreground
func OpenInEditor(pref string, filePaths ...string) tea.Cmd {
	ed, args := FindPreferredEditor(pref)
	allArgs := append(args, filePaths...)
	c := exec.Command(ed, allArgs...)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return EditorFinishedMsg{Err: err}
	})
}
