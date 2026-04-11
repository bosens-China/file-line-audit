#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
output_dir="$repo_root/.agents/skills/file-line-audit/scripts"

mkdir -p "$output_dir"

build_target() {
  local goos="$1"
  local goarch="$2"
  local output="$3"
  GOOS="$goos" GOARCH="$goarch" go build -o "$output_dir/$output" ./cmd/line-audit
}

cd "$repo_root"

build_target windows amd64 line-audit-windows-amd64.exe
build_target linux amd64 line-audit-linux-amd64
build_target linux arm64 line-audit-linux-arm64
build_target darwin amd64 line-audit-darwin-amd64
build_target darwin arm64 line-audit-darwin-arm64
