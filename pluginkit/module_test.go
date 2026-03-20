package pluginkit

import (
	"testing"
	"time"
)

func makeTestConfig() *PluginConfig {
	return &PluginConfig{
		Name:        "test-plugin",
		Description: "A test plugin",
		Version:     "1.0.0",
		Handler: HandlerConfig{
			Type:    "subprocess",
			Command: "echo hello",
		},
		Tools: []PluginToolDef{
			{
				Name:        "test_tool",
				Description: "A test tool",
				InputSchema: map[string]interface{}{
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "The name",
						},
					},
					"required": []interface{}{"name"},
				},
			},
			{
				Name:        "another_tool",
				Description: "Another test tool",
			},
		},
	}
}

func TestPluginModuleNameDescription(t *testing.T) {
	m := NewPluginModule(makeTestConfig())
	if m.Name() != "test-plugin" {
		t.Errorf("expected name 'test-plugin', got %q", m.Name())
	}
	if m.Description() != "A test plugin" {
		t.Errorf("expected description 'A test plugin', got %q", m.Description())
	}
}

func TestPluginModuleTools(t *testing.T) {
	m := NewPluginModule(makeTestConfig())
	tools := m.Tools()
	if len(tools) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(tools))
	}
	if tools[0].Tool.Name != "test_tool" {
		t.Errorf("expected first tool 'test_tool', got %q", tools[0].Tool.Name)
	}
	if tools[1].Tool.Name != "another_tool" {
		t.Errorf("expected second tool 'another_tool', got %q", tools[1].Tool.Name)
	}
}

func TestPluginModuleTimeout(t *testing.T) {
	// Default timeout (no timeout string)
	cfg := makeTestConfig()
	m := NewPluginModule(cfg)
	if m.handler.Timeout != 30*time.Second {
		t.Errorf("expected default timeout 30s, got %v", m.handler.Timeout)
	}

	// Custom timeout
	cfg2 := makeTestConfig()
	cfg2.Handler.Timeout = "10s"
	m2 := NewPluginModule(cfg2)
	if m2.handler.Timeout != 10*time.Second {
		t.Errorf("expected timeout 10s, got %v", m2.handler.Timeout)
	}

	// Invalid timeout falls back to default
	cfg3 := makeTestConfig()
	cfg3.Handler.Timeout = "not-a-duration"
	m3 := NewPluginModule(cfg3)
	if m3.handler.Timeout != 30*time.Second {
		t.Errorf("expected fallback timeout 30s for invalid duration, got %v", m3.handler.Timeout)
	}
}

func TestPluginModuleInputSchema(t *testing.T) {
	m := NewPluginModule(makeTestConfig())
	tools := m.Tools()

	schema := tools[0].Tool.InputSchema
	if schema.Type != "object" {
		t.Errorf("expected schema type 'object', got %q", schema.Type)
	}
	if _, ok := schema.Properties["name"]; !ok {
		t.Error("expected 'name' property in schema")
	}
	if len(schema.Required) != 1 || schema.Required[0] != "name" {
		t.Errorf("expected required=['name'], got %v", schema.Required)
	}
}

func TestPluginModuleToolsNoSchema(t *testing.T) {
	m := NewPluginModule(makeTestConfig())
	tools := m.Tools()

	// Second tool has no InputSchema
	schema := tools[1].Tool.InputSchema
	if schema.Type != "object" {
		t.Errorf("expected schema type 'object', got %q", schema.Type)
	}
	if len(schema.Properties) != 0 {
		t.Errorf("expected empty properties for tool without schema, got %v", schema.Properties)
	}
	if len(schema.Required) != 0 {
		t.Errorf("expected no required fields for tool without schema, got %v", schema.Required)
	}
}
