# ADR-0004: Native MCP and Headless Agent Interface

## Status
Accepted

## Context
Coding agents (Claude Code, Gemini CLI, Antigravity CLI, OpenClaw, OpenCode) are rapidly becoming standard pair-programming operators. To manipulate presentations effectively, agents must not be forced to write fragile bash regexes or parse terminal ANSI strings.

## Decision
Provide two first-class interfaces for AI agents:
1. **Headless CLI Subcommands (`deckforge agent ...`)**:
   - Machine-readable JSON output for slide inspection, modification, reordering, and token auditing.
2. **Embedded Model Context Protocol (MCP) Server**:
   - Implements the MCP tool specification over stdio and SSE at `/mcp`.
   - Exposes tools: `deckforge_list_slides`, `deckforge_get_slide`, `deckforge_update_slide`, `deckforge_audit_tokens`, `deckforge_apply_token_fix`, `deckforge_insert_component`.
3. **Exportable Agent Skill (`assets/agent/SKILL.md`)**:
   - Provides system prompts, 1080p canvas rules, and CSS class dictionaries ready to install into agent work environments.

## Consequences
- **Positive**: Any modern AI agent can inspect and build decks with full semantic understanding and zero hallucination of geometry.
