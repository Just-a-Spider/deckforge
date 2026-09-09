package theme

import (
	"deckforge/internal/models"
)

func academicCrimsonPreset() models.Theme {
	return models.Theme{
		Tokens: models.ThemeTokens{
			Name:        "academic-crimson",
			DisplayName: "Academic Crimson",
			Vibe:        "Institutional, authoritative, Canva ledger aesthetic",
			Description: "Deep Peruvian crimson with midnight navy headings on clean canvas",
			Palette: models.ThemePalette{
				BgCanvas:          "#f8f9fc",
				BgSurface:         "#ffffff",
				BgSurfaceElevated: "#f1f5f9",
				TextPrimary:       "#0f172a",
				TextSecondary:     "#334155",
				TextMuted:         "#64748b",
				AccentPrimary:     "#820024",
				AccentSecondary:   "#0f172a",
				AccentGlow:        "rgba(130, 0, 36, 0.12)",
				AccentDanger:      "#dc2626",
				BorderColor:       "#820024",
				BorderFaint:       "rgba(15, 23, 42, 0.12)",
			},
			Fonts: models.ThemeFonts{
				Display:   "'Poppins', system-ui, -apple-system, sans-serif",
				Body:      "'Plus Jakarta Sans', system-ui, -apple-system, sans-serif",
				Mono:      "'JetBrains Mono', monospace",
				ImportURL: "https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;700&family=Plus+Jakarta+Sans:wght@400;500;600;700&family=Poppins:ital,wght@0,600;0,700;0,800;1,700&display=swap",
			},
			Geometry: models.ThemeGeometry{
				RadiusCard:  "8px",
				RadiusSm:    "4px",
				RadiusLg:    "12px",
				BorderWidth: "1.5px",
				BorderStyle: "solid",
			},
			Shadows: models.ThemeShadows{
				Card:     "0 4px 16px rgba(15, 23, 42, 0.08)",
				Elevated: "0 8px 24px rgba(15, 23, 42, 0.12)",
				Glow:     "0 0 16px rgba(130, 0, 36, 0.15)",
			},
			Surfaces: models.ThemeSurfaces{
				BackdropBlur:   "0px",
				SurfaceOpacity: "1.0",
			},
			Spacing: models.ThemeSpacing{
				CardPadding: "28px",
			},
			Motion: models.ThemeMotion{
				SlideTransition: "fade",
				AnimationSpeed:  "0.25s",
			},
		},
		TypographyCSS: `
.slide-topbar .topbar-title { font-weight: 800; letter-spacing: -0.02em; }
.panel-headline { font-weight: 700; letter-spacing: -0.01em; }
`,
		SurfacesCSS: `
.glass-panel {
    border: 1.5px solid var(--border-color);
    box-shadow: 0 4px 16px rgba(15, 23, 42, 0.08);
}
.glass-panel.lime-edge { border-left: 5px solid var(--accent-primary); }
`,
		BackdropCSS: `
.stage-backdrop {
    background-color: var(--bg-canvas);
    background-image:
        linear-gradient(rgba(15, 23, 42, 0.03) 1px, transparent 1px),
        linear-gradient(90deg, rgba(15, 23, 42, 0.03) 1px, transparent 1px),
        radial-gradient(circle at 92% 8%, rgba(130, 0, 36, 0.04) 0%, transparent 45%),
        radial-gradient(circle at 8% 92%, rgba(15, 23, 42, 0.04) 0%, transparent 45%);
    background-size: 48px 48px, 48px 48px, 100% 100%, 100% 100%;
}
`,
	}
}

