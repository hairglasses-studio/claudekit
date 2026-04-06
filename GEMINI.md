# claudekit — Gemini CLI Instructions

This repo uses [AGENTS.md](AGENTS.md) as the canonical instruction file. Treat this file as compatibility guidance for Gemini-specific workflows.

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


## Shared Research Repository

Cross-project research lives at `~/hairglasses-studio/docs/` (git: hairglasses-studio/docs). When launching research agents, check existing docs first and write reusable research outputs back to the shared repo rather than local docs/.
