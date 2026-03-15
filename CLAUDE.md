# claudekit

Claude Code terminal customization framework. Built on [mcpkit](https://github.com/hairglasses-studio/mcpkit).

## Build

```
make check   # vet + test + build
make build   # compile all packages
make test    # run tests
```

## Package Map

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

## MCP Tools

34 tools across 10 modules + dynamic plugin tools:
- **fonts**: `font_status`, `font_install`, `font_configure`
- **theme**: `theme_apply`, `theme_list`
- **statusline**: `statusline_install`
- **env**: `env_status`, `env_snapshot`
- **skills**: `skill_list`, `skill_install`
- **ralph**: `ralph_start`, `ralph_stop`, `ralph_status`
- **roadmap**: `roadmap_read`, `roadmap_update`, `roadmap_gaps`, `roadmap_next_phase`
- **finops**: `finops_status`, `finops_reset`
- **memory**: `memory_get`, `memory_set`, `memory_list`, `memory_search`
- **rdcycle**: `rdcycle_scan`, `rdcycle_plan`, `rdcycle_verify`, `rdcycle_artifacts`, `rdcycle_commit`, `rdcycle_report`, `rdcycle_schedule`, `rdcycle_notes`, `rdcycle_improve`
- **workflow**: `workflow_run`, `workflow_list`

## Key Patterns

- `ToolModule` interface: `Name()`, `Description()`, `Tools() []ToolDefinition`
- Typed inputs/outputs with `jsonschema` tags for MCP tools
- Tests alongside source files (`_test.go`)
- CLI: `os.Args` routing with `parseFlag(key, fallback)` helper
- Context-aware `exec.CommandContext` for all shell commands
- Font fallback: MonaspiceNe → MonaspaceNeon → Menlo
- Theme export: one Catppuccin palette → multiple targets
- Dotfile management: `# claudekit:begin` / `# claudekit:end` markers

## External Dependencies

- [mcpkit](https://github.com/hairglasses-studio/mcpkit) — MCP server framework (local replace `../mcpkit`)
- [Monaspace](https://github.com/githubnext/monaspace) — GitHub's variable font family
- [Monaspice](https://github.com/aaronliu0130/monaspice) — Nerd Font patched Monaspace
- [Catppuccin](https://github.com/catppuccin/catppuccin) — Soothing pastel theme
- [Starship](https://starship.rs/) — Cross-shell prompt
- [mise](https://mise.jdx.dev/) — Tool version manager
- [bat](https://github.com/sharkdp/bat) — Cat clone with syntax highlighting
- [delta](https://github.com/dandavison/delta) — Git diff viewer
