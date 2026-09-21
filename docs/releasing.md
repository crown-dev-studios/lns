# Releasing LNS

LNS releases are built from immutable semantic-version tags. The first public release is `v0.1.0`; the executable reports `0.1.0`.

## Distribution channels

One tag drives every channel:

- GitHub Releases receives macOS and Linux binaries, checksums, SBOMs, and provenance attestations.
- `crown-dev-studios/homebrew-tap` receives an automatically generated source formula. Homebrew downloads the exact tagged commit, verifies its checksum, builds LNS, and installs the `caddy` formula as a runtime dependency.
- The public Go module makes `go install github.com/crown-dev-studios/lns/cmd/lns@v0.1.0` available from the same tag. Nothing is uploaded to a separate Go registry.

Homebrew is the recommended macOS path because it installs Caddy too. Direct archives and `go install` install only LNS; `lns doctor` verifies their local Caddy installation.

## One-time repository setup

Create the public `crown-dev-studios/homebrew-tap` repository with a `main` branch. Give the release workflow a fine-grained token named `HOMEBREW_TAP_TOKEN` with Contents read/write access to that repository only. This is the only release secret. Stable releases require it; release candidates publish GitHub prereleases without updating Homebrew.

The Homebrew path builds from source, so it does not require an Apple Developer certificate or notarization credentials. The release workflow never removes macOS quarantine metadata.

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
bash scripts/test-homebrew-formula.sh
```

The snapshot requires Syft because the release contains archive SBOMs. CI pins Go `1.26.1`, GoReleaser `v2.18.2`, and Syft `v1.52.0`; use those versions when reproducing CI locally. The normal CI matrix also tests the module's minimum Go 1.24 line.

Run the host binary under `build/release` and confirm that `lns version` reports the snapshot version, commit, and commit date. The formula renderer test proves that stable versions produce valid Ruby with the pinned source checksum, release metadata, Go build dependency, and Caddy runtime dependency.

## Publishing `v0.1.0`

After the release commit is merged to `main`, create and push an annotated tag:

```bash
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

The tag starts `.github/workflows/release.yml`. That workflow:

1. Validates the tag.
2. Runs the same macOS/Linux and Go 1.24/1.26 verification workflow required by normal CI.
3. Builds all four snapshot archives on matching native runners, extracts each archive, and smoke-tests `lns version` plus a fixture-backed `lns plan`.
4. Builds macOS and Linux release archives for amd64 and arm64.
5. Generates checksums and SBOMs.
6. Publishes and attests the GitHub release archives.
7. For a stable tag, hashes the exact release commit, renders the Homebrew formula, and updates the tap.
8. Independently verifies the published checksums and installs the stable formula on a clean macOS runner.

The Homebrew verification builds LNS from source, confirms that Caddy was installed, runs the formula test and `lns doctor`, and proves installation created no LNS runtime state. A tag such as `v0.1.0-rc1` skips the tap and this stable-channel check.

Do not move or replace a published tag. If a release is wrong, fix the problem and publish the next patch version, such as `v0.1.1`.

## Post-release verification

Verify all three install paths from clean environments:

```bash
brew uninstall lns 2>/dev/null || true
brew install crown-dev-studios/tap/lns
lns version
lns doctor

go install github.com/crown-dev-studios/lns/cmd/lns@v0.1.0
lns version

gh release download v0.1.0 --repo crown-dev-studios/lns
shasum -a 256 --check checksums.txt
```

The Homebrew installation must not start Caddy or create `~/.lns`; those happen only when the user runs LNS.
