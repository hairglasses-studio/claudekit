package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hairglasses-studio/claudekit/envkit"
	"github.com/hairglasses-studio/claudekit/fontkit"
	"github.com/hairglasses-studio/claudekit/mcpserver"
	"github.com/hairglasses-studio/claudekit/pluginkit"
	"github.com/hairglasses-studio/claudekit/skillkit"
	"github.com/hairglasses-studio/claudekit/statusline"
	"github.com/hairglasses-studio/claudekit/themekit"
	"github.com/hairglasses-studio/mcpkit/ralph"
	"github.com/hairglasses-studio/mcpkit/registry"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	group := os.Args[1]
	cmd := ""
	if len(os.Args) > 2 {
		cmd = os.Args[2]
	}

	var err error
	switch group {
	case "fonts":
		err = runFonts(ctx, cmd)
	case "theme":
		err = runTheme(cmd)
	case "env":
		err = runEnv(ctx, cmd)
	case "mcp":
		err = runMCP(ctx, cmd)
	case "statusline":
		err = runStatusline(cmd)
	case "plugin":
		err = runPlugin(cmd)
	case "skill":
		err = runSkill(cmd)
	case "ralph":
		err = runRalph(ctx, cmd)
	case "help", "--help", "-h":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command group: %s\n", group)
		usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Println(`claudekit — Claude Code terminal customization toolkit

Usage:
  claudekit fonts status          Detect installed fonts and terminal
  claudekit fonts install         Install Monaspice via Homebrew
  claudekit fonts configure       Write terminal font configuration
  claudekit fonts configure --terminal ghostty
  claudekit fonts preview         Show ligatures, Nerd Font glyphs, fallback tiers
  claudekit fonts setup           All-in-one: detect → install → configure

  claudekit theme apply           Apply Catppuccin theme to terminal
  claudekit theme apply --flavor mocha --terminal ghostty
  claudekit theme sync            Apply theme to terminal + bat + delta
  claudekit theme preview         Preview all Catppuccin flavors

  claudekit env status            Show mise + shell + dotfile info
  claudekit env snapshot          Capture current config files
  claudekit env mise              Install/configure mise

  claudekit mcp tools             List registered MCP tools
  claudekit mcp publish           Publish to MCP Registry
  claudekit mcp serve             Start WebMCP HTTP server (default :8080)

  claudekit statusline install    Install the Claude Code statusline
  claudekit statusline preview    Preview statusline with sample data

  claudekit plugin list          List installed plugins
  claudekit plugin add <path>    Install a plugin from YAML file

  claudekit skill list           List installed and available skills
  claudekit skill install <name> Install a skill from the marketplace
  claudekit skill remove <name>  Remove an installed skill

  claudekit ralph tail <file>    Watch ralph progress file in real time
  claudekit ralph status <file>  Show current ralph progress snapshot
  claudekit ralph status <file> --json  Output progress as JSON`)
}

// hasFlag returns true if --key is present in os.Args (with no value).
func hasFlag(key string) bool {
	flag := "--" + key
	for _, arg := range os.Args {
		if arg == flag {
			return true
		}
	}
	return false
}

// parseFlag returns the value of --key from os.Args, or fallback if not found.
func parseFlag(key, fallback string) string {
	prefix := "--" + key
	for i, arg := range os.Args {
		if arg == prefix && i+1 < len(os.Args) {
			return os.Args[i+1]
		}
		if strings.HasPrefix(arg, prefix+"=") {
			return strings.TrimPrefix(arg, prefix+"=")
		}
	}
	return fallback
}

func runFonts(ctx context.Context, cmd string) error {
	switch cmd {
	case "status":
		return fontsStatus(ctx)
	case "install":
		return fontsInstall(ctx)
	case "configure":
		return fontsConfigure()
	case "preview":
		return fontsPreview()
	case "setup":
		return fontsSetup(ctx)
	default:
		return fmt.Errorf("unknown fonts command: %s (try: status, install, configure, preview, setup)", cmd)
	}
}

func fontsStatus(ctx context.Context) error {
	_, err := printFontStatus(ctx)
	return err
}

