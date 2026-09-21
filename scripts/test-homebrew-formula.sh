#!/usr/bin/env bash
set -euo pipefail

script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT

formula="$tmp_dir/Formula/lns.rb"
checksums="$tmp_dir/checksums.txt"

cat > "$checksums" <<'EOF'
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa  lns_0.1.0_darwin_arm64.tar.gz
bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb  lns_0.1.0_darwin_amd64.tar.gz
cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc  lns_0.1.0_linux_arm64.tar.gz
dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd  lns_0.1.0_linux_amd64.tar.gz
EOF

bash "$script_dir/render-homebrew-formula.sh" \
  0.1.0 \
  "$checksums" \
  "$formula"

ruby -c "$formula" >/dev/null
grep -Fq 'url "https://github.com/crown-dev-studios/lns/releases/download/v0.1.0/lns_0.1.0_darwin_arm64.tar.gz"' "$formula"
grep -Fq 'sha256 "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"' "$formula"
grep -Fq 'url "https://github.com/crown-dev-studios/lns/releases/download/v0.1.0/lns_0.1.0_darwin_amd64.tar.gz"' "$formula"
grep -Fq 'sha256 "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"' "$formula"
grep -Fq 'url "https://github.com/crown-dev-studios/lns/releases/download/v0.1.0/lns_0.1.0_linux_arm64.tar.gz"' "$formula"
grep -Fq 'sha256 "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"' "$formula"
grep -Fq 'url "https://github.com/crown-dev-studios/lns/releases/download/v0.1.0/lns_0.1.0_linux_amd64.tar.gz"' "$formula"
grep -Fq 'sha256 "dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"' "$formula"
if grep -Fq 'depends_on "go"' "$formula"; then
  echo "Generated formula still depends on Go" >&2
  exit 1
fi
grep -Fq 'depends_on "caddy"' "$formula"
grep -Fq 'bin.install "lns"' "$formula"

if bash "$script_dir/render-homebrew-formula.sh" \
  0.1.0-rc1 \
  "$checksums" \
  "$tmp_dir/prerelease.rb" >/dev/null 2>&1; then
  echo "Formula renderer accepted a prerelease version" >&2
  exit 1
fi

sed '/linux_amd64/d' "$checksums" > "$tmp_dir/incomplete-checksums.txt"
if bash "$script_dir/render-homebrew-formula.sh" \
  0.1.0 \
  "$tmp_dir/incomplete-checksums.txt" \
  "$tmp_dir/incomplete.rb" >/dev/null 2>&1; then
  echo "Formula renderer accepted a missing archive checksum" >&2
  exit 1
fi

sed 's/^aaaa/zzzz/' "$checksums" > "$tmp_dir/invalid-checksums.txt"
if bash "$script_dir/render-homebrew-formula.sh" \
  0.1.0 \
  "$tmp_dir/invalid-checksums.txt" \
  "$tmp_dir/invalid.rb" >/dev/null 2>&1; then
  echo "Formula renderer accepted an invalid checksum" >&2
  exit 1
fi
