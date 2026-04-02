# claudekit — Gemini CLI Instructions

Claude Code terminal customization framework. Built on [mcpkit](https://github.com/hairglasses-studio/mcpkit).

## Build & Test

```
make check   # vet + test + build
make build   # compile all packages
make test    # run tests
```

## Architecture


## Key Conventions

- `ToolModule` interface: `Name()`, `Description()`, `Tools() []ToolDefinition`
- Typed inputs/outputs with `jsonschema` tags for MCP tools
- Tests alongside source files (`_test.go`)
- CLI: `os.Args` routing with `parseFlag(key, fallback)` helper
- Context-aware `exec.CommandContext` for all shell commands
- Font fallback: MonaspiceNe → MonaspaceNeon → Menlo
- Theme export: one Catppuccin palette → multiple targets
- Dotfile management: `# claudekit:begin` / `# claudekit:end` markers

