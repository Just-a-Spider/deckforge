package models

// CanvasRules defines the deterministic 1080p canvas constraints
type CanvasRules struct {
	StageWidth   int      `json:"stage_width"`
	StageHeight  int      `json:"stage_height"`
	AspectRatio  string   `json:"aspect_ratio"`
	ScaleFormula string   `json:"scale_formula"`
	Forbidden    []string `json:"forbidden"`
	Required     []string `json:"required"`
}

// ClassItem defines a single semantic CSS class with its purpose and example
type ClassItem struct {
	Class       string `json:"class"`
	Description string `json:"description"`
	Example     string `json:"example,omitempty"`
}

// ClassCategory groups semantic classes by role
type ClassCategory struct {
	Category string      `json:"category"`
	Classes  []ClassItem `json:"classes"`
}

// TokenItem defines a CSS custom property variable
type TokenItem struct {
	Variable    string `json:"variable"`
	Role        string `json:"role"`
	DefaultDark string `json:"default_dark,omitempty"`
}

// SpacingItem defines a standardized spacing scale step
type SpacingItem struct {
	Level    int    `json:"level"`
	Variable string `json:"variable"`
	Value    string `json:"value"`
	Utility  string `json:"utility"`
}

// StylingGuide provides external agents and MCP tools with complete style context
type StylingGuide struct {
	CanvasRules    CanvasRules     `json:"canvas_rules"`
	Categories     []ClassCategory `json:"class_categories"`
	SpacingScale   []SpacingItem   `json:"spacing_scale"`
	ColorTokens    []TokenItem     `json:"color_tokens"`
	QuickReference string          `json:"quick_reference"`
}

