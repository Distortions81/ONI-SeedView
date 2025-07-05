#!/bin/bash
set -euo pipefail

# Convert all PNG textures in biomes/ to JPEG format.
# Requires Go to be installed.

go run "$(dirname "$0")/convert_biomes.go"
