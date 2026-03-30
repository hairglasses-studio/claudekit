---
name: 12hr perpetual loop runbook
description: Step-by-step instructions for running a $100 budget 12-hour perpetual R&D cycle
type: project
---

## 12-Hour Perpetual R&D Loop — $100 Budget Runbook

### Pre-flight

1. **Ensure .env has the API key** (work-api key with $10K credits):
   ```
   ANTHROPIC_API_KEY=<from 1Password: "Anthropic API Key (Work - 10K credits)">
   ```

2. **Verify roadmap.json exists** in the project root with active work items.

3. **Build both repos**:
   ```bash
   cd ~/hairglasses-studio/mcpkit && make check
   cd ~/hairglasses-studio/claudekit && make check
   ```

### Budget Profile

A custom profile `overnight-100.json` is saved at `claudekit/rdcycle/profiles/overnight-100.json`:
- **$8/cycle** (allows ~12 cycles within $100)
- **$100/day cap** (hard stop)
- **200 iterations/cycle** (generous per-cycle limit)
- **10M token budget** per cycle
- Uses opus for plan/implement, sonnet for scan/verify/reflect/report/schedule

### Governor Settings

The perpetual orchestrator creates:
- **Circuit breaker**: threshold=3 (opens after 3 no-progress cycles), cooldown=30min
- **Cost velocity governor**: window=5, alarm=$8/cycle avg, unproductive cap=3

These are hardcoded defaults in `perpetual.go`. For the $100 run, pass `alarm_per_cycle=8.0` to `rdcycle_perpetual_start`.

### Launch Sequence

```bash
# Terminal 1: Start the MCP server with the overnight profile
cd ~/hairglasses-studio/claudekit
CLAUDEKIT_BUDGET_PROFILE=rdcycle/profiles/overnight-100.json ./bin/claudekit-mcp
```

Or via Claude Code (preferred — it handles sampler wiring):
```
# In Claude Code, call:
rdcycle_perpetual_start {"max_cycles": 0, "alarm_per_cycle": 8.0, "improve_cadence": 4}
```

Parameters:
- `max_cycles: 0` — unlimited (governor + breaker provide safety)
- `alarm_per_cycle: 8.0` — matches the $8/cycle budget profile
- `improve_cadence: 4` — meta-improvement every 4 cycles (3x in a 12-cycle run)

### Monitoring

```
# Check status at any time:
rdcycle_perpetual_status
```

Returns: `running`, `cycle_num`, `breaker_state`, `total_cost`.

**Key checkpoints:**
- After cycle 1: verify scan pulled real GitHub data and spec was synthesized
- After cycle 4: first meta-improvement should fire
- If `breaker_state` = "open": 3 consecutive no-progress cycles, waiting 30min cooldown
- If `total_cost` > $80: approaching budget limit, consider stopping

### Stop

```
rdcycle_perpetual_stop
```

Finishes the current cycle gracefully, then halts.

### Post-Run

1. **Check artifacts**: `rdcycle/artifacts/*.json` has all scan/plan/verify artifacts
2. **Check specs**: `rdcycle/specs/cycle-*.json` has generated specs
3. **Check notes**: `rdcycle/notes/improvement_log.json` has all retrospectives
4. **Review git log**: each cycle should have committed changes via `rdcycle_commit`
5. **Run `rdcycle_improve`** manually for a final analysis of the full run

### Failure Modes

| Symptom | Cause | Fix |
|---------|-------|-----|
| Breaker opens immediately | No roadmap items ready | Add planned items to roadmap.json |
| Governor halts early | Cost too high per cycle | Reduce model usage: more sonnet, less opus |
| Scan returns 0 commits/issues | gh auth expired | Run `gh auth login` |
| RalphStarter fails | Missing sampler | Check ANTHROPIC_API_KEY is set |
| Stuck at same cycle | Verify fails repeatedly | Check `make check` passes manually first |

### Safety Nets

The system has 4 layers of protection:
1. **Per-cycle dollar budget** ($8) — CostPolicy stops individual ralph loop
2. **Circuit breaker** — opens after 3 no-progress cycles
3. **Cost velocity governor** — halts if rolling avg exceeds $8/cycle
4. **Daily dollar cap** ($100) — WindowedTracker hard stop

**Why:** Prevents runaway API spend. A single stuck cycle burns at most $8. The governor catches patterns where cycles complete but waste money.

**How to apply:** Before any perpetual run, verify the profile's `dollar_budget * expected_cycles ≤ daily_dollar_cap`. For this run: $8 × 12 = $96 ≤ $100.
