package theme

import (
	"deckforge/internal/models"
)

func cyberDarkPreset() models.Theme {
	return models.Theme{
		Tokens: models.ThemeTokens{
			Name:        "cyber-dark",
			DisplayName: "Cyber Dark",
			Vibe:        "Futuristic, DevSecOps, high-contrast OLED",
			Description: "Deep OLED canvas with neon emerald and purple electric accents",
			Palette: models.ThemePalette{
				BgCanvas:          "#090d16",
				BgSurface:         "rgba(15, 23, 42, 0.85)",
				BgSurfaceElevated: "#0f172a",
				TextPrimary:       "#f8fafc",
				TextSecondary:     "#94a3b8",
				TextMuted:         "#64748b",
				AccentPrimary:     "#10b981",
				AccentSecondary:   "#818cf8",
				AccentGlow:        "rgba(16, 185, 129, 0.25)",
				AccentDanger:      "#ef4444",
				BorderColor:       "#10b981",
				BorderFaint:       "rgba(255, 255, 255, 0.12)",
			},
			Fonts: models.ThemeFonts{
				Display:   "'Poppins', system-ui, -apple-system, sans-serif",
				Body:      "'Plus Jakarta Sans', system-ui, -apple-system, sans-serif",
				Mono:      "'JetBrains Mono', monospace",
				ImportURL: "https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;700&family=Plus+Jakarta+Sans:wght@400;500;600;700&family=Poppins:ital,wght@0,600;0,700;0,800;1,700&display=swap",
			},
			Geometry: models.ThemeGeometry{
				RadiusCard:  "10px",
				RadiusSm:    "4px",
				RadiusLg:    "16px",
				BorderWidth: "1px",
				BorderStyle: "solid",
			},
			Shadows: models.ThemeShadows{
				Card:     "0 10px 30px rgba(0, 0, 0, 0.5)",
				Elevated: "0 20px 40px rgba(0, 0, 0, 0.6)",
				Glow:     "0 0 25px rgba(16, 185, 129, 0.25)",
			},
			Surfaces: models.ThemeSurfaces{
				BackdropBlur:   "16px",
				SurfaceOpacity: "0.85",
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
.slide-topbar .topbar-title { font-weight: 800; letter-spacing: -0.02em; }
.panel-headline { font-weight: 700; color: #f8fafc; }
`,
		SurfacesCSS: `
.glass-panel {
    background: rgba(15, 23, 42, 0.85);
    border: 1px solid rgba(16, 185, 129, 0.25);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
    backdrop-filter: blur(12px);
    -webkit-backdrop-filter: blur(12px);
}
.glass-panel.lime-edge { border-left: 4px solid var(--accent-primary); }
`,
		BackdropCSS: `
.stage-backdrop {
    background-color: var(--bg-canvas);
    background-image:
        radial-gradient(circle at 10% 20%, rgba(16, 185, 129, 0.08) 0%, transparent 40%),
        radial-gradient(circle at 90% 80%, rgba(129, 140, 248, 0.08) 0%, transparent 40%),
        linear-gradient(rgba(255, 255, 255, 0.02) 1px, transparent 1px),
        linear-gradient(90deg, rgba(255, 255, 255, 0.02) 1px, transparent 1px);
    background-size: 100% 100%, 100% 100%, 40px 40px, 40px 40px;
}
`,
		AnimationsCSS: `
@keyframes neonPulse {
    0%, 100% { box-shadow: 0 0 15px rgba(16, 185, 129, 0.2); }
    50% { box-shadow: 0 0 30px rgba(16, 185, 129, 0.45); }
}
.panel-bordered:hover { animation: neonPulse 2s infinite ease-in-out; }
`,
	}
}

func neonCyberPreset() models.Theme {
	return models.Theme{
		Tokens: models.ThemeTokens{
			Name:        "neon-cyber",
			DisplayName: "Neon Cyber",
			Vibe:        "High-tech, confident, neon nightscape",
			Description: "Deep navy background with glowing cyan and electric magenta accents",
			Palette: models.ThemePalette{
				BgCanvas:          "#0a0f1c",
				BgSurface:         "#111827",
				BgSurfaceElevated: "#1f2937",
				TextPrimary:       "#ffffff",
				TextSecondary:     "#9ca3af",
				TextMuted:         "#6b7280",
				AccentPrimary:     "#00ffcc",
				AccentSecondary:   "#ff00aa",
				AccentGlow:        "rgba(0, 255, 204, 0.25)",
				AccentDanger:      "#ff3366",
				BorderColor:       "#00ffcc",
				BorderFaint:       "rgba(255, 255, 255, 0.12)",
			},
			Fonts: models.ThemeFonts{
				Display:   "'Poppins', sans-serif",
				Body:      "'Plus Jakarta Sans', sans-serif",
				Mono:      "'JetBrains Mono', monospace",
				ImportURL: "https://fonts.googleapis.com/css2?family=Poppins:wght@700;800;900&family=Plus+Jakarta+Sans:wght@400;500;600;700&family=JetBrains+Mono:wght@500;700&display=swap",
			},
			Geometry: models.ThemeGeometry{
				RadiusCard:  "12px",
				RadiusSm:    "6px",
				RadiusLg:    "20px",
				BorderWidth: "1.5px",
				BorderStyle: "solid",
			},
			Shadows: models.ThemeShadows{
				Card:     "0 0 25px rgba(0, 255, 204, 0.18)",
				Elevated: "0 0 35px rgba(255, 0, 170, 0.25)",
				Glow:     "0 0 30px rgba(0, 255, 204, 0.35)",
			},
			Surfaces: models.ThemeSurfaces{
				BackdropBlur:   "20px",
				SurfaceOpacity: "0.80",
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
.slide-topbar .topbar-title { font-weight: 800; letter-spacing: -0.02em; }
`,
		SurfacesCSS: `
.glass-panel {
    border: 1px solid rgba(0, 255, 204, 0.3);
    box-shadow: 0 0 20px rgba(0, 255, 204, 0.12);
}
.glass-panel.lime-edge { border-left: 4px solid var(--accent-primary); }
`,
		BackdropCSS: `
.stage-backdrop {
    background-color: var(--bg-canvas);
    background-image:
        radial-gradient(circle at 10% 10%, rgba(0, 255, 204, 0.08) 0%, transparent 40%),
        radial-gradient(circle at 90% 90%, rgba(255, 0, 170, 0.08) 0%, transparent 40%);
}
`,
	}
}

func terminalGreenPreset() models.Theme {
	return models.Theme{
		Tokens: models.ThemeTokens{
			Name:        "terminal-green",
			DisplayName: "Terminal Green",
			Vibe:        "Developer, hacker, phosphor terminal",
			Description: "GitHub dark background with crisp matrix green monospace styling",
			Palette: models.ThemePalette{
				BgCanvas:          "#0d1117",
				BgSurface:         "#161b22",
				BgSurfaceElevated: "#21262d",
				TextPrimary:       "#39d353",
				TextSecondary:     "#8b949e",
				TextMuted:         "#484f58",
				AccentPrimary:     "#39d353",
				AccentSecondary:   "#56d364",
				AccentGlow:        "rgba(57, 211, 83, 0.2)",
				AccentDanger:      "#f85149",
				BorderColor:       "#39d353",
				BorderFaint:       "#30363d",
			},
			Fonts: models.ThemeFonts{
				Display:   "'JetBrains Mono', monospace",
				Body:      "'JetBrains Mono', monospace",
				Mono:      "'JetBrains Mono', monospace",
				ImportURL: "https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600;700;800&display=swap",
			},
			Geometry: models.ThemeGeometry{
				RadiusCard:  "4px",
				RadiusSm:    "2px",
				RadiusLg:    "6px",
				BorderWidth: "1px",
				BorderStyle: "solid",
			},
			Shadows: models.ThemeShadows{
				Card:     "0 4px 20px rgba(0, 0, 0, 0.6)",
				Elevated: "0 8px 30px rgba(0, 0, 0, 0.8)",
				Glow:     "0 0 15px rgba(57, 211, 83, 0.3)",
			},
			Surfaces: models.ThemeSurfaces{
				BackdropBlur:   "0px",
				SurfaceOpacity: "0.95",
			},
			Spacing: models.ThemeSpacing{
				CardPadding: "24px",
			},
			Motion: models.ThemeMotion{
				SlideTransition: "none",
				AnimationSpeed:  "0.1s",
			},
		},
		TypographyCSS: `
.slide-topbar .topbar-title { font-family: 'JetBrains Mono', monospace; font-weight: 800; }
`,
		SurfacesCSS: `
.glass-panel {
    border: 1px solid var(--border-color);
    border-radius: 4px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
}
.glass-panel.lime-edge { border-left: 4px solid var(--accent-primary); }
`,
		BackdropCSS: `
.stage-backdrop {
    background-color: var(--bg-canvas);
    background-image:
        linear-gradient(rgba(57, 211, 83, 0.03) 1px, transparent 1px),
        linear-gradient(90deg, rgba(57, 211, 83, 0.03) 1px, transparent 1px);
    background-size: 32px 32px;
}
`,
	}
}
