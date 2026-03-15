package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hairglasses-studio/claudekit/fontkit"
	"github.com/hairglasses-studio/claudekit/statusline"
	"github.com/hairglasses-studio/claudekit/themekit"
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
	case "statusline":
		err = runStatusline(ctx, cmd)
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
  claudekit theme preview         Preview all Catppuccin flavors

  claudekit statusline install    Install the Claude Code statusline
  claudekit statusline preview    Preview statusline with sample data`)
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
		return err
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
	return nil
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
	default:
		return fmt.Errorf("unsupported terminal: %s (supported: auto, iterm2, ghostty)", termFlag)
	}

	var path string
	var err error
	switch terminal {
	case fontkit.TerminalITerm2:
		path, err = fontkit.ConfigureITerm2(fontkit.ITerm2Opts{})
	case fontkit.TerminalGhostty:
		path, err = fontkit.ConfigureGhostty(fontkit.GhosttyOpts{})
	default:
		return fmt.Errorf("unsupported terminal: %s (supported: iTerm2, Ghostty)", terminal)
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
	if err := fontsStatus(ctx); err != nil {
		return err
	}

	status, err := fontkit.Detect(ctx)
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
	case "preview":
		return themePreview()
	default:
		return fmt.Errorf("unknown theme command: %s (try: apply, preview)", cmd)
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
	default:
		return fmt.Errorf("unsupported terminal: %s (supported: auto, iterm2, ghostty)", termStr)
	}

	var path string
	var err error
	switch terminal {
	case fontkit.TerminalITerm2:
		path, err = themekit.ExportITerm2(p)
	case fontkit.TerminalGhostty:
		path, err = themekit.ExportGhostty(p)
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

func runStatusline(ctx context.Context, cmd string) error {
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
