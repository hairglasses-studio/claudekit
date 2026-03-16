# Cross-Repo Ralph Loop Learnings (2026-03-16)

Source: hg-mcp, mesmer, ralphglasses, ralph-claude-code, mcpkit parallel research

## Critical Fixes (Apply Everywhere)

### 1. Nested Session Prevention
All repos hit `Error: Claude Code cannot be launched inside another Claude Code session`. Fix:
```bash
CLAUDECODE= claude --dangerously-skip-permissions ...
```
Must be in the invocation itself. `env -u` from parent shell is unreliable.

### 2. ANTHROPIC_API_KEY Conflict (ralphglasses discovery)
Claude Code uses OAuth now. Setting ANTHROPIC_API_KEY overrides OAuth with stale key. Fix:
```bash
unset ANTHROPIC_API_KEY  # Claude now uses OAuth, not keys
```

### 3. macOS-to-Linux Compatibility
- `sed -i '' 's/old/new/'` fails on Linux (use `sed -i 's/old/new/'`)
- `df -g` is macOS-only (use `df -BG . 2>/dev/null || df -g .`)
- `brew` not available on WSL (fontkit tests use `t.Skip`)

## Architecture Patterns Worth Adopting

### From hg-mcp: Shell-based Library Modules
Modular bash libraries in `.ralph/lib/`:
- `circuit_breaker.sh` — 3-state (CLOSED/HALF_OPEN/OPEN) with auto-recovery
- `cost_governor.sh` — rolling velocity alarm + unproductive streak
- `improvement_journal.sh` — retrospective extraction + meta-improvement cycles
- `model_selector.sh` — 4-layer selection (forced > task-pattern > budget-aware > opus-window)
- `response_analyzer.sh` — JSON+text signal extraction from Claude output

### From mesmer: Adaptive Quality Gates by Role
Skip unnecessary gates based on phase role:
- Strategist (markdown-only): skip all code gates
- Reconciler: skip full test suite
- Builder: all gates required
- ~60% time savings on non-code phases

### From mesmer: Task Batching (3x speedup)
Batch 2-3 similar tasks per loop iteration. Config:
```
BATCH_SIMILAR_TASKS=true
MAX_TASKS_PER_BATCH=3
MAX_LINES_PER_BATCH=600
```

### From ralphglasses: Marathon Supervisor
External supervisor script (366 lines) with:
- Git checkpoints every 3 hours (tags + commits with metrics)
- 90% budget headroom protection
- Graceful SIGTERM→SIGKILL shutdown
- 30-second polling loop for duration/budget/health checks

### From ralph-claude-code: Exponential Backoff
```
RETRY_BACKOFF_INITIAL=30s
RETRY_BACKOFF_MAX=300s
RETRY_BACKOFF_MULTIPLIER=2
```
Plus stale call counter detection: reset when TIMESTAMP_FILE hour is >1 hour old.

### From mcpkit: Task DAG with LLM-driven Selection
Dependency DAG where ready tasks are chosen by LLM sampler. Enables smarter prioritization than sequential task lists.

### From claudekit: APISamplingClient + SetSampler() Deferred Wiring
When MCP sampling isn't supported, call Anthropic API directly. Register modules before server boot, wire sampler after init.

## Budget & Cost Control Comparison

| Repo | Session Budget | Per-Loop | Model Strategy | Downgrade |
|------|---------------|----------|----------------|-----------|
| hg-mcp | $100/12h | $1.50 max | sonnet default, opus window (30/5h) | 60% budget |
| mesmer | $100/12h | ~$1.80 | sonnet default, role-based | 70% budget |
| ralphglasses | $100/12h | $0.15 est | sonnet, 5 calls/hr | 90% headroom |
| ralph-claude-code | varies | --max-cost | configurable | CB threshold |
| claudekit | $100/day, $8/cycle | ~$0.15-0.50 | opus default, sonnet for light tasks | CostPolicy auto-stop |

## Key Improvement: claudekit-specific

claudekit's perpetual loop uses Go-based MCP tools (not shell scripts), which means:
- Improvements live in mcpkit Go code, not .ralph/ bash scripts
- The orchestrator is typed and testable (vs. shell script fragility)
- Cost tracking is integrated via finops middleware (automatic, not manual)
- But: less flexible for quick iteration than shell scripts
- Consider: hybrid approach — shell launcher + Go MCP tools
