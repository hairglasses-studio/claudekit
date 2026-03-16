#!/usr/bin/env bash
set -euo pipefail

# Perpetual R&D loop launcher.
# Invokes Claude Code with a prompt that starts the perpetual orchestrator,
# monitors it, and lets it run autonomously.

PROFILE="${1:-rdcycle/profiles/overnight-100.json}"
ALARM="${2:-8.0}"
CADENCE="${3:-4}"
MAX_CYCLES="${4:-0}"

PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$PROJECT_DIR"

# Resolve absolute path for profile if relative.
if [[ ! "$PROFILE" = /* ]]; then
    PROFILE="$PROJECT_DIR/$PROFILE"
fi

if [[ ! -f "$PROFILE" ]]; then
    echo "error: profile not found: $PROFILE" >&2
    exit 1
fi

# Validate API key.
if [[ -z "${ANTHROPIC_API_KEY:-}" ]]; then
    ENV_FILE="$PROJECT_DIR/.env"
    if [[ -f "$ENV_FILE" ]]; then
        eval "$(grep -E '^ANTHROPIC_API_KEY=' "$ENV_FILE" | head -1)"
        export ANTHROPIC_API_KEY
    fi
fi

if [[ -z "${ANTHROPIC_API_KEY:-}" ]]; then
    echo "error: ANTHROPIC_API_KEY not set and not found in .env" >&2
    exit 1
fi

# Build if needed.
if [[ ! -x bin/claudekit-mcp ]]; then
    echo "Building claudekit-mcp..."
    make build
fi

echo "=== Perpetual R&D Loop ==="
echo "  profile:       $PROFILE"
echo "  alarm/cycle:   \$$ALARM"
echo "  improve every: $CADENCE cycles"
echo "  max cycles:    $([ "$MAX_CYCLES" = "0" ] && echo "unlimited" || echo "$MAX_CYCLES")"
echo "  project:       $PROJECT_DIR"
echo ""
echo "Launching Claude Code with perpetual loop prompt..."
echo "To monitor: open another Claude Code session and call rdcycle_perpetual_status"
echo "To stop:    call rdcycle_perpetual_stop from any Claude Code session"
echo ""

export CLAUDEKIT_BUDGET_PROFILE="$PROFILE"

PROMPT="Start the perpetual R&D cycle orchestrator with these settings:
- max_cycles: $MAX_CYCLES
- alarm_per_cycle: $ALARM
- improve_cadence: $CADENCE

Call rdcycle_perpetual_start with those parameters. Then monitor with rdcycle_perpetual_status every 2 minutes until it completes or I interrupt. Report each status check concisely: cycle number, breaker state, total cost."

exec claude --dangerously-skip-permissions -p "$PROMPT"
