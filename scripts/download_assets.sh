#!/bin/bash
set -euo pipefail

ASSET_DIR="$(dirname "$0")/../assets"
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
find "$ASSET_SRC" -name '*.png' -exec cp {} "$ASSET_DIR" \;
find "$ASSET_SRC" -name '*.webp' -print0 | while IFS= read -r -d '' img; do
    out="$ASSET_DIR/$(basename "${img%.webp}").png"
    dwebp "$img" -o "$out" >/dev/null
done

# copy local icons if present
cp "$(dirname "$0")/../icons"/*.png "$ASSET_DIR" 2>/dev/null || true

echo "Cleaning up..."
rm -rf "$TEMP_DIR"

echo "Downloaded assets to $ASSET_DIR"