// GetStylingGuide returns the comprehensive, self-describing CSS and token guide
func GetStylingGuide() StylingGuide {
	return StylingGuide{
		CanvasRules: CanvasRules{
			StageWidth:   1920,
			StageHeight:  1080,
			AspectRatio:  "16:9",
			ScaleFormula: "Math.min(window.innerWidth / 1920, window.innerHeight / 1080)",
			Forbidden: []string{
				"NEVER use dynamic viewport units (100vw, 100vh, vmin, vmax) in slide HTML",
				"NEVER allow cards or columns to wrap unpredictably (avoid unstructured float/inline)",
				"NEVER hardcode raw hex colors (#ffffff, #000000) for text or surfaces",
				"NEVER exceed 860px total content height inside .slide-content-area (causes footer collisions)",
			},
			Required: []string{
				"ALWAYS use fixed pixel values (34px, 70px) or percentage grid tracks relative to 1920x1080 stage",
				"ALWAYS wrap slide content inside <section class=\"slide\"><div class=\"slide-frame\">...</div></section>",
				"ALWAYS reference semantic tokens via var(--token-name)",
				"ALWAYS budget vertical height: available content area is 860px. When combining grid + footer, size grid (e.g. 620px) and footer (e.g. 120px) proportionally or use .row-fill and .row-auto",
			},
		},
		Categories: []ClassCategory{
			{
				Category: "Containers",
				Classes: []ClassItem{
					{Class: ".slide", Description: "Top-level slide container pinned to 1920x1080 viewport"},
					{Class: ".slide-frame", Description: "Inner padded boundary with flex layout and safe margins"},
					{Class: ".slide-content-area", Description: "Content body area below standard slide topbar (max-height: 860px)"},
					{Class: ".slide.chapter-break", Description: "Centered act/chapter divider with glowing focal point"},
					{Class: ".row-fill", Description: "Flex-fill row inside .slide-content-area taking remaining height"},
					{Class: ".row-auto", Description: "Auto-height row inside .slide-content-area (e.g. summary bars, metric footers)"},
				},
			},
			{
				Category: "Layouts",
				Classes: []ClassItem{
					{Class: ".grid-2col", Description: "Two equal-width columns (1fr 1fr) with fixed gap"},
					{Class: ".grid-3col", Description: "Three equal-width columns (1fr 1fr 1fr) for card decks"},
					{Class: ".grid-4col", Description: "Four equal-width columns for metric cards or architecture nodes"},
					{Class: ".grid-split-hero", Description: "Split hero layout: 1.2fr left narrative, 0.8fr right media/card"},
					{Class: ".grid-metrics", Description: "Compact auto-fit grid specifically for KPI callouts"},
				},
			},
			{
				Category: "Panels & Cards",
				Classes: []ClassItem{
					{Class: ".panel-bordered", Description: "Clean balanced card with 1.5px continuous border in theme accent/border color"},
					{Class: ".panel-top-accent", Description: "Modern card with bold 4px accent line across the top border"},
					{Class: ".panel-elevated", Description: "Minimalist elevated card on elevated surface with soft shadow, no harsh stripes"},
					{Class: ".panel-outline", Description: "Subtle transparent wireframe card with faint border"},
					{Class: ".glass-panel", Description: "Translucent backdrop-filtered card with subtle border and surface elevation"},
					{Class: ".lime-edge / .accent-edge", Description: "Left accent border stripe for key callout cards (use selectively for visual rhythm)"},
				},
			},
			{
				Category: "Typography & Headers",
				Classes: []ClassItem{
					{Class: ".slide-topbar", Description: "Standardized flex header row displaying title and tags"},
					{Class: ".topbar-title", Description: "Slide category or section kicker in topbar"},
					{Class: ".topbar-tag", Description: "Monospace tracking tag or phase badge in topbar"},
					{Class: ".micro-kicker", Description: "Tiny uppercase section kicker (font-mono, tracking-wide)"},
					{Class: ".panel-headline", Description: "Display heading inside cards and panels (font-display)"},
					{Class: ".card-tag", Description: "Monospace tag badge inside cards and nodes"},
					{Class: ".bullet-list", Description: "High-contrast clean bulleted list with custom bullets"},
					{Class: ".chapter-kicker", Description: "Oversized section identifier on chapter divider slides"},
					{Class: ".chapter-num", Description: "Gigantic display numeral or act identifier (font-display)"},
					{Class: ".chapter-title", Description: "Hero chapter headline with high-contrast text"},
					{Class: ".chapter-summary", Description: "Executive summary paragraph for chapter dividers"},
					{Class: ".chapter-divider", Description: "Accent horizontal divider line"},
				},
			},
			{
				Category: "Highlights & Badges",
				Classes: []ClassItem{
					{Class: ".lime-hl", Description: "Primary accent color text highlight"},
					{Class: ".accent-hl", Description: "Secondary accent color text highlight"},
					{Class: ".lime-badge", Description: "Pill badge with primary accent background/border"},
					{Class: ".purple-badge", Description: "Pill badge with secondary accent background/border"},
					{Class: ".badge-danger", Description: "Critical / alert status pill badge"},
				},
			},
			{
				Category: "Spacing Utilities",
				Classes: []ClassItem{
					{Class: ".p-1 to .p-6", Description: "Standard padding (p-1: 8px, p-2: 16px, p-3: 24px, p-4: 32px, p-5: 48px, p-6: 64px)"},
					{Class: ".px-1 to .px-6", Description: "Horizontal padding utility"},
					{Class: ".py-1 to .py-6", Description: "Vertical padding utility"},
					{Class: ".gap-1 to .gap-6", Description: "Grid and flex gap utility (gap-1: 8px, gap-2: 16px, gap-3: 24px, gap-4: 32px, gap-5: 48px, gap-6: 64px)"},
					{Class: ".m-auto", Description: "Center alignment margin auto"},
					{Class: ".mt-1 to .mt-6", Description: "Top margin spacing step"},
					{Class: ".mb-1 to .mb-6", Description: "Bottom margin spacing step"},
				},
			},
		},
		SpacingScale: []SpacingItem{
			{Level: 1, Variable: "--space-1", Value: "8px", Utility: ".p-1 / .gap-1 / .mt-1"},
			{Level: 2, Variable: "--space-2", Value: "16px", Utility: ".p-2 / .gap-2 / .mt-2"},
			{Level: 3, Variable: "--space-3", Value: "24px", Utility: ".p-3 / .gap-3 / .mt-3"},
			{Level: 4, Variable: "--space-4", Value: "32px", Utility: ".p-4 / .gap-4 / .mt-4"},
			{Level: 5, Variable: "--space-5", Value: "48px", Utility: ".p-5 / .gap-5 / .mt-5"},
			{Level: 6, Variable: "--space-6", Value: "64px", Utility: ".p-6 / .gap-6 / .mt-6"},
		},
		ColorTokens: []TokenItem{
			{Variable: "var(--bg-canvas)", Role: "Slide stage background color", DefaultDark: "#090d16"},
			{Variable: "var(--bg-surface)", Role: "Standard card surface background", DefaultDark: "rgba(15, 23, 42, 0.85)"},
			{Variable: "var(--bg-surface-elevated)", Role: "Elevated modal/inspector surface", DefaultDark: "rgba(30, 41, 59, 0.95)"},
			{Variable: "var(--text-primary)", Role: "Highest contrast primary copy", DefaultDark: "#f8fafc"},
			{Variable: "var(--text-secondary)", Role: "Subordinate body copy (min 4.5:1 AA)", DefaultDark: "#94a3b8"},
			{Variable: "var(--text-muted)", Role: "Faint labels, kickers, timestamps", DefaultDark: "#64748b"},
			{Variable: "var(--accent-primary)", Role: "Main visual focal accent (crimson, emerald, blue, amber)", DefaultDark: "#10b981"},
			{Variable: "var(--accent-secondary)", Role: "Complementary accent (indigo, violet, sky)", DefaultDark: "#818cf8"},
			{Variable: "var(--accent-danger)", Role: "Penalties, errors, security warnings", DefaultDark: "#ef4444"},
			{Variable: "var(--border-color)", Role: "Primary border color matching theme palette", DefaultDark: "#10b981"},
			{Variable: "var(--border-faint)", Role: "Subtle borders and dividers", DefaultDark: "rgba(255, 255, 255, 0.12)"},
			{Variable: "var(--font-display)", Role: "Primary headline font family", DefaultDark: "Poppins, Outfit, sans-serif"},
			{Variable: "var(--font-body)", Role: "Body paragraph font family", DefaultDark: "Plus Jakarta Sans, Inter, sans-serif"},
			{Variable: "var(--font-mono)", Role: "Code, kickers, metrics, metadata", DefaultDark: "JetBrains Mono, monospace"},
			{Variable: "var(--radius-card)", Role: "Corner radius for panels and cards", DefaultDark: "10px"},
			{Variable: "var(--radius-sm)", Role: "Corner radius for badges, tags, pills", DefaultDark: "4px"},
			{Variable: "var(--border-width)", Role: "Border stroke thickness", DefaultDark: "1px"},
			{Variable: "var(--border-style)", Role: "Border stroke style (solid, dashed)", DefaultDark: "solid"},
			{Variable: "var(--shadow-card)", Role: "Default card elevation and drop shadow", DefaultDark: "0 10px 30px rgba(0, 0, 0, 0.5)"},
			{Variable: "var(--shadow-elevated)", Role: "Modal and inspector elevated shadow", DefaultDark: "0 20px 40px rgba(0, 0, 0, 0.6)"},
			{Variable: "var(--shadow-glow)", Role: "Optical glow around focal elements", DefaultDark: "0 0 25px rgba(16, 185, 129, 0.25)"},
			{Variable: "var(--backdrop-blur)", Role: "Translucent backdrop blur filter radius", DefaultDark: "16px"},
			{Variable: "var(--card-padding)", Role: "Standard inner card padding", DefaultDark: "28px"},
			{Variable: "var(--animation-speed)", Role: "Slide transition and hover speed", DefaultDark: "0.2s"},
		},
		QuickReference: "To scaffold a new slide: use .slide > .slide-frame > (.slide-topbar + layout container). Vary card styles (.panel-bordered, .panel-top-accent, .panel-elevated, .glass-panel). Cards automatically inherit active theme tokens (radius, shadow, blur, borders, padding). Choose a theme tailored to the subject (e.g. swiss-minimal, paper-ink, academic-crimson, terminal-green, electric-studio). Respect the 860px vertical budget: size multi-row content proportionally to prevent footer collisions.",
	}
}
