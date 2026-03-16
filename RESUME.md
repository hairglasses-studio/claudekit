# Resume Instructions — Perpetual R&D Loop

## What was done

Implemented perpetual task synthesis for Ralph loops across 8 commits:

**mcpkit** (`feature/rdcycle-auto-2` branch):
1. Circuit breaker — 3-state halt control for no-progress cycles
2. FileArtifactStore — persistent JSON artifact storage
3. Real GitHub scanner — `gh` CLI integration with fallback
4. Task synthesizer — PlanOutput → ralph.Spec conversion
5. Cost velocity governor — rolling-window budget monitoring
6. Orchestrator — perpetual scan→plan→synthesize→ralph→notes→improve loop
7. Perpetual MCP tools — `rdcycle_perpetual_start/stop/status`

**claudekit** (`main` branch):
8. Wiring — FileArtifactStore, RalphStarter closure, budget profile

## Setup on new machine

```bash
# 1. Clone both repos
git clone git@github.com:hairglasses-studio/claudekit.git
git clone git@github.com:hairglasses-studio/mcpkit.git

# 2. Switch mcpkit to the feature branch
cd mcpkit
git checkout feature/rdcycle-auto-2

# 3. Retrieve secrets from 1Password
#    Account: my.1password.com | Vault: Personal
op item get "Anthropic API Key (Work - 10K credits)" --account my.1password.com --vault Personal --fields credential --reveal
op item get "AFTRS MCP - Claude Max Personal API Key" --account my.1password.com --vault Personal --fields credential --reveal
op item get "AFTRS MCP - mcpkit GitHub PAT" --account my.1password.com --vault Personal --fields credential --reveal

# 4. Create .env in claudekit root
cd ../claudekit
cat > .env << 'ENVEOF'
ANTHROPIC_API_KEY=<work key from step 3>
PERSONAL_CLAUDE_MAX_ANTHROPIC_API_KEY=<personal key from step 3>
MCPKIT_TOKEN=<github pat from step 3>
ENVEOF

# 5. Verify both repos build
cd ../mcpkit && make check
cd ../claudekit && make check

# 6. Verify gh is authenticated
gh auth status
```

## Run the 12-hour perpetual loop

```bash
cd claudekit
./scripts/perpetual-loop.sh
```

Or with custom args:
```bash
./scripts/perpetual-loop.sh rdcycle/profiles/overnight-100.json 8.0 4 0
#                           ^profile                           ^$/cyc ^cadence ^max(0=∞)
```

## Budget math

- Profile: `overnight-100.json` — $8/cycle, $100/day cap
- Expected: ~12 cycles × $8 = $96
- Safety: circuit breaker (3 no-progress), governor (rolling avg), daily cap

## 1Password items for this project

| Item | .env Key |
|------|----------|
| Anthropic API Key (Work - 10K credits) | ANTHROPIC_API_KEY |
| AFTRS MCP - Claude Max Personal API Key | PERSONAL_CLAUDE_MAX_ANTHROPIC_API_KEY |
| AFTRS MCP - mcpkit GitHub PAT | MCPKIT_TOKEN |

## WSL/Linux notes

- **No Homebrew**: fontkit install tests auto-skip via `t.Skip`
- **1Password CLI**: Use `op.exe` (Windows bridge), not `op` (Linux). Field name is `credential`, not `password`.
  ```bash
  op.exe item get "Anthropic API Key (Work - 10K credits)" --account my.1password.com --vault Personal --fields credential --reveal
  ```
- **Nesting guard**: `perpetual-loop.sh` will refuse to run inside a Claude Code session (`CLAUDECODE=1`). Always launch from a raw terminal.

## Key files

| File | Purpose |
|------|---------|
| `scripts/perpetual-loop.sh` | Launcher script with configurable args |
| `rdcycle/profiles/overnight-100.json` | Budget profile for 12hr/$100 run |
| `mcpserver/rdcyclemodule.go` | Wires FileArtifactStore + returns module |
| `cmd/claudekit-mcp/main.go` | RalphStarter closure + profile loading |
| `rdcycle/notes/improvement_log.json` | Cross-cycle learning log |
| `roadmap.json` | Machine-readable roadmap (tier-1 through tier-8) |
