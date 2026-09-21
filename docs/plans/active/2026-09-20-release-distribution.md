---
title: Release and Distribution Workflow
date: 2026-09-20
status: complete
---

# Release and Distribution Workflow

## Outcome

One semantic-version tag publishes LNS through GitHub Releases, Homebrew, and the public Go module.

GoReleaser builds and publishes versioned archives for macOS and Linux on amd64 and arm64. Each binary reports its version, commit, and build date.

The Homebrew formula installs the matching archive and declares Caddy as its runtime dependency. Direct archive and `go install` users install Caddy separately and verify it with `lns doctor`.

## Release contract

- Git tags are the release version source.
- Published tags are immutable; corrections use a new patch version.
- Every archive has a checksum, SBOM, and build attestation.
- Stable releases update `crown-dev-studios/homebrew-tap`.
- Release candidates publish GitHub prereleases without updating Homebrew.
- Installation never starts Caddy, creates routes, edits a project, or changes Docker configuration.
- Windows remains unsupported until its runtime behavior has a defined contract.

## Components

- `.goreleaser.yaml` defines targets, archive names, checksums, SBOMs, changelogs, and build metadata.
- `packaging/homebrew/lns.rb.tmpl` defines the Homebrew formula.
- `scripts/render-homebrew-formula.sh` fills the formula with release archive checksums.
- `.github/workflows/verify.yml` tests supported Go versions and release archives.
- `.github/workflows/release.yml` publishes tagged releases and updates the Homebrew tap.
- `docs/releasing.md` is the maintainer runbook.

## Release workflow

1. Validate the semantic-version tag.
2. Run the normal verification workflow.
3. Build macOS and Linux archives for amd64 and arm64.
4. Generate checksums and SBOMs.
5. Publish and attest the GitHub release archives.
6. For a stable release, render and publish the Homebrew formula.
7. Verify published checksums and the Homebrew installation.

## Verification

- `go test -race ./...`
- `go vet ./...`
- `goreleaser check`
- `goreleaser release --snapshot --clean`
- `bash scripts/test-homebrew-formula.sh`
- Run `lns version` from an extracted archive.
- Install the published formula on a clean macOS runner and run `brew test` and `lns doctor`.
