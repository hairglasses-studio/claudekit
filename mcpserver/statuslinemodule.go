package mcpserver

import (
	"context"

	"github.com/hairglasses-studio/claudekit/statusline"
	"github.com/hairglasses-studio/mcpkit/handler"
	"github.com/hairglasses-studio/mcpkit/registry"
)

// StatuslineModule exposes statusline management tools via MCP.
type StatuslineModule struct{}

func (m *StatuslineModule) Name() string        { return "statusline" }
func (m *StatuslineModule) Description() string { return "Claude Code statusline installation and management" }

type StatuslineInstallInput struct {
	Style string `json:"style,omitempty" jsonschema:"description=Statusline style (full|compact|minimal),enum=full,enum=compact,enum=minimal"`
}

type StatuslineInstallOutput struct {
	ScriptPath   string `json:"script_path"`
	SettingsPath string `json:"settings_path"`
	BackedUp     string `json:"backed_up,omitempty"`
	Message      string `json:"message"`
}

func (m *StatuslineModule) Tools() []registry.ToolDefinition {
	return []registry.ToolDefinition{
		handler.TypedHandler[StatuslineInstallInput, StatuslineInstallOutput](
			"statusline_install",
			"Install the Claude Code statusline script and update settings. The statusline shows model info, context usage, cost, and duration with Nerd Font icon support.",
			func(_ context.Context, input StatuslineInstallInput) (StatuslineInstallOutput, error) {
				style := statusline.StyleFull
				switch input.Style {
				case "compact":
					style = statusline.StyleCompact
				case "minimal":
					style = statusline.StyleMinimal
				}

				result, err := statusline.Install(statusline.InstallOpts{Style: style})
				if err != nil {
					return StatuslineInstallOutput{}, err
				}

				return StatuslineInstallOutput{
					ScriptPath:   result.ScriptPath,
					SettingsPath: result.SettingsPath,
					BackedUp:     result.BackedUp,
					Message:      "statusline installed — restart Claude Code to activate",
				}, nil
			},
		),
	}
}
