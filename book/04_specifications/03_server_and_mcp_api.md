# Server REST and MCP API Specification

## REST API Endpoints

### `GET /api/deck`
Returns metadata and discovered slide files for the currently active deck.

### `POST /api/slides/update`
Receives `{ "index": int, "html": string }`. Sanitizes HTML (strips runtime studio attributes) and updates `slides/NN_*.html`.

### `POST /api/slides/reorder`
Receives `{ "newOrder": []int }` (e.g. `[1, 3, 2, 4]`). Safely re-indexes slide files sequentially on disk.

### `POST /api/slides/create`
Receives `{ "title": string, "layout": string }`. Generates a new slide file and appends to deck.

### `DELETE /api/slides/:index`
Removes specified slide from disk and re-indexes subsequent slides.

### `GET /api/themes`
Returns array of all available presets and custom themes with name, displayName, description, and vibe.

### `POST /api/deck/theme`
Receives `{ "theme": string }`. Updates `deck.json`'s theme property, rebuilds deck, and broadcasts `reload` event.

### `POST /api/deck/create`
Receives `{ "name": string, "theme": string, "slides": int, "path": string }`. Scaffolds a new presentation, generates initial build, and returns deck metadata.

### `POST /api/tokens/update`
Receives `{ "tokens": ThemeTokens }`. Updates `tokens.json` in active theme directory and triggers recompile.

### `POST /api/tokens/autofix`
Applies recommended WCAG contrast fixes to the active theme.

### `GET /api/catalog`
Returns complete array of component definitions directly from embedded memory.

### `GET /api/styling-guide`
Returns complete self-describing styling guide including canvas rules, CSS classes, spacing scale, and color token variables.

### `GET /api/events`
Server-Sent Events (SSE) stream broadcasting:
- `event: reload`: signals that compiler has rebuilt `dist/index.html`.
- `event: slide_updated`: signals that slide index has changed.

---

## Model Context Protocol (MCP) Tools

- `deckforge_get_catalog`: Returns full component catalog with input schemas, slots, and snippets.
- `deckforge_get_styling_guide`: Returns comprehensive styling rules, semantic CSS dictionary, spacing scale, and token variables.
- `deckforge_create_deck`: Scaffolds a new presentation deck with chosen theme and slide count.
- `deckforge_list_themes`: Lists all available themes with displayName, vibe, and custom status.
- `deckforge_set_theme`: Updates the deck's theme in `deck.json` and recompiles the presentation.
- `deckforge_create_theme`: Scaffolds a new theme skeleton or cloned preset in workspace or global scope with CSS or SCSS modular files.
- `deckforge_list_slides`: Lists all slides in active deck with indices, filenames, and titles.
- `deckforge_get_slide`: Retrieves complete HTML and extracted text slots of slide `index`.
- `deckforge_update_slide`: Overwrites content of slide `index` with validated HTML.
- `deckforge_reorder_slides`: Re-indexes deck according to `order` array.
- `deckforge_audit_tokens`: Runs WCAG 2.1 contrast evaluation across all theme token pairs.
- `deckforge_apply_token_fix`: Automatically patches failing contrast tokens in `tokens.json`.
- `deckforge_insert_component`: Renders component template `selector` with `props` into target slide.

