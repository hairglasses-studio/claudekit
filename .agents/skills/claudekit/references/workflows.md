# claudekit Workflow Catalog

## Environment

- Inspect the current shell, plugin manager, mise installation, and managed config inventory before proposing changes.
- Treat snapshot and restore operations as change-safety steps, not the first action.
- Call out when a user needs to reload a shell or restart a terminal app for settings to take effect.

## Fonts

- Confirm the current terminal target first because iTerm2, Ghostty, and WezTerm each need different configuration paths.
- Prefer the documented fallback chain instead of assuming a single installed font family.
- Explain whether the change is an install, a config write, or both.

## Themes

- Choose a Catppuccin flavor intentionally and list the affected targets before writing files.
- Keep terminal color changes separate from auxiliary sync surfaces like `bat`, `delta`, and `starship`.
- Note restart or reload requirements for iTerm2, Ghostty, and WezTerm.

## Statusline And Skills

- Treat statusline installation as a separate follow-up after the environment, font, and theme are stable.
- Skill install or removal work should mention the skill marketplace behavior and any compatibility caveats.
