package scaffold

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"deckforge/internal/models"
	"deckforge/internal/workspace"
)

var slugRegex = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(text string) string {
	clean := strings.ToLower(strings.TrimSpace(text))
	clean = slugRegex.ReplaceAllString(clean, "_")
	clean = strings.Trim(clean, "_")
	if len(clean) > 40 {
		clean = clean[:40]
	}
	if clean == "" {
		clean = "slide"
	}
	return clean
}

// ScaffoldDeck creates a new presentation with the specified slide count and theme
func ScaffoldDeck(targetDir, deckName, themeName string, slideCount int) (*models.Deck, error) {
	if slideCount <= 0 {
		slideCount = 5
	}
	if themeName == "" {
		themeName = "academic-crimson"
	}

	cleanDir, err := filepath.Abs(targetDir)
	if err != nil {
		cleanDir = targetDir
	}

	slidesDir := filepath.Join(cleanDir, "slides")
	if err := os.MkdirAll(slidesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create slides directory: %w", err)
	}

	// Create deck.json
	conf := models.DeckConfig{
		Title:    strings.Title(strings.ReplaceAll(deckName, "_", " ")),
		Subtitle: "High-Fidelity 1080p Modular Slide Deck",
		Theme:    themeName,
	}
	confData, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(cleanDir, "deck.json"), confData, 0644); err != nil {
		return nil, err
	}

	// Generate N sequential slides
	for i := 1; i <= slideCount; i++ {
		var filename string
		var content string

		if i == 1 {
			filename = "01_title.html"
			content = titleSlideHTML(1, conf.Title, conf.Subtitle)
		} else if i == 2 && slideCount > 2 {
			filename = "02_agenda.html"
			content = agendaSlideHTML(2, "Agenda and Executive Roadmap")
		} else if i == slideCount && slideCount > 1 {
			filename = fmt.Sprintf("%02d_conclusions.html", i)
			content = conclusionSlideHTML(i, "Strategic Conclusions and Next Steps")
		} else {
			// Alternate among layouts
			switch i % 4 {
			case 3:
				filename = fmt.Sprintf("%02d_architecture.html", i)
				content = twoColSlideHTML(i, fmt.Sprintf("Module %02d: Core Architecture", i))
			case 0:
				filename = fmt.Sprintf("%02d_feature_matrix.html", i)
				content = threeColSlideHTML(i, fmt.Sprintf("Module %02d: Component Matrix", i))
			case 1:
				filename = fmt.Sprintf("%02d_metrics_impact.html", i)
				content = splitHeroSlideHTML(i, fmt.Sprintf("Module %02d: Key Metrics & Performance", i))
			case 2:
				filename = fmt.Sprintf("%02d_technical_specs.html", i)
				content = codeSpecSlideHTML(i, fmt.Sprintf("Module %02d: Implementation Protocols", i))
			}
		}

		filePath := filepath.Join(slidesDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return nil, err
		}
	}

	return workspace.InspectDeck(cleanDir)
}

// AddSlide appends a new slide with layout to an existing deck
func AddSlide(deckPath, title, layout string) (*models.SlideInfo, error) {
	slidesDir := filepath.Join(deckPath, "slides")
	if err := os.MkdirAll(slidesDir, 0755); err != nil {
		return nil, err
	}

	deck, err := workspace.InspectDeck(deckPath)
	if err != nil {
		return nil, err
	}

	nextNum := len(deck.Slides) + 1
	slug := slugify(title)
	filename := fmt.Sprintf("%02d_%s.html", nextNum, slug)
	filePath := filepath.Join(slidesDir, filename)

	var content string
	switch layout {
	case "3col":
		content = threeColSlideHTML(nextNum, title)
	case "hero":
		content = splitHeroSlideHTML(nextNum, title)
	case "code":
		content = codeSpecSlideHTML(nextNum, title)
	default:
		content = twoColSlideHTML(nextNum, title)
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return nil, err
	}

	return &models.SlideInfo{
		Index:    nextNum,
		Filename: filename,
		Path:     filePath,
		Title:    title,
		Slug:     slug,
	}, nil
}

func titleSlideHTML(num int, title, subtitle string) string {
	return fmt.Sprintf(`<section class="slide active" data-slide="%d">
    <div class="slide-frame" style="justify-content: center; align-items: center; text-align: center;">
        <div style="max-width: 1200px;">
            <div class="badge" style="margin-bottom: 24px; font-size: 13px; padding: 8px 18px;">DECKFORGE • 1080P MASTER CANVAS</div>
            <h1 style="font-family: var(--font-display); font-size: 64px; font-weight: 800; line-height: 1.1; margin-bottom: 24px; color: var(--text-primary); letter-spacing: -0.03em;">
                %s
            </h1>
            <p style="font-size: 24px; line-height: 1.5; color: var(--text-secondary); max-width: 800px; margin: 0 auto 36px auto;">
                %s
            </p>
            <div style="font-family: var(--font-mono); font-size: 14px; color: var(--accent-primary); letter-spacing: 0.1em; text-transform: uppercase;">
                Modular Presentation Architecture
            </div>
        </div>
    </div>
</section>
`, num, title, subtitle)
}