// printFontStatus detects and prints font status, returning the result for reuse.
func printFontStatus(ctx context.Context) (*fontkit.FontStatus, error) {
	term := fontkit.DetectTerminal()
	fmt.Printf("Terminal: %s", term.Name)
	if term.Version != "" {
		fmt.Printf(" %s", term.Version)
	}
	fmt.Println()
	fmt.Printf("Font config support: %v\n", term.Terminal.SupportsConfig())
	fmt.Println()

	status, err := fontkit.Detect(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Homebrew: %v\n", status.BrewAvail)
	fmt.Printf("Best available font: %s\n\n", status.BestFont)

	if len(status.Installed) > 0 {
		fmt.Println("Installed:")
		for _, f := range status.Installed {
			nerd := ""
			if f.HasNerdGlyphs {
				nerd = " (Nerd Font)"
			}
			fmt.Printf("  ✓ %s%s\n", f.Name, nerd)
		}
		fmt.Println()
	}

	if len(status.NotInstalled) > 0 {
		fmt.Println("Not installed:")
		for _, f := range status.NotInstalled {
			fmt.Printf("  ✗ %s\n", f.Name)
		}
		fmt.Println()
	}
	return status, nil
}

func fontsInstall(ctx context.Context) error {
	fmt.Println("Installing Monaspice Nerd Font via Homebrew...")
	result, err := fontkit.Install(ctx, fontkit.InstallOpts{NerdFont: true})
	if err != nil {
		return err
	}
	if result.AlreadyOK {
		fmt.Printf("Already installed: %s\n", result.Cask)
	} else {
		fmt.Printf("Installed: %s\n", result.Cask)
	}
	if result.Output != "" {
		fmt.Println(result.Output)
	}
	return nil
}

func fontsConfigure() error {
	termFlag := parseFlag("terminal", "auto")

	var terminal fontkit.Terminal
	switch termFlag {
	case "auto":
		terminal = fontkit.DetectTerminal().Terminal
	case "iterm2":
		terminal = fontkit.TerminalITerm2
	case "ghostty":
		terminal = fontkit.TerminalGhostty
	case "wezterm":
		terminal = fontkit.TerminalWezTerm
	default:
		return fmt.Errorf("unsupported terminal: %s (supported: auto, iterm2, ghostty, wezterm)", termFlag)
	}

	var path string
	var err error
	switch terminal {
	case fontkit.TerminalITerm2:
		path, err = fontkit.ConfigureITerm2(fontkit.ITerm2Opts{})
	case fontkit.TerminalGhostty:
		path, err = fontkit.ConfigureGhostty(fontkit.GhosttyOpts{})
	case fontkit.TerminalWezTerm:
		path, err = fontkit.ConfigureWezTerm(fontkit.WezTermOpts{})
	default:
		return fmt.Errorf("unsupported terminal: %s (supported: iTerm2, Ghostty, WezTerm)", terminal)
	}
	if err != nil {
		return err
	}
	fmt.Printf("Wrote font config: %s\n", path)

	if terminal == fontkit.TerminalITerm2 {
		fmt.Println("\nTo activate: iTerm2 → Preferences → Profiles → select \"Claudekit Monaspace\"")
	}
	return nil
}

func fontsPreview() error {
	fmt.Print(fontkit.Preview(fontkit.PreviewOpts{ShowAll: true}))
	return nil
}

func fontsSetup(ctx context.Context) error {
	fmt.Println("=== Font Status ===")
	status, err := printFontStatus(ctx)
	if err != nil {
		return err
	}

	hasMonaspice := false
	for _, f := range status.Installed {
		if f.HasNerdGlyphs {
			hasMonaspice = true
			break
		}
	}

	if !hasMonaspice {
		fmt.Println("=== Installing Monaspice ===")
		if err := fontsInstall(ctx); err != nil {
			return err
		}
	}

	term := fontkit.DetectTerminal()
	if term.Terminal.SupportsConfig() {
		fmt.Println("=== Configuring Terminal ===")
		if err := fontsConfigure(); err != nil {
			return err
		}
	}

	fmt.Println("\n=== Setup Complete ===")
	return nil
}

func runTheme(cmd string) error {
	switch cmd {
	case "apply":
		return themeApply()
	case "sync":
		return themeSync()
	case "preview":
		return themePreview()
	default:
		return fmt.Errorf("unknown theme command: %s (try: apply, sync, preview)", cmd)
	}
}

func themeApply() error {
	flavorStr := parseFlag("flavor", "mocha")
	termStr := parseFlag("terminal", "auto")

	var flavor themekit.Flavor
	switch flavorStr {
	case "latte":
		flavor = themekit.Latte
	case "frappe":
		flavor = themekit.Frappe
	case "macchiato":
		flavor = themekit.Macchiato
	case "mocha":
		flavor = themekit.Mocha
	default:
		return fmt.Errorf("unknown flavor: %s (try: mocha, macchiato, frappe, latte)", flavorStr)
	}

	p := themekit.Catppuccin(flavor)

	var terminal fontkit.Terminal
	switch termStr {
	case "auto":
		terminal = fontkit.DetectTerminal().Terminal
	case "iterm2":
		terminal = fontkit.TerminalITerm2
	case "ghostty":
		terminal = fontkit.TerminalGhostty
	case "wezterm":
		terminal = fontkit.TerminalWezTerm
	default:
		return fmt.Errorf("unsupported terminal: %s (supported: auto, iterm2, ghostty, wezterm)", termStr)
	}

	var path string
	var err error
	switch terminal {
	case fontkit.TerminalITerm2:
		path, err = themekit.ExportITerm2(p)
	case fontkit.TerminalGhostty:
		path, err = themekit.ExportGhostty(p)
	case fontkit.TerminalWezTerm:
		path, err = themekit.ExportWezTerm(p)
	default:
		return fmt.Errorf("unsupported terminal: %s", terminal)
	}
	if err != nil {
		return err
	}

	fmt.Printf("Applied %s theme: %s\n", p.Name, path)
	if terminal == fontkit.TerminalGhostty {
		fmt.Printf("\nAdd to your Ghostty config: theme = claudekit-%s\n", strings.ReplaceAll(strings.ToLower(p.Name), " ", "-"))
	}
	if terminal == fontkit.TerminalWezTerm {
		fmt.Println("\nAdd to your wezterm.lua: config.colors = require(\"claudekit-colors\")")
	}
	return nil
}

func themeSync() error {
	flavorStr := parseFlag("flavor", "mocha")

	var flavor themekit.Flavor
	switch flavorStr {
	case "latte":
		flavor = themekit.Latte
	case "frappe":
		flavor = themekit.Frappe
	case "macchiato":
		flavor = themekit.Macchiato
	case "mocha":
		flavor = themekit.Mocha
	default:
		return fmt.Errorf("unknown flavor: %s", flavorStr)
	}

	p := themekit.Catppuccin(flavor)
	fmt.Printf("Syncing %s across all targets...\n\n", p.Name)

	// Terminal theme
	term := fontkit.DetectTerminal()
	if term.Terminal.SupportsConfig() {
		var path string
		var err error
		switch term.Terminal {
		case fontkit.TerminalITerm2:
			path, err = themekit.ExportITerm2(p)
		case fontkit.TerminalGhostty:
			path, err = themekit.ExportGhostty(p)
		case fontkit.TerminalWezTerm:
			path, err = themekit.ExportWezTerm(p)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ Terminal: %v\n", err)
		} else {
			fmt.Printf("  ✓ Terminal (%s): %s\n", term.Name, path)
		}
	}

	// bat
	if path, err := themekit.ExportBat(p); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ bat: %v\n", err)
	} else {
		fmt.Printf("  ✓ bat: %s\n", path)
	}

	// delta
	if path, err := themekit.ExportDelta(p); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ delta: %v\n", err)
	} else {
		fmt.Printf("  ✓ delta: %s\n", path)
	}

	// starship
	if path, err := themekit.ExportStarship(p); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ starship: %v\n", err)
	} else {
		fmt.Printf("  ✓ starship: %s\n", path)
	}

	fmt.Printf("\nTheme sync complete. Set CLAUDEKIT_THEME=%s for statusline colors.\n", flavorStr)
	return nil
}

