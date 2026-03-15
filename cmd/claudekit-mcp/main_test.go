package main

import (
	"testing"

	"github.com/hairglasses-studio/mcpkit/rdcycle"
)

func TestResolveProfileDefault(t *testing.T) {
	// With no env or flag, should return personal.
	t.Setenv("CLAUDEKIT_BUDGET_PROFILE", "")
	p := resolveProfile()
	expected := rdcycle.PersonalProfile()
	if p.Name != expected.Name {
		t.Errorf("default profile name = %q, want %q", p.Name, expected.Name)
	}
}

func TestResolveProfilePersonalExplicit(t *testing.T) {
	t.Setenv("CLAUDEKIT_BUDGET_PROFILE", "personal")
	p := resolveProfile()
	expected := rdcycle.PersonalProfile()
	if p.Name != expected.Name {
		t.Errorf("profile name = %q, want %q", p.Name, expected.Name)
	}
}

func TestResolveProfileWorkAPI(t *testing.T) {
	t.Setenv("CLAUDEKIT_BUDGET_PROFILE", "work-api")
	p := resolveProfile()
	expected := rdcycle.WorkAPIProfile()
	if p.Name != expected.Name {
		t.Errorf("profile name = %q, want %q", p.Name, expected.Name)
	}
	if p.DollarBudget != expected.DollarBudget {
		t.Errorf("DollarBudget = %.2f, want %.2f", p.DollarBudget, expected.DollarBudget)
	}
}

func TestResolveProfileInvalidFallsBack(t *testing.T) {
	t.Setenv("CLAUDEKIT_BUDGET_PROFILE", "/nonexistent/profile.json")
	p := resolveProfile()
	// Should fall back to personal.
	expected := rdcycle.PersonalProfile()
	if p.Name != expected.Name {
		t.Errorf("fallback profile name = %q, want %q", p.Name, expected.Name)
	}
}
