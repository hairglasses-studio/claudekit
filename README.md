# claudekit

Claude Code terminal customization toolkit — fonts, themes, statusline, environment bootstrap, and MCP integration.

Built with [mcpkit](https://github.com/hairglasses-studio/mcpkit) for the MCP server layer.

## Quick Start

```bash
# Install fonts + configure terminal + install statusline
claudekit fonts setup
claudekit statusline install

# Apply Catppuccin Mocha to everything
claudekit theme sync --flavor mocha

# Check environment
claudekit env status
```

## Installation

```bash
go install github.com/hairglasses-studio/claudekit/cmd/claudekit@latest
```

Or build from source:

```bash
git clone https://github.com/hairglasses-studio/claudekit.git
cd claudekit
make build
```

## Claude Code MCP Integration

Add to your project's `.mcp.json`:

```json
{
  "mcpServers": {
    "claudekit": {
      "command": "go",
      "args": ["run", "./cmd/claudekit-mcp"]
    }
  }
}
```

Or with gateway for aggregating external MCP servers:

```json
{
  "mcpServers": {
    "claudekit": {
      "command": "go",
      "args": ["run", "./cmd/claudekit-mcp", "--gateway=github=http://localhost:8080"]
    }
  }
}
```

## CLI Reference

### Fonts

| Command | Description |
|---------|-------------|
| `claudekit fonts status` | Detect installed fonts and terminal emulator |
| `claudekit fonts install` | Install Monaspice Nerd Font via Homebrew |
| `claudekit fonts configure` | Write font config for detected terminal |
| `claudekit fonts configure --terminal ghostty` | Write font config for specific terminal |
| `claudekit fonts preview` | Show ligatures, Nerd Font glyphs, fallback tiers |
| `claudekit fonts setup` | All-in-one: detect, install, configure |

Supported terminals: **iTerm2** (Dynamic Profiles), **Ghostty** (config file), **WezTerm** (Lua module).

Font families:
- [Monaspace](https://github.com/githubnext/monaspace) — 5 subfamilies (Argon, Neon, Xenon, Radon, Krypton)
- [Monaspice](https://github.com/aaronliu0130/monaspice) — Nerd Font patched variants with icon glyphs

### Themes

| Command | Description |
|---------|-------------|
| `claudekit theme apply` | Apply Catppuccin theme to terminal |
| `claudekit theme apply --flavor macchiato --terminal ghostty` | Apply specific flavor to specific terminal |
| `claudekit theme sync --flavor mocha` | Apply theme to terminal + bat + delta + starship |
| `claudekit theme preview` | Preview all Catppuccin flavors with color swatches |

Catppuccin flavors: **Mocha** (dark), **Macchiato** (dark), **Frappe** (dark), **Latte** (light).

Sync targets: terminal (iTerm2/Ghostty/WezTerm), [bat](https://github.com/sharkdp/bat), [delta](https://github.com/dandavison/delta), [Starship](https://starship.rs/).

### Statusline

| Command | Description |
|---------|-------------|
| `claudekit statusline install` | Install Claude Code statusline script + settings |
| `claudekit statusline preview` | Preview statusline with sample data |
| `claudekit statusline render` | Render statusline from stdin JSON (used by the script) |

The statusline shows: model name, CWD slug, context usage bar, cost, and duration.

Set `CLAUDEKIT_THEME=mocha` for Catppuccin-colored statusline output.

Font fallback tiers:
1. **Nerd Font** — full icons (requires Monaspice)
2. **Unicode** — geometric symbols (any modern font)
3. **ASCII** — pure text (any monospace font)

### Environment

| Command | Description |
|---------|-------------|
| `claudekit env status` | Show mise, shell, and managed config info |
| `claudekit env snapshot` | Capture current config files |
| `claudekit env mise` | Install and configure [mise](https://mise.jdx.dev/) |

### Plugins

| Command | Description |
|---------|-------------|
| `claudekit plugin list` | List installed plugins |
| `claudekit plugin add <path>` | Install a plugin from YAML file |

Plugins are YAML files in `~/.claudekit/plugins/`. Each plugin defines tools that are invoked via subprocess:

```yaml
name: my-plugin
description: Example plugin
version: "1.0.0"
handler:
  type: subprocess
  command: my-handler
  timeout: "30s"
tools:
  - name: my_tool
    description: Does something useful
    input_schema:
      type: object
      properties:
        name:
          type: string
```

The plugin subprocess receives `{"method":"<tool>","params":<input>}` on stdin and returns `{"result":...}` on stdout.

### MCP Management

| Command | Description |
|---------|-------------|
| `claudekit mcp tools` | List all registered MCP tools |
| `claudekit mcp publish` | Publish to [MCP Registry](https://registry.modelcontextprotocol.io/) |
| `claudekit mcp serve` | Start WebMCP HTTP server |
| `claudekit mcp serve --addr :9090` | Start on custom port |

## MCP Tools

| Tool | Description | Writes? |
|------|-------------|---------|
| `font_status` | Detect fonts + terminal | No |
| `font_install` | Install fonts via Homebrew | Yes |
| `font_configure` | Write terminal font config | Yes |
| `theme_apply` | Apply Catppuccin theme | Yes |
| `theme_list` | List available flavors | No |
| `statusline_install` | Install Claude Code statusline | Yes |
| `env_status` | Check mise + shell info | No |
| `env_snapshot` | Capture managed config files | No |
| `ralph_start` | Start autonomous task loop | Yes |
| `ralph_stop` | Stop running loop | No |
| `ralph_status` | Get loop status | No |

## Architecture

```
claudekit/
├── fontkit/        Font detection, installation, terminal config (pure Go)
├── themekit/       Catppuccin palettes, iTerm2/Ghostty/WezTerm/bat/delta/Starship export
├── statusline/     Claude Code statusline with 3-tier font fallback
├── envkit/         Mise integration, shell detection, dotfile management
├── pluginkit/      YAML plugin loading, subprocess handler, ToolModule bridge
├── mcpserver/      MCP tool modules + gateway + ralph + discovery + WebMCP
├── cmd/claudekit/  CLI entrypoint
└── cmd/claudekit-mcp/  MCP stdio server entrypoint
```

All packages follow [mcpkit](https://github.com/hairglasses-studio/mcpkit) patterns: `ToolModule` interface, typed handlers, middleware chains.

## Development

```bash
make check   # go vet + go test + go build
make build   # compile all packages
make test    # run tests
make clean   # remove binaries
```

## License

Private — [hairglasses-studio](https://github.com/hairglasses-studio)
