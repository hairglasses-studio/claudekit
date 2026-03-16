# Community Claude Code Statusline Projects

Reference list of community-built statusline tools for Claude Code CLI.

## Full-Featured

| Project | Language | Install | Key Differentiator |
|---------|----------|---------|-------------------|
| [rz1989s/claude-code-statusline](https://github.com/rz1989s/claude-code-statusline) | Bash | `curl \| bash` | 35 atomic components, TOML config, 9-line layouts, catppuccin theme |
| [sirmalloc/ccstatusline](https://github.com/sirmalloc/ccstatusline) | TypeScript (React/Ink) | `npm install` | tok/sec speed widgets, custom command widget, thinking effort, powerline separators |
| [npow/oh-my-claude](https://github.com/npow/oh-my-claude) | Shell | `curl \| bash` | Plugin framework (`omc create`), 87 plugins, CI status widget, color-shifting context bar |
| [Haleclipse/CCometixLine](https://github.com/Haleclipse/CCometixLine) | Rust | `npm i -g @cometix/ccline` | Cross-platform binaries, TUI config, fast execution |

## Lightweight

| Project | Language | Install | Key Differentiator |
|---------|----------|---------|-------------------|
| [Owloops/claude-powerline](https://github.com/Owloops/claude-powerline) | TypeScript | `npx` | Vim-style powerline, auto-download |
| [chongdashu/cc-statusline](https://github.com/chongdashu/cc-statusline) | TypeScript | `npm install` | Interactive setup wizard, session time tracking |
| [kamranahmedse/claude-statusline](https://github.com/kamranahmedse/claude-statusline) | TypeScript | `npx` | Minimal and opinionated, one-command install |
| [ssenart/oh-my-claude](https://github.com/ssenart/oh-my-claude) | Shell | manual | Oh My Posh rendering, dynamic git status colors (clean/dirty/ahead/behind/diverged) |

## Specialized

| Project | Language | Install | Key Differentiator |
|---------|----------|---------|-------------------|
| [gabriel-dehan/claude_monitor_statusline](https://github.com/gabriel-dehan/claude_monitor_statusline) | Ruby | manual | Token and message usage with plan limits |
| [ryoppippi/ccusage](https://github.com/ryoppippi/ccusage) | TypeScript | `npm install` | Cost-focused, real-time burn rates from local JSONL, offline mode |

## Unique Widgets by Tool

Widgets not available in rz1989s (our installed tool) that exist in other projects:

| Widget | Available In | Description |
|--------|-------------|-------------|
| Input/Output Speed (tok/sec) | ccstatusline | Real-time token throughput with configurable window |
| Custom Command | ccstatusline | Execute any shell command and display output inline |
| Thinking Effort | ccstatusline | Shows current thinking effort level |
| Skills | ccstatusline | Last-used skill, count, list with hide-when-empty |
| Session Name | ccstatusline | From `/rename` command |
| Flex Separator | ccstatusline | Adaptive width spacing between components |
| CI Status | oh-my-claude | CI/CD pipeline status |
| Plugin System | oh-my-claude | Create custom widgets via `omc create my-plugin` |
| Dynamic Git Colors | oh-my-claude (ssenart) | Repo state color-coded (clean/dirty/ahead/behind) |

## Most Requested Features (Not Available Anywhere Yet)

- Rate limit utilization % ([anthropics/claude-code#30784](https://github.com/anthropics/claude-code/issues/30784))
- Background task count ([anthropics/claude-code#33310](https://github.com/anthropics/claude-code/issues/33310))
- Timer-based auto-refresh ([anthropics/claude-code#5685](https://github.com/anthropics/claude-code/issues/5685))
- Extended thinking mode status ([anthropics/claude-code#26279](https://github.com/anthropics/claude-code/issues/26279))

## Our Current Setup

Using **rz1989s/claude-code-statusline v2.21.5** with catppuccin theme, 2-line layout:

```
Line 1: repo_info │ model_info │ cache_efficiency │ code_productivity │ version_info
Line 2: cost_live │ cost_repo │ burn_rate │ context_window │ block_projection │ usage_reset
```

Config: `~/.claude/statusline/Config.toml`
