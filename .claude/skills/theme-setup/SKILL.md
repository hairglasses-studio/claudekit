# Theme Setup Skill

Activates when the user mentions themes, colors, Catppuccin, dark mode, or terminal appearance.

## Workflow

1. Call `theme_list` to show available [Catppuccin](https://github.com/catppuccin/catppuccin) flavors
2. Ask which flavor (default: Mocha for dark, Latte for light)
3. Call `theme_apply` with their choice
4. Explain what was configured:
   - iTerm2: Dynamic Profile with ANSI color mapping
   - Ghostty: theme file in `~/.config/ghostty/themes/`
5. Suggest `claudekit theme sync` for bat + delta + Starship too
6. Suggest `CLAUDEKIT_THEME=<flavor>` env var for statusline colors

## Flavors

| Flavor | Type | Base Color |
|--------|------|-----------|
| Mocha | Dark | `#1e1e2e` |
| Macchiato | Dark | `#24273a` |
| Frappe | Dark | `#303446` |
| Latte | Light | `#eff1f5` |

## Sync Targets

| Target | Config Path | Tool |
|--------|------------|------|
| iTerm2 | `~/Library/Application Support/iTerm2/DynamicProfiles/` | [iTerm2](https://iterm2.com/) |
| Ghostty | `~/.config/ghostty/themes/` | [Ghostty](https://ghostty.org/) |
| bat | `~/.config/bat/config` | [bat](https://github.com/sharkdp/bat) |
| delta | `~/.config/delta/catppuccin.gitconfig` | [delta](https://github.com/dandavison/delta) |
| Starship | `~/.config/starship.toml` | [Starship](https://starship.rs/) |
