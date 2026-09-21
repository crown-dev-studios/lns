#!/usr/bin/env bash
set -euo pipefail

script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT

checksum=0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
commit=0123456789abcdef0123456789abcdef01234567
build_date=2026-09-20T12:34:56Z
formula="$tmp_dir/Formula/lns.rb"

bash "$script_dir/render-homebrew-formula.sh" \
  0.1.0 \
  "$checksum" \
  "$commit" \
  "$build_date" \
  "$formula"

ruby -c "$formula" >/dev/null
grep -Fq "url \"https://github.com/crown-dev-studios/lns/archive/$commit.tar.gz\"" "$formula"
grep -Fq 'version "0.1.0"' "$formula"
grep -Fq "sha256 \"$checksum\"" "$formula"
grep -Fq 'depends_on "go" => :build' "$formula"
grep -Fq 'depends_on "caddy"' "$formula"
grep -Fq -- "-X main.commit=$commit" "$formula"
grep -Fq -- "-X main.buildDate=$build_date" "$formula"

if bash "$script_dir/render-homebrew-formula.sh" \
  0.1.0-rc1 \
  "$checksum" \
  "$commit" \
  "$build_date" \
  "$tmp_dir/prerelease.rb" >/dev/null 2>&1; then
  echo "Formula renderer accepted a prerelease version" >&2
  exit 1
fi

if bash "$script_dir/render-homebrew-formula.sh" \
  0.1.0 \
  invalid \
  "$commit" \
  "$build_date" \
  "$tmp_dir/invalid.rb" >/dev/null 2>&1; then
  echo "Formula renderer accepted an invalid checksum" >&2
  exit 1
fi
