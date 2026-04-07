---
name: reconnect
description: Reconnect all MCP servers and verify their status. Use when MCP tools stop responding, servers are disconnected, or after rebuilding an MCP server binary.
---

Reconnect all MCP servers. For each configured MCP server:
1. Check if it's currently connected
2. If disconnected, reconnect it
3. Report the final status of all servers

Use `/mcp` to access the MCP management interface and reconnect any disconnected servers.
