package theme

import (
	"bytes"
	"fmt"
	"os/exec"
)

// CompileSCSS attempts to transpile SCSS content to standard CSS using local sass tools.
// If no sass compiler is available on the system, it passes the raw content through,
// allowing modern browser CSS native nesting and features to execute directly.
func CompileSCSS(content string) (string, error) {
	if content == "" {
		return "", nil
	}

	// 1. Check for native standalone sass or dart-sass
	sassBin, err := exec.LookPath("sass")
	if err != nil {
		sassBin, err = exec.LookPath("dart-sass")
	}

	if err == nil {
		cmd := exec.Command(sassBin, "--stdin", "--no-source-map", "--style=expanded")
		cmd.Stdin = bytes.NewBufferString(content)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err == nil {
			return out.String(), nil
		}
	}

	// 2. Check for npx if node is installed
	if npxBin, err := exec.LookPath("npx"); err == nil {
		cmd := exec.Command(npxBin, "-y", "sass", "--stdin", "--no-source-map", "--style=expanded")
		cmd.Stdin = bytes.NewBufferString(content)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err == nil {
			return out.String(), nil
		}
	}

	// 3. Fallback: return raw content (modern native CSS handles nesting, variables, etc.)
	return content, fmt.Errorf("no sass compiler found; using native modern CSS")
}
