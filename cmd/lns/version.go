package main

import (
	"fmt"
	"regexp"
	"runtime/debug"
	"strings"
)

var (
	version   = "dev"
	commit    = ""
	buildDate = ""

	taggedVersionPattern = regexp.MustCompile(`^v?(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-[0-9A-Za-z]+([.-][0-9A-Za-z]+)*)?$`)
	pseudoVersionPattern = regexp.MustCompile(`[.-][0-9]{14}-[0-9a-f]{12}(\+.*)?$`)
)

func versionString() string {
	resolvedVersion, resolvedCommit, resolvedDate := version, commit, buildDate
	if info, ok := debug.ReadBuildInfo(); ok {
		if resolvedVersion == "dev" {
			if tagged := taggedModuleVersion(info.Main.Version); tagged != "" {
				resolvedVersion = tagged
			}
		}
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				if resolvedCommit == "" {
					resolvedCommit = setting.Value
				}
			case "vcs.time":
				if resolvedDate == "" {
					resolvedDate = setting.Value
				}
			}
		}
	}

	resolvedVersion = strings.TrimPrefix(resolvedVersion, "v")
	metadata := make([]string, 0, 2)
	if resolvedCommit != "" {
		if len(resolvedCommit) > 12 {
			resolvedCommit = resolvedCommit[:12]
		}
		metadata = append(metadata, "commit "+resolvedCommit)
	}
	if resolvedDate != "" {
		metadata = append(metadata, "built "+resolvedDate)
	}
	if len(metadata) == 0 {
		return "lns version " + resolvedVersion
	}
	return fmt.Sprintf("lns version %s (%s)", resolvedVersion, strings.Join(metadata, ", "))
}

func taggedModuleVersion(candidate string) string {
	if !taggedVersionPattern.MatchString(candidate) || pseudoVersionPattern.MatchString(candidate) {
		return ""
	}
	return strings.TrimPrefix(candidate, "v")
}
