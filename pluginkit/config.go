package pluginkit

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// PluginConfig represents a YAML-defined plugin.
type PluginConfig struct {
	Name        string          `yaml:"name"`
	Description string          `yaml:"description"`
	Version     string          `yaml:"version"`
	Handler     HandlerConfig   `yaml:"handler"`
	Tools       []PluginToolDef `yaml:"tools"`
}

// HandlerConfig describes how to invoke the plugin.
type HandlerConfig struct {
	Type    string `yaml:"type"` // "subprocess"
	Command string `yaml:"command"`
	Timeout string `yaml:"timeout"` // e.g. "30s"
}

// PluginToolDef describes a tool exposed by the plugin.
type PluginToolDef struct {
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	InputSchema map[string]interface{} `yaml:"input_schema"`
}

// LoadPlugin reads and parses a plugin config from a YAML file.
func LoadPlugin(path string) (*PluginConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read plugin %s: %w", path, err)
	}

	var cfg PluginConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse plugin %s: %w", path, err)
	}

	if cfg.Name == "" {
		return nil, fmt.Errorf("plugin %s: name is required", path)
	}

	return &cfg, nil
}

// LoadPlugins scans a directory for .yaml/.yml plugin files and loads them all.
func LoadPlugins(dir string) ([]*PluginConfig, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read plugin dir %s: %w", dir, err)
	}

	var plugins []*PluginConfig
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := filepath.Ext(entry.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}
		cfg, err := LoadPlugin(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		plugins = append(plugins, cfg)
	}
	return plugins, nil
}

// DefaultPluginDir returns the default plugin directory (~/.claudekit/plugins/).
func DefaultPluginDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".claudekit", "plugins"), nil
}
