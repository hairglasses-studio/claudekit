package mcpserver

import (
	"testing"
)

func TestParseUpstream(t *testing.T) {
	tests := []struct {
		spec    string
		name    string
		url     string
		wantErr bool
	}{
		{"github=http://localhost:8080", "github", "http://localhost:8080", false},
		{"slack=https://mcp.slack.com/v1", "slack", "https://mcp.slack.com/v1", false},
		{"name=url=with=equals", "name", "url=with=equals", false},
		{"noequals", "", "", true},
		{"=noname", "", "", true},
		{"nourl=", "", "", true},
		{"", "", "", true},
	}

	for _, tt := range tests {
		name, url, err := parseUpstream(tt.spec)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseUpstream(%q): err=%v, wantErr=%v", tt.spec, err, tt.wantErr)
			continue
		}
		if name != tt.name {
			t.Errorf("parseUpstream(%q): name=%q, want %q", tt.spec, name, tt.name)
		}
		if url != tt.url {
			t.Errorf("parseUpstream(%q): url=%q, want %q", tt.spec, url, tt.url)
		}
	}
}

func TestSetupGatewayEmptyUpstreams(t *testing.T) {
	// Empty upstreams should succeed and return a valid gateway.
	gw, dynReg, err := SetupGateway(nil, nil)
	if err != nil {
		t.Fatalf("SetupGateway with empty upstreams: %v", err)
	}
	if gw == nil {
		t.Fatal("expected non-nil gateway")
	}
	if dynReg == nil {
		t.Fatal("expected non-nil dynamic registry")
	}
	defer func() { _ = gw.Close() }()

	upstreams := gw.ListUpstreams()
	if len(upstreams) != 0 {
		t.Errorf("expected 0 upstreams, got %d", len(upstreams))
	}
}

func TestSetupGatewayInvalidFormat(t *testing.T) {
	_, _, err := SetupGateway(nil, []string{"bad-format"})
	if err == nil {
		t.Fatal("expected error for invalid upstream format")
	}
}
