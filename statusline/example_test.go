package statusline_test

import (
	"fmt"
	"strings"

	"github.com/hairglasses-studio/claudekit/statusline"
)

func ExampleRender() {
	input := `{"model":"claude-sonnet-4-20250514","working_directory":"/tmp/demo","session_id":"abc","cost_usd":0.42,"total_tokens":5000,"max_tokens":200000,"duration_ms":12000}`
	output := statusline.Render(strings.NewReader(input))
	// Output contains the model name and cost.
	// Exact output depends on font tier and theme env vars.
	if strings.Contains(output, "Sonnet") {
		fmt.Println("contains model name")
	}
	if strings.Contains(output, "0.42") {
		fmt.Println("contains cost")
	}
	// Output:
	// contains model name
	// contains cost
}

func ExampleScriptContent() {
	script := statusline.ScriptContent()
	// The script starts with a shebang line.
	lines := strings.SplitN(script, "\n", 2)
	fmt.Println(lines[0])
	// Output:
	// #!/bin/bash
}