func themePreview() error {
	reset := "\033[0m"
	for _, flavor := range themekit.AllFlavors() {
		p := themekit.Catppuccin(flavor)
		fmt.Printf("=== %s ===\n", p.Name)

		// Show key colors
		names := []string{"base", "text", "red", "green", "yellow", "blue", "pink", "teal", "mauve", "peach"}
		for _, name := range names {
			col := p.Get(name)
			fmt.Printf("  %s██%s %s (#%s)\n", col.ANSI(), reset, col.Name, col.Hex)
		}
		fmt.Println()
	}
	return nil
}

func runStatusline(cmd string) error {
	switch cmd {
	case "install":
		return statuslineInstall()
	case "preview":
		return statuslinePreview()
	case "render":
		return statuslineRender()
	default:
		return fmt.Errorf("unknown statusline command: %s (try: install, preview, render)", cmd)
	}
}

func statuslineInstall() error {
	result, err := statusline.Install(statusline.InstallOpts{Style: statusline.StyleFull})
	if err != nil {
		return err
	}
	fmt.Printf("Statusline script: %s\n", result.ScriptPath)
	fmt.Printf("Settings updated: %s\n", result.SettingsPath)
	if result.BackedUp != "" {
		fmt.Printf("Previous config backed up to: %s\n", result.BackedUp)
	}
	fmt.Println("\nRestart Claude Code to see the statusline.")
	return nil
}

