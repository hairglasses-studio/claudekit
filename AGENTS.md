# claudekit — Agent Instructions

Claude Code terminal customization framework. Built on [mcpkit](https://github.com/hairglasses-studio/mcpkit).

## Build & Test

```
make check   # vet + test + build
make build   # compile all packages
make test    # run tests
```

## Architecture

| Package | Purpose | Dependencies |
|---------|---------|-------------|
| `fontkit` | Font detection, installation, terminal config | None (pure Go) |
| `themekit` | Catppuccin palettes + theme export (iTerm2, Ghostty, WezTerm, bat, delta, Starship) | None (pure Go) |
| `envkit` | Mise integration, shell detection, dotfile snapshot/restore | None (pure Go) |
| `statusline` | Claude Code statusline script + installer | `themekit` |
| `pluginkit` | YAML plugin loading, subprocess handler, ToolModule bridge | `mcpkit/registry`, `yaml.v3` |
| `skillkit` | Claude Code skill marketplace — discovery, install, remove | None (pure Go) |
| `mcpserver` | MCP tool modules (font, theme, env, statusline, skill, ralph, roadmap, finops, memory, rdcycle, workflow, gateway, discovery, webmcp) | `fontkit`, `themekit`, `envkit`, `statusline`, `skillkit`, mcpkit |
| `cmd/claudekit` | CLI entrypoint | all packages |
| `cmd/claudekit-mcp` | MCP server entrypoint | `mcpserver`, `pluginkit` |

## Key Conventions

- `ToolModule` interface: `Name()`, `Description()`, `Tools() []ToolDefinition`
- Typed inputs/outputs with `jsonschema` tags for MCP tools
- Tests alongside source files (`_test.go`)
- CLI: `os.Args` routing with `parseFlag(key, fallback)` helper
- Context-aware `exec.CommandContext` for all shell commands
- Font fallback: MonaspiceNe → MonaspaceNeon → Menlo
- Theme export: one Catppuccin palette → multiple targets
- Dotfile management: `# claudekit:begin` / `# claudekit:end` markers


## Shared Research Repository

Cross-project research lives at `~/hairglasses-studio/docs/` (git: hairglasses-studio/docs). When launching research agents, check existing docs first and write reusable research outputs back to the shared repo rather than local docs/.
