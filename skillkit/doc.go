// Package skillkit provides Claude Code skill discovery, installation,
// and management.
//
// A skill is a SKILL.md file that guides Claude Code's behavior for a
// specific domain (e.g., font setup, theme configuration, environment
// bootstrap). Skills live in .claude/skills/<name>/SKILL.md within a
// project directory.
//
// # Built-in Skills
//
// [BuiltinIndex] returns the claudekit-authored skill marketplace with
// entries for font-setup, theme-setup, and env-setup. Each entry
// includes the full SKILL.md content and the MCP tools it uses.
//
// # Lifecycle
//
// [ListInstalled] scans a project's .claude/skills/ directory for
// installed skills. [Install] writes a SKILL.md into the project.
// [Remove] deletes a skill directory. [AvailableSkills] returns index
// entries that are not yet installed.
//
// # Parsing
//
// [ParseSkill] extracts the title (H1 heading) and description (first
// paragraph) from a SKILL.md file.
package skillkit