func statuslineRender() error {
	output := statusline.Render(os.Stdin)
	fmt.Print(output)
	return nil
}

func statuslinePreview() error {
	sampleJSON := `{
		"model": "claude-opus-4-6",
		"working_directory": "/Users/dev/hairglasses-studio/claudekit",
		"session_id": "abc123",
		"cost_usd": 0.42,
		"total_tokens": 15000,
		"max_tokens": 200000,
		"duration_ms": 120000
	}`
	output := statusline.Render(strings.NewReader(sampleJSON))
	fmt.Print(output)
	return nil
}

// --- env commands ---

func runEnv(ctx context.Context, cmd string) error {
	switch cmd {
	case "status":
		return envStatus(ctx)
	case "snapshot":
		return envSnapshot()
	case "mise":
		return envMise(ctx)
	default:
		return fmt.Errorf("unknown env command: %s (try: status, snapshot, mise)", cmd)
	}
}

func envStatus(ctx context.Context) error {
	mise, err := envkit.MiseStatus(ctx)
	if err != nil {
		return err
	}

	fmt.Println("=== Mise ===")
	if mise.Installed {
		fmt.Printf("Installed: yes (%s)\n", mise.Version)
		if len(mise.Tools) > 0 {
			fmt.Println("Active tools:")
			for name, version := range mise.Tools {
				fmt.Printf("  %s: %s\n", name, version)
			}
		}
	} else {
		fmt.Println("Installed: no")
	}
	fmt.Println()

	shell := envkit.DetectShell()
	fmt.Println("=== Shell ===")
	fmt.Println(shell.ShellSummary())
	fmt.Println()

	snap, err := envkit.Snapshot()
	if err != nil {
		return err
	}
	fmt.Println("=== Managed Configs ===")
	fmt.Print(envkit.SnapshotSummary(snap))

	return nil
}

func envSnapshot() error {
	snap, err := envkit.Snapshot()
	if err != nil {
		return err
	}

	fmt.Print(envkit.SnapshotSummary(snap))
	if len(snap.Files) > 0 {
		fmt.Println("\nFile contents captured in memory. Use MCP env_snapshot tool to access programmatically.")
	}
	return nil
}

func envMise(ctx context.Context) error {
	dryRun := hasFlag("dry-run")

	fmt.Println("Configuring mise tool version manager...")
	result, err := envkit.MiseInstall(ctx, envkit.MiseOpts{DryRun: dryRun})
	if err != nil {
		return err
	}

	if result.Installed {
		fmt.Println("Mise installed successfully.")
	}
	if result.ConfigPath != "" {
		fmt.Printf("Config written: %s\n", result.ConfigPath)
	}
	if result.Output != "" {
		fmt.Println(result.Output)
	}
	return nil
}

// --- mcp commands ---

