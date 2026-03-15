// Package mcpserver provides MCP tool modules for claudekit.
package mcpserver

import (
	"context"

	"github.com/hairglasses-studio/claudekit/fontkit"
	"github.com/hairglasses-studio/mcpkit/handler"
	"github.com/hairglasses-studio/mcpkit/registry"
)

// FontModule exposes font management tools via MCP.
type FontModule struct{}

func (m *FontModule) Name() string        { return "fonts" }
func (m *FontModule) Description() string { return "Font detection, installation, and terminal configuration" }

// --- font_status ---

type FontStatusInput struct{}

type FontStatusOutput struct {
	Terminal     string   `json:"terminal"`
	TermVersion  string   `json:"terminal_version"`
	ConfigSupport bool   `json:"config_support"`
	BestFont     string   `json:"best_font"`
	BrewAvail    bool     `json:"brew_available"`
	Installed    []string `json:"installed"`
	NotInstalled []string `json:"not_installed"`
}

// --- font_install ---

type FontInstallInput struct {
	NerdFont bool `json:"nerd_font,omitempty" jsonschema:"description=Install Monaspice Nerd Font variant (default true)"`
	DryRun   bool `json:"dry_run,omitempty" jsonschema:"description=Preview what would be installed without doing it"`
}

type FontInstallOutput struct {
	Cask      string `json:"cask"`
	DryRun    bool   `json:"dry_run"`
	Output    string `json:"output"`
	AlreadyOK bool   `json:"already_installed"`
}

// --- font_configure ---

type FontConfigureInput struct {
	Terminal string `json:"terminal,omitempty" jsonschema:"description=Terminal to configure (auto|iterm2|ghostty|wezterm),enum=auto,enum=iterm2,enum=ghostty,enum=wezterm"`
	FontSize int    `json:"font_size,omitempty" jsonschema:"description=Font size in points (default 15)"`
}

type FontConfigureOutput struct {
	ConfigPath string `json:"config_path"`
	Terminal   string `json:"terminal"`
	Message    string `json:"message"`
}

func (m *FontModule) Tools() []registry.ToolDefinition {
	return []registry.ToolDefinition{
		handler.TypedHandler[FontStatusInput, FontStatusOutput](
			"font_status",
			"Detect installed Monaspace/Monaspice fonts and terminal emulator. Returns which fonts are available, which terminal is running, and the best font to use.",
			func(ctx context.Context, _ FontStatusInput) (FontStatusOutput, error) {
				term := fontkit.DetectTerminal()
				status, err := fontkit.Detect(ctx)
				if err != nil {
					return FontStatusOutput{}, err
				}

				out := FontStatusOutput{
					Terminal:      string(term.Terminal),
					TermVersion:   term.Version,
					ConfigSupport: term.Terminal.SupportsConfig(),
					BestFont:      status.BestFont,
					BrewAvail:     status.BrewAvail,
				}
				for _, f := range status.Installed {
					out.Installed = append(out.Installed, f.Name)
				}
				for _, f := range status.NotInstalled {
					out.NotInstalled = append(out.NotInstalled, f.Name)
				}
				return out, nil
			},
		),
		handler.TypedHandler[FontInstallInput, FontInstallOutput](
			"font_install",
			"Install Monaspace or Monaspice fonts via Homebrew. Defaults to Monaspice Nerd Font for full icon support. Supports dry-run mode.",
			func(ctx context.Context, input FontInstallInput) (FontInstallOutput, error) {
				result, err := fontkit.Install(ctx, fontkit.InstallOpts{
					NerdFont: true, // Always install Nerd Font variant
					DryRun:   input.DryRun,
				})
				if err != nil {
					return FontInstallOutput{}, err
				}

				return FontInstallOutput{
					Cask:      result.Cask,
					DryRun:    result.DryRun,
					Output:    result.Output,
					AlreadyOK: result.AlreadyOK,
				}, nil
			},
		),
		handler.TypedHandler[FontConfigureInput, FontConfigureOutput](
			"font_configure",
			"Write font configuration for the detected or specified terminal emulator. Creates iTerm2 Dynamic Profile or Ghostty config with Monaspice + Menlo fallback.",
			func(_ context.Context, input FontConfigureInput) (FontConfigureOutput, error) {
				terminal := input.Terminal
				if terminal == "" || terminal == "auto" {
					term := fontkit.DetectTerminal()
					terminal = string(term.Terminal)
				}

				fontSize := input.FontSize

				var path string
				var err error
				switch fontkit.Terminal(terminal) {
				case fontkit.TerminalITerm2:
					path, err = fontkit.ConfigureITerm2(fontkit.ITerm2Opts{FontSize: fontSize})
				case fontkit.TerminalGhostty:
					path, err = fontkit.ConfigureGhostty(fontkit.GhosttyOpts{FontSize: fontSize})
				case fontkit.TerminalWezTerm:
					path, err = fontkit.ConfigureWezTerm(fontkit.WezTermOpts{FontSize: fontSize})
				default:
					return FontConfigureOutput{
						Terminal: terminal,
						Message:  "unsupported terminal — only iTerm2, Ghostty, and WezTerm are supported",
					}, nil
				}
				if err != nil {
					return FontConfigureOutput{}, err
				}

				msg := "font configuration written"
				if fontkit.Terminal(terminal) == fontkit.TerminalITerm2 {
					msg = "iTerm2 Dynamic Profile created — switch to \"Claudekit Monaspace\" in Preferences → Profiles"
				}

				return FontConfigureOutput{
					ConfigPath: path,
					Terminal:   terminal,
					Message:    msg,
				}, nil
			},
		),
	}
}
