package theme

import (
	"deckforge/internal/models"
)

func darkBotanicalPreset() models.Theme {
	return models.Theme{
		Tokens: models.ThemeTokens{
			Name:        "dark-botanical",
			DisplayName: "Dark Botanical",
			Vibe:        "Elegant, sophisticated, artistic, premium editorial",
			Description: "Charcoal noir canvas with warm gold, terracotta, and soft blush accents",
			Palette: models.ThemePalette{
				BgCanvas:          "#0f0f0f",
				BgSurface:         "#1a1a1a",
				BgSurfaceElevated: "#262626",
				TextPrimary:       "#e8e4df",
				TextSecondary:     "#b5b0a8",
				TextMuted:         "#7a756e",
				AccentPrimary:     "#d4a574",
				AccentSecondary:   "#e8b4b8",
				AccentGlow:        "rgba(212, 165, 116, 0.15)",
				AccentDanger:      "#ef4444",
				BorderColor:       "#d4a574",
				BorderFaint:       "rgba(255, 255, 255, 0.08)",
			},
			Fonts: models.ThemeFonts{
				Display:   "'Cormorant', serif",
				Body:      "'IBM Plex Sans', sans-serif",
				Mono:      "'JetBrains Mono', monospace",
				ImportURL: "https://fonts.googleapis.com/css2?family=Cormorant:ital,wght@0,600;0,700;1,600&family=IBM+Plex+Sans:wght@300;400;500;600&family=JetBrains+Mono:wght@400;600&display=swap",
			},
			Geometry: models.ThemeGeometry{
				RadiusCard:  "6px",
				RadiusSm:    "3px",
				RadiusLg:    "10px",
				BorderWidth: "1px",
				BorderStyle: "solid",
			},
			Shadows: models.ThemeShadows{
				Card:     "0 8px 30px rgba(0, 0, 0, 0.4)",
				Elevated: "0 14px 40px rgba(0, 0, 0, 0.6)",
				Glow:     "0 0 18px rgba(212, 165, 116, 0.2)",
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
.slide-topbar .topbar-title { font-family: 'Cormorant', serif; font-size: 38px; font-weight: 700; }
`,
		SurfacesCSS: `
.glass-panel {
    border: 1px solid rgba(212, 165, 116, 0.25);
    border-radius: 6px;
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.4);
}
.glass-panel.lime-edge { border-left: 4px solid var(--accent-primary); }
`,
		BackdropCSS: `
.stage-backdrop {
    background-color: var(--bg-canvas);
    background-image:
        radial-gradient(circle at 85% 15%, rgba(212, 165, 116, 0.07) 0%, transparent 45%),
        radial-gradient(circle at 15% 85%, rgba(232, 180, 184, 0.05) 0%, transparent 45%);
}
`,
	}
}

func notebookTabsPreset() models.Theme {
	return models.Theme{
		Tokens: models.ThemeTokens{
			Name:        "notebook-tabs",
			DisplayName: "Notebook Tabs",
			Vibe:        "Editorial, organized, tactile paper feeling",
			Description: "Warm cream paper sheets with pastel tabs on dark charcoal outer",
			Palette: models.ThemePalette{
				BgCanvas:          "#2d2d2d",
				BgSurface:         "#f8f6f1",
				BgSurfaceElevated: "#eae6dc",
				TextPrimary:       "#1a1a1a",
				TextSecondary:     "#444444",
				TextMuted:         "#777777",
				AccentPrimary:     "#98d4bb",
				AccentSecondary:   "#c7b8ea",
				AccentGlow:        "rgba(152, 212, 187, 0.2)",
				AccentDanger:      "#e11d48",
				BorderColor:       "#444444",
				BorderFaint:       "rgba(0, 0, 0, 0.12)",
			},
			Fonts: models.ThemeFonts{
				Display:   "'Bodoni Moda', serif",
				Body:      "'DM Sans', sans-serif",
				Mono:      "'JetBrains Mono', monospace",
				ImportURL: "https://fonts.googleapis.com/css2?family=Bodoni+Moda:ital,wght@0,600;0,700;1,600&family=DM+Sans:wght@400;500;700&family=JetBrains+Mono:wght@400;600&display=swap",
			},
			Geometry: models.ThemeGeometry{
				RadiusCard:  "8px",
				RadiusSm:    "4px",
				RadiusLg:    "12px",
				BorderWidth: "1px",
				BorderStyle: "solid",
			},
			Shadows: models.ThemeShadows{
				Card:     "0 4px 20px rgba(0, 0, 0, 0.15)",
				Elevated: "0 8px 28px rgba(0, 0, 0, 0.25)",
				Glow:     "0 0 14px rgba(152, 212, 187, 0.25)",
			},
			Surfaces: models.ThemeSurfaces{
				BackdropBlur:   "0px",
				SurfaceOpacity: "1.0",
			},
			Spacing: models.ThemeSpacing{
				CardPadding: "28px",
			},
			Motion: models.ThemeMotion{
				SlideTransition: "slide",
				AnimationSpeed:  "0.25s",
			},
		},
		TypographyCSS: `
.slide-topbar .topbar-title { font-family: 'Bodoni Moda', serif; font-size: 36px; font-weight: 700; }
`,
		SurfacesCSS: `
.glass-panel {
    background: #f8f6f1;
    border: 1px solid rgba(0, 0, 0, 0.12);
    border-radius: 8px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}
.glass-panel.lime-edge { border-left: 6px solid var(--accent-primary); }
`,
		BackdropCSS: `
.stage-backdrop {
    background: #2d2d2d;
}
`,
	}
}

func pastelGeometryPreset() models.Theme {
	return models.Theme{
		Tokens: models.ThemeTokens{
			Name:        "pastel-geometry",
			DisplayName: "Pastel Geometry",
			Vibe:        "Approachable, clean, friendly modern cards",
			Description: "Soft pastel sky background with clean cards and violet pill tags",
			Palette: models.ThemePalette{
				BgCanvas:          "#c8d9e6",
				BgSurface:         "#faf9f7",
				BgSurfaceElevated: "#f0ede6",
				TextPrimary:       "#1e293b",
				TextSecondary:     "#475569",
				TextMuted:         "#64748b",
				AccentPrimary:     "#7c6aad",
				AccentSecondary:   "#f0b4d4",
				AccentGlow:        "rgba(124, 106, 173, 0.15)",
				AccentDanger:      "#e11d48",
				BorderColor:       "#7c6aad",
				BorderFaint:       "rgba(0, 0, 0, 0.08)",
			},
			Fonts: models.ThemeFonts{
				Display:   "'Plus Jakarta Sans', sans-serif",
				Body:      "'Plus Jakarta Sans', sans-serif",
				Mono:      "'JetBrains Mono', monospace",
				ImportURL: "https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=JetBrains+Mono:wght@400;600&display=swap",
			},
			Geometry: models.ThemeGeometry{
				RadiusCard:  "14px",
				RadiusSm:    "6px",
				RadiusLg:    "20px",
				BorderWidth: "1px",
				BorderStyle: "solid",
			},
			Shadows: models.ThemeShadows{
				Card:     "0 8px 24px rgba(71, 85, 105, 0.12)",
				Elevated: "0 12px 32px rgba(71, 85, 105, 0.18)",
				Glow:     "0 0 16px rgba(124, 106, 173, 0.18)",
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
.slide-topbar .topbar-title { font-weight: 800; }
`,
		SurfacesCSS: `
.glass-panel {
    border-radius: 14px;
    box-shadow: 0 8px 24px rgba(71, 85, 105, 0.12);
}
.glass-panel.lime-edge { border-left: 5px solid var(--accent-primary); }
`,
		BackdropCSS: `
.stage-backdrop {
    background: #c8d9e6;
}
`,
	}
}

func splitPastelPreset() models.Theme {
	return models.Theme{
		Tokens: models.ThemeTokens{
			Name:        "split-pastel",
			DisplayName: "Split Pastel",
			Vibe:        "Playful, creative, modern split design",
			Description: "Two-color peach and lavender panels with mint accents",
			Palette: models.ThemePalette{
				BgCanvas:          "#f5e6dc",
				BgSurface:         "#ffffff",
				BgSurfaceElevated: "#e4dff0",
				TextPrimary:       "#1a1a1a",
				TextSecondary:     "#4a4a4a",
				TextMuted:         "#777777",
				AccentPrimary:     "#4a7c59",
				AccentSecondary:   "#c8f0d8",
				AccentGlow:        "rgba(74, 124, 89, 0.15)",
				AccentDanger:      "#e11d48",
				BorderColor:       "#4a7c59",
				BorderFaint:       "rgba(0, 0, 0, 0.08)",
			},
			Fonts: models.ThemeFonts{
				Display:   "'Outfit', sans-serif",
				Body:      "'Outfit', sans-serif",
				Mono:      "'JetBrains Mono', monospace",
				ImportURL: "https://fonts.googleapis.com/css2?family=Outfit:wght@400;500;600;700;800&family=JetBrains+Mono:wght@400;600&display=swap",
			},
			Geometry: models.ThemeGeometry{
				RadiusCard:  "12px",
				RadiusSm:    "6px",
				RadiusLg:    "18px",
				BorderWidth: "1px",
				BorderStyle: "solid",
			},
			Shadows: models.ThemeShadows{
				Card:     "0 6px 20px rgba(0, 0, 0, 0.06)",
				Elevated: "0 10px 28px rgba(0, 0, 0, 0.1)",
				Glow:     "0 0 14px rgba(74, 124, 89, 0.18)",
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
`,
		SurfacesCSS: `
.glass-panel {
    border-radius: 12px;
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.06);
}
.glass-panel.lime-edge { border-left: 5px solid var(--accent-primary); }
`,
		BackdropCSS: `
.stage-backdrop {
    background: linear-gradient(90deg, #f5e6dc 50%, #e4dff0 50%);
}
`,
	}
}