func agendaSlideHTML(num int, title string) string {
	return fmt.Sprintf(`<section class="slide" data-slide="%d">
    <div class="slide-frame">
        <div class="slide-topbar">
            <div class="topbar-title">%s</div>
            <div class="topbar-tag">ROADMAP • SLIDE %02d</div>
        </div>
        <div class="slide-content-area">
            <div class="grid-3col">
                <div class="glass-panel lime-edge" style="padding: 28px;">
                    <div class="card-tag">PHASE 01</div>
                    <div class="panel-headline">Foundation & Context</div>
                    <p style="font-size: 16px; line-height: 1.6; color: var(--text-secondary); margin-bottom: 16px;">
                        Baseline operational objectives, architectural invariants, and domain analysis.
                    </p>
                    <div class="badge">DISCOVERY</div>
                </div>
                <div class="glass-panel lime-edge" style="padding: 28px;">
                    <div class="card-tag">PHASE 02</div>
                    <div class="panel-headline">Technical Execution</div>
                    <p style="font-size: 16px; line-height: 1.6; color: var(--text-secondary); margin-bottom: 16px;">
                        Core protocol design, service topology, and data verification workflows.
                    </p>
                    <div class="badge">ARCHITECTURE</div>
                </div>
                <div class="glass-panel lime-edge" style="padding: 28px;">
                    <div class="card-tag">PHASE 03</div>
                    <div class="panel-headline">Impact & Deliverables</div>
                    <p style="font-size: 16px; line-height: 1.6; color: var(--text-secondary); margin-bottom: 16px;">
                        Performance metrics, strategic conclusions, and rollout roadmap.
                    </p>
                    <div class="badge">RESULTS</div>
                </div>
            </div>
        </div>
    </div>
</section>
`, num, title, num)
}

func twoColSlideHTML(num int, title string) string {
	return fmt.Sprintf(`<section class="slide" data-slide="%d">
    <div class="slide-frame">
        <div class="slide-topbar">
            <div class="topbar-title">%s</div>
            <div class="topbar-tag">MODULE • SLIDE %02d</div>
        </div>
        <div class="slide-content-area">
            <div class="grid-2col">
                <div class="glass-panel lime-edge" style="padding: 28px;">
                    <div class="card-tag">COMPONENT SPECIFICATION</div>
                    <div class="panel-headline">Functional Architecture</div>
                    <p style="font-size: 16px; line-height: 1.6; color: var(--text-secondary); margin-bottom: 20px;">
                        Deep-dive breakdown of the module mechanics, interface contracts, and execution guarantees.
                    </p>
                    <ul class="bullet-list">
                        <li>Zero external runtime dependencies for standalone distribution.</li>
                        <li>Deterministic 1920x1080 canvas scaling without viewport distortion.</li>
                        <li>Unified keyboard controller with instant 4K UHD capture.</li>
                    </ul>
                </div>
                <div class="glass-panel lime-edge" style="padding: 28px;">
                    <div class="card-tag">VERIFICATION & METRICS</div>
                    <div class="panel-headline">System Invariants</div>
                    <p style="font-size: 16px; line-height: 1.6; color: var(--text-secondary); margin-bottom: 20px;">
                        Operational boundary checks, fault tolerance metrics, and deployment criteria.
                    </p>
                    <ul class="bullet-list">
                        <li>Fixed aspect ratio maintained across all display factors.</li>
                        <li>Segmented CSS architecture for targeted theme adjustments.</li>
                        <li>Embedded HTTP server with real-time file watching.</li>
                    </ul>
                </div>
            </div>
        </div>
    </div>
</section>
`, num, title, num)
}

func threeColSlideHTML(num int, title string) string {
	return fmt.Sprintf(`<section class="slide" data-slide="%d">
    <div class="slide-frame">
        <div class="slide-topbar">
            <div class="topbar-title">%s</div>
            <div class="topbar-tag">ANALYSIS • SLIDE %02d</div>
        </div>
        <div class="slide-content-area">
            <div class="grid-3col">
                <div class="glass-panel lime-edge" style="padding: 24px;">
                    <div class="card-tag">PILLAR A</div>
                    <div class="panel-headline">Core Engine</div>
                    <p style="font-size: 15px; line-height: 1.55; color: var(--text-secondary);">
                        Pure Go standalone binary orchestrates build, serve, and presentation lifecycles.
                    </p>
                </div>
                <div class="glass-panel lime-edge" style="padding: 24px;">
                    <div class="card-tag">PILLAR B</div>
                    <div class="panel-headline">Segmented Themes</div>
                    <p style="font-size: 15px; line-height: 1.55; color: var(--text-secondary);">
                        Dismantles monolithic stylesheets into Palette, Typography, Surfaces, and Backdrop.
                    </p>
                </div>
                <div class="glass-panel lime-edge" style="padding: 24px;">
                    <div class="card-tag">PILLAR C</div>
                    <div class="panel-headline">External Editor</div>
                    <p style="font-size: 15px; line-height: 1.55; color: var(--text-secondary);">
                        Seamless integration with $EDITOR for surgical editing of single or batch slides.
                    </p>
                </div>
            </div>
        </div>
    </div>
</section>
`, num, title, num)
}

