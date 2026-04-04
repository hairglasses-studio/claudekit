package pluginkit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

// SubprocessHandler calls a plugin by spawning a subprocess.
type SubprocessHandler struct {
	Command string
	Timeout time.Duration
}

type subprocessRequest struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type subprocessResponse struct {
	Result json.RawMessage `json:"result,omitempty"`
	Error  string          `json:"error,omitempty"`
}

// Call invokes a tool by writing a JSON request to the subprocess stdin
// and reading a JSON response from stdout.
func (h *SubprocessHandler) Call(ctx context.Context, toolName string, input json.RawMessage) (json.RawMessage, error) {
	timeout := h.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req := subprocessRequest{
		Method: toolName,
		Params: input,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	reqBytes = append(reqBytes, '\n')

	// Trust boundary: plugins are user-installed from local YAML files.
	// The command is executed via shell to support pipes and redirects.
	// LoadPlugin warns on shell metacharacters for auditability.
	cmd := exec.CommandContext(ctx, "sh", "-c", h.Command)
	cmd.Stdin = bytes.NewReader(reqBytes)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("plugin timeout after %s", timeout)
		}
		return nil, fmt.Errorf("plugin exec: %w: %s", err, stderr.String())
	}

	var resp subprocessResponse
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		raw := stdout.String()
		if len(raw) > 256 {
			raw = raw[:256] + "..."
		}
		return nil, fmt.Errorf("parse plugin response: %w (raw: %s)", err, raw)
	}

	if resp.Error != "" {
		return nil, fmt.Errorf("plugin error: %s", resp.Error)
	}

	return resp.Result, nil
}
