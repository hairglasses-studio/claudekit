package mcpserver

import (
	"context"
	"fmt"
	"os"

	"github.com/hairglasses-studio/mcpkit/discovery"
	"github.com/hairglasses-studio/mcpkit/registry"
)

const (
	claudekitName        = "claudekit"
	claudekitDescription = "Claude Code terminal customization toolkit — fonts, themes, statusline"
)

// GenerateMetadata creates ServerMetadata from the current tool registry
// for publishing to the MCP Registry.
func GenerateMetadata(reg *registry.ToolRegistry) discovery.ServerMetadata {
	return discovery.MetadataFromRegistry(
		claudekitName,
		claudekitDescription,
		reg,
		[]discovery.TransportInfo{{Type: "stdio"}},
	)
}

// Publish registers or updates claudekit in the MCP Registry.
// Requires CLAUDEKIT_REGISTRY_TOKEN env var.
func Publish(ctx context.Context, reg *registry.ToolRegistry) error {
	token := os.Getenv("CLAUDEKIT_REGISTRY_TOKEN")
	if token == "" {
		return fmt.Errorf("CLAUDEKIT_REGISTRY_TOKEN environment variable is required to publish to the MCP Registry")
	}

	pub, err := discovery.NewPublisher(discovery.PublisherConfig{
		Token: token,
	})
	if err != nil {
		return fmt.Errorf("creating publisher: %w", err)
	}

	meta := GenerateMetadata(reg)
	_, err = pub.Register(ctx, meta)
	if err != nil {
		return fmt.Errorf("publishing to MCP Registry: %w", err)
	}
	return nil
}
