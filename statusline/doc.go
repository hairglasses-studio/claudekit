// Package statusline provides a Claude Code statusline renderer and
// installer.
//
// The statusline displays model name, working directory, context usage
// progress bar, cost, and duration. It supports three icon tiers
// (Nerd Font, Unicode, ASCII) and optional Catppuccin color theming
// via the CLAUDEKIT_THEME environment variable.
//
// # Rendering
//
// [Render] reads session JSON from an [io.Reader] and produces styled
// terminal output with ANSI color codes. The output format adapts to
// the detected [FontTier] and active theme.
//
// # Installation
//
// [Install] writes a shell script to ~/.claude/statusline.sh and
// updates ~/.claude/settings.json to register it as the active
// statusline command. If a statusline configuration already exists,
// the previous settings are backed up.
//
// # Script
//
// [ScriptContent] returns the shell script that Claude Code executes.
// When the claudekit binary is available, it delegates to the Go
// renderer; otherwise it falls back to an inline bash implementation.
package statusline
