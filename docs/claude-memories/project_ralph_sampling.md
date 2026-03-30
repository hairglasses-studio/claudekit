---
name: Ralph sampling limitation
description: Claude Code doesn't support MCP sampling — Ralph needs an API-based sampler for autonomous loops
type: project
---

Claude Code does not implement MCP sampling/createMessage as a client capability. The ServerSamplingClient wiring in claudekit-mcp is architecturally correct but non-functional because Claude Code won't respond to sampling requests.

**Why:** MCP sampling is an optional client capability. Claude Code's MCP client implementation only supports tools, resources, and prompts — not sampling. This means ralph_start always fails with "sampler is required" when the sampler tries to call RequestSampling.

**How to apply:** Ralph autonomous loops cannot run via MCP server + Claude Code. Need one of:
1. APISamplingClient — call Anthropic API directly with ANTHROPIC_API_KEY
2. Claude Agent SDK integration as alternative loop runtime
3. Wait for Claude Code to add sampling support
Running rdcycle tools manually (scan → plan → verify → notes → report → schedule) works fine as a workaround.
