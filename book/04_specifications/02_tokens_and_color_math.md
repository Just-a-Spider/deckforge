# Tokens & WCAG Color Mathematics

## Relative Luminance Formula (WCAG 2.1)

For an sRGB color with 8-bit channels $R_{8bit}, G_{8bit}, B_{8bit} \in [0, 255]$:

1. Normalize channels to $C_{srgb} = C_{8bit} / 255.0$.
2. Convert to linear RGB:
   $$C_{linear} = \begin{cases} \frac{C_{srgb}}{12.92} & \text{if } C_{srgb} \le 0.04045 \\ \left(\frac{C_{srgb} + 0.055}{1.055}\right)^{2.4} & \text{otherwise} \end{cases}$$
3. Compute relative luminance $L$:
   $$L = 0.2126 R_{linear} + 0.7152 G_{linear} + 0.0722 B_{linear}$$

## Contrast Ratio Formula

Given two colors with relative luminances $L_1$ and $L_2$ where $L_1 \ge L_2$:

$$\text{Ratio} = \frac{L_1 + 0.05}{L_2 + 0.05}$$

- **WCAG AA (Normal Text)**: Ratio $\ge 4.5:1$
- **WCAG AA (Large Text $\ge 24\text{px}$)**: Ratio $\ge 3.0:1$
- **WCAG AAA (Enhanced)**: Ratio $\ge 7.0:1$

## Recommended Action Binary Search Algorithm

When a token pair fails WCAG AA:
1. Convert the target foreground or surface color from hex to HSL ($H, S, L$).
2. Determine direction:
   - If background is dark ($L_{bg} < 0.5$), increase $L_{fg} \to 1.0$.
   - If background is light ($L_{bg} \ge 0.5$), decrease $L_{fg} \to 0.0$.
3. Execute binary search over $L \in [0.0, 1.0]$ with tolerance $\epsilon = 0.01$ until $\text{Ratio} \ge 4.55$ (safety margin above AA threshold).
4. Convert optimal $(H, S, L_{opt})$ back to hexadecimal color string.
5. Return structured `RecommendedAction` record.

---

## Full-Spectrum Theme Token Dictionary

Themes are defined via `tokens.json` and compiled into `:root` CSS custom properties:

| Category | Token Variable | Role / Purpose | Fallback Default |
|---|---|---|---|
| **Palette** | `--bg-canvas` | Primary slide backdrop fill | `#090d16` |
| | `--bg-surface` | Default card & panel surface | `rgba(15, 23, 42, 0.85)` |
| | `--bg-surface-elevated`| Modals, drawers, elevated surfaces | `#0f172a` |
| | `--text-primary` | Main titles and high-contrast copy | `#f8fafc` |
| | `--text-secondary` | Body paragraphs (min 4.5:1 AA) | `#94a3b8` |
| | `--text-muted` | Kickers, timestamps, footnotes | `#64748b` |
| | `--accent-primary` | Primary visual anchor & focus | `#10b981` |
| | `--accent-secondary` | Secondary visual contrast | `#818cf8` |
| | `--accent-glow` | Translucent glow for focal items | `rgba(16, 185, 129, 0.18)` |
| | `--accent-danger` | Critical penalties & error warnings | `#ef4444` |
| | `--border-color` | Primary continuous card borders | `#10b981` |
| | `--border-faint` | Hairlines and structural dividers | `rgba(255, 255, 255, 0.12)` |
| **Typography** | `--font-display` | Slide kickers & card headlines | `'Poppins', sans-serif` |
| | `--font-body` | Paragraphs and narrative text | `'Plus Jakarta Sans', sans-serif`|
| | `--font-mono` | Code blocks, metrics, tags | `'JetBrains Mono', monospace` |
| **Geometry** | `--radius-card` | Panels and major cards | `12px` (or `0px` in Swiss) |
| | `--radius-sm` | Badges, pills, tags, chips | `4px` |
| | `--radius-lg` | Outer containers & hero cards | `16px` |
| | `--border-width` | Border thickness across components| `1.5px` |
| | `--border-style` | Border stroke style | `solid` |
| **Elevation** | `--shadow-card` | Standard card drop shadow | `0 10px 30px rgba(0,0,0,0.25)` |
| | `--shadow-elevated` | Modal/inspector elevation shadow | `0 20px 40px rgba(0,0,0,0.35)` |
| | `--shadow-glow` | Optical radiant neon aura | `0 0 20px var(--accent-glow)` |
| **Materials** | `--backdrop-blur` | Glassmorphic blur radius | `12px` (or `0px` in print) |
| | `--surface-opacity` | Surface translucency coefficient | `0.85` |
| **Spacing** | `--card-padding` | Inner card padding boundary | `28px` |
| **Motion** | `--slide-transition` | Slide transition mode | `fade` / `slide` / `none` |
| | `--animation-speed` | Speed of transitions & hovers | `0.25s` |

---

## Modular Stylesheets & SCSS Pipeline

Themes support open modular stylesheets in their folder:
1. `typography.css` / `.scss`: Custom font weights, headings, display kickers.
2. `surfaces.css` / `.scss`: Card hover lifts, glass treatments, panel variants.
3. `backdrop.css` / `.scss`: Canvas gradients, geometric meshes, grid patterns.
4. `components.css` / `.scss`: Angular archetype overrides (`df-metric-card`, etc.).
5. `animations.css` / `.scss`: Keyframe animations, pulse glows, entrance effects.
6. `styles/*.css` / `*.scss`: Arbitrary granular modular files stitched in alphabetical order.
