package models

import (
	"fmt"
	"sort"
	"strings"
)

// Theme represents a complete segmented and modular presentation theme
type Theme struct {
	Tokens        ThemeTokens
	TypographyCSS string
	SurfacesCSS   string
	BackdropCSS   string
	ComponentsCSS string            // modular component / archetype overrides
	AnimationsCSS string            // modular animation keyframes and transitions
	ExtraCSS      map[string]string // custom stylesheets discovered from styles/ or theme root
	Dir           string            // directory where theme files reside, if on disk
	IsCustom      bool
	Scope         string            // "builtin", "global", "workspace"
}

// CompileFullCSS stitches all theme segments and custom stylesheets together in deterministic cascade order
func (t *Theme) CompileFullCSS() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("/* Theme: %s */\n%s\n", t.Tokens.Name, t.Tokens.RootCSSVariables()))

	if strings.TrimSpace(t.TypographyCSS) != "" {
		sb.WriteString(fmt.Sprintf("\n/* Typography Segment */\n%s\n", t.TypographyCSS))
	}
	if strings.TrimSpace(t.SurfacesCSS) != "" {
		sb.WriteString(fmt.Sprintf("\n/* Surfaces Segment */\n%s\n", t.SurfacesCSS))
	}
	if strings.TrimSpace(t.BackdropCSS) != "" {
		backdrop := t.BackdropCSS
		if strings.Contains(backdrop, ".stage-backdrop") && !strings.Contains(backdrop, ".slide-backdrop") {
			backdrop = strings.ReplaceAll(backdrop, ".stage-backdrop", ".stage-backdrop, .slide-backdrop")
		}
		sb.WriteString(fmt.Sprintf("\n/* Backdrop Segment */\n%s\n", backdrop))
	}
	if strings.TrimSpace(t.ComponentsCSS) != "" {
		sb.WriteString(fmt.Sprintf("\n/* Components Segment */\n%s\n", t.ComponentsCSS))
	}
	if strings.TrimSpace(t.AnimationsCSS) != "" {
		sb.WriteString(fmt.Sprintf("\n/* Animations Segment */\n%s\n", t.AnimationsCSS))
	}

	if len(t.ExtraCSS) > 0 {
		var extraKeys []string
		for k := range t.ExtraCSS {
			extraKeys = append(extraKeys, k)
		}
		sort.Strings(extraKeys)
		for _, k := range extraKeys {
			content := strings.TrimSpace(t.ExtraCSS[k])
			if content != "" {
				sb.WriteString(fmt.Sprintf("\n/* Custom Segment: %s */\n%s\n", k, content))
			}
		}
	}

	return sb.String()
}
