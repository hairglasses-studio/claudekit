package mcpserver

import (
	"context"

	"github.com/hairglasses-studio/claudekit/envkit"
	"github.com/hairglasses-studio/mcpkit/handler"
	"github.com/hairglasses-studio/mcpkit/registry"
)

// EnvModule exposes environment bootstrap tools via MCP.
type EnvModule struct{}

func (m *EnvModule) Name() string        { return "env" }
func (m *EnvModule) Description() string { return "Environment bootstrap — mise, shell, and dotfile management" }

// --- env_status ---

type EnvStatusInput struct{}

type EnvStatusOutput struct {
	Mise  *envkit.MiseInfo  `json:"mise"`
	Shell envkit.ShellInfo  `json:"shell"`
}

// --- env_snapshot ---

type EnvSnapshotInput struct{}

type EnvSnapshotOutput struct {
	Files   map[string]string `json:"files"`
	Summary string            `json:"summary"`
}

func (m *EnvModule) Tools() []registry.ToolDefinition {
	return []registry.ToolDefinition{
		handler.TypedHandler[EnvStatusInput, EnvStatusOutput](
			"env_status",
			"Check environment status: mise tool versions and shell configuration. Returns mise installation info, active tool versions, detected shell, and plugin manager.",
			func(ctx context.Context, _ EnvStatusInput) (EnvStatusOutput, error) {
				mise, err := envkit.MiseStatus(ctx)
				if err != nil {
					return EnvStatusOutput{}, err
				}
				shell := envkit.DetectShell()

				return EnvStatusOutput{
					Mise:  mise,
					Shell: shell,
				}, nil
			},
		),
		handler.TypedHandler[EnvSnapshotInput, EnvSnapshotOutput](
			"env_snapshot",
			"Capture a snapshot of all claudekit-managed config files (starship, ghostty, bat, delta, iTerm2 profiles). Returns file paths and contents.",
			func(_ context.Context, _ EnvSnapshotInput) (EnvSnapshotOutput, error) {
				snap, err := envkit.Snapshot()
				if err != nil {
					return EnvSnapshotOutput{}, err
				}

				return EnvSnapshotOutput{
					Files:   snap.Files,
					Summary: envkit.SnapshotSummary(snap),
				}, nil
			},
		),
	}
}
