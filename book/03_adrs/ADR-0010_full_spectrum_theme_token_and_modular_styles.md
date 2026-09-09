# ADR-0010: Full-Spectrum Theme Tokens, Modular Stylesheets, and SCSS Preprocessing

## Status
Accepted

## Context
Prior to this architectural evolution, DeckForge themes only managed color palettes (`ThemePalette`) and font families (`ThemeFonts`). Critical design attributes such as card corner radii, drop shadows, backdrop-filter blur, border widths, inner padding, and motion timings were hardcoded directly in `engine/core/components.css` (e.g. 12px border radius, 12px blur).

Consequently, disparate design aesthetics (such as brutalist monochrome `swiss-minimal`, delicate literary `paper-ink`, and glowing cyberpunk `neon-cyber`) were constrained to the same rounded glassmorphic card silhouettes. Furthermore, developers and LLM agents were confined to only 3 static CSS files (`typography.css`, `surfaces.css`, `backdrop.css`) without support for dedicated archetype stylesheets (`components.css`), animations (`animations.css`), granular styles (`styles/*.css`), or SCSS preprocessor workflows. Additionally, monolithic Go source files (`presets.go` at 715 lines) hindered long-term maintainability.

## Decision
1. **Decouple Platform Boundaries from Theme Authority**:
   - **Engine Scope (Platform Basics)**: 1920x1080 fixed stage coordinate engine, viewport scaling (`transform: scale(...)`), slide bounding box and vertical content budgeting (860px max height), slide lifecycle, Master Index drawer, and Studio WYSIWYG editor.
   - **Theme Scope (Design Authority)**: All visual design tokens and presentation styling are delegated entirely to the theme.
2. **Full-Spectrum Token Architecture**:
   - Expand `ThemeTokens` in `internal/models/theme_tokens.go` to include:
     - `Geometry`: `RadiusCard`, `RadiusSm`, `RadiusLg`, `BorderWidth`, `BorderStyle`.
     - `Shadows`: `Card`, `Elevated`, `Glow`, `Flat`.
     - `Surfaces`: `BackdropBlur`, `SurfaceOpacity`.
     - `Spacing`: `CardPadding`.
     - `Motion`: `SlideTransition`, `AnimationSpeed`.
   - Emit standard CSS variables: `--radius-card`, `--radius-sm`, `--border-width`, `--border-style`, `--shadow-card`, `--shadow-elevated`, `--shadow-glow`, `--backdrop-blur`, `--card-padding`, `--slide-transition`, `--animation-speed`.
   - Update `engine/core/components.css` and `engine/core/stage.css` to use these variables with backward-compatible fallbacks.
3. **Modular Stylesheet Architecture & SCSS Preprocessing**:
   - Expand theme loading in `internal/theme/loader.go` to support:
     - Canonical segments: `typography.css`, `surfaces.css`, `backdrop.css`, `components.css`, `animations.css`.
     - Open directory scanning for arbitrary modular stylesheets in `styles/*.css` or `styles/*.scss`.
   - Implement `internal/theme/scss.go` to automatically compile `.scss` when `sass`, `dart-sass`, or `npx sass` is present, while preserving native modern CSS syntax as a zero-dependency default.
4. **Go Codebase Modularization**:
   - Fragment `internal/models/theme.go` into `theme_tokens.go` (token schemas and variable emitter) and `theme.go` (composition and consolidation).
   - Fragment `internal/theme/presets.go` into modular packages by aesthetic cluster: `presets_cyber.go`, `presets_editorial.go`, `presets_modern.go`, and `presets_pastel.go`.
   - Separate file I/O and scanning into `internal/theme/loader.go`, leaving `internal/theme/manager.go` as a clean coordinator.

## Consequences
- **Positive**: Complete visual freedom for theme authors and LLM agents; brutalist themes render sharp 0px corners and flat drop shadows while cyber themes render glowing translucent cards; styles can be split across granular modular files and SCSS; Go codebase is clean and maintainable.
- **Negative**: Preprocessor requires `sass` or `npx` if developers choose to use SCSS features beyond modern CSS native nesting.
