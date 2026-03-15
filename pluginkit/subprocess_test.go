package pluginkit

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestSubprocessHandlerCall(t *testing.T) {
	h := &SubprocessHandler{
		Command: `echo '{"result":{"greeting":"hello"}}'`,
		Timeout: 5 * time.Second,
	}

	result, err := h.Call(context.Background(), "greet", json.RawMessage(`{"name":"world"}`))
	if err != nil {
		t.Fatal(err)
	}

	var out map[string]string
	if err := json.Unmarshal(result, &out); err != nil {
		t.Fatal(err)
	}
	if out["greeting"] != "hello" {
		t.Errorf("greeting = %q, want hello", out["greeting"])
	}
}

func TestSubprocessHandlerError(t *testing.T) {
	h := &SubprocessHandler{
		Command: `echo '{"error":"something went wrong"}'`,
		Timeout: 5 * time.Second,
	}

	_, err := h.Call(context.Background(), "test", json.RawMessage(`{}`))
	if err == nil {
		t.Error("expected error")
	}
}

func TestSubprocessHandlerTimeout(t *testing.T) {
	h := &SubprocessHandler{
		Command: "sleep 10",
		Timeout: 100 * time.Millisecond,
	}

	_, err := h.Call(context.Background(), "test", json.RawMessage(`{}`))
	if err == nil {
		t.Error("expected timeout error")
	}
}

func TestSubprocessHandlerMalformed(t *testing.T) {
	h := &SubprocessHandler{
		Command: "echo 'not json'",
		Timeout: 5 * time.Second,
	}

	_, err := h.Call(context.Background(), "test", json.RawMessage(`{}`))
	if err == nil {
		t.Error("expected parse error")
	}
}
