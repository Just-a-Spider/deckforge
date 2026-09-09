package models

import (
	"fmt"
	"strings"
)

// ThemePalette defines core semantic color tokens
type ThemePalette struct {
	BgCanvas          string `json:"bgCanvas"`
	BgSurface         string `json:"bgSurface"`
	BgSurfaceElevated string `json:"bgSurfaceElevated,omitempty"`
	TextPrimary       string `json:"textPrimary"`
	TextSecondary     string `json:"textSecondary"`
	TextMuted         string `json:"textMuted,omitempty"`
	AccentPrimary     string `json:"accentPrimary"`
	AccentSecondary   string `json:"accentSecondary,omitempty"`
	AccentGlow        string `json:"accentGlow,omitempty"`
	AccentDanger      string `json:"accentDanger,omitempty"`
	BorderColor       string `json:"borderColor"`
	BorderFaint       string `json:"borderFaint,omitempty"`
}

// ThemeFonts defines typography tokens
type ThemeFonts struct {
	Display   string `json:"display"`
	Body      string `json:"body"`
	Mono      string `json:"mono"`
	ImportURL string `json:"importUrl,omitempty"`
}

// ThemeGeometry defines corner radii and border metrics
type ThemeGeometry struct {
	RadiusCard  string `json:"radiusCard,omitempty"`  // e.g. "12px", "0px", "4px"
	RadiusSm    string `json:"radiusSm,omitempty"`    // e.g. "4px", "0px"
	RadiusLg    string `json:"radiusLg,omitempty"`    // e.g. "16px", "0px"
	BorderWidth string `json:"borderWidth,omitempty"` // e.g. "1.5px", "2px", "1px"
	BorderStyle string `json:"borderStyle,omitempty"` // e.g. "solid", "dashed"
}

// ThemeShadows defines elevation levels and optical glow effects
type ThemeShadows struct {
	Card     string `json:"card,omitempty"`     // standard card shadow
	Elevated string `json:"elevated,omitempty"` // modal / dropdown elevation
	Glow     string `json:"glow,omitempty"`     // focal element glow
	Flat     string `json:"flat,omitempty"`     // hard brutalist shadow
}

// ThemeSurfaces defines material properties (blur, translucency)
type ThemeSurfaces struct {
	BackdropBlur   string `json:"backdropBlur,omitempty"`   // e.g. "12px", "0px", "24px"
	SurfaceOpacity string `json:"surfaceOpacity,omitempty"` // e.g. "0.85", "1.0"
}

// ThemeSpacing defines component-level density and inner margins
type ThemeSpacing struct {
	CardPadding string `json:"cardPadding,omitempty"` // e.g. "28px", "20px", "32px"
}

// ThemeMotion defines transitions and keyframe speed
type ThemeMotion struct {
	SlideTransition string `json:"slideTransition,omitempty"` // "fade", "slide", "none"
	AnimationSpeed  string `json:"animationSpeed,omitempty"`  // "0.25s", "0.15s"
}

// ThemeTokens defines the structured JSON metadata for a theme
type ThemeTokens struct {
	Name        string        `json:"name"`
	DisplayName string        `json:"displayName"`
	Description string        `json:"description,omitempty"`
	Vibe        string        `json:"vibe,omitempty"`
	Palette     ThemePalette  `json:"palette"`
	Fonts       ThemeFonts    `json:"fonts"`
	Geometry    ThemeGeometry `json:"geometry,omitempty"`
	Shadows     ThemeShadows  `json:"shadows,omitempty"`
	Surfaces    ThemeSurfaces `json:"surfaces,omitempty"`
	Spacing     ThemeSpacing  `json:"spacing,omitempty"`
	Motion      ThemeMotion   `json:"motion,omitempty"`
}

