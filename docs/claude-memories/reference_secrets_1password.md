---
name: Secrets and 1Password management
description: How to manage .env secrets, GitHub repo secrets, and 1Password vault for AFTRS MCP project
type: reference
---

## 1Password Accounts
- Personal: `my.1password.com` (mixellburk@gmail.com)
- Work: `team-galileo.1password.com` (mitch@galileo.ai)

## AFTRS Secrets Location
- Vault: `Personal` on `my.1password.com`
- Naming: "AFTRS MCP - <Service Name>" with tags "api,<service>,aftrs"
- Category: "API Credential" for API keys, "Login" for user/pass pairs

## Existing Secrets

| 1Password Item | .env Key(s) |
|---|---|
| Anthropic API Key (Work - 10K credits) | ANTHROPIC_API_KEY, PERSONAL_CLAUDE_MAX_ANTHROPIC_API_KEY |
| AFTRS MCP - mcpkit GitHub PAT | MCPKIT_TOKEN |
| AFTRS MCP - Notion API Key | NOTION_API_KEY |
| AFTRS MCP - AWS cr8 Profile | AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY |
| AFTRS MCP - Discord Bot | DISCORD_BOT_TOKEN, DISCORD_CHANNEL_ID, DISCORD_GUILD_ID |
| AFTRS MCP - Google Drive API | GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, GOOGLE_REFRESH_TOKEN |

## Rules
- Never commit .env — must be in .gitignore
- Always use `--account my.1password.com --vault Personal` for AFTRS secrets
- ANTHROPIC_API_KEY is a work key with $10K credits — API-billed work only
- Config flags (DISABLE_METRICS, AWS_REGION, etc.) go in .env only, not 1Password

## Commands
- Retrieve: `op item get "<title>" --account my.1password.com --vault Personal --fields password --reveal`
- Create: `op item create --category "API Credential" --title "AFTRS MCP - <Name>" --vault Personal --account my.1password.com "credential=<value>" --tags "api,<service>,aftrs"`
- GitHub: `gh secret set <KEY> --body "<value>"`
