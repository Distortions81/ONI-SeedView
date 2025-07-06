#!/bin/bash
set -euo pipefail

REPO_ROOT="$(dirname "$0")/.."
ASSET_DIR="$REPO_ROOT/objects"
TEMP_DIR=$(mktemp -d)

if [ -n "$(ls -A "$ASSET_DIR" 2>/dev/null | grep -v '.gitkeep' || true)" ]; then
    echo "Assets already exist in $ASSET_DIR"
    exit 0
fi

if ! command -v dwebp >/dev/null; then
    echo "dwebp command not found. Please install the 'webp' package." >&2
    exit 1
fi

echo "Cloning oni-seed-browser repository..."
git clone --depth 1 https://github.com/MapsNotIncluded/oni-seed-browser "$TEMP_DIR"

echo "Copying assets..."
mkdir -p "$ASSET_DIR"
ASSET_SRC="$TEMP_DIR/app/src/commonMain/composeResources/drawable"
USED_NAMES=$(grep -ho '[A-Za-z0-9_]*\.png' "$REPO_ROOT"/display.go "$REPO_ROOT"/main.go | sort -u)
for name in $USED_NAMES; do
    if [ -f "$ASSET_SRC/$name" ]; then
        cp "$ASSET_SRC/$name" "$ASSET_DIR/$name"
    elif [ -f "$ASSET_SRC/${name%.png}.webp" ]; then
        dwebp "$ASSET_SRC/${name%.png}.webp" -o "$ASSET_DIR/$name" >/dev/null
    else
        echo "Missing asset $name" >&2
    fi
done

echo "Cleaning up..."
rm -rf "$TEMP_DIR"

echo "Downloaded assets to $ASSET_DIR"
