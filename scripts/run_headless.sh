#!/bin/bash
set -euo pipefail

# Run the viewer using a virtual framebuffer.
# Any arguments passed to this script are forwarded to the Go program.

if ! command -v Xvfb > /dev/null; then
  echo "Xvfb is required for headless mode" >&2
  exit 1
fi

xvfb-run --auto-servernum --server-args='-screen 0 1280x720x24' go run . "$@"