func newToolRegistry() *registry.ToolRegistry {
	reg := registry.NewToolRegistry()
	reg.RegisterModule(&mcpserver.FontModule{})
	reg.RegisterModule(&mcpserver.ThemeModule{})
	reg.RegisterModule(&mcpserver.StatuslineModule{})
	reg.RegisterModule(&mcpserver.EnvModule{})
	reg.RegisterModule(&mcpserver.SkillModule{ProjectDir: projectDir()})

	// Load plugins
	if pluginDir, err := pluginkit.DefaultPluginDir(); err == nil {
		if plugins, err := pluginkit.LoadPlugins(pluginDir); err == nil {
			for _, cfg := range plugins {
				reg.RegisterModule(pluginkit.NewPluginModule(cfg))
			}
		}
	}

	return reg
}

func runMCP(ctx context.Context, cmd string) error {
	switch cmd {
	case "tools":
		return mcpTools()
	case "publish":
		return mcpPublish(ctx)
	case "serve":
		return mcpServe()
	default:
		return fmt.Errorf("unknown mcp command: %s (try: tools, publish, serve)", cmd)
	}
}

func mcpTools() error {
	reg := newToolRegistry()
	meta := mcpserver.GenerateMetadata(reg)
	fmt.Printf("claudekit MCP tools (%d registered):\n\n", len(meta.Tools))
	for _, tool := range meta.Tools {
		fmt.Printf("  %-24s %s\n", tool.Name, tool.Description)
	}
	return nil
}

func mcpPublish(ctx context.Context) error {
	reg := newToolRegistry()
	fmt.Println("Publishing claudekit to MCP Registry...")
	if err := mcpserver.Publish(ctx, reg); err != nil {
		return err
	}
	fmt.Println("Published successfully.")
	return nil
}

func mcpServe() error {
	reg := newToolRegistry()
	addr := parseFlag("addr", ":8080")
	handler := mcpserver.WebMCPHandler(reg)
	fmt.Printf("Starting WebMCP server on %s\n", addr)
	fmt.Println("Endpoints:")
	fmt.Println("  GET /tools   — list registered MCP tools")
	fmt.Println("  GET /skills  — list available skills")
	fmt.Println("  GET /health  — health check")
	return http.ListenAndServe(addr, handler)
}

// --- plugin commands ---

func runPlugin(cmd string) error {
	switch cmd {
	case "list":
		return pluginList()
	case "add":
		if len(os.Args) < 4 {
			return fmt.Errorf("usage: claudekit plugin add <path>")
		}
		return pluginAdd(os.Args[3])
	default:
		return fmt.Errorf("unknown plugin command: %s (try: list, add)", cmd)
	}
}

func pluginList() error {
	dir, err := pluginkit.DefaultPluginDir()
	if err != nil {
		return err
	}

	plugins, err := pluginkit.LoadPlugins(dir)
	if err != nil {
		return err
	}

	if len(plugins) == 0 {
		fmt.Printf("No plugins installed. (plugin dir: %s)\n", dir)
		return nil
	}

	fmt.Printf("Installed plugins (%s):\n\n", dir)
	for _, p := range plugins {
		fmt.Printf("  %-24s v%-8s %d tool(s)\n", p.Name, p.Version, len(p.Tools))
	}
	return nil
}

func pluginAdd(path string) error {
	// Validate the plugin file first
	cfg, err := pluginkit.LoadPlugin(path)
	if err != nil {
		return err
	}

	dir, err := pluginkit.DefaultPluginDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create plugin dir: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	dest := filepath.Join(dir, filepath.Base(path))
	if err := os.WriteFile(dest, data, 0o644); err != nil {
		return err
	}

	fmt.Printf("Installed plugin: %s v%s (%d tools) → %s\n", cfg.Name, cfg.Version, len(cfg.Tools), dest)
	return nil
}

// --- skill commands ---

func projectDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}

func runSkill(cmd string) error {
	switch cmd {
	case "list":
		return skillList()
	case "install":
		if len(os.Args) < 4 {
			return fmt.Errorf("usage: claudekit skill install <name>")
		}
		return skillInstall(os.Args[3])
	case "remove":
		if len(os.Args) < 4 {
			return fmt.Errorf("usage: claudekit skill remove <name>")
		}
		return skillRemove(os.Args[3])
	default:
		return fmt.Errorf("unknown skill command: %s (try: list, install, remove)", cmd)
	}
}

