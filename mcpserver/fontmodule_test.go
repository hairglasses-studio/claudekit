package mcpserver

import "testing"

func TestFontModuleInterface(t *testing.T) {
	m := &FontModule{}
	if m.Name() != "fonts" {
		t.Errorf("Name() = %q, want fonts", m.Name())
	}
	if m.Description() == "" {
		t.Error("Description() should not be empty")
	}
	tools := m.Tools()
	if len(tools) != 3 {
		t.Fatalf("Tools() = %d tools, want 3", len(tools))
	}

	names := map[string]bool{}
	for _, td := range tools {
		names[td.Tool.Name] = true
	}
	for _, want := range []string{"font_status", "font_install", "font_configure"} {
		if !names[want] {
			t.Errorf("missing tool %q", want)
		}
	}
}

func TestStatuslineModuleInterface(t *testing.T) {
	m := &StatuslineModule{}
	if m.Name() != "statusline" {
		t.Errorf("Name() = %q, want statusline", m.Name())
	}
	tools := m.Tools()
	if len(tools) != 1 {
		t.Fatalf("Tools() = %d tools, want 1", len(tools))
	}
	if tools[0].Tool.Name != "statusline_install" {
		t.Errorf("tool name = %q, want statusline_install", tools[0].Tool.Name)
	}
}
