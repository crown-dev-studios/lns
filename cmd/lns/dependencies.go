package main

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	"github.com/crown-dev-studios/lns/internal/caddy"
)

func requireCaddy(ctx context.Context) (caddy.Installation, error) {
	installation, err := caddy.Probe(ctx)
	if err == nil {
		return installation, nil
	}
	if errors.Is(err, caddy.ErrNotFound) {
		if runtime.GOOS == "darwin" {
			return caddy.Installation{}, fmt.Errorf("Caddy v2 is required; install it with `brew install caddy`")
		}
		return caddy.Installation{}, fmt.Errorf("Caddy v2 is required; install it from https://caddyserver.com/docs/install")
	}
	return caddy.Installation{}, fmt.Errorf("Caddy v2 check failed: %w", err)
}