func skillList() error {
	dir := projectDir()

	installed, err := skillkit.ListInstalled(dir)
	if err != nil {
		return err
	}

	if len(installed) > 0 {
		fmt.Println("Installed skills:")
		for _, s := range installed {
			fmt.Printf("  %-24s %s\n", s.Name, s.Description)
		}
		fmt.Println()
	}

	available, err := skillkit.AvailableSkills(dir)
	if err != nil {
		return err
	}

	if len(available) > 0 {
		fmt.Println("Available from marketplace:")
		for _, e := range available {
			fmt.Printf("  %-24s %s\n", e.Name, e.Description)
		}
		fmt.Println()
		fmt.Println("Install with: claudekit skill install <name>")
	}

	if len(installed) == 0 && len(available) == 0 {
		fmt.Println("No skills found.")
	}
	return nil
}

func skillInstall(name string) error {
	entry := skillkit.FindInIndex(name)
	if entry == nil {
		return fmt.Errorf("skill %q not found in marketplace (try: claudekit skill list)", name)
	}

	dir := projectDir()
	path, err := skillkit.Install(dir, entry.Name, entry.Content)
	if err != nil {
		return err
	}

	fmt.Printf("Installed skill: %s → %s\n", entry.Title, path)
	return nil
}

func skillRemove(name string) error {
	dir := projectDir()
	if err := skillkit.Remove(dir, name); err != nil {
		return err
	}
	fmt.Printf("Removed skill: %s\n", name)
	return nil
}

// --- ralph commands ---

func runRalph(ctx context.Context, cmd string) error {
	switch cmd {
	case "tail":
		if len(os.Args) < 4 {
			return fmt.Errorf("usage: claudekit ralph tail <progress-file>")
		}
		return ralphTail(ctx, os.Args[3])
	case "status":
		if len(os.Args) < 4 {
			return fmt.Errorf("usage: claudekit ralph status <progress-file>")
		}
		return ralphStatus(os.Args[3])
	default:
		return fmt.Errorf("unknown ralph command: %s (try: tail, status)", cmd)
	}
}

func ralphStatus(path string) error {
	p, err := ralph.LoadProgress(path)
	if err != nil {
		return fmt.Errorf("load progress: %w", err)
	}
	if parseFlag("json", "") == "true" || parseFlag("json", "") == "" && hasFlag("json") {
		data, err := json.MarshalIndent(p, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}
	printProgress(p)
	return nil
}

func ralphTail(ctx context.Context, path string) error {
	interval := parseFlag("interval", "2")
	dur, err := time.ParseDuration(interval + "s")
	if err != nil {
		dur = 2 * time.Second
	}

	var lastIter int
	fmt.Printf("Watching %s (Ctrl+C to stop)\n\n", path)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		p, err := ralph.LoadProgress(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  (waiting for progress file...)\r")
			time.Sleep(dur)
			continue
		}

		if p.Iteration != lastIter {
			fmt.Printf("\033[2K") // clear line
			printProgress(p)

			// Print last log entry if new
			if len(p.Log) > 0 {
				last := p.Log[len(p.Log)-1]
				fmt.Printf("  Last: [%s] %s\n", last.TaskID, truncate(last.Result, 80))
			}
			fmt.Println()
			lastIter = p.Iteration
		}

		if p.Status == ralph.StatusCompleted || p.Status == ralph.StatusFailed || p.Status == ralph.StatusStopped {
			fmt.Printf("Loop finished: %s\n", p.Status)
			return nil
		}

		time.Sleep(dur)
	}
}

func printProgress(p ralph.Progress) {
	fmt.Printf("Status: %s | Iteration: %d | Completed: %v\n", p.Status, p.Iteration, p.CompletedIDs)
	if p.SpecFile != "" {
		fmt.Printf("Spec: %s\n", p.SpecFile)
	}
	if !p.StartedAt.IsZero() {
		elapsed := time.Since(p.StartedAt).Truncate(time.Second)
		fmt.Printf("Elapsed: %s\n", elapsed)
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}
