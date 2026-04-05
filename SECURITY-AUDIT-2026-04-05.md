# Security Audit — claudekit — 2026-04-05

Pre-release audit for Wave 3 public release.

## Audit Scope

| Check | Result |
|-------|--------|
| Hardcoded secrets in source | PASS |
| Committed .env files | PASS |
| .env files in git history | PASS |
| Personal paths (`/home/hg`) | PASS |
| Local `replace` directives in go.mod | PASS |
| Employer/personal names in source | PASS |
| LICENSE file present | PASS |
| README.md present | PASS |
| Author emails in git history | **FAIL** |

## Findings

### FAIL: Non-hairglasses author emails in git history

**Severity**: Medium (reputation/privacy, not a credential leak)

The git log contains commits authored with emails that should not appear in a public repo:

| Email | Commit Count | Issue |
|-------|-------------|-------|
| `mitch@galileo.ai` | 63 | Employer email — leaks employer association |
| `<redacted-personal-email>` | 11 | Personal gmail — should be project email |
| `mitch@hairglasses.studio` | 11 | Correct |

**Remediation options** (pick one before public push):

1. **Rewrite history** (`git filter-repo --mailmap`). Cleanest result but rewrites all commit hashes. Acceptable since the repo is not yet public.
   ```
   # Create mailmap file:
   # mitch@hairglasses.studio <mitch@galileo.ai>
   # mitch@hairglasses.studio <<redacted-personal-email>>
   git filter-repo --mailmap mailmap.txt
   ```
2. **Accept and move forward**. The emails are not secrets per se, but `galileo.ai` ties the repo to an employer, which violates the org convention "no employer names in personal repos."

**Recommendation**: Option 1 (rewrite) before the first public push. The repo has no external forks or CI references to preserve.

### PASS Notes

- **Secrets grep**: All "token" matches are legitimate code (env var reads via `os.Getenv("CLAUDEKIT_REGISTRY_TOKEN")`, struct fields for token-usage tracking, function parameters). No hardcoded API keys, passwords, or service account credentials found.
- **.gitignore**: Properly excludes `.env`, `.env.*`, `.env.local`, `.envrc`, with explicit allowlist for `.env.example` and `.env.*.example`.
- **go.mod**: Uses published `github.com/hairglasses-studio/mcpkit v0.2.0` — no local `replace` directives.
- **Personal paths**: Zero matches for `/home/hg` across all `.go`, `.md`, and `.json` files.

## Checklist Status

- [x] MIT LICENSE file present (Copyright 2024-2026 hairglasses-studio)
- [x] .env in .gitignore and not in git history
- [x] No hardcoded secrets (sk-ant, sk-svcacct, AIzaSy, ghp_, passwords)
- [x] No personal paths (/home/hg, /Users/)
- [x] No local replace directives in go.mod
- [x] README.md exists
- [ ] **Git author emails need rewrite** (63 galileo.ai + 11 gmail commits)

## Conclusion

The claudekit repo is **clear for public release** pending one remediation: rewrite git author emails to `mitch@hairglasses.studio` using `git filter-repo --mailmap` before the first public push. No secrets, credentials, personal paths, or local dependency overrides were found.
