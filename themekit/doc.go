// Package themekit provides terminal color theme management with full
// Catppuccin palette support.
//
// It defines the complete Catppuccin color palettes (Mocha, Macchiato,
// Frappe, Latte) and exports theme configurations for multiple terminal
// emulators and CLI tools.
//
// # Palettes
//
// [Catppuccin] returns the [Palette] for a given [Flavor]. Each palette
// contains 26 named colors (rosewater through crust) with hex values
// and RGB components. [Color.ANSI] and [Color.ANSIBg] produce 24-bit
// escape sequences for direct terminal use.
//
// # Export Targets
//
// Theme export functions write config files for specific tools:
//
//   - [ExportGhostty]: writes a theme file to ~/.config/ghostty/themes/
//   - [ExportITerm2]: writes an iTerm2 Dynamic Profile with ANSI mapping
//   - [ExportWezTerm]: writes a Lua color table for WezTerm
//   - [ExportBat]: updates ~/.config/bat/config with the theme name
//   - [ExportStarship]: writes palette colors to ~/.config/starship.toml
package themekit
