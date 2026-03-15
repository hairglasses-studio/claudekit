# Environment Setup Skill

Activates when the user mentions environment, mise, tool versions, dotfiles, shell setup, or dev tools.

## Workflow

1. Call `env_status` to check current state:
   - [mise](https://mise.jdx.dev/) installation and active tool versions
   - Shell type (zsh/bash/fish) and plugin manager (oh-my-zsh/zinit)
   - Managed config file inventory
2. If mise not installed, offer to set it up
3. If config files are out of sync, suggest `claudekit theme sync`
4. For dotfile backup, explain `env_snapshot` captures all managed configs

## MCP Tools

| Tool | Description |
|------|-------------|
| `env_status` | Returns mise info, shell info, managed config paths |
| `env_snapshot` | Captures content of all claudekit-managed config files |

## Managed Configs

Files tracked by claudekit's dotfile system:
- `~/.config/starship.toml` — Starship prompt
- `~/.config/ghostty/config` — Ghostty terminal
- `~/.config/bat/config` — bat syntax highlighter
- `~/.config/delta/catppuccin.gitconfig` — delta git pager
- `~/Library/Application Support/iTerm2/DynamicProfiles/claudekit-*.json` — iTerm2 profiles

## Tools Managed by mise

Default `.mise.toml` configuration:
```toml
[tools]
go = "latest"
node = "lts"
python = "latest"
```
