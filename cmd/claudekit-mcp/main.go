// Command claudekit-mcp is the MCP server entrypoint for Claude Code integration.
package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hairglasses-studio/claudekit/mcpserver"
	"github.com/hairglasses-studio/claudekit/pluginkit"
	"github.com/hairglasses-studio/mcpkit/ralph"
	"github.com/hairglasses-studio/mcpkit/rdcycle"
	"github.com/hairglasses-studio/mcpkit/registry"
	"github.com/hairglasses-studio/mcpkit/sampling"
)

func main() {
	loadDotenv(".env")

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

	// Register file tools for ralph autonomous loops.
	reg.RegisterModule(&ralph.FileToolModule{Root: wd})

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

	// Persistent artifact store for rdcycle.
	artDir := filepath.Join(wd, "rdcycle", "artifacts")
	fileStore, err := rdcycle.NewFileArtifactStore(artDir)
	if err != nil {
		log.Printf("warning: file artifact store: %v (using in-memory)", err)
		fileStore = nil
	}
	rdcycleMod := mcpserver.SetupRDCycle(reg, wd, fileStore)

	// Wire cost reader so the perpetual orchestrator can track per-cycle cost.
	if costPolicy != nil {
		rdcycleMod.SetCostReader(costPolicy.TotalCost)
	}

	// Wire RalphStarter so perpetual orchestrator can launch ralph loops.
	rdcycleMod.SetRalphStarter(func(ctx context.Context, specPath string) error {
		maxIter := profile.MaxIterations
		config := ralph.Config{
			SpecFile:      specPath,
			ProjectRoot:   wd,
			MaxIterations: maxIter,
			MaxTokens:     profile.MaxTokensPerReq,
			ToolRegistry:  reg,
			Sampler:       ralphMod.CurrentSampler(),
			CostTracker:   tracker,
		}
		if modelTiers.Default != "" {
			config.ModelSelector = modelTiers.Selector()
		}
		loop, err := ralph.NewLoop(config)
		if err != nil {
			return err
		}
		return loop.Run(ctx)
	})

	reg.RegisterModule(mcpserver.NewWorkflowModule())
	mcpserver.SetupMemory(reg)

	// Load plugins.
	if pluginDir, err := pluginkit.DefaultPluginDir(); err == nil {
		plugins, err := pluginkit.LoadPlugins(pluginDir)
		if err != nil {
			log.Printf("warning: failed to load plugins from %s: %v", pluginDir, err)
		}
		for _, cfg := range plugins {
			reg.RegisterModule(pluginkit.NewPluginModule(cfg))
		}
	}

	// Wire ralph's sampler. Select API key based on budget profile:
	//   personal  → PERSONAL_CLAUDE_MAX_ANTHROPIC_API_KEY (Claude Max billing)
	//   work-api  → ANTHROPIC_API_KEY (API console billing, $10K credits)
	// Falls back to ANTHROPIC_API_KEY if profile-specific key is not set.
	apiKey := resolveAPIKey(profile.Name)

	s := registry.NewMCPServer("claudekit", "0.1.0")
	reg.RegisterWithServer(s)

	if apiKey != "" {
		log.Printf("[ralph] using API sampler (profile=%s)", profile.Name)
		ralphMod.SetSampler(&sampling.APISamplingClient{
			APIKey:       apiKey,
			DefaultModel: "claude-sonnet-4-6",
		})
	} else {
		// Fallback to MCP server sampling (works if client supports it).
		ralphMod.SetSampler(&sampling.ServerSamplingClient{Server: s})
	}

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

// loadDotenv reads KEY=VALUE pairs from a file and sets them as env vars.
// Silently skips if the file doesn't exist. Does not override existing env vars.
func loadDotenv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.Trim(strings.TrimSpace(v), "'\"")
		if os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}
}

// resolveAPIKey selects the Anthropic API key based on budget profile name.
// Personal profile uses the Claude Max key; work-api uses the console key.
func resolveAPIKey(profileName string) string {
	switch profileName {
	case "personal":
		if key := os.Getenv("PERSONAL_CLAUDE_MAX_ANTHROPIC_API_KEY"); key != "" {
			return key
		}
	case "work-api":
		if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
			return key
		}
	}
	// Fallback: try ANTHROPIC_API_KEY for any profile.
	return os.Getenv("ANTHROPIC_API_KEY")
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