func splitHeroSlideHTML(num int, title string) string {
	return fmt.Sprintf(`<section class="slide" data-slide="%d">
    <div class="slide-frame">
        <div class="slide-topbar">
            <div class="topbar-title">%s</div>
            <div class="topbar-tag">METRICS • SLIDE %02d</div>
        </div>
        <div class="slide-content-area">
            <div class="grid-split-hero">
                <div class="glass-panel lime-edge" style="padding: 32px;">
                    <div class="card-tag">PERFORMANCE PROFILE</div>
                    <div class="panel-headline" style="font-size: 28px;">High-Speed Compilation</div>
                    <p style="font-size: 17px; line-height: 1.6; color: var(--text-secondary); margin-bottom: 24px;">
                        Native Go assembler packages multi-slide HTML decks in under 15ms. Zero NodeJS, zero npm bundles, zero external network requests needed at runtime.
                    </p>
                    <div class="table-container">
                        <table>
                            <thead>
                                <tr><th>Metric</th><th>Target</th><th>Benchmark</th></tr>
                            </thead>
                            <tbody>
                                <tr><td>Build Latency</td><td>&lt; 50ms</td><td>12ms</td></tr>
                                <tr><td>Memory Footprint</td><td>&lt; 20MB</td><td>8.4MB</td></tr>
                                <tr><td>Stage Resolution</td><td>1920x1080</td><td>Deterministic</td></tr>
                            </tbody>
                        </table>
                    </div>
                </div>
                <div class="glass-panel" style="padding: 32px; justify-content: center; align-items: center; text-align: center;">
                    <div style="font-size: 72px; font-weight: 800; color: var(--accent-primary); font-family: var(--font-display);">
                        1080p
                    </div>
                    <div style="font-family: var(--font-mono); font-size: 14px; letter-spacing: 0.1em; color: var(--text-secondary); margin-bottom: 16px;">
                        PRISTINE STAGE GEOMETRY
                    </div>
                    <p style="font-size: 15px; line-height: 1.5; color: var(--text-muted);">
                        Scaled uniformly to any screen with Math.min(W/1920, H/1080).
                    </p>
                </div>
            </div>
        </div>
    </div>
</section>
`, num, title, num)
}

func codeSpecSlideHTML(num int, title string) string {
	return fmt.Sprintf(`<section class="slide" data-slide="%d">
    <div class="slide-frame">
        <div class="slide-topbar">
            <div class="topbar-title">%s</div>
            <div class="topbar-tag">PROTOCOL • SLIDE %02d</div>
        </div>
        <div class="slide-content-area">
            <div class="grid-2col">
                <div class="glass-panel lime-edge" style="padding: 28px;">
                    <div class="card-tag">SYSTEM INTERACTION</div>
                    <div class="panel-headline">Process Suspension Protocol</div>
                    <p style="font-size: 16px; line-height: 1.6; color: var(--text-secondary); margin-bottom: 16px;">
                        When editing externally, Bubble Tea suspends terminal raw mode, spawns $EDITOR, and smoothly restores UI state upon exit.
                    </p>
                    <div class="code-block">
func (m Model) editSlide(path string) tea.Cmd {
    c := exec.Command(editor, path)
    return tea.ExecProcess(c, func(err error) tea.Msg {
        return slideEditedMsg{path, err}
    })
}
                    </div>
                </div>
                <div class="glass-panel lime-edge" style="padding: 28px;">
                    <div class="card-tag">BENEFITS</div>
                    <div class="panel-headline">Operational Highlights</div>
                    <ul class="bullet-list">
                        <li>Instant hot-reload during presentation rehearsals.</li>
                        <li>Full keyboard navigation: Arrow keys, Space, M (Index), E (Edit).</li>
                        <li>Super-sampled 4K UHD capture with P key.</li>
                    </ul>
                </div>
            </div>
        </div>
    </div>
</section>
`, num, title, num)
}

func conclusionSlideHTML(num int, title string) string {
	return fmt.Sprintf(`<section class="slide" data-slide="%d">
    <div class="slide-frame" style="justify-content: center; align-items: center; text-align: center;">
        <div style="max-width: 960px;">
            <div class="badge" style="margin-bottom: 20px;">SUMMARY & NEXT STEPS</div>
            <h2 style="font-family: var(--font-display); font-size: 52px; font-weight: 800; line-height: 1.15; margin-bottom: 24px; color: var(--text-primary); letter-spacing: -0.02em;">
                %s
            </h2>
            <p style="font-size: 20px; line-height: 1.6; color: var(--text-secondary); margin-bottom: 36px;">
                DeckForge provides deterministic 1080p modular presentation authoring with zero runtime overhead.
            </p>
            <div style="display: flex; justify-content: center; gap: 16px;">
                <div class="badge" style="padding: 10px 20px; font-size: 13px;">Q & A SESSION</div>
                <div class="badge" style="padding: 10px 20px; font-size: 13px;">DOCUMENTATION</div>
            </div>
        </div>
    </div>
</section>
`, num, title)
}
