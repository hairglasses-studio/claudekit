package mcpserver

import (
	"path/filepath"

	"github.com/hairglasses-studio/mcpkit/rdcycle"
	"github.com/hairglasses-studio/mcpkit/registry"
)

// NewRDCycleModule creates an R&D cycle module configured for the claudekit ecosystem.
func NewRDCycleModule(projectDir string) *rdcycle.Module {
	return rdcycle.NewModule(rdcycle.CycleConfig{
		RoadmapPath: filepath.Join(projectDir, "roadmap.json"),
		GitRoot:     projectDir,
		ScanRepos: []string{
			"hairglasses-studio/mcpkit",
			"hairglasses-studio/claudekit",
			"githubnext/monaspace",
			"catppuccin/catppuccin",
		},
	})
}

// SetupRDCycle registers the R&D cycle module with the registry.
func SetupRDCycle(reg *registry.ToolRegistry, projectDir string) {
	reg.RegisterModule(NewRDCycleModule(projectDir))
}
