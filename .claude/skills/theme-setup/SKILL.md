# Theme Setup Skill

Activates when the user mentions themes, colors, Catppuccin, dark mode, or terminal appearance.

## Workflow

1. Call `theme_list` to show available Catppuccin flavors
2. Ask the user which flavor they prefer (default: Mocha for dark terminals, Latte for light)
3. Call `theme_apply` with their chosen flavor
4. Explain what was configured:
   - For iTerm2: Dynamic Profile with full ANSI color mapping
   - For Ghostty: Theme file in ~/.config/ghostty/themes/
5. Suggest running `claudekit theme sync` to also apply to bat and delta
6. Suggest setting `CLAUDEKIT_THEME=<flavor>` for statusline colors

## Key Details

- **Catppuccin** has 4 flavors: Mocha (dark), Macchiato (dark), Frappé (dark), Latte (light)
- Each flavor has 26 semantic colors (base, text, 14 accent colors, surface layers)
- Theme sync applies to: terminal + bat (syntax highlighter) + delta (git pager)
- The statusline picks up theme colors via the `CLAUDEKIT_THEME` env var
