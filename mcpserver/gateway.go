package mcpserver

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hairglasses-studio/mcpkit/gateway"
	"github.com/hairglasses-studio/mcpkit/registry"
)

// SetupGateway creates a Gateway that aggregates claudekit tools with external MCP servers.
// upstreams is a list of "name=url" pairs (e.g., "github=http://localhost:8080").
// Returns the gateway and its DynamicRegistry for the caller to manage lifecycle.
func SetupGateway(reg *registry.ToolRegistry, upstreams []string) (*gateway.Gateway, *registry.DynamicRegistry, error) {
	gw, dynReg := gateway.NewGateway()

	for _, spec := range upstreams {
		name, url, err := parseUpstream(spec)
		if err != nil {
			_ = gw.Close()
			return nil, nil, fmt.Errorf("invalid upstream %q: %w", spec, err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, err = gw.AddUpstream(ctx, gateway.UpstreamConfig{
			Name: name,
			URL:  url,
		})
		cancel()
		if err != nil {
			_ = gw.Close()
			return nil, nil, fmt.Errorf("adding upstream %q: %w", name, err)
		}
	}

	return gw, dynReg, nil
}

// parseUpstream splits a "name=url" string into its components.
func parseUpstream(spec string) (name, url string, err error) {
	parts := strings.SplitN(spec, "=", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("expected format \"name=url\", got %q", spec)
	}
	return parts[0], parts[1], nil
}
