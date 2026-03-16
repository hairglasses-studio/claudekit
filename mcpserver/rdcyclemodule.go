package mcpserver

import (
	"path/filepath"

	"github.com/hairglasses-studio/mcpkit/rdcycle"
	"github.com/hairglasses-studio/mcpkit/registry"
)

// NewRDCycleModule creates an R&D cycle module configured for the claudekit ecosystem.
// When store is non-nil it is used for persistent artifact storage; otherwise
// the module defaults to an in-memory store.
func NewRDCycleModule(projectDir string, store rdcycle.ArtifactStore) *rdcycle.Module {
	var opts []rdcycle.ModuleOption
	if store != nil {
		opts = append(opts, rdcycle.WithArtifactStore(store))
	}
	return rdcycle.NewModule(rdcycle.CycleConfig{
		RoadmapPath: filepath.Join(projectDir, "roadmap.json"),
		GitRoot:     projectDir,
		ScanRepos: []string{
			"hairglasses-studio/mcpkit",
			"hairglasses-studio/claudekit",
			"githubnext/monaspace",
			"catppuccin/catppuccin",
		},
	}, opts...)
}

// SetupRDCycle registers the R&D cycle module with the registry and returns
// the module so callers can wire a RalphStarter for perpetual loops.
func SetupRDCycle(reg *registry.ToolRegistry, projectDir string, store rdcycle.ArtifactStore) *rdcycle.Module {
	mod := NewRDCycleModule(projectDir, store)
	reg.RegisterModule(mod)
	return mod
}
