package main

import (
	"context"
	"os"
	"path/filepath"
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

func TestDoctorReportsCaddyProblems(t *testing.T) {
	tests := []struct {
		name       string
		executable string
		want       string
	}{
		{name: "missing", want: "Caddy v2 is required"},
		{name: "unavailable", executable: "#!/bin/sh\necho 'broken installation' >&2\nexit 1\n", want: "Caddy v2 check failed: caddy executable is unavailable"},
		{name: "incompatible", executable: "#!/bin/sh\necho 'v1.0.5'\n", want: "Caddy v2 check failed: incompatible caddy executable"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bin := t.TempDir()
			if test.executable != "" {
				path := filepath.Join(bin, "caddy")
				if err := os.WriteFile(path, []byte(test.executable), 0755); err != nil {
					t.Fatal(err)
				}
			}
			t.Setenv("PATH", bin)

			err := runDoctor()
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("expected doctor error containing %q, got %v", test.want, err)
			}
		})
	}
}

func TestDoctorAcceptsUsableCaddy(t *testing.T) {
	executable := writeExecutable(t, "#!/bin/sh\nif [ \"$1\" = version ]; then echo 'v2.10.2'; fi\n")
	t.Setenv("PATH", filepath.Dir(executable))

	if err := runDoctor(); err != nil {
		t.Fatalf("expected doctor to accept usable Caddy, got %v", err)
	}
}

func runDoctor() error {
	command := *doctorCmd
	command.SetContext(context.Background())
	return command.RunE(&command, nil)
}
