# claudekit

Claude Code terminal customization framework. Follows mcpkit patterns.

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
| `statusline` | Claude Code statusline script + installer | `fontkit` |
| `mcpserver` | MCP tool modules | `fontkit`, `statusline`, mcpkit |
| `cmd/claudekit` | CLI entrypoint | `fontkit`, `statusline` |
| `cmd/claudekit-mcp` | MCP server entrypoint | `mcpserver` |

## Conventions

- Typed inputs/outputs with `jsonschema` tags for MCP tools
- `ToolModule` interface: `Name()`, `Description()`, `Tools()`
- Tests alongside source files (`_test.go`)
- No heavy CLI frameworks — `os.Args` routing
- Context-aware exec for all shell commands
