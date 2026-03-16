#!/usr/bin/env bash
set -euo pipefail

# Perpetual R&D loop launcher.
# Invokes Claude Code with a prompt that starts the perpetual orchestrator,
# monitors it, and lets it run autonomously.

PROFILE="${1:-rdcycle/profiles/overnight-100.json}"
ALARM="${2:-8.0}"
CADENCE="${3:-4}"
MAX_CYCLES="${4:-0}"

# Guard: cannot launch inside another Claude Code session.
if [[ -n "${CLAUDECODE:-}" ]]; then
    echo "error: cannot launch inside a Claude Code session (CLAUDECODE is set)." >&2
    echo "Run this script from a regular terminal." >&2
    exit 1
fi

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

# Load all env vars from .env (does not override existing).
ENV_FILE="$PROJECT_DIR/.env"
if [[ -f "$ENV_FILE" ]]; then
    while IFS= read -r line || [[ -n "$line" ]]; do
        line="${line%%#*}"          # strip comments
        line="${line#"${line%%[![:space:]]*}"}"  # trim leading whitespace
        line="${line%"${line##*[![:space:]]}"}"  # trim trailing whitespace
        [[ -z "$line" ]] && continue
        key="${line%%=*}"
        val="${line#*=}"
        val="${val#[\"\']}"        # strip leading quote
        val="${val%[\"\']}"        # strip trailing quote
        if [[ -z "${!key:-}" ]]; then
            export "$key=$val"
        fi
    done < "$ENV_FILE"
fi

# Validate API key.
if [[ -z "${ANTHROPIC_API_KEY:-}" ]]; then
    echo "error: ANTHROPIC_API_KEY not set and not found in .env" >&2
    exit 1
fi

# Verify build compiles (the MCP server runs via 'go run' so no binary is needed).
echo "Verifying build..."
go build ./... || { echo "error: build failed" >&2; exit 1; }

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
