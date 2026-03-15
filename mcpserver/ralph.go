package mcpserver

import (
	"github.com/hairglasses-studio/mcpkit/ralph"
	"github.com/hairglasses-studio/mcpkit/registry"
	"github.com/hairglasses-studio/mcpkit/sampling"
)

// SetupRalph registers the ralph autonomous loop module with the registry.
// sampler can be nil if sampling is not yet available.
func SetupRalph(reg *registry.ToolRegistry, sampler sampling.SamplingClient) {
	mod := ralph.NewModule(reg, sampler)
	reg.RegisterModule(mod)
}
