#!/bin/bash
set -euo pipefail

ASSET_DIR="$(dirname "$0")/../assets"
TEMP_DIR=$(mktemp -d)

if [ -n "$(ls -A "$ASSET_DIR" 2>/dev/null | grep -v '.gitkeep' || true)" ]; then
    echo "Assets already exist in $ASSET_DIR"
    exit 0
fi

echo "Cloning oni-seed-browser repository..."
 git clone --depth 1 https://github.com/MapsNotIncluded/oni-seed-browser "$TEMP_DIR"

echo "Copying assets..."
mkdir -p "$ASSET_DIR"
cp "$TEMP_DIR"/app/src/commonMain/composeResources/drawable/* "$ASSET_DIR"/

echo "Cleaning up..."
rm -rf "$TEMP_DIR"

echo "Downloaded assets to $ASSET_DIR"
