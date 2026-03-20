package mcpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/hairglasses-studio/mcpkit/handler"
	"github.com/hairglasses-studio/mcpkit/registry"
)

// PromptModule exposes the Anthropic prompt improver API via MCP.
type PromptModule struct {
	APIKey     string       // injected from resolveAPIKey()
	BaseURL    string       // default "https://api.anthropic.com", overridable for tests
	HTTPClient *http.Client // optional, defaults to 5min timeout
}

func (m *PromptModule) Name() string        { return "prompt" }
func (m *PromptModule) Description() string { return "Prompt improvement using the Anthropic API" }

func (m *PromptModule) httpClient() *http.Client {
	if m.HTTPClient != nil {
		return m.HTTPClient
	}
	return &http.Client{Timeout: 5 * time.Minute}
}

func (m *PromptModule) baseURL() string {
	if m.BaseURL != "" {
		return m.BaseURL
	}
	return "https://api.anthropic.com"
}

type PromptImproveInput struct {
	Prompt      string `json:"prompt" jsonschema:"description=The prompt to improve — may contain {{variable}} placeholders,required"`
	System      string `json:"system,omitempty" jsonschema:"description=System prompt to improve alongside the user prompt"`
	Feedback    string `json:"feedback,omitempty" jsonschema:"description=Guidance for the improver (e.g. add more examples)"`
	TargetModel string `json:"target_model,omitempty" jsonschema:"description=Model the prompt targets (default claude-sonnet-4-6)"`
}

type PromptImproveOutput struct {
	ImprovedPrompt string `json:"improved_prompt"`
	ImprovedSystem string `json:"improved_system,omitempty"`
	TargetModel    string `json:"target_model"`
}

// promptAPIRequest is the request body for /v1/experimental/improve_prompt.
type promptAPIRequest struct {
	Prompt      promptMessage  `json:"prompt"`
	System      string         `json:"system,omitempty"`
	Feedback    string         `json:"feedback,omitempty"`
	TargetModel string         `json:"target_model"`
}

type promptMessage struct {
	Messages []promptMsg `json:"messages"`
}

type promptMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// promptAPIResponse is the response from the improve_prompt endpoint.
type promptAPIResponse struct {
	Prompt promptAPIResponsePrompt `json:"prompt"`
	System string                  `json:"system,omitempty"`
	Error  *promptAPIError         `json:"error,omitempty"`
}

type promptAPIResponsePrompt struct {
	Messages []promptMsg `json:"messages"`
}

type promptAPIError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// promptRateLimitError identifies 429 responses for retry logic.
type promptRateLimitError struct {
	status int
	body   string
}

func (e *promptRateLimitError) Error() string {
	return fmt.Sprintf("prompt improver: rate limited (status %d): %s", e.status, e.body)
}

func isPromptRateLimitError(err error) bool {
	_, ok := err.(*promptRateLimitError)
	return ok
}

func (m *PromptModule) Tools() []registry.ToolDefinition {
	return []registry.ToolDefinition{
		handler.TypedHandler[PromptImproveInput, PromptImproveOutput](
			"prompt_improve",
			"Improve a prompt using the Anthropic prompt improver. Rewrites the prompt for clarity, specificity, and effectiveness.",
			func(ctx context.Context, input PromptImproveInput) (PromptImproveOutput, error) {
				if m.APIKey == "" {
					return PromptImproveOutput{}, fmt.Errorf("prompt improver: no API key configured — set ANTHROPIC_API_KEY")
				}

				targetModel := input.TargetModel
				if targetModel == "" {
					targetModel = "claude-sonnet-4-6"
				}

				body := promptAPIRequest{
					Prompt: promptMessage{
						Messages: []promptMsg{
							{Role: "user", Content: input.Prompt},
						},
					},
					TargetModel: targetModel,
				}
				if input.System != "" {
					body.System = input.System
				}
				if input.Feedback != "" {
					body.Feedback = input.Feedback
				}

				resp, err := m.doImproveRequest(ctx, body)
				if err != nil {
					return PromptImproveOutput{}, err
				}

				out := PromptImproveOutput{
					TargetModel: targetModel,
				}

				// Extract improved prompt from response messages.
				for _, msg := range resp.Prompt.Messages {
					if msg.Role == "user" {
						out.ImprovedPrompt = msg.Content
						break
					}
				}

				if resp.System != "" {
					out.ImprovedSystem = resp.System
				}

				return out, nil
			},
		),
	}
}

func (m *PromptModule) doImproveRequest(ctx context.Context, body promptAPIRequest) (*promptAPIResponse, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("prompt improver: marshal request: %w", err)
	}

	url := m.baseURL() + "/v1/experimental/improve_prompt"

	var resp *promptAPIResponse
	var lastErr error
	for attempt := 0; attempt < 4; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		resp, lastErr = m.doHTTP(ctx, url, payload)
		if lastErr == nil {
			break
		}
		if !isPromptRateLimitError(lastErr) {
			return nil, lastErr
		}
	}
	if lastErr != nil {
		return nil, fmt.Errorf("prompt improver: all retries exhausted: %w", lastErr)
	}

	return resp, nil
}

func (m *PromptModule) doHTTP(ctx context.Context, url string, payload []byte) (*promptAPIResponse, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("prompt improver: create request: %w", err)
	}
	httpReq.Header.Set("x-api-key", m.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	httpReq.Header.Set("anthropic-beta", "prompt-tools-2025-04-02")
	httpReq.Header.Set("content-type", "application/json")

	httpResp, err := m.httpClient().Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("prompt improver: http error: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("prompt improver: read response: %w", err)
	}

	if httpResp.StatusCode == http.StatusTooManyRequests {
		return nil, &promptRateLimitError{status: httpResp.StatusCode, body: string(respBody)}
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prompt improver: API error (status %d): %s", httpResp.StatusCode, string(respBody))
	}

	var resp promptAPIResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("prompt improver: unmarshal response: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("prompt improver: API error: %s: %s", resp.Error.Type, resp.Error.Message)
	}

	return &resp, nil
}
