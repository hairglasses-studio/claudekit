// Command claudekit-mcp is the MCP server entrypoint for Claude Code integration.
package main

import (
	"log"

	"github.com/hairglasses-studio/claudekit/mcpserver"
	"github.com/hairglasses-studio/mcpkit/registry"
)

func main() {
	reg := registry.NewToolRegistry()
	reg.RegisterModule(&mcpserver.FontModule{})
	reg.RegisterModule(&mcpserver.ThemeModule{})
	reg.RegisterModule(&mcpserver.StatuslineModule{})

	s := registry.NewMCPServer("claudekit", "0.1.0")
	reg.RegisterWithServer(s)

	if err := registry.ServeStdio(s); err != nil {
		log.Fatal(err)
	}
}
