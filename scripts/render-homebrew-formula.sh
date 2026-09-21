#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 5 ]]; then
  echo "usage: $0 VERSION SHA256 COMMIT BUILD_DATE OUTPUT" >&2
  exit 2
fi

version=$1
checksum=$2
commit=$3
build_date=$4
output_path=$5

if [[ ! "$version" =~ ^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]]; then
  echo "Homebrew formula version must be a stable semantic version, for example 0.1.0" >&2
  exit 1
fi
if [[ ! "$checksum" =~ ^[0-9a-f]{64}$ ]]; then
  echo "Homebrew formula checksum must be a lowercase SHA-256" >&2
  exit 1
fi
if [[ ! "$commit" =~ ^[0-9a-f]{40}$ ]]; then
  echo "Homebrew formula commit must be a full Git commit SHA" >&2
  exit 1
fi
if [[ ! "$build_date" =~ ^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(\.[0-9]+)?(Z|[+-][0-9]{2}:[0-9]{2})$ ]]; then
  echo "Homebrew formula build date must be RFC 3339" >&2
  exit 1
fi

script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
template_path="$script_dir/../packaging/homebrew/lns.rb.tmpl"
mkdir -p "$(dirname -- "$output_path")"

sed \
  -e "s|@@VERSION@@|$version|g" \
  -e "s|@@SHA256@@|$checksum|g" \
  -e "s|@@COMMIT@@|$commit|g" \
  -e "s|@@BUILD_DATE@@|$build_date|g" \
  "$template_path" > "$output_path"
