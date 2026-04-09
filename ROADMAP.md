# Claudekit Roadmap

## Tier 1: Font Management + Statusline
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
- [x] Dotfile management (snapshot/restore)
- [x] Tool version management (mise integration)
- [x] Shell plugin detection (oh-my-zsh, zinit)

## Tier 4: Advanced MCP Integration
- [x] Gateway pattern — aggregate multiple claudekit MCP modules
- [x] Ralph loop — autonomous terminal setup verification
- [x] WebMCP bridge — HTTP server with /tools and /health endpoints
- [x] MCP Registry publishing — discoverable via official registry

## Tier 5: Ecosystem
- [x] Plugin system for community-contributed terminal configs (YAML-driven, subprocess handler)
- [x] Claude Code skill marketplace integration (skillkit package, CLI + MCP + WebMCP)
- [x] Multi-terminal sync (one config → iTerm2 + Ghostty + WezTerm)
- [x] CI/CD for dotfile validation (GitHub Actions)

## Tier 6: Automated R&D
- [x] Machine-readable roadmap with MCP tools (roadmap_read/update/gaps/next_phase)
- [x] FinOps token tracking middleware + budget enforcement (finops_status/reset)
- [x] Agent memory (get/set/list/search) for cross-step persistence
- [x] R&D cycle tools (scan/plan/verify/artifacts/commit/report/schedule/notes/improve)
- [x] Workflow engine with predefined setup graphs (workflow_run/list)

## Tier 7: Hardening & Autonomous Ops
- [x] Wire MCP sampling client into ralph for fully autonomous loops
- [x] CostPolicy breach triggers ralph_stop automatically
- [x] rdcycle_improve feeds suggestions back into next cycle spec
- [x] CLI command to tail ralph progress in a parallel terminal (ralph tail/status)

## Tier 8: Cross-Platform & Loop Reliability
- [x] WSL/Linux compatibility (skip brew-dependent tests, op.exe bridge docs)
- [ ] Fix verify task stuck-loop: mark_done immediately when make check passes with no changes
- [ ] Orchestrator maintenance mode: skip plan/synthesize when all phases complete
- [ ] Deduplicate orchestrator notes (don't write hollow "Completed cycle N")
- [x] mcpserver test coverage to 51%+ with handler-level tests
- [x] cmd/claudekit test coverage to 20%+ with routing and helper tests
- [x] CLAUDECODE nesting guard in perpetual-loop.sh
- [x] Wire CostReader for real per-cycle dollar tracking in governor
- [x] Consecutive sampler failure limit (5) for fast-fail on API outages
- [x] Widen scan repo list (anthropics/anthropic-sdk-go, mark3labs/mcp-go)

<!-- whiteclaw-rollout:start -->
## Whiteclaw-Derived Overhaul (2026-04-08)

This tranche applies the highest-value whiteclaw findings that fit this repo's real surface: engineer briefs, bounded skills/runbooks, searchable provenance, scoped MCP packaging, and explicit verification ladders.

### Strategic Focus
- This repo is publicly tagged but functionally deprecated into `mcpkit`, so the roadmap should emphasize surface honesty and migration clarity.
- Only keep whiteclaw patterns that reduce confusion for operators: accurate manifests, a minimal engineer brief, and clear verification/migration guidance.
- Avoid broad new autonomy work here unless the repo regains a real standalone MCP surface.

### Recommended Work
- [ ] [MCP contract] Either restore a minimal discovery-only self server or stop advertising live MCP parity in docs, manifests, and instructions.
- [ ] [Deprecation] Make the deprecation and migration path to `mcpkit` explicit in README, roadmap, and any repo-local manifests.
- [ ] [Verification] Keep a thin smoke path proving the remaining CLI behavior still works and matches the published guidance.
- [ ] [Public docs] Trim instructions and skills to the small surface that is still intentionally maintained.

### Rationale Snapshot
- Tier / lifecycle: `tier-2` / `active`
- Language profile: `Go`
- Visibility / sensitivity: `PUBLIC` / `public`
- Surface baseline: AGENTS=yes, skills=yes, codex=yes, mcp_manifest=empty, ralph=yes, roadmap=yes
- Whiteclaw transfers in scope: surface honesty, migration to mcpkit, thin verification, minimal engineer brief
- Live repo notes: AGENTS, skills, Codex config, empty .mcp.json, .ralph

<!-- whiteclaw-rollout:end -->
