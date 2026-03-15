# Font Setup Skill

Activates when the user mentions fonts, terminal appearance, Monaspace, Monaspice, or ligatures.

## Workflow

1. Call `font_status` to detect the current terminal and installed fonts
2. Present findings to the user:
   - Which terminal they're running
   - Which fonts (if any) are already installed
   - The recommended font (MonaspiceNe Nerd Font)
3. If fonts need to be installed:
   - Explain what will be installed (Monaspice Nerd Font via Homebrew)
   - Ask for confirmation before proceeding
   - Call `font_install` with `nerd_font: true`
4. If the terminal supports configuration:
   - Call `font_configure` to write the appropriate config
   - For iTerm2: explain they need to switch to the "Claudekit Monaspace" profile
   - For Ghostty: the config takes effect on restart
5. Suggest installing the statusline with `statusline_install` to showcase the fonts

## Key Details

- **Monaspace** = upstream GitHub font family (5 subfamilies: Argon, Neon, Xenon, Radon, Krypton)
- **Monaspice** = Nerd Font patched variant with icon glyphs (5 variants: Ar, Ne, Xe, Rn, Kr)
- Default recommendation: **MonaspiceNe Nerd Font** (Neon subfamily + icons)
- Fallback chain: MonaspiceNe → MonaspaceNeon → Menlo (system)
- Both are installed via Homebrew casks
