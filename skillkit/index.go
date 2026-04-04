package skillkit

// IndexEntry describes a skill available in the marketplace.
type IndexEntry struct {
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tools       []string `json:"tools"`   // MCP tools this skill uses
	Content     string   `json:"content"` // Full SKILL.md content
}

// BuiltinIndex returns the built-in skill marketplace entries.
// These are the claudekit-authored skills that ship with the binary.
func BuiltinIndex() []IndexEntry {
	return []IndexEntry{
		{
			Name:        "font-setup",
			Title:       "Font Setup Skill",
			Description: "Activates when the user mentions fonts, terminal appearance, Monaspace, Monaspice, or ligatures.",
			Tools:       []string{"font_status", "font_install", "font_configure"},
			Content:     fontSetupSkill,
		},
		{
			Name:        "theme-setup",
			Title:       "Theme Setup Skill",
			Description: "Activates when the user mentions themes, colors, Catppuccin, dark mode, or terminal appearance.",
			Tools:       []string{"theme_apply", "theme_list"},
			Content:     themeSetupSkill,
		},
		{
			Name:        "env-setup",
			Title:       "Environment Setup Skill",
			Description: "Activates when the user mentions environment, mise, tool versions, dotfiles, shell setup, or dev tools.",
			Tools:       []string{"env_status", "env_snapshot"},
			Content:     envSetupSkill,
		},
	}
}

// FindInIndex searches the builtin index by name.
func FindInIndex(name string) *IndexEntry {
	for _, e := range BuiltinIndex() {
		if e.Name == name {
			return &e
		}
	}
	return nil
}

// AvailableSkills returns index entries that are not yet installed.
func AvailableSkills(projectDir string) ([]IndexEntry, error) {
	installed, err := ListInstalled(projectDir)
	if err != nil {
		return nil, err
	}

	installedSet := make(map[string]bool, len(installed))
	for _, s := range installed {
		installedSet[s.Name] = true
	}

	var available []IndexEntry
	for _, e := range BuiltinIndex() {
		if !installedSet[e.Name] {
			available = append(available, e)
		}
	}
	return available, nil
}

const fontSetupSkill = `# Font Setup Skill

Activates when the user mentions fonts, terminal appearance, Monaspace, Monaspice, or ligatures.

## Workflow

1. Call ` + "`font_status`" + ` to detect the current terminal and installed fonts
2. Present findings:
   - Terminal: iTerm2 / Ghostty / WezTerm / Apple Terminal / other
   - Installed font families (if any)
   - Best available font from the fallback chain
3. If fonts need installing:
   - Explain: Monaspice Nerd Font = Monaspace + icon glyphs, installed via Homebrew
   - Ask for confirmation
   - Call ` + "`font_install`" + ` with ` + "`nerd_font: true`" + `
4. If terminal supports configuration (iTerm2, Ghostty, or WezTerm):
   - Call ` + "`font_configure`" + ` to write config with fallback
   - iTerm2: creates Dynamic Profile "Claudekit Monaspace" — user switches manually
   - Ghostty: writes font-family directives, takes effect on restart
   - WezTerm: writes Lua module, user adds require("claudekit").apply(config, wezterm)
5. Suggest ` + "`statusline_install`" + ` and ` + "`theme_apply`" + ` for the full experience

## Font Families

| Family | Source | Nerd Glyphs | Brew Cask |
|--------|--------|-------------|-----------|
| [Monaspace](https://github.com/githubnext/monaspace) | GitHub | No | ` + "`font-monaspace`" + ` |
| [Monaspice](https://github.com/aaronliu0130/monaspice) | Community | Yes | ` + "`font-monaspice-nerd-font`" + ` |

Subfamilies: Argon, Neon, Xenon, Radon, Krypton (Monaspace) / Ar, Ne, Xe, Rn, Kr (Monaspice)

Default: **MonaspiceNe Nerd Font** (Neon + icons). Fallback: MonaspiceNe → MonaspaceNeon → Menlo.
`

