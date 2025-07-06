#!/usr/bin/env bash
set -euo pipefail

# Extract existing MAJOR.MINOR.PATCH from const.go
current=$(grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' const.go)
version="${current}-$(date -u +%y%m%d%H%M)"
# Replace the ClientVersion constant in const.go
sed -i -E "s/ClientVersion\s*=\s*\"[^\"]+\"/ClientVersion    = \"$version\"/" const.go

echo "Updated ClientVersion to $version"
