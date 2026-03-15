#!/usr/bin/env bash
set -euo pipefail

REPO="claudekit"
echo "=== $REPO ==="
echo ""

# Resolve repo root (script lives in scripts/).
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
PARENT_DIR="$(cd "$REPO_ROOT/.." && pwd)"

# 1. Verify mcpkit sibling exists (required by go.mod replace directive).
if [ ! -d "$PARENT_DIR/mcpkit" ]; then
    echo "mcpkit not found at $PARENT_DIR/mcpkit"
    echo "Cloning mcpkit..."
    git clone https://github.com/hairglasses-studio/mcpkit.git "$PARENT_DIR/mcpkit"
else
    echo "mcpkit found at $PARENT_DIR/mcpkit"
fi

cd "$REPO_ROOT"

# 2. Download Go dependencies.
echo ""
echo "Downloading dependencies..."
go mod download

# 3. Vet.
echo ""
echo "Running go vet..."
go vet ./...

# 4. Test.
echo ""
echo "Running tests..."
go test ./...

# 5. Build.
echo ""
echo "Building..."
go build ./...

# 6. Print budget profile info.
echo ""
echo "Budget profiles available:"
echo "  CLAUDEKIT_BUDGET_PROFILE=personal  (default — conservative, for Claude Max)"
echo "  CLAUDEKIT_BUDGET_PROFILE=work-api  (higher limits, for API credits)"
echo "  --budget=<path>                    (custom profile JSON)"

# 7. Print run command.
echo ""
echo "To start the MCP server:"
echo "  go run ./cmd/claudekit-mcp"
echo ""

echo "=== $REPO ==="
