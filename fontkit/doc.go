// Package fontkit provides font detection, installation, and terminal
// configuration for the Monaspace and Monaspice (Nerd Font) font families.
//
// It supports three terminal emulators: iTerm2, Ghostty, and WezTerm.
// Each terminal has its own configuration writer that generates the
// appropriate config format (iTerm2 Dynamic Profiles JSON, Ghostty
// key=value config, WezTerm Lua modules).
//
// # Font Families
//
// Monaspace is the upstream GitHub monospace font with five subfamilies
// (Argon, Neon, Xenon, Radon, Krypton). Monaspice is the community
// Nerd Font patched variant with icon glyphs (Ar, Ne, Xe, Rn, Kr).
//
// [DefaultFont] is MonaspiceNe (Neon + Nerd Font icons).
// [FallbackChain] returns the ordered font preference for terminal config.
//
// # Detection and Installation
//
// [Detect] scans system font directories for installed families.
// [Install] uses Homebrew to install missing font casks.
//
// # Terminal Configuration
//
// [ConfigureGhostty], [ConfigureITerm2], and [ConfigureWezTerm] write
// font configuration for the respective terminals.
package fontkit
