# Contributing

This project supports development with **Claude Code**, **Gemini CLI**, and **OpenAI Codex CLI**. Any provider can lead development.

## Development Setup

### 1. Clone and build

```bash
git clone https://github.com/hairglasses-studio/claudekit
cd claudekit
make build   # or: go build ./...
make test    # or: go test ./... -count=1
```

### 2. Verify

```bash
make check   # build + vet + test
```

## Making Changes

1. Create a branch: `git checkout -b feat/my-change`
2. Make your changes
3. Run checks: `make check`
4. Commit with a descriptive message
5. Push and open a PR

## Code Style

- **Go**: `gofmt` formatting, `go vet` clean, golangci-lint passing
- **Node.js**: ESLint/Prettier where configured
- **Python**: ruff/black formatting

Editor settings are in `.editorconfig` — most editors pick this up automatically.

## Pre-commit Hooks

Install with:

```bash
make install-hooks
```

This runs vet + fast tests before each commit.

## CI

All PRs trigger CI automatically. The pipeline runs lint, test, and build checks.

## Questions?

Open an issue or tag `@hairglasses` in your PR.
