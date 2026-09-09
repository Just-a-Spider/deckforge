# Component Schema Specification

## Schema Definition (JSON)

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "DeckForgeComponent",
  "type": "object",
  "required": ["selector", "name", "category", "inputs", "templateSnippet"],
  "properties": {
    "selector": {
      "type": "string",
      "pattern": "^df-[a-z0-9-]+$",
      "description": "Unique component custom tag name (e.g. df-metric-card)"
    },
    "name": {
      "type": "string",
      "description": "Human-readable component name"
    },
    "category": {
      "type": "string",
      "enum": ["metrics", "layout", "typography", "code", "cards", "media"]
    },
    "inputs": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["name", "type", "default"],
        "properties": {
          "name": { "type": "string" },
          "type": { "type": "string", "enum": ["string", "number", "boolean"] },
          "default": { "type": ["string", "number", "boolean"] },
          "description": { "type": "string" }
        }
      }
    },
    "slots": {
      "type": "array",
      "items": { "type": "string" },
      "description": "Named projection slots (e.g. body, footer, notes)"
    },
    "classes": {
      "type": "array",
      "items": { "type": "string" },
      "description": "Default CSS utility classes attached to root"
    },
    "templateSnippet": {
      "type": "string",
      "description": "Go template snippet producing clean HTML"
    }
  }
}
```

## Built-in Components Catalog

1. `df-metric-card`: Stat number, label, delta badge, optional sparkline.
2. `df-comparison-column`: Pros/cons, protocol comparison, before/after.
3. `df-architecture-node`: DevSecOps/cloud microservice node with badge and tag.
4. `df-code-spec-box`: Code syntax snippet with language badge and path header.
5. `df-executive-quote`: Bordered editorial quote with author attribution.
6. `df-timeline-step`: Numbered sequence step with connector line.
7. `df-feature-pill-grid`: Compact matrix of badge pills.

## Spacing Primitives & Scale

Standardized atomic spacing variables defined in `:root`:
- `--space-1`: `8px` (`.p-1`, `.gap-1`, `.mt-1`, `.mb-1`)
- `--space-2`: `16px` (`.p-2`, `.gap-2`, `.mt-2`, `.mb-2`)
- `--space-3`: `24px` (`.p-3`, `.gap-3`, `.mt-3`, `.mb-3`)
- `--space-4`: `32px` (`.p-4`, `.gap-4`, `.mt-4`, `.mb-4`)
- `--space-5`: `48px` (`.p-5`, `.gap-5`, `.mt-5`, `.mb-5`)
- `--space-6`: `64px` (`.p-6`, `.gap-6`, `.mt-6`, `.mb-6`)

Directional utilities:
- Horizontal padding: `.px-1` to `.px-6`
- Vertical padding: `.py-1` to `.py-6`
- Centering: `.m-auto`

## Chapter & Act Divider Layout

High-contrast centered section divider slides:
```html
<section class="slide chapter-break" data-slide="N">
    <div class="slide-frame">
        <div class="chapter-container">
            <div class="chapter-kicker">ACT II // SYSTEM ARCHITECTURE</div>
            <div class="chapter-num">02</div>
            <h1 class="chapter-title">The Triad Execution Engine</h1>
            <div class="chapter-divider"></div>
            <p class="chapter-summary">
                Deep dive into synchronous runtime compilation, headless CLI mechanics, and deterministic canvas geometry.
            </p>
        </div>
    </div>
</section>
```

