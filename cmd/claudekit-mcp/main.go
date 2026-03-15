// Command claudekit-mcp is the MCP server entrypoint for Claude Code integration.
package main

import (
	"log"
	"os"
	"strings"

	"github.com/hairglasses-studio/claudekit/mcpserver"
	"github.com/hairglasses-studio/claudekit/pluginkit"
	"github.com/hairglasses-studio/mcpkit/rdcycle"
	"github.com/hairglasses-studio/mcpkit/registry"
	"github.com/hairglasses-studio/mcpkit/sampling"
)

func main() {
	reg := registry.NewToolRegistry()
	wd, _ := os.Getwd()

	// Resolve budget profile from env or CLI flag.
	profile := resolveProfile()

	// FinOps middleware must be registered first so it wraps all subsequent tool handlers.
	tracker, costPolicy, _ := mcpserver.SetupFinOpsFromProfile(reg, profile)

	reg.RegisterModule(&mcpserver.FontModule{})
	reg.RegisterModule(&mcpserver.ThemeModule{})
	reg.RegisterModule(&mcpserver.StatuslineModule{})
	reg.RegisterModule(&mcpserver.EnvModule{})

	// Register skill marketplace module.
	reg.RegisterModule(&mcpserver.SkillModule{ProjectDir: wd})

	// Model tier config: opus for plan/implement, sonnet for everything else.
	modelTiers := rdcycle.ModelTierConfig{
		Default: "claude-opus-4-6",
		TaskOverrides: map[string]string{
			"scan":     "claude-sonnet-4-6",
			"verify":   "claude-sonnet-4-6",
			"reflect":  "claude-sonnet-4-6",
			"report":   "claude-sonnet-4-6",
			"schedule": "claude-sonnet-4-6",
		},
	}

	// Register ralph autonomous loop module with budget-aware config.
	// Sampler is wired after MCP server creation (see below).
	ralphMod := mcpserver.SetupRalph(reg, nil, mcpserver.RalphConfig{
		Profile:     profile,
		ModelTiers:  modelTiers,
		CostTracker: tracker,
		CostPolicy:  costPolicy,
	})

	// Register Tier 6: Automated R&D modules.
	mcpserver.SetupRoadmap(reg, wd)
	mcpserver.SetupRDCycle(reg, wd)
	reg.RegisterModule(mcpserver.NewWorkflowModule())
	mcpserver.SetupMemory(reg)

	// Load plugins.
	if pluginDir, err := pluginkit.DefaultPluginDir(); err == nil {
		if plugins, err := pluginkit.LoadPlugins(pluginDir); err == nil {
			for _, cfg := range plugins {
				reg.RegisterModule(pluginkit.NewPluginModule(cfg))
			}
		}
	}

	s := registry.NewMCPServer("claudekit", "0.1.0")
	reg.RegisterWithServer(s)

	// Wire the MCP server as ralph's sampling client so loops can run autonomously.
	ralphMod.SetSampler(&sampling.ServerSamplingClient{Server: s})

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

// resolveProfile returns the budget profile based on CLAUDEKIT_BUDGET_PROFILE env var
// or --budget CLI flag. Defaults to personal profile.
func resolveProfile() rdcycle.BudgetProfile {
	name := os.Getenv("CLAUDEKIT_BUDGET_PROFILE")

	// CLI flag overrides env var.
	for i, arg := range os.Args[1:] {
		if arg == "--budget" && i+1 < len(os.Args[1:]) {
			name = os.Args[i+2]
			break
		}
		if strings.HasPrefix(arg, "--budget=") {
			name = strings.TrimPrefix(arg, "--budget=")
			break
		}
	}

	switch name {
	case "work-api":
		return rdcycle.WorkAPIProfile()
	case "", "personal":
		return rdcycle.PersonalProfile()
	default:
		// Try loading from file path.
		p, err := rdcycle.LoadProfile(name)
		if err != nil {
			log.Printf("failed to load budget profile %q, using personal: %v", name, err)
			return rdcycle.PersonalProfile()
		}
		return p
	}
}