// RootCSSVariables generates the :root CSS variable block from tokens
func (t *ThemeTokens) RootCSSVariables() string {
	bgElev := t.Palette.BgSurfaceElevated
	if bgElev == "" {
		bgElev = t.Palette.BgSurface
	}
	textMuted := t.Palette.TextMuted
	if textMuted == "" {
		textMuted = t.Palette.TextSecondary
	}
	accSec := t.Palette.AccentSecondary
	if accSec == "" {
		accSec = t.Palette.AccentPrimary
	}
	accGlow := t.Palette.AccentGlow
	if accGlow == "" {
		accGlow = "rgba(0, 0, 0, 0.1)"
	}
	accDanger := t.Palette.AccentDanger
	if accDanger == "" {
		accDanger = "#ef4444"
	}
	borderFaint := t.Palette.BorderFaint
	if borderFaint == "" {
		borderFaint = "rgba(255, 255, 255, 0.12)"
	}

	// Geometry defaults
	radCard := t.Geometry.RadiusCard
	if radCard == "" {
		radCard = "12px"
	}
	radSm := t.Geometry.RadiusSm
	if radSm == "" {
		radSm = "4px"
	}
	radLg := t.Geometry.RadiusLg
	if radLg == "" {
		radLg = "16px"
	}
	bWidth := t.Geometry.BorderWidth
	if bWidth == "" {
		bWidth = "1.5px"
	}
	bStyle := t.Geometry.BorderStyle
	if bStyle == "" {
		bStyle = "solid"
	}

	// Shadow defaults
	shCard := t.Shadows.Card
	if shCard == "" {
		shCard = "0 10px 30px rgba(0, 0, 0, 0.25)"
	}
	shElev := t.Shadows.Elevated
	if shElev == "" {
		shElev = "0 20px 40px rgba(0, 0, 0, 0.35)"
	}
	shGlow := t.Shadows.Glow
	if shGlow == "" {
		shGlow = "0 0 20px " + accGlow
	}

	// Surface defaults
	blur := t.Surfaces.BackdropBlur
	if blur == "" {
		blur = "12px"
	}
	opacity := t.Surfaces.SurfaceOpacity
	if opacity == "" {
		opacity = "0.85"
	}

	// Spacing defaults
	cardPadding := t.Spacing.CardPadding
	if cardPadding == "" {
		cardPadding = "28px"
	}

	// Motion defaults
	transition := t.Motion.SlideTransition
	if transition == "" {
		transition = "fade"
	}
	speed := t.Motion.AnimationSpeed
	if speed == "" {
		speed = "0.25s"
	}

	importCSS := ""
	if t.Fonts.ImportURL != "" {
		importCSS = fmt.Sprintf("@import url('%s');\n\n", t.Fonts.ImportURL)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`%s:root {
    --bg-canvas: %s;
    --bg-surface: %s;
    --bg-surface-elevated: %s;
    --text-primary: %s;
    --text-secondary: %s;
    --text-muted: %s;
    --accent-primary: %s;
    --accent-secondary: %s;
    --accent-glow: %s;
    --accent-danger: %s;
    --border-color: %s;
    --border-faint: %s;
    --font-display: %s;
    --font-body: %s;
    --font-mono: %s;
    --radius-card: %s;
    --radius-sm: %s;
    --radius-lg: %s;
    --border-width: %s;
    --border-style: %s;
    --shadow-card: %s;
    --shadow-elevated: %s;
    --shadow-glow: %s;
    --backdrop-blur: %s;
    --surface-opacity: %s;
    --card-padding: %s;
    --slide-transition: %s;
    --animation-speed: %s;
}`,
		importCSS,
		t.Palette.BgCanvas,
		t.Palette.BgSurface,
		bgElev,
		t.Palette.TextPrimary,
		t.Palette.TextSecondary,
		textMuted,
		t.Palette.AccentPrimary,
		accSec,
		accGlow,
		accDanger,
		t.Palette.BorderColor,
		borderFaint,
		t.Fonts.Display,
		t.Fonts.Body,
		t.Fonts.Mono,
		radCard,
		radSm,
		radLg,
		bWidth,
		bStyle,
		shCard,
		shElev,
		shGlow,
		blur,
		opacity,
		cardPadding,
		transition,
		speed,
	))

	return sb.String()
}
