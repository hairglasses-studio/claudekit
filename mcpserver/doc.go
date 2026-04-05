// Package mcpserver implements the claudekit MCP server, exposing terminal
// customization tools over the Model Context Protocol.
//
// It provides mcpkit [registry.ToolModule] implementations for each
// claudekit domain: fonts, themes, environment, skills, statusline,
// workflows, R&D cycles, FinOps, memory, and roadmap management.
//
// # Modules
//
// Each tool domain is implemented as a separate module that registers
// its tools with the mcpkit registry:
//
//   - EnvModule: shell detection, mise status, dotfile snapshots
//   - FontModule: font detection, installation, terminal configuration
//   - ThemeModule: Catppuccin theme listing and application
//   - SkillModule: skill marketplace, installation, removal
//   - StatuslineModule: statusline script installation
//   - WorkflowModule: multi-step workflow orchestration
//   - RDCycleModule: R&D cycle tracking
//   - FinOpsModule: cost and usage analytics
//   - MemoryModule: cross-session key-value memory
//   - RoadmapModule: roadmap status and planning
//
// # Gateway
//
// [SetupGateway] creates a gateway that aggregates claudekit tools with
// external upstream MCP servers specified as "name=url" pairs.
//
// # Discovery
//
// [GenerateMetadata] and [Publish] handle MCP Registry integration for
// tool discovery.
package mcpserver
