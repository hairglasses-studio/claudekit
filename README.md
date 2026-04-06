# claudekit

[![CI](https://github.com/hairglasses-studio/claudekit/actions/workflows/ci.yml/badge.svg)](https://github.com/hairglasses-studio/claudekit/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/hairglasses-studio/claudekit.svg)](https://pkg.go.dev/github.com/hairglasses-studio/claudekit)
[![Go Report Card](https://goreportcard.com/badge/github.com/hairglasses-studio/claudekit)](https://goreportcard.com/report/github.com/hairglasses-studio/claudekit)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Claude Code terminal customization and autonomous R&D framework. Combines environment management (fonts, themes, statusline) with autonomous loop orchestration, budget-aware model routing, and self-improvement cycles.

37 MCP tools across 10 modules: environment setup, roadmap automation, Ralph autonomous loops with budget profiles, and rdcycle self-improving R&D workflows. Built with [mcpkit](https://github.com/hairglasses-studio/mcpkit).

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

### Skills (Marketplace)

| Command | Description |
|---------|-------------|
| `claudekit skill list` | List installed and available skills |
| `claudekit skill install <name>` | Install a skill from the marketplace |
| `claudekit skill remove <name>` | Remove an installed skill |

Skills are Claude Code behavior guides (SKILL.md files) that teach Claude how to use claudekit MCP tools for specific workflows. The built-in marketplace includes: `font-setup`, `theme-setup`, `env-setup`.

### Ralph (Autonomous Loop)

| Command | Description |
|---------|-------------|
| `claudekit ralph tail <file>` | Watch ralph progress file in real time |
| `claudekit ralph status <file>` | Show current progress snapshot |
| `claudekit ralph status <file> --json` | Output progress as JSON |

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
| `skill_list` | List installed and available skills | No |
| `skill_install` | Install a skill from marketplace | Yes |
| `ralph_start` | Start autonomous task loop | Yes |
| `ralph_stop` | Stop running loop | No |
| `ralph_status` | Get loop status | No |
| `roadmap_read` | Read project roadmap | No |
| `roadmap_update` | Update roadmap item status | Yes |
| `roadmap_gaps` | Find incomplete roadmap items | No |
| `roadmap_next_phase` | Get next actionable phase | No |
| `finops_status` | Get token usage summary | No |
| `finops_reset` | Reset token usage counters | Yes |
| `memory_get` | Retrieve value from agent memory | No |
| `memory_set` | Store value in agent memory | Yes |
| `memory_list` | List memory entries with filters | No |
| `memory_search` | Search memory by text query | No |
| `rdcycle_scan` | Scan ecosystem repos for changes | No |
| `rdcycle_plan` | Plan next work from roadmap gaps | No |
| `rdcycle_verify` | Run build/test verification | No |
| `rdcycle_artifacts` | List R&D cycle artifacts | No |
| `rdcycle_commit` | Stage files and commit on feature branch | Yes |
| `rdcycle_report` | Generate RESEARCH-*.md report | Yes |
| `rdcycle_schedule` | Write next cycle's Ralph spec file | Yes |
| `rdcycle_notes` | Record improvement notes per cycle | Yes |
| `rdcycle_improve` | Analyze notes for patterns and cost trends | No |
| `workflow_run` | Execute a predefined workflow | Yes |
| `workflow_list` | List available workflows | No |

## Automated R&D (Tiers 6–7)

Claudekit includes self-improving R&D tools powered by [mcpkit](https://github.com/hairglasses-studio/mcpkit):

- **Roadmap** — Machine-readable `roadmap.json` with gap analysis and phase planning
- **FinOps** — Profile-aware budget tracking with dollar costs, token budgets, and per-tool breakdown
- **Memory** — In-memory key/value store with tiers (episodic/semantic/procedural), tags, and text search
- **R&D Cycle** — Ecosystem scanning, roadmap-driven planning, build verification, artifact tracking, improvement notes, and self-analysis
- **Workflow** — Predefined multi-step graphs (e.g., `full-setup`: detect → font install → theme → statusline → env snapshot)
- **Ralph** — Autonomous loop with budget profiles, model tier selection, cost-breach auto-stop, and MCP sampling

### Budget Profiles

```bash
# Conservative (Claude Max subscription): $5/cycle, 50 iterations
CLAUDEKIT_BUDGET_PROFILE=personal go run ./cmd/claudekit-mcp

# Higher limits (API credits): $50/cycle, 200 iterations
CLAUDEKIT_BUDGET_PROFILE=work-api go run ./cmd/claudekit-mcp

# Custom profile from JSON file
go run ./cmd/claudekit-mcp --budget=my-profile.json
```

### Model Tier Selection

Ralph uses Opus for planning/implementation and Sonnet for lighter tasks:

| Task | Model |
|------|-------|
| plan, implement | claude-opus-4-6 |
| scan, verify, reflect, report, schedule | claude-sonnet-4-6 |

## Architecture

```
claudekit/
├── fontkit/        Font detection, installation, terminal config (pure Go)
├── themekit/       Catppuccin palettes, iTerm2/Ghostty/WezTerm/bat/delta/Starship export
├── statusline/     Claude Code statusline with 3-tier font fallback
├── envkit/         Mise integration, shell detection, dotfile management
├── pluginkit/      YAML plugin loading, subprocess handler, ToolModule bridge
├── skillkit/       Claude Code skill marketplace — discovery, install, manage
├── mcpserver/      MCP tool modules + gateway + ralph + skills + discovery + WebMCP
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

[MIT](LICENSE)
