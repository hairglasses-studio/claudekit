// Package pluginkit provides a YAML-based plugin system for extending
// claudekit with external tools.
//
// Plugins are defined as YAML files in ~/.claudekit/plugins/ and are
// loaded into mcpkit [registry.ToolModule] instances. Each plugin
// specifies a subprocess command that is invoked via JSON-over-stdin/stdout
// for each tool call.
//
// # Plugin Configuration
//
// A plugin YAML file declares the plugin name, handler command, and a
// list of tools with their input schemas:
//
//	name: my-plugin
//	description: Example plugin
//	handler:
//	  type: subprocess
//	  command: my-plugin-binary
//	  timeout: 30s
//	tools:
//	  - name: my_tool
//	    description: Does something
//	    input_schema:
//	      properties:
//	        query: { type: string }
//
// # Loading Plugins
//
// [LoadPlugin] reads a single YAML file. [LoadPlugins] scans a directory
// for all .yaml/.yml files. [DefaultPluginDir] returns ~/.claudekit/plugins/.
//
// # Module Adapter
//
// [NewPluginModule] adapts a [PluginConfig] into an mcpkit ToolModule
// that can be registered alongside built-in claudekit modules.
package pluginkit