func vintageEditorialPreset() models.Theme {
	return models.Theme{
		Tokens: models.ThemeTokens{
			Name:        "vintage-editorial",
			DisplayName: "Vintage Editorial",
			Vibe:        "Witty, confident, personality-driven publication",
			Description: "Warm cream background with distinctive Fraunces serif and bordered cards",
			Palette: models.ThemePalette{
				BgCanvas:          "#f5f3ee",
				BgSurface:         "#ffffff",
				BgSurfaceElevated: "#ede9df",
				TextPrimary:       "#1a1a1a",
				TextSecondary:     "#4f4f4f",
				TextMuted:         "#808080",
				AccentPrimary:     "#c25e00",
				AccentSecondary:   "#e8d4c0",
				AccentGlow:        "rgba(194, 94, 0, 0.12)",
				AccentDanger:      "#b91c1c",
				BorderColor:       "#1a1a1a",
				BorderFaint:       "#ded8cb",
			},
			Fonts: models.ThemeFonts{
				Display:   "'Fraunces', serif",
				Body:      "'Work Sans', sans-serif",
				Mono:      "'JetBrains Mono', monospace",
				ImportURL: "https://fonts.googleapis.com/css2?family=Fraunces:ital,opsz,wght@0,9..144,700;0,9..144,900;1,9..144,700&family=Work+Sans:wght@400;500;600&family=JetBrains+Mono:wght@400;600&display=swap",
			},
			Geometry: models.ThemeGeometry{
				RadiusCard:  "4px",
				RadiusSm:    "2px",
				RadiusLg:    "6px",
				BorderWidth: "2px",
				BorderStyle: "solid",
			},
			Shadows: models.ThemeShadows{
				Card:     "4px 4px 0 #1a1a1a",
				Elevated: "6px 6px 0 #1a1a1a",
				Flat:     "3px 3px 0 #1a1a1a",
			},
			Surfaces: models.ThemeSurfaces{
				BackdropBlur:   "0px",
				SurfaceOpacity: "1.0",
			},
			Spacing: models.ThemeSpacing{
				CardPadding: "28px",
			},
			Motion: models.ThemeMotion{
				SlideTransition: "fade",
				AnimationSpeed:  "0.25s",
			},
		},
		TypographyCSS: `
.slide-topbar .topbar-title { font-family: 'Fraunces', serif; font-size: 38px; font-weight: 900; }
`,
		SurfacesCSS: `
.glass-panel {
    border: 2px solid #1a1a1a;
    border-radius: 4px;
    box-shadow: 4px 4px 0 #1a1a1a;
}
.glass-panel.lime-edge { border-left: 6px solid var(--accent-primary); }
`,
		BackdropCSS: `
.stage-backdrop {
    background: #f5f3ee;
}
`,
	}
}

func paperAndInkPreset() models.Theme {
	return models.Theme{
		Tokens: models.ThemeTokens{
			Name:        "paper-ink",
			DisplayName: "Paper & Ink",
			Vibe:        "Literary, thoughtful, classical publication",
			Description: "Warm cream background with rich charcoal ink and crimson highlight",
			Palette: models.ThemePalette{
				BgCanvas:          "#faf9f7",
				BgSurface:         "#ffffff",
				BgSurfaceElevated: "#f2efe9",
				TextPrimary:       "#1a1a1a",
				TextSecondary:     "#3b3b3b",
				TextMuted:         "#6e6e6e",
				AccentPrimary:     "#c41e3a",
				AccentSecondary:   "#1a1a1a",
				AccentGlow:        "rgba(196, 30, 58, 0.12)",
				AccentDanger:      "#991b1b",
				BorderColor:       "#c41e3a",
				BorderFaint:       "#e5e0d8",
			},
			Fonts: models.ThemeFonts{
				Display:   "'Cormorant Garamond', serif",
				Body:      "'Source Serif 4', serif",
				Mono:      "'JetBrains Mono', monospace",
				ImportURL: "https://fonts.googleapis.com/css2?family=Cormorant+Garamond:ital,wght@0,600;0,700;1,600&family=Source+Serif+4:wght@400;600&family=JetBrains+Mono:wght@400;600&display=swap",
			},
			Geometry: models.ThemeGeometry{
				RadiusCard:  "4px",
				RadiusSm:    "2px",
				RadiusLg:    "8px",
				BorderWidth: "1px",
				BorderStyle: "solid",
			},
			Shadows: models.ThemeShadows{
				Card:     "0 4px 16px rgba(0, 0, 0, 0.05)",
				Elevated: "0 8px 24px rgba(0, 0, 0, 0.08)",
				Glow:     "0 0 12px rgba(196, 30, 58, 0.15)",
			},
			Surfaces: models.ThemeSurfaces{
				BackdropBlur:   "0px",
				SurfaceOpacity: "1.0",
			},
			Spacing: models.ThemeSpacing{
				CardPadding: "32px",
			},
			Motion: models.ThemeMotion{
				SlideTransition: "fade",
				AnimationSpeed:  "0.25s",
			},
		},
		TypographyCSS: `
.slide-topbar .topbar-title { font-family: 'Cormorant Garamond', serif; font-size: 38px; font-weight: 700; }
`,
		SurfacesCSS: `
.glass-panel {
    border: 1px solid var(--border-faint);
    border-radius: 6px;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
}
.glass-panel.lime-edge { border-left: 4px solid var(--accent-primary); }
`,
		BackdropCSS: `
.stage-backdrop {
    background: #faf9f7;
}
`,
	}
}
