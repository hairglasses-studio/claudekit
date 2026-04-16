# claudekit — Agent Instructions

> **DEPRECATED.** This repo has been retired. Its functionality was consolidated into
> [`mcpkit`](https://github.com/hairglasses-studio/mcpkit) and
> [`dotfiles/mcp/`](https://github.com/hairglasses-studio/dotfiles).

## Status

Both `cmd/claudekit` and `cmd/claudekit-mcp` are deprecation stubs that print a
retirement message and exit. The packages advertised in the previous architecture
table (`fontkit`, `themekit`, `envkit`, `statusline`, `pluginkit`, `skillkit`,
`mcpserver`) have been removed. The repo's `.mcp.json` is `{}`.

Do not add new features here. Any ongoing work should land in one of the canonical
downstream homes:

- **MCP server framework** → [`mcpkit`](https://github.com/hairglasses-studio/mcpkit)
- **Claude Code skill marketplace / skill surface** → [`dotfiles/.agents/skills`](https://github.com/hairglasses-studio/dotfiles) and the canonical skill surface under `.agents/skills/`
- **Terminal/font/theme customization** → [`dotfiles`](https://github.com/hairglasses-studio/dotfiles) theme pipeline
- **Statusline** → `dotfiles/scripts/hg-statusline.sh` or equivalent

## Historical context

The v3 llm-ops audit (docs@a5079e0e, wave-2-05) flagged claudekit as a
declaration/runtime-inverted shadow surface: `workspace/manifest.json` still
classified it `active_first_party` / `gateway` despite the code being stub-only.
This AGENTS.md rewrite brings the repo-owned documentation in line with the
runtime state. A matching manifest update moves claudekit to `lifecycle: "deprecated"`
and `mcp_surface_class: "none"`.