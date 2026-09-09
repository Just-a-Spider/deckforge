package theme

import (
	"deckforge/internal/models"
)

func swissMinimalPreset() models.Theme {
	return models.Theme{
		Tokens: models.ThemeTokens{
			Name:        "swiss-minimal",
			DisplayName: "Swiss Minimal",
			Vibe:        "Stark Bauhaus, typographic hierarchy, brutalist discipline",
			Description: "High-contrast monochrome with bold red focus accents",
			Palette: models.ThemePalette{
				BgCanvas:          "#ffffff",
				BgSurface:         "#ffffff",
				BgSurfaceElevated: "#f4f4f5",
				TextPrimary:       "#000000",
				TextSecondary:     "#3f3f46",
				TextMuted:         "#71717a",
				AccentPrimary:     "#ff3300",
				AccentSecondary:   "#000000",
				AccentGlow:        "rgba(255, 51, 0, 0.12)",
				AccentDanger:      "#ef4444",
				BorderColor:       "#000000",
				BorderFaint:       "#e4e4e7",
			},
			Fonts: models.ThemeFonts{
				Display:   "'Archivo', sans-serif",
				Body:      "'Nunito', sans-serif",
				Mono:      "'JetBrains Mono', monospace",
				ImportURL: "https://fonts.googleapis.com/css2?family=Archivo:wght@700;800;900&family=Nunito:wght@400;500;600;700&family=JetBrains+Mono:wght@500;700&display=swap",
			},
			Geometry: models.ThemeGeometry{
				RadiusCard:  "0px",
				RadiusSm:    "0px",
				RadiusLg:    "0px",
				BorderWidth: "2px",
				BorderStyle: "solid",
			},
			Shadows: models.ThemeShadows{
				Card:     "6px 6px 0 #000000",
				Elevated: "10px 10px 0 #000000",
				Flat:     "4px 4px 0 #000000",
			},
			Surfaces: models.ThemeSurfaces{
				BackdropBlur:   "0px",
				SurfaceOpacity: "1.0",
			},
			Spacing: models.ThemeSpacing{
				CardPadding: "24px",
			},
			Motion: models.ThemeMotion{
				SlideTransition: "none",
				AnimationSpeed:  "0.15s",
			},
		},
		TypographyCSS: `
.slide-topbar .topbar-title { font-weight: 900; letter-spacing: -0.04em; }
.panel-headline { font-weight: 800; letter-spacing: -0.02em; }
`,
		SurfacesCSS: `
.glass-panel {
    border: 2px solid #000000;
    border-radius: 0;
    box-shadow: 6px 6px 0 #000000;
}
.glass-panel.lime-edge { border-left: 6px solid var(--accent-primary); }
`,
		BackdropCSS: `
.stage-backdrop {
    background-color: var(--bg-canvas);
    background-image:
        linear-gradient(#000000 1px, transparent 1px),
        linear-gradient(90deg, #000000 1px, transparent 1px);
    background-size: 60px 60px;
    opacity: 0.04;
}
`,
	}
}

func boldSignalPreset() models.Theme {
	return models.Theme{
		Tokens: models.ThemeTokens{
			Name:        "bold-signal",
			DisplayName: "Bold Signal",
			Vibe:        "Confident, bold, modern, high-impact",
			Description: "Vibrant orange focal card on dark gradient canvas",
			Palette: models.ThemePalette{
				BgCanvas:          "#1a1a1a",
				BgSurface:         "#262626",
				BgSurfaceElevated: "#333333",
				TextPrimary:       "#ffffff",
				TextSecondary:     "#d4d4d4",
				TextMuted:         "#a3a3a3",
				AccentPrimary:     "#ff5722",
				AccentSecondary:   "#ff8a65",
				AccentGlow:        "rgba(255, 87, 34, 0.2)",
				AccentDanger:      "#ef4444",
				BorderColor:       "#ff5722",
				BorderFaint:       "rgba(255, 255, 255, 0.15)",
			},
			Fonts: models.ThemeFonts{
				Display:   "'Archivo Black', sans-serif",
				Body:      "'Space Grotesk', sans-serif",
				Mono:      "'JetBrains Mono', monospace",
				ImportURL: "https://fonts.googleapis.com/css2?family=Archivo+Black&family=Space+Grotesk:wght@400;500;600;700&family=JetBrains+Mono:wght@500;700&display=swap",
			},
			Geometry: models.ThemeGeometry{
				RadiusCard:  "12px",
				RadiusSm:    "6px",
				RadiusLg:    "16px",
				BorderWidth: "1.5px",
				BorderStyle: "solid",
			},
			Shadows: models.ThemeShadows{
				Card:     "0 10px 30px rgba(0, 0, 0, 0.35)",
				Elevated: "0 16px 40px rgba(0, 0, 0, 0.5)",
				Glow:     "0 0 20px rgba(255, 87, 34, 0.25)",
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
				AnimationSpeed:  "0.2s",
			},
		},
		TypographyCSS: `
.slide-topbar .topbar-title { font-family: 'Archivo Black', sans-serif; text-transform: uppercase; }
.panel-headline { font-weight: 700; letter-spacing: -0.01em; }
`,
		SurfacesCSS: `
.glass-panel {
    background: #262626;
    border: 1.5px solid var(--border-faint);
    border-radius: 12px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.35);
}
.glass-panel.lime-edge { border-left: 6px solid var(--accent-primary); }
`,
		BackdropCSS: `
.stage-backdrop {
    background: linear-gradient(135deg, #141414 0%, #222222 50%, #141414 100%);
}
`,
	}
}

