package pluginkit

import (
	"os"
	"path/filepath"
	"testing"
)

const samplePlugin = `name: hello-world
description: A test plugin
version: "1.0.0"
handler:
  type: subprocess
  command: echo
  timeout: "10s"
tools:
  - name: hello_greet
    description: Greet someone
    input_schema:
      type: object
      properties:
        name:
          type: string
          description: Name to greet
`

func TestLoadPlugin(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "hello.yaml")
	os.WriteFile(path, []byte(samplePlugin), 0o644)

	cfg, err := LoadPlugin(path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Name != "hello-world" {
		t.Errorf("name = %q, want hello-world", cfg.Name)
	}
	if cfg.Version != "1.0.0" {
		t.Errorf("version = %q, want 1.0.0", cfg.Version)
	}
	if cfg.Handler.Type != "subprocess" {
		t.Errorf("handler.type = %q, want subprocess", cfg.Handler.Type)
	}
	if cfg.Handler.Command != "echo" {
		t.Errorf("handler.command = %q, want echo", cfg.Handler.Command)
	}
	if len(cfg.Tools) != 1 {
		t.Fatalf("tools count = %d, want 1", len(cfg.Tools))
	}
	if cfg.Tools[0].Name != "hello_greet" {
		t.Errorf("tool name = %q, want hello_greet", cfg.Tools[0].Name)
	}
}

func TestLoadPluginMissingName(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "bad.yaml")
	os.WriteFile(path, []byte("description: no name\n"), 0o644)

	_, err := LoadPlugin(path)
	if err == nil {
		t.Error("expected error for missing name")
	}
}

func TestLoadPlugins(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "a.yaml"), []byte(samplePlugin), 0o644)
	os.WriteFile(filepath.Join(tmpDir, "b.yml"), []byte("name: second\nversion: \"0.1\"\n"), 0o644)
	os.WriteFile(filepath.Join(tmpDir, "ignore.txt"), []byte("not a plugin"), 0o644)

	plugins, err := LoadPlugins(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(plugins) != 2 {
		t.Errorf("plugin count = %d, want 2", len(plugins))
	}
}

func TestLoadPluginsEmptyDir(t *testing.T) {
	tmpDir := t.TempDir()

	plugins, err := LoadPlugins(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(plugins) != 0 {
		t.Errorf("expected 0 plugins, got %d", len(plugins))
	}
}

func TestLoadPluginsNonexistent(t *testing.T) {
	plugins, err := LoadPlugins("/tmp/nonexistent-dir-12345")
	if err != nil {
		t.Fatal(err)
	}
	if plugins != nil {
		t.Error("expected nil for nonexistent dir")
	}
}

func TestDefaultPluginDir(t *testing.T) {
	dir, err := DefaultPluginDir()
	if err != nil {
		t.Fatal(err)
	}
	if dir == "" {
		t.Error("plugin dir should not be empty")
	}
}
