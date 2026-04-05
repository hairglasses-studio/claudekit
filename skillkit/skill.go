package skillkit

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Skill represents a Claude Code skill (a SKILL.md behavior guide).
type Skill struct {
	Name        string // Directory name (e.g. "font-setup")
	Title       string // H1 title from SKILL.md
	Description string // First paragraph after the title
	Path        string // Absolute path to SKILL.md
	Installed   bool   // Whether this skill is installed locally
}

// ListInstalled scans a .claude/skills/ directory for installed skills.
func ListInstalled(projectDir string) ([]Skill, error) {
	skillsDir := filepath.Join(projectDir, ".claude", "skills")
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read skills dir: %w", err)
	}

	var skills []Skill
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		skillPath := filepath.Join(skillsDir, entry.Name(), "SKILL.md")
		if _, err := os.Stat(skillPath); err != nil {
			continue
		}
		s, err := ParseSkill(skillPath)
		if err != nil {
			continue
		}
		s.Name = entry.Name()
		s.Path = skillPath
		s.Installed = true
		skills = append(skills, s)
	}
	return skills, nil
}

// ParseSkill extracts metadata from a SKILL.md file.
func ParseSkill(path string) (Skill, error) {
	f, err := os.Open(path)
	if err != nil {
		return Skill{}, err
	}
	defer f.Close()

	var s Skill
	scanner := bufio.NewScanner(f)
	foundTitle := false

	for scanner.Scan() {
		line := scanner.Text()

		// Extract H1 title
		if !foundTitle && strings.HasPrefix(line, "# ") {
			s.Title = strings.TrimPrefix(line, "# ")
			foundTitle = true
			continue
		}

		// First non-empty line after title is the description
		if foundTitle && s.Description == "" && strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "#") {
			s.Description = strings.TrimSpace(line)
			break
		}
	}

	return s, scanner.Err()
}

// Install writes a SKILL.md file into the project's .claude/skills/ directory.
func Install(projectDir, name, content string) (string, error) {
	skillDir := filepath.Join(projectDir, ".claude", "skills", name)
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		return "", fmt.Errorf("create skill dir: %w", err)
	}

	skillPath := filepath.Join(skillDir, "SKILL.md")
	if err := os.WriteFile(skillPath, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("write skill: %w", err)
	}

	return skillPath, nil
}

// Remove deletes a skill directory from the project.
func Remove(projectDir, name string) error {
	skillDir := filepath.Join(projectDir, ".claude", "skills", name)
	if _, err := os.Stat(skillDir); os.IsNotExist(err) {
		return fmt.Errorf("skill %q not installed", name)
	}
	return os.RemoveAll(skillDir)
}
