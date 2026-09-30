package tui

import (
	"os/exec"

	"deckforge/internal/editor"

	tea "github.com/charmbracelet/bubbletea"
)

// EditorFinishedMsg is sent after the external editor process terminates
type EditorFinishedMsg struct {
	Err error
}

// FindPreferredEditor resolves editor executable and initial flags based on preference
func FindPreferredEditor(pref string) (string, []string) {
	return editor.ResolveEditorCommand(pref)
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
