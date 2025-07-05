#!/usr/bin/env bash
set -euo pipefail

version="v0.0.0-$(date -u +%y%m%d%H%M)"
# Replace the ClientVersion constant in const.go
sed -i -E "s/ClientVersion\s*=\s*\"[^\"]+\"/ClientVersion    = \"$version\"/" const.go

echo "Updated ClientVersion to $version"