const themeSetupSkill = `# Theme Setup Skill

Activates when the user mentions themes, colors, Catppuccin, dark mode, or terminal appearance.

## Workflow

1. Call ` + "`theme_list`" + ` to show available [Catppuccin](https://github.com/catppuccin/catppuccin) flavors
2. Ask which flavor (default: Mocha for dark, Latte for light)
3. Call ` + "`theme_apply`" + ` with their choice
4. Explain what was configured:
   - iTerm2: Dynamic Profile with ANSI color mapping
   - Ghostty: theme file in ` + "`~/.config/ghostty/themes/`" + `
   - WezTerm: Lua color table in ` + "`~/.config/wezterm/claudekit-colors.lua`" + `
5. Suggest ` + "`claudekit theme sync`" + ` for bat + delta + Starship too
6. Suggest ` + "`CLAUDEKIT_THEME=<flavor>`" + ` env var for statusline colors

## Flavors

| Flavor | Type | Base Color |
|--------|------|-----------|
| Mocha | Dark | ` + "`#1e1e2e`" + ` |
| Macchiato | Dark | ` + "`#24273a`" + ` |
| Frappe | Dark | ` + "`#303446`" + ` |
| Latte | Light | ` + "`#eff1f5`" + ` |

## Sync Targets

| Target | Config Path | Tool |
|--------|------------|------|
| iTerm2 | ` + "`~/Library/Application Support/iTerm2/DynamicProfiles/`" + ` | [iTerm2](https://iterm2.com/) |
| Ghostty | ` + "`~/.config/ghostty/themes/`" + ` | [Ghostty](https://ghostty.org/) |
| WezTerm | ` + "`~/.config/wezterm/claudekit-colors.lua`" + ` | [WezTerm](https://wezfurlong.org/wezterm/) |
| bat | ` + "`~/.config/bat/config`" + ` | [bat](https://github.com/sharkdp/bat) |
| delta | ` + "`~/.config/delta/catppuccin.gitconfig`" + ` | [delta](https://github.com/dandavison/delta) |
| Starship | ` + "`~/.config/starship.toml`" + ` | [Starship](https://starship.rs/) |
`

const envSetupSkill = `# Environment Setup Skill

Activates when the user mentions environment, mise, tool versions, dotfiles, shell setup, or dev tools.

## Workflow

1. Call ` + "`env_status`" + ` to check current state:
   - [mise](https://mise.jdx.dev/) installation and active tool versions
   - Shell type (zsh/bash/fish) and plugin manager (oh-my-zsh/zinit)
   - Managed config file inventory
2. If mise not installed, offer to set it up
3. If config files are out of sync, suggest ` + "`claudekit theme sync`" + `
4. For dotfile backup, explain ` + "`env_snapshot`" + ` captures all managed configs

## MCP Tools

| Tool | Description |
|------|-------------|
| ` + "`env_status`" + ` | Returns mise info, shell info, managed config paths |
| ` + "`env_snapshot`" + ` | Captures content of all claudekit-managed config files |

## Managed Configs

Files tracked by claudekit's dotfile system:
- ` + "`~/.config/starship.toml`" + ` — Starship prompt
- ` + "`~/.config/ghostty/config`" + ` — Ghostty terminal
- ` + "`~/.config/wezterm/claudekit.lua`" + ` — WezTerm font config
- ` + "`~/.config/wezterm/claudekit-colors.lua`" + ` — WezTerm color theme
- ` + "`~/.config/bat/config`" + ` — bat syntax highlighter
- ` + "`~/.config/delta/catppuccin.gitconfig`" + ` — delta git pager
- ` + "`~/Library/Application Support/iTerm2/DynamicProfiles/claudekit-*.json`" + ` — iTerm2 profiles

## Tools Managed by mise

Default ` + "`.mise.toml`" + ` configuration:
` + "```toml" + `
[tools]
go = "latest"
node = "lts"
python = "latest"
` + "```" + `
`
