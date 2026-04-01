# claudekit — Gemini CLI Instructions

Claude Code terminal customization framework built on mcpkit. 37 MCP tools across 10 modules for fonts, themes, env, statusline, skills, and autonomous R&D.

## Build & Test

```bash
make check   # vet + test + build
make build   # compile all packages
make test    # run tests
```

## Architecture

- `fontkit/`, `themekit/`, `envkit/` — Pure Go utility packages
- `statusline/`, `pluginkit/`, `skillkit/` — Higher-level features
- `mcpserver/` — MCP tool modules (depends on all packages + mcpkit)
- `cmd/claudekit/` — CLI, `cmd/claudekit-mcp/` — MCP server

## Key Patterns

- `ToolModule` interface: `Name()`, `Description()`, `Tools() []ToolDefinition`
- Typed inputs/outputs with `jsonschema` tags
- Tests alongside source files (`_test.go`)
- Dotfile markers: `# claudekit:begin` / `# claudekit:end`
- Local mcpkit dependency via `replace` directive in go.mod
