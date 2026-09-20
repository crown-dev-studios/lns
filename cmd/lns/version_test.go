package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCommandReportsReleaseMetadata(t *testing.T) {
	originalVersion, originalCommit, originalDate := version, commit, buildDate
	t.Cleanup(func() {
		version, commit, buildDate = originalVersion, originalCommit, originalDate
		versionCmd.SetOut(nil)
	})
	version = "v0.1.0"
	commit = "0123456789abcdef"
	buildDate = "2026-09-20T12:00:00Z"

	var output bytes.Buffer
	versionCmd.SetOut(&output)
	versionCmd.Run(versionCmd, nil)

	for _, want := range []string{"lns version 0.1.0", "commit 0123456789ab", "built 2026-09-20T12:00:00Z"} {
		if !strings.Contains(output.String(), want) {
			t.Fatalf("version output %q does not contain %q", output.String(), want)
		}
	}
}

func TestVersionCommandKeepsLocalBuildsAtDev(t *testing.T) {
	originalVersion, originalCommit, originalDate := version, commit, buildDate
	t.Cleanup(func() {
		version, commit, buildDate = originalVersion, originalCommit, originalDate
		versionCmd.SetOut(nil)
	})
	version, commit, buildDate = "dev", "", ""

	var output bytes.Buffer
	versionCmd.SetOut(&output)
	versionCmd.Run(versionCmd, nil)

	if !strings.HasPrefix(output.String(), "lns version dev") {
		t.Fatalf("local build must report dev, got %q", output.String())
	}
}

func TestTaggedModuleVersionRejectsLocalPseudoVersion(t *testing.T) {
	for _, candidate := range []string{
		"v0.0.0-20260829200420-d69d313d34c3+dirty",
		"v0.1.1-0.20260829200420-d69d313d34c3",
		"(devel)",
	} {
		if taggedModuleVersion(candidate) != "" {
			t.Fatalf("expected %q to remain a dev build", candidate)
		}
	}
	if got := taggedModuleVersion("v0.1.0"); got != "0.1.0" {
		t.Fatalf("expected tagged module version 0.1.0, got %q", got)
	}
}
