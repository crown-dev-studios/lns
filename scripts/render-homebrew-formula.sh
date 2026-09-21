#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 3 ]]; then
  echo "usage: $0 VERSION CHECKSUMS_FILE OUTPUT" >&2
  exit 2
fi

version=$1
checksums_path=$2
output_path=$3

if [[ ! "$version" =~ ^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]]; then
  echo "Homebrew formula version must be a stable semantic version, for example 0.1.0" >&2
  exit 1
fi
if [[ ! -f "$checksums_path" ]]; then
  echo "Homebrew checksums file does not exist: $checksums_path" >&2
  exit 1
fi

checksum_for() {
  local archive=$1
  local checksum
  local matches

  read -r matches checksum < <(
    awk -v archive="$archive" \
      '$2 == archive { count++; checksum = $1 } END { print count + 0, checksum }' \
      "$checksums_path"
  )
  if [[ "$matches" -ne 1 ]]; then
    echo "Expected one checksum for $archive, found $matches" >&2
    exit 1
  fi

  if [[ ! "$checksum" =~ ^[0-9a-f]{64}$ ]]; then
    echo "Checksum for $archive must be a lowercase SHA-256" >&2
    exit 1
  fi

  printf '%s' "$checksum"
}

darwin_arm64_checksum=$(checksum_for "lns_${version}_darwin_arm64.tar.gz")
darwin_amd64_checksum=$(checksum_for "lns_${version}_darwin_amd64.tar.gz")
linux_arm64_checksum=$(checksum_for "lns_${version}_linux_arm64.tar.gz")
linux_amd64_checksum=$(checksum_for "lns_${version}_linux_amd64.tar.gz")

script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
template_path="$script_dir/../packaging/homebrew/lns.rb.tmpl"
mkdir -p "$(dirname -- "$output_path")"

sed \
  -e "s|@@VERSION@@|$version|g" \
  -e "s|@@DARWIN_ARM64_SHA256@@|$darwin_arm64_checksum|g" \
  -e "s|@@DARWIN_AMD64_SHA256@@|$darwin_amd64_checksum|g" \
  -e "s|@@LINUX_ARM64_SHA256@@|$linux_arm64_checksum|g" \
  -e "s|@@LINUX_AMD64_SHA256@@|$linux_amd64_checksum|g" \
  "$template_path" > "$output_path"
