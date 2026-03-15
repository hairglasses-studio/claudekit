# Changelog

## v0.1.0 — 2026-03-15

Initial release: 7 tiers complete, 34 MCP tools, 126 tests.

### Tier 1: Terminal Fonts
- `font_status`, `font_install`, `font_configure`
- MonaspiceNe Nerd Font detection and installation
- Terminal-specific configuration (iTerm2, Ghostty, WezTerm)

### Tier 2: Theme Engine
- `theme_apply`, `theme_list`
- Catppuccin palette export to iTerm2, Ghostty, WezTerm, bat, delta, Starship

### Tier 3: Statusline
- `statusline_install`
- Claude Code statusline script with Catppuccin colors

### Tier 4: Environment Management
- `env_status`, `env_snapshot`
- Mise integration, shell detection, dotfile snapshot/restore

### Tier 5: Skill Marketplace
- `skill_list`, `skill_install`
- Claude Code skill discovery and installation

### Tier 6: Automated R&D
- `rdcycle_scan`, `rdcycle_plan`, `rdcycle_verify`, `rdcycle_artifacts`
- `rdcycle_commit`, `rdcycle_report`, `rdcycle_schedule`, `rdcycle_notes`, `rdcycle_improve`
- `roadmap_read`, `roadmap_update`, `roadmap_gaps`, `roadmap_next_phase`
- `workflow_run`, `workflow_list`

### Tier 7: Autonomous Loop (Ralph)
- `ralph_start`, `ralph_stop`, `ralph_status`
- `finops_status`, `finops_reset`
- `memory_get`, `memory_set`, `memory_list`, `memory_search`
- Budget profiles (personal/$5, work-api/$50)
- Model tier selection (Opus for plan/implement, Sonnet for scan/verify/reflect)
- Cost tracking with auto-stop on budget breach
