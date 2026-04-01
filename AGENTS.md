# claudekit — Agent Instructions

Claude Code terminal customization framework built on mcpkit. Provides font detection, theme export, env management, statusline, skill marketplace, and 37 MCP tools across 10 modules.

## Build & Test

```bash
make check   # vet + test + build
make build   # compile all packages
make test    # run tests
```

## Architecture

```
claudekit/
├── fontkit/         # Font detection, installation, terminal config (pure Go)
├── themekit/        # Catppuccin palettes + theme export (iTerm2, Ghostty, WezTerm, bat, delta, Starship)
├── envkit/          # Mise integration, shell detection, dotfile snapshot/restore
├── statusline/      # Claude Code statusline script + installer (depends: themekit)
├── pluginkit/       # YAML plugin loading, subprocess handler, ToolModule bridge
├── skillkit/        # Claude Code skill marketplace (pure Go)
├── mcpserver/       # MCP tool modules (depends: all packages + mcpkit)
├── cmd/claudekit/   # CLI entrypoint
└── cmd/claudekit-mcp/ # MCP server entrypoint
```

## Key Patterns

- `ToolModule` interface: `Name()`, `Description()`, `Tools() []ToolDefinition`
- Typed inputs/outputs with `jsonschema` tags for MCP tools
- Tests alongside source files (`_test.go`)
- CLI: `os.Args` routing with `parseFlag(key, fallback)` helper
- Context-aware `exec.CommandContext` for all shell commands
- Font fallback: MonaspiceNe -> MonaspaceNeon -> Menlo
- Theme export: one Catppuccin palette -> multiple targets
- Dotfile management: `# claudekit:begin` / `# claudekit:end` markers

## MCP Tools (37 tools, 10 modules)

- fonts: `font_status`, `font_install`, `font_configure`
- theme: `theme_apply`, `theme_list`
- statusline: `statusline_install`
- env: `env_status`, `env_snapshot`
- skills: `skill_list`, `skill_install`
- ralph: `ralph_start`, `ralph_stop`, `ralph_status`
- roadmap: `roadmap_read`, `roadmap_update`, `roadmap_gaps`, `roadmap_next_phase`
- finops: `finops_status`, `finops_reset`
- memory: `memory_get`, `memory_set`, `memory_list`, `memory_search`
- rdcycle: `rdcycle_scan`, `rdcycle_plan`, `rdcycle_verify`, `rdcycle_artifacts`, `rdcycle_commit`, `rdcycle_report`, `rdcycle_schedule`, `rdcycle_notes`, `rdcycle_improve`, `rdcycle_perpetual_start`, `rdcycle_perpetual_stop`, `rdcycle_perpetual_status`
- workflow: `workflow_run`, `workflow_list`

## Dependencies

- mcpkit (local replace `../mcpkit`) — MCP server framework
- Monaspace / Monaspice — Font families
- Catppuccin — Theme palette
- Starship, mise, bat, delta — External tools

## Coding Conventions

- Middleware signature: `func(name string, td registry.ToolDefinition, next registry.ToolHandlerFunc) registry.ToolHandlerFunc`
- Error handling: return `ErrorResult(err), nil` for expected errors
- Tests: `*_test.go` in same package
