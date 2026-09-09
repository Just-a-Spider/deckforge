# ADR-0003: Angular-Inspired Declarative Component Specification

## Status
Accepted

## Context
DeckForge needs a library of reusable slide components (e.g. Stat Cards, Comparison Columns, Architecture Nodes). If components are merely hardcoded HTML strings, they cannot be validated, typed, or transformed into an internal component framework in the future. Conversely, adopting heavy web frameworks (like Angular/React npm bundles) violates DeckForge's zero-dependency pure Go philosophy.

## Decision
Define components using a declarative schema inspired by Angular's `@Component` metadata:
- **Selector**: Canonical identifier (e.g. `df-metric-card`).
- **Inputs/Props**: Typed parameters with defaults (e.g. `value: string`, `accentEdge: boolean`).
- **Slots**: Named projection targets (e.g. `notes`, `footer`).
- **Template Snippet**: HTML template parameterized by inputs and slots.
- **Classes**: Encapsulated utility classes applied to the root element.

In the initial implementation, inserting a component renders its instantiated HTML snippet directly into the slide file. In future phases, this schema serves as the foundation for an internal Go-based AST component compiler that supports `<df-metric-card ...>` tags directly in slide source.

## Consequences
- **Positive**: Clean HTML in slide files today; seamless evolution to custom component compilation tomorrow.
- **Positive**: LLMs and Studio UIs can introspect input schemas and generate components reliably.
