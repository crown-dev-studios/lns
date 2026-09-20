package caddy

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestProbeReportsMissingCaddy(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	_, err := Probe(context.Background())
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestProbeReportsUsableCaddyVersion(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "caddy")
	if err := os.WriteFile(path, []byte("#!/bin/sh\necho 'v2.10.2 h1:example'\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)

	installation, err := Probe(context.Background())
	if err != nil {
		t.Fatalf("probe usable Caddy: %v", err)
	}
	if installation.Path != path || installation.Version != "v2.10.2" {
		t.Fatalf("unexpected installation: %#v", installation)
	}
}

func TestProbeReportsCaddyThatCannotRun(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "caddy")
	if err := os.WriteFile(path, []byte("#!/bin/sh\necho 'broken installation' >&2\nexit 1\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)

	_, err := Probe(context.Background())
	if !errors.Is(err, ErrUnavailable) {
		t.Fatalf("expected ErrUnavailable, got %v", err)
	}
}

func TestProbeRejectsIncompatibleCaddy(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "caddy")
	if err := os.WriteFile(path, []byte("#!/bin/sh\necho 'v1.0.5'\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)

	_, err := Probe(context.Background())
	if !errors.Is(err, ErrIncompatible) {
		t.Fatalf("expected ErrIncompatible, got %v", err)
	}
}

func TestProbeRejectsCaddyWithoutRequiredAdapter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "caddy")
	script := "#!/bin/sh\nif [ \"$1\" = version ]; then echo 'v2.10.2'; exit 0; fi\necho 'adapt unavailable' >&2\nexit 1\n"
	if err := os.WriteFile(path, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)

	_, err := Probe(context.Background())
	if !errors.Is(err, ErrIncompatible) {
		t.Fatalf("expected ErrIncompatible, got %v", err)
	}
}
