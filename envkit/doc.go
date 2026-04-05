// Package envkit provides environment bootstrap utilities for Claude Code
// terminal customization.
//
// It detects the current shell and plugin manager, manages dotfile snapshots
// for claudekit-controlled config files (Ghostty, Starship, bat, delta,
// iTerm2 Dynamic Profiles), and integrates with mise for development tool
// version management.
//
// # Shell Detection
//
// [DetectShell] identifies the active shell (zsh, bash, fish) and any
// plugin manager (oh-my-zsh, zinit).
//
// # Dotfile Snapshots
//
// [Snapshot] and [Restore] capture and restore the content of all
// claudekit-managed config files, enabling safe before/after comparisons
// during theme and font changes.
//
// # mise Integration
//
// [MiseStatus] checks if mise is installed and lists active tool versions.
// [MiseInstall] bootstraps mise and writes a .mise.toml with recommended
// tool versions for Claude Code development (Go, Node, Python).
package envkit
