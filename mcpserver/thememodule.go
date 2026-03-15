package mcpserver

import (
	"context"
	"fmt"

	"github.com/hairglasses-studio/claudekit/fontkit"
	"github.com/hairglasses-studio/claudekit/themekit"
	"github.com/hairglasses-studio/mcpkit/handler"
	"github.com/hairglasses-studio/mcpkit/registry"
)

// ThemeModule exposes terminal theme tools via MCP.
type ThemeModule struct{}

func (m *ThemeModule) Name() string        { return "theme" }
func (m *ThemeModule) Description() string { return "Terminal color theme management with Catppuccin support" }

type ThemeApplyInput struct {
	Flavor   string `json:"flavor,omitempty" jsonschema:"description=Catppuccin flavor (mocha|macchiato|frappe|latte),enum=mocha,enum=macchiato,enum=frappe,enum=latte"`
	Terminal string `json:"terminal,omitempty" jsonschema:"description=Terminal to configure (auto|iterm2|ghostty|wezterm),enum=auto,enum=iterm2,enum=ghostty,enum=wezterm"`
}

type ThemeApplyOutput struct {
	ConfigPath string `json:"config_path"`
	Terminal   string `json:"terminal"`
	Flavor     string `json:"flavor"`
	Message    string `json:"message"`
}

type ThemeListInput struct{}

type ThemeListOutput struct {
	Flavors []ThemeFlavorInfo `json:"flavors"`
}

type ThemeFlavorInfo struct {
	Name   string `json:"name"`
	Flavor string `json:"flavor"`
	Type   string `json:"type"` // "dark" or "light"
}

func (m *ThemeModule) Tools() []registry.ToolDefinition {
	return []registry.ToolDefinition{
		handler.TypedHandler[ThemeApplyInput, ThemeApplyOutput](
			"theme_apply",
			"Apply a Catppuccin color theme to the terminal. Generates iTerm2 Dynamic Profile colors or Ghostty theme file. Defaults to Mocha flavor.",
			func(_ context.Context, input ThemeApplyInput) (ThemeApplyOutput, error) {
				flavorStr := input.Flavor
				if flavorStr == "" {
					flavorStr = "mocha"
				}

				var flavor themekit.Flavor
				switch flavorStr {
				case "mocha":
					flavor = themekit.Mocha
				case "macchiato":
					flavor = themekit.Macchiato
				case "frappe":
					flavor = themekit.Frappe
				case "latte":
					flavor = themekit.Latte
				default:
					return ThemeApplyOutput{}, fmt.Errorf("unknown flavor: %s", flavorStr)
				}

				p := themekit.Catppuccin(flavor)

				terminal := input.Terminal
				if terminal == "" || terminal == "auto" {
					term := fontkit.DetectTerminal()
					terminal = string(term.Terminal)
				}

				var path string
				var err error
				switch fontkit.Terminal(terminal) {
				case fontkit.TerminalITerm2:
					path, err = themekit.ExportITerm2(p)
				case fontkit.TerminalGhostty:
					path, err = themekit.ExportGhostty(p)
				case fontkit.TerminalWezTerm:
					path, err = themekit.ExportWezTerm(p)
				default:
					return ThemeApplyOutput{
						Terminal: terminal,
						Flavor:   flavorStr,
						Message:  "unsupported terminal — only iTerm2, Ghostty, and WezTerm are supported",
					}, nil
				}
				if err != nil {
					return ThemeApplyOutput{}, err
				}

				return ThemeApplyOutput{
					ConfigPath: path,
					Terminal:   terminal,
					Flavor:     flavorStr,
					Message:    fmt.Sprintf("%s theme applied", p.Name),
				}, nil
			},
		),
		handler.TypedHandler[ThemeListInput, ThemeListOutput](
			"theme_list",
			"List available Catppuccin theme flavors with their characteristics.",
			func(_ context.Context, _ ThemeListInput) (ThemeListOutput, error) {
				return ThemeListOutput{
					Flavors: []ThemeFlavorInfo{
						{Name: "Catppuccin Mocha", Flavor: "mocha", Type: "dark"},
						{Name: "Catppuccin Macchiato", Flavor: "macchiato", Type: "dark"},
						{Name: "Catppuccin Frappé", Flavor: "frappe", Type: "dark"},
						{Name: "Catppuccin Latte", Flavor: "latte", Type: "light"},
					},
				}, nil
			},
		),
	}
}
