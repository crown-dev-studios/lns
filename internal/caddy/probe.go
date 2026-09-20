package caddy

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var (
	ErrNotFound     = errors.New("caddy executable not found")
	ErrUnavailable  = errors.New("caddy executable is unavailable")
	ErrIncompatible = errors.New("incompatible caddy executable")
)

type Installation struct {
	Path    string
	Version string
}

func Probe(ctx context.Context) (Installation, error) {
	path, err := exec.LookPath("caddy")
	if err != nil {
		return Installation{}, fmt.Errorf("%w in PATH", ErrNotFound)
	}
	output, err := exec.CommandContext(ctx, path, "version").CombinedOutput()
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return Installation{}, ctxErr
		}
		detail := strings.TrimSpace(string(output))
		if detail == "" {
			detail = err.Error()
		}
		return Installation{}, fmt.Errorf("%w: %s: %s", ErrUnavailable, path, detail)
	}
	fields := strings.Fields(string(output))
	version := ""
	if len(fields) > 0 {
		version = fields[0]
	}
	if !strings.HasPrefix(version, "v2.") {
		return Installation{}, fmt.Errorf("%w: %s reports %q; LNS requires Caddy v2", ErrIncompatible, path, version)
	}
	adapt := exec.CommandContext(ctx, path, "adapt", "--config", "-")
	adapt.Stdin = strings.NewReader("{\n\tadmin off\n}\n")
	if output, err := adapt.CombinedOutput(); err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return Installation{}, ctxErr
		}
		detail := strings.TrimSpace(string(output))
		if detail == "" {
			detail = err.Error()
		}
		return Installation{}, fmt.Errorf("%w: %s cannot adapt Caddyfiles: %s", ErrIncompatible, path, detail)
	}
	return Installation{Path: path, Version: version}, nil
}
