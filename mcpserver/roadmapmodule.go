package mcpserver

import (
	"path/filepath"

	"github.com/hairglasses-studio/mcpkit/registry"
	"github.com/hairglasses-studio/mcpkit/roadmap"
)

// NewRoadmapModule creates a roadmap tool module pointing at the project's roadmap.json.
func NewRoadmapModule(projectDir string) *roadmap.Module {
	return roadmap.NewModule(roadmap.Config{
		RoadmapPath: filepath.Join(projectDir, "roadmap.json"),
	})
}

// SetupRoadmap registers the roadmap module with the registry.
func SetupRoadmap(reg *registry.ToolRegistry, projectDir string) {
	reg.RegisterModule(NewRoadmapModule(projectDir))
}
