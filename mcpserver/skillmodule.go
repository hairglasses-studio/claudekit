package mcpserver

import (
	"context"
	"fmt"

	"github.com/hairglasses-studio/claudekit/skillkit"
	"github.com/hairglasses-studio/mcpkit/handler"
	"github.com/hairglasses-studio/mcpkit/registry"
)

// SkillModule exposes skill marketplace tools via MCP.
type SkillModule struct {
	ProjectDir string // Project directory to install skills into
}

func (m *SkillModule) Name() string        { return "skills" }
func (m *SkillModule) Description() string { return "Claude Code skill marketplace — discover, install, and manage skills" }

type SkillListInput struct {
	Filter string `json:"filter,omitempty" jsonschema:"description=Filter: installed|available|all (default all),enum=installed,enum=available,enum=all"`
}

type SkillListOutput struct {
	Installed []SkillInfo `json:"installed"`
	Available []SkillInfo `json:"available,omitempty"`
}

type SkillInfo struct {
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tools       []string `json:"tools,omitempty"`
	Installed   bool     `json:"installed"`
}

type SkillInstallInput struct {
	Name string `json:"name" jsonschema:"description=Skill name to install from the marketplace (e.g. font-setup),required"`
}

type SkillInstallOutput struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Message string `json:"message"`
}

func (m *SkillModule) Tools() []registry.ToolDefinition {
	projectDir := m.ProjectDir

	return []registry.ToolDefinition{
		handler.TypedHandler[SkillListInput, SkillListOutput](
			"skill_list",
			"List installed and available Claude Code skills. Skills are behavior guides that teach Claude how to use claudekit MCP tools effectively.",
			func(_ context.Context, input SkillListInput) (SkillListOutput, error) {
				filter := input.Filter
				if filter == "" {
					filter = "all"
				}

				var out SkillListOutput

				if filter == "installed" || filter == "all" {
					installed, err := skillkit.ListInstalled(projectDir)
					if err != nil {
						return SkillListOutput{}, err
					}
					for _, s := range installed {
						out.Installed = append(out.Installed, SkillInfo{
							Name:        s.Name,
							Title:       s.Title,
							Description: s.Description,
							Installed:   true,
						})
					}
				}

				if filter == "available" || filter == "all" {
					available, err := skillkit.AvailableSkills(projectDir)
					if err != nil {
						return SkillListOutput{}, err
					}
					for _, e := range available {
						out.Available = append(out.Available, SkillInfo{
							Name:        e.Name,
							Title:       e.Title,
							Description: e.Description,
							Tools:       e.Tools,
							Installed:   false,
						})
					}
				}

				return out, nil
			},
		),
		handler.TypedHandler[SkillInstallInput, SkillInstallOutput](
			"skill_install",
			"Install a Claude Code skill from the marketplace into the current project. Skills teach Claude how to use claudekit tools for specific workflows.",
			func(_ context.Context, input SkillInstallInput) (SkillInstallOutput, error) {
				entry := skillkit.FindInIndex(input.Name)
				if entry == nil {
					return SkillInstallOutput{}, fmt.Errorf("skill %q not found in marketplace", input.Name)
				}

				path, err := skillkit.Install(projectDir, entry.Name, entry.Content)
				if err != nil {
					return SkillInstallOutput{}, err
				}

				return SkillInstallOutput{
					Name:    entry.Name,
					Path:    path,
					Message: fmt.Sprintf("Installed %q skill — Claude will now use it when relevant topics come up", entry.Title),
				}, nil
			},
		),
	}
}
