package mcpserver

import (
	"testing"
)

func TestSkillModuleNameDescription(t *testing.T) {
	mod := &SkillModule{ProjectDir: t.TempDir()}

	if mod.Name() != "skills" {
		t.Errorf("Name() = %q, want %q", mod.Name(), "skills")
	}
	if mod.Description() == "" {
		t.Error("Description() should not be empty")
	}
}

func TestSkillModuleToolCount(t *testing.T) {
	mod := &SkillModule{ProjectDir: t.TempDir()}
	tools := mod.Tools()

	if len(tools) != 2 {
		t.Errorf("got %d tools, want 2", len(tools))
	}

	wantNames := map[string]bool{
		"skill_list":    false,
		"skill_install": false,
	}
	for _, td := range tools {
		if _, ok := wantNames[td.Tool.Name]; ok {
			wantNames[td.Tool.Name] = true
		}
	}
	for name, found := range wantNames {
		if !found {
			t.Errorf("missing tool %q", name)
		}
	}
}

func TestSkillListHandler(t *testing.T) {
	mod := &SkillModule{ProjectDir: t.TempDir()}
	tools := mod.Tools()

	td := findTool(tools, "skill_list")
	if td == nil {
		t.Fatal("skill_list tool not found")
	}

	result := callTool(t, td, map[string]interface{}{
		"filter": "available",
	})

	var out SkillListOutput
	extractJSON(t, result, &out)

	// The builtin index has at least one entry.
	if len(out.Available) == 0 {
		t.Error("expected at least one available skill from the builtin index")
	}

	// Verify each available skill has required fields.
	for _, s := range out.Available {
		if s.Name == "" {
			t.Error("available skill has empty Name")
		}
		if s.Title == "" {
			t.Errorf("skill %q has empty Title", s.Name)
		}
	}
}
