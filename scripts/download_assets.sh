#!/bin/bash
set -euo pipefail

ASSET_DIR="$(dirname "$0")/../assets"
TEMP_DIR=$(mktemp -d)

if ! command -v dwebp >/dev/null; then
    echo "dwebp command not found. Please install the 'webp' package." >&2
    exit 1
fi

echo "Cloning ONITraitFinder repository..."
git clone --depth 1 https://github.com/MapsNotIncluded/ONITraitFinder "$TEMP_DIR"

echo "Copying assets..."
mkdir -p "$ASSET_DIR"
find "$TEMP_DIR"/docs/images -name '*.png' -exec cp {} "$ASSET_DIR" \;
find "$TEMP_DIR"/docs/images -name '*.webp' -print0 | while IFS= read -r -d '' img; do
    out="$ASSET_DIR/$(basename "${img%.webp}").png"
    dwebp "$img" -o "$out" >/dev/null
done

echo "Cleaning up..."
rm -rf "$TEMP_DIR"

echo "Downloaded assets to $ASSET_DIR"