func electricStudioPreset() models.Theme {
	return models.Theme{
		Tokens: models.ThemeTokens{
			Name:        "electric-studio",
			DisplayName: "Electric Studio",
			Vibe:        "Bold, clean, professional, high-contrast studio",
			Description: "High-contrast dark canvas with sharp electric royal blue",
			Palette: models.ThemePalette{
				BgCanvas:          "#0a0a0a",
				BgSurface:         "#171717",
				BgSurfaceElevated: "#262626",
				TextPrimary:       "#ffffff",
				TextSecondary:     "#e5e5e5",
				TextMuted:         "#a3a3a3",
				AccentPrimary:     "#4361ee",
				AccentSecondary:   "#3a0ca3",
				AccentGlow:        "rgba(67, 97, 238, 0.25)",
				AccentDanger:      "#f72585",
				BorderColor:       "#4361ee",
				BorderFaint:       "rgba(255, 255, 255, 0.12)",
			},
			Fonts: models.ThemeFonts{
				Display:   "'Manrope', sans-serif",
				Body:      "'Manrope', sans-serif",
				Mono:      "'JetBrains Mono', monospace",
				ImportURL: "https://fonts.googleapis.com/css2?family=Manrope:wght@400;500;600;700;800&family=JetBrains+Mono:wght@500;700&display=swap",
			},
			Geometry: models.ThemeGeometry{
				RadiusCard:  "8px",
				RadiusSm:    "4px",
				RadiusLg:    "14px",
				BorderWidth: "1px",
				BorderStyle: "solid",
			},
			Shadows: models.ThemeShadows{
				Card:     "0 4px 24px rgba(0, 0, 0, 0.3)",
				Elevated: "0 10px 32px rgba(0, 0, 0, 0.45)",
				Glow:     "0 0 20px rgba(67, 97, 238, 0.25)",
			},
			Surfaces: models.ThemeSurfaces{
				BackdropBlur:   "8px",
				SurfaceOpacity: "0.90",
			},
			Spacing: models.ThemeSpacing{
				CardPadding: "28px",
			},
			Motion: models.ThemeMotion{
				SlideTransition: "fade",
				AnimationSpeed:  "0.2s",
			},
		},
		TypographyCSS: `
.slide-topbar .topbar-title { font-weight: 800; letter-spacing: -0.03em; }
`,
		SurfacesCSS: `
.glass-panel {
    border: 1px solid rgba(67, 97, 238, 0.25);
    border-radius: 8px;
    box-shadow: 0 4px 24px rgba(0, 0, 0, 0.3);
}
.glass-panel.lime-edge { border-left: 5px solid var(--accent-primary); }
`,
		BackdropCSS: `
.stage-backdrop {
    background-color: var(--bg-canvas);
    background-image: radial-gradient(circle at 85% 15%, rgba(67, 97, 238, 0.12) 0%, transparent 50%);
}
`,
	}
}

func creativeVoltagePreset() models.Theme {
	return models.Theme{
		Tokens: models.ThemeTokens{
			Name:        "creative-voltage",
			DisplayName: "Creative Voltage",
			Vibe:        "Bold, creative, energetic, retro-modern",
			Description: "Electric cobalt with vivid neon yellow callouts",
			Palette: models.ThemePalette{
				BgCanvas:          "#1a1a2e",
				BgSurface:         "#16213e",
				BgSurfaceElevated: "#0f3460",
				TextPrimary:       "#ffffff",
				TextSecondary:     "#e2e8f0",
				TextMuted:         "#94a3b8",
				AccentPrimary:     "#d4ff00",
				AccentSecondary:   "#0066ff",
				AccentGlow:        "rgba(212, 255, 0, 0.22)",
				AccentDanger:      "#ff0055",
				BorderColor:       "#d4ff00",
				BorderFaint:       "rgba(255, 255, 255, 0.15)",
			},
			Fonts: models.ThemeFonts{
				Display:   "'Syne', sans-serif",
				Body:      "'Space Mono', monospace",
				Mono:      "'Space Mono', monospace",
				ImportURL: "https://fonts.googleapis.com/css2?family=Syne:wght@700;800&family=Space+Mono:wght@400;700&display=swap",
			},
			Geometry: models.ThemeGeometry{
				RadiusCard:  "10px",
				RadiusSm:    "4px",
				RadiusLg:    "16px",
				BorderWidth: "2px",
				BorderStyle: "solid",
			},
			Shadows: models.ThemeShadows{
				Card:     "4px 4px 0 var(--accent-primary)",
				Elevated: "6px 6px 0 var(--accent-primary)",
				Flat:     "3px 3px 0 var(--accent-primary)",
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
				AnimationSpeed:  "0.2s",
			},
		},
		TypographyCSS: `
.slide-topbar .topbar-title { font-family: 'Syne', sans-serif; font-weight: 800; }
`,
		SurfacesCSS: `
.glass-panel {
    border: 2px solid var(--accent-secondary);
    border-radius: 10px;
    box-shadow: 4px 4px 0 var(--accent-primary);
}
.glass-panel.lime-edge { border-left: 6px solid var(--accent-primary); }
`,
		BackdropCSS: `
.stage-backdrop {
    background-color: var(--bg-canvas);
    background-image:
        radial-gradient(circle at 10% 90%, rgba(0, 102, 255, 0.2) 0%, transparent 45%),
        radial-gradient(circle at 90% 10%, rgba(212, 255, 0, 0.1) 0%, transparent 40%);
}
`,
	}
}
