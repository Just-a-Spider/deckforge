# ADR-0005: WCAG 2.1 Contrast Validator with Recommended Actions

## Status
Accepted

## Context
Presentation themes often use vibrant brand colors that fail accessibility standards when projected in conference rooms or viewed on low-contrast screens. Passive linter warnings are easily ignored by developers and designers.

## Decision
Implement an intelligent token engine in Go with:
1. Exact relative luminance calculation ($L = 0.2126 R + 0.7152 G + 0.0722 B$) and WCAG 2.1 contrast ratio $((L_1 + 0.05) / (L_2 + 0.05))$.
2. Automated evaluation of critical token pairs (`textPrimary` vs `bgCanvas`, `textPrimary` vs `bgSurface`, `textSecondary` vs `bgSurface`, `accentPrimary` vs `bgCanvas`).
3. **Automated Recommended Actions**: A binary search algorithm that adjusts the HSL lightness of the foreground or surface token to determine the minimal shift required to meet 4.5:1 (AA) or 7.0:1 (AAA) compliance.
4. One-click execution in Studio, TUI shortcut, and CLI auto-fix command.

## Consequences
- **Positive**: Proactive accessibility compliance built directly into the development cycle.
