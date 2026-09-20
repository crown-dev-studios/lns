package main

import (
	"context"
	"strings"
	"testing"
)

func TestRequireCaddyGivesInstallationGuidance(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	_, err := requireCaddy(context.Background())
	if err == nil || !strings.Contains(err.Error(), "Caddy v2 is required") {
		t.Fatalf("expected actionable Caddy guidance, got %v", err)
	}
}
