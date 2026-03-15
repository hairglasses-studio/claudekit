# Font Setup Skill

Activates when the user mentions fonts, terminal appearance, Monaspace, Monaspice, or ligatures.

## Workflow

1. Call `font_status` to detect the current terminal and installed fonts
2. Present findings:
   - Terminal: iTerm2 / Ghostty / Apple Terminal / other
   - Installed font families (if any)
   - Best available font from the fallback chain
3. If fonts need installing:
   - Explain: Monaspice Nerd Font = Monaspace + icon glyphs, installed via Homebrew
   - Ask for confirmation
   - Call `font_install` with `nerd_font: true`
4. If terminal supports configuration (iTerm2 or Ghostty):
   - Call `font_configure` to write config with fallback
   - iTerm2: creates Dynamic Profile "Claudekit Monaspace" — user switches manually
   - Ghostty: writes font-family directives, takes effect on restart
5. Suggest `statusline_install` and `theme_apply` for the full experience

## Font Families

| Family | Source | Nerd Glyphs | Brew Cask |
|--------|--------|-------------|-----------|
| [Monaspace](https://github.com/githubnext/monaspace) | GitHub | No | `font-monaspace` |
| [Monaspice](https://github.com/aaronliu0130/monaspice) | Community | Yes | `font-monaspice-nerd-font` |

Subfamilies: Argon, Neon, Xenon, Radon, Krypton (Monaspace) / Ar, Ne, Xe, Rn, Kr (Monaspice)

Default: **MonaspiceNe Nerd Font** (Neon + icons). Fallback: MonaspiceNe → MonaspaceNeon → Menlo.
