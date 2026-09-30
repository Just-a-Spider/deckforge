package models

import (
	"bytes"
	"fmt"
	"text/template"
)

// ComponentInput specifies a configurable parameter on a component
type ComponentInput struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"` // "string", "number", "boolean"
	Default     interface{} `json:"default"`
	Description string      `json:"description,omitempty"`
}

// ComponentDefinition models an Angular-inspired declarative component schema
type ComponentDefinition struct {
	Selector        string           `json:"selector"`
	Name            string           `json:"name"`
	Category        string           `json:"category"`
	Description     string           `json:"description,omitempty"`
	Inputs          []ComponentInput `json:"inputs"`
	Slots           []string         `json:"slots,omitempty"`
	Classes         []string         `json:"classes,omitempty"`
	Styles          string           `json:"styles,omitempty"`
	TemplateSnippet string           `json:"templateSnippet"`
	Scope           string           `json:"scope,omitempty"`      // "builtin", "global", "workspace"
	SourcePath      string           `json:"sourcePath,omitempty"` // Path on disk if custom
}

// Render compiles the component template with supplied inputs and slots
func (c *ComponentDefinition) Render(inputs map[string]interface{}, slots map[string]string) (string, error) {
	tmplData := make(map[string]interface{})

	// Populate defaults
	for _, in := range c.Inputs {
		tmplData[in.Name] = in.Default
	}

	// Override with provided inputs
	for k, v := range inputs {
		tmplData[k] = v
	}

	// Populate slots
	for _, s := range c.Slots {
		if val, exists := slots[s]; exists {
			tmplData["slot_"+s] = val
		} else {
			tmplData["slot_"+s] = ""
		}
	}

	tmpl, err := template.New(c.Selector).Parse(c.TemplateSnippet)
	if err != nil {
		return "", fmt.Errorf("template parse error in %s: %w", c.Selector, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tmplData); err != nil {
		return "", fmt.Errorf("template execute error in %s: %w", c.Selector, err)
	}

	return buf.String(), nil
}

// GetBuiltinComponents returns the curated library of DeckForge components
func GetBuiltinComponents() []ComponentDefinition {
	return []ComponentDefinition{
		{
			Selector:    "df-metric-card",
			Name:        "Metric Impact Card",
			Category:    "metrics",
			Description: "High-contrast statistic card with value, label, and delta tag",
			Inputs: []ComponentInput{
				{Name: "value", Type: "string", Default: "+99.9%", Description: "Primary metric value"},
				{Name: "label", Type: "string", Default: "Operational Uptime", Description: "Metric description"},
				{Name: "delta", Type: "string", Default: "SLO TARGET", Description: "Badge kicker"},
				{Name: "accentEdge", Type: "boolean", Default: true, Description: "Show primary accent left border"},
			},
			Slots:   []string{"notes"},
			Classes: []string{"glass-panel", "lime-edge"},
			TemplateSnippet: `<div class="glass-panel{{if .accentEdge}} lime-edge{{end}}" style="padding: 24px;">
    <div class="card-tag">{{.delta}}</div>
    <div style="font-family: var(--font-display); font-size: 52px; font-weight: 800; line-height: 1; margin-bottom: 8px; color: var(--accent-primary);">{{.value}}</div>
    <div class="panel-headline" style="margin-bottom: 4px;">{{.label}}</div>
    {{if .slot_notes}}<div style="font-size: 13px; color: var(--text-secondary); margin-top: 8px;">{{.slot_notes}}</div>{{end}}
</div>`,
		},
		{
			Selector:    "df-comparison-column",
			Name:        "Comparison Column Panel",
			Category:    "cards",
			Description: "Structured protocol comparison panel with header and bullet items",
			Inputs: []ComponentInput{
				{Name: "title", Type: "string", Default: "Architecture A", Description: "Column header"},
				{Name: "badge", Type: "string", Default: "RECOMMENDED", Description: "Top status badge"},
			},
			Slots:   []string{"items"},
			Classes: []string{"glass-panel"},
			TemplateSnippet: `<div class="glass-panel" style="padding: 24px;">
    <div class="card-tag">{{.badge}}</div>
    <div class="panel-headline">{{.title}}</div>
    <ul class="bullet-list" style="margin-top: 16px;">
        {{if .slot_items}}{{.slot_items}}{{else}}<li>Deterministic 16:9 stage geometry</li><li>Zero-runtime HTML output</li>{{end}}
    </ul>
</div>`,
		},
		{
			Selector:    "df-architecture-node",
			Name:        "Architecture Flow Node",
			Category:    "cards",
			Description: "Microservice / DevSecOps protocol node with kicker and description",
			Inputs: []ComponentInput{
				{Name: "stepNumber", Type: "string", Default: "01", Description: "Step sequence identifier"},
				{Name: "title", Type: "string", Default: "Ingress Gateway", Description: "Node title"},
				{Name: "desc", Type: "string", Default: "Mutual TLS 1.3 verification and token authentication", Description: "Node description"},
			},
			Classes: []string{"glass-panel"},
			TemplateSnippet: `<div class="glass-panel" style="padding: 20px; border-left: 3px solid var(--accent-primary);">
    <div style="font-family: var(--font-mono); font-size: 11px; font-weight: 700; color: var(--accent-primary); margin-bottom: 6px;">NODE // {{.stepNumber}}</div>
    <div style="font-family: var(--font-display); font-size: 18px; font-weight: 700; color: var(--text-primary); margin-bottom: 8px;">{{.title}}</div>
    <div style="font-size: 14px; line-height: 1.45; color: var(--text-secondary);">{{.desc}}</div>
</div>`,
		},
		{
			Selector:    "df-code-spec-box",
			Name:        "Code Specification Box",
			Category:    "code",
			Description: "Syntax snippet container with file path header and language tag",
			Inputs: []ComponentInput{
				{Name: "filename", Type: "string", Default: "auth_service.go", Description: "Code file path"},
				{Name: "lang", Type: "string", Default: "GO 1.23", Description: "Language pill"},
				{Name: "code", Type: "string", Default: "func VerifyToken(ctx context.Context, token string) error {\n    return verifier.Check(token)\n}", Description: "Source code"},
			},
			Classes: []string{"glass-panel"},
			TemplateSnippet: `<div class="glass-panel" style="padding: 0; overflow: hidden; background: #0b0f19;">
    <div style="display: flex; justify-content: space-between; align-items: center; padding: 10px 16px; background: rgba(255,255,255,0.05); border-bottom: 1px solid var(--border-color);">
        <span style="font-family: var(--font-mono); font-size: 12px; color: var(--text-primary);">{{.filename}}</span>
        <span class="card-tag" style="margin: 0;">{{.lang}}</span>
    </div>
    <pre style="margin: 0; padding: 18px; font-family: var(--font-mono); font-size: 13px; line-height: 1.5; color: #a5f3fc; overflow-x: auto;"><code>{{.code}}</code></pre>
</div>`,
		},
		{
			Selector:    "df-executive-quote",
			Name:        "Executive Callout Quote",
			Category:    "typography",
			Description: "Bordered editorial quote with prominent text and attribution",
			Inputs: []ComponentInput{
				{Name: "quote", Type: "string", Default: "Simplicity is prerequisite for reliability.", Description: "Quote body"},
				{Name: "author", Type: "string", Default: "Edsger W. Dijkstra", Description: "Author / Source"},
				{Name: "role", Type: "string", Default: "Turing Award Laureate", Description: "Author title"},
			},
			Classes: []string{"glass-panel", "lime-edge"},
			TemplateSnippet: `<div class="glass-panel lime-edge" style="padding: 28px 32px;">
    <div style="font-family: var(--font-display); font-size: 26px; font-style: italic; font-weight: 600; line-height: 1.35; color: var(--text-primary); margin-bottom: 16px;">
        "{{.quote}}"
    </div>
    <div style="font-family: var(--font-mono); font-size: 12px;">
        <span style="color: var(--accent-primary); font-weight: 700;">{{.author}}</span>
        <span style="color: var(--text-secondary); margin-left: 8px;">// {{.role}}</span>
    </div>
</div>`,
		},
		{
			Selector:    "df-timeline-step",
			Name:        "Horizontal Timeline Step",
			Category:    "layout",
			Description: "Sequential phase box with step kicker and milestone objectives",
			Inputs: []ComponentInput{
				{Name: "phase", Type: "string", Default: "PHASE 01", Description: "Phase tag"},
				{Name: "title", Type: "string", Default: "Discovery & Blueprint", Description: "Milestone title"},
				{Name: "eta", Type: "string", Default: "Q3 2026", Description: "Timeline indicator"},
			},
			Slots:   []string{"details"},
			Classes: []string{"glass-panel"},
			TemplateSnippet: `<div class="glass-panel" style="padding: 20px;">
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">
        <span class="card-tag" style="margin: 0;">{{.phase}}</span>
        <span style="font-family: var(--font-mono); font-size: 11px; color: var(--text-muted);">{{.eta}}</span>
    </div>
    <div class="panel-headline" style="font-size: 18px; margin-bottom: 8px;">{{.title}}</div>
    <div style="font-size: 13px; color: var(--text-secondary); line-height: 1.45;">
        {{if .slot_details}}{{.slot_details}}{{else}}Establish architectural requirements and contracts.{{end}}
    </div>
</div>`,
		},
		{
			Selector:    "df-feature-pill-grid",
			Name:        "Feature Badge Pills Matrix",
			Category:    "layout",
			Description: "Compact cluster of capability badges with icons/indicators",
			Inputs: []ComponentInput{
				{Name: "tag1", Type: "string", Default: "DETERMINISTIC 1080P"},
				{Name: "tag2", Type: "string", Default: "ZERO DEPENDENCIES"},
				{Name: "tag3", Type: "string", Default: "WCAG AAA AUDIT"},
				{Name: "tag4", Type: "string", Default: "MCP AGENT READY"},
			},
			TemplateSnippet: `<div style="display: flex; flex-wrap: wrap; gap: 10px;">
    <span class="lime-badge">{{.tag1}}</span>
    <span class="purple-badge">{{.tag2}}</span>
    <span class="lime-badge">{{.tag3}}</span>
    <span class="purple-badge">{{.tag4}}</span>
</div>`,
		},
	}
}
