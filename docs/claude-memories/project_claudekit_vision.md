---
name: claudekit-vision
description: Claudekit is a Claude Code terminal customization bootstrapper modeled after mcpkit architecture
type: project
---

Claudekit is a Go-based Claude Code extension framework for bootstrapping terminal customizations (fonts, themes, statusline, environment).

**Why:** No existing tool manages Claude Code's terminal appearance holistically. The user wants Monaspace/Monaspice font support with fallback as the first feature, followed by a rich statusline, then broader theming/environment tools.

**How to apply:** All code should follow mcpkit's architectural patterns (module-based, middleware chains, typed handlers, layered deps). Deliver as CLI + MCP server + Claude Code skill triple. Design for open source distinction.

**Key sibling project:** ~/hairglasses-studio/mcpkit — the reference architecture for claudekit's MCP server layer.
