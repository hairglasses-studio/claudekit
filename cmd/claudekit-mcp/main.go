// Command claudekit-mcp is the MCP server entrypoint for Claude Code integration.
package main

import (
	"log"
	"os"
	"strings"

	"github.com/hairglasses-studio/claudekit/mcpserver"
	"github.com/hairglasses-studio/mcpkit/registry"
)

func main() {
	reg := registry.NewToolRegistry()
	reg.RegisterModule(&mcpserver.FontModule{})
	reg.RegisterModule(&mcpserver.ThemeModule{})
	reg.RegisterModule(&mcpserver.StatuslineModule{})
	reg.RegisterModule(&mcpserver.EnvModule{})

	// Register ralph autonomous loop module.
	mcpserver.SetupRalph(reg, nil)

	s := registry.NewMCPServer("claudekit", "0.1.0")
	reg.RegisterWithServer(s)

	// Optional gateway for aggregating external MCP servers.
	var gatewayUpstreams []string
	for i, arg := range os.Args[1:] {
		if arg == "--gateway" && i+1 < len(os.Args[1:]) {
			gatewayUpstreams = strings.Split(os.Args[i+2], ",")
			break
		}
		if strings.HasPrefix(arg, "--gateway=") {
			gatewayUpstreams = strings.Split(strings.TrimPrefix(arg, "--gateway="), ",")
			break
		}
	}

	if len(gatewayUpstreams) > 0 {
		gw, dynReg, err := mcpserver.SetupGateway(reg, gatewayUpstreams)
		if err != nil {
			log.Fatalf("gateway setup: %v", err)
		}
		defer gw.Close()
		dynReg.RegisterWithServer(s)
	}

	if err := registry.ServeStdio(s); err != nil {
		log.Fatal(err)
	}
}
