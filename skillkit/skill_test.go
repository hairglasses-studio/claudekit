package skillkit

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSkill(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "SKILL.md")
	os.WriteFile(path, []byte("# My Skill\n\nDoes cool stuff when triggered.\n\n## Workflow\n\n1. Do things\n"), 0o644)

	s, err := ParseSkill(path)
	if err != nil {
		t.Fatal(err)
	}

	if s.Title != "My Skill" {
		t.Errorf("title = %q, want My Skill", s.Title)
	}
	if s.Description != "Does cool stuff when triggered." {
		t.Errorf("description = %q, want 'Does cool stuff when triggered.'", s.Description)
	}
}

func TestListInstalled(t *testing.T) {
	tmpDir := t.TempDir()
	skillsDir := filepath.Join(tmpDir, ".claude", "skills", "test-skill")
	os.MkdirAll(skillsDir, 0o755)
	os.WriteFile(filepath.Join(skillsDir, "SKILL.md"), []byte("# Test\n\nA test skill.\n"), 0o644)

	skills, err := ListInstalled(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(skills) != 1 {
		t.Fatalf("got %d skills, want 1", len(skills))
	}
	if skills[0].Name != "test-skill" {
		t.Errorf("name = %q, want test-skill", skills[0].Name)
	}
	if !skills[0].Installed {
		t.Error("expected Installed = true")
	}
}

func TestListInstalledEmpty(t *testing.T) {
	tmpDir := t.TempDir()

	skills, err := ListInstalled(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if skills != nil {
		t.Errorf("expected nil, got %d skills", len(skills))
	}
}

func TestInstallAndRemove(t *testing.T) {
	tmpDir := t.TempDir()

	path, err := Install(tmpDir, "my-skill", "# My Skill\n\nDoes things.\n")
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "# My Skill\n\nDoes things.\n" {
		t.Error("content mismatch")
	}

	// Verify it shows up in list
	skills, _ := ListInstalled(tmpDir)
	if len(skills) != 1 {
		t.Fatalf("expected 1 installed skill, got %d", len(skills))
	}

	// Remove
	if err := Remove(tmpDir, "my-skill"); err != nil {
		t.Fatal(err)
	}

	skills, _ = ListInstalled(tmpDir)
	if len(skills) != 0 {
		t.Errorf("expected 0 skills after remove, got %d", len(skills))
	}
}

func TestRemoveNonexistent(t *testing.T) {
	tmpDir := t.TempDir()
	err := Remove(tmpDir, "nope")
	if err == nil {
		t.Error("expected error for nonexistent skill")
	}
}

func TestBuiltinIndex(t *testing.T) {
	index := BuiltinIndex()
	if len(index) != 3 {
		t.Errorf("builtin index has %d entries, want 3", len(index))
	}

	for _, e := range index {
		if e.Name == "" {
			t.Error("index entry has empty name")
		}
		if e.Content == "" {
			t.Error("index entry has empty content")
		}
		if len(e.Tools) == 0 {
			t.Errorf("index entry %q has no tools", e.Name)
		}
	}
}

func TestFindInIndex(t *testing.T) {
	e := FindInIndex("font-setup")
	if e == nil {
		t.Fatal("expected to find font-setup")
	}
	if e.Title != "Font Setup Skill" {
		t.Errorf("title = %q", e.Title)
	}

	if FindInIndex("nonexistent") != nil {
		t.Error("expected nil for nonexistent skill")
	}
}

func TestAvailableSkills(t *testing.T) {
	tmpDir := t.TempDir()

	// Nothing installed — all should be available
	available, err := AvailableSkills(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(available) != 3 {
		t.Errorf("expected 3 available, got %d", len(available))
	}

	// Install one
	Install(tmpDir, "font-setup", "# Font Setup Skill\n\nTest.\n")

	available, err = AvailableSkills(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(available) != 2 {
		t.Errorf("expected 2 available after installing one, got %d", len(available))
	}
}
