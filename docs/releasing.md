# Releasing LNS

LNS releases are built from immutable semantic-version tags. The first public release is `v0.1.0`; the executable reports `0.1.0`.

## Distribution channels

One tag drives every channel:

- GitHub Releases receives signed and notarized macOS binaries, Linux binaries, checksums, SBOMs, and provenance attestations.
- `crown-dev-studios/homebrew-tap` receives an automatically generated LNS cask. The cask declares the Homebrew `caddy` formula as a dependency.
- The public Go module makes `go install github.com/crown-dev-studios/lns/cmd/lns@v0.1.0` available from the same tag. Nothing is uploaded to a separate Go registry.

Homebrew is the recommended macOS path because it installs Caddy too. Direct archives and `go install` install only LNS; `lns doctor` verifies their local Caddy installation.

## One-time repository setup

Create the public `crown-dev-studios/homebrew-tap` repository with a `main` branch. Give the release workflow a fine-grained token named `HOMEBREW_TAP_TOKEN` with Contents read/write access to that repository only.

Add these Apple release secrets to the LNS repository:

- `MACOS_SIGN_P12`: base64-encoded Developer ID Application certificate.
- `MACOS_SIGN_PASSWORD`: password for the certificate.
- `MACOS_NOTARY_KEY`: base64-encoded App Store Connect API key.
- `MACOS_NOTARY_KEY_ID`: App Store Connect key ID.
- `MACOS_NOTARY_ISSUER_ID`: App Store Connect issuer UUID.

The release workflow fails before publishing if the tap token or any Apple credential is absent. It never falls back to an unsigned cask or removes macOS quarantine metadata.

Protect tags matching `v*` so only maintainers can create or delete them. Protect `main` with the CI workflow required.

## Before tagging

From a clean checkout of the commit intended for release:

```bash
go mod tidy
git diff --exit-code -- go.mod go.sum
go fmt ./...
git diff --exit-code
go test -race ./...
go vet ./...
goreleaser check
goreleaser release --snapshot --clean
```

The snapshot requires Syft because the release contains archive SBOMs. CI pins Go `1.26.1`, GoReleaser `v2.18.2`, and Syft `v1.52.0`; use those versions when reproducing CI locally. The normal CI matrix also tests the module's minimum Go 1.24 line.

Inspect the generated cask at `build/release/homebrew/Casks/lns.rb` and verify that it declares the `caddy` formula dependency. Run the host binary under `build/release` and confirm that `lns version` reports the snapshot version, commit, and commit date.

## Publishing `v0.1.0`

After the release commit is merged to `main`, create and push an annotated tag:

```bash
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

The tag starts `.github/workflows/release.yml`. That workflow:

1. Validates the tag and module files.
2. Runs the race-enabled test suite and vet.
3. Signs and notarizes macOS binaries.
4. Builds macOS and Linux archives for amd64 and arm64.
5. Generates checksums and SBOMs.
6. Publishes the GitHub release and updates the Homebrew tap.
7. Attests the archives and independently downloads the release to verify checksums.

Do not move or replace a published tag. If a release is wrong, fix the problem and publish the next patch version, such as `v0.1.1`.

## Post-release verification

Verify all three install paths from clean environments:

```bash
brew uninstall --cask lns 2>/dev/null || true
brew install --cask crown-dev-studios/tap/lns
lns version
lns doctor

go install github.com/crown-dev-studios/lns/cmd/lns@v0.1.0
lns version

gh release download v0.1.0 --repo crown-dev-studios/lns
shasum -a 256 --check checksums.txt
```

On macOS, also verify the released executable's Developer ID signature and notarization before announcing the release. The Homebrew installation must not start Caddy or create `~/.lns`; those happen only when the user runs LNS.
