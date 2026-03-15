# Claudekit Roadmap

## Tier 1: Font Management + Statusline (Current)
- [x] Monaspace/Monaspice detection and installation
- [x] iTerm2 Dynamic Profile with font fallback
- [x] Ghostty config with font fallback
- [x] CLI + MCP server + Claude Code skill
- [x] Font preview command (ligatures, Nerd Font icons)
- [x] Claude Code statusline with Monaspace multi-font fallback

## Tier 2: Terminal Theming
- [x] Catppuccin theme support (one palette → multiple targets)
- [x] Starship prompt configuration
- [x] Terminal color scheme management
- [x] Syntax highlighting theme sync (bat, delta, etc.)

## Tier 3: Environment Bootstrap
- [ ] Dotfile management (chezmoi-style templates)
- [ ] Tool version management (mise integration)
- [ ] Shell plugin management (oh-my-zsh, zinit)

## Tier 4: Advanced MCP Integration
- [x] Gateway pattern — aggregate multiple claudekit MCP modules
- [x] Ralph loop — autonomous terminal setup verification
- [ ] WebMCP bridge — browser-based config previewer
- [ ] MCP Registry publishing — discoverable via official registry

## Tier 5: Ecosystem
- [ ] Plugin system for community-contributed terminal configs
- [ ] Claude Code skill marketplace integration
- [ ] Multi-terminal sync (one config → iTerm2 + Ghostty + WezTerm)
- [ ] CI/CD for dotfile validation
