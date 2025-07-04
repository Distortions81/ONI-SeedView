#!/bin/bash
set -euo pipefail
mkdir -p dist
env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o dist/oni-view-linux-amd64 .
#env CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o dist/oni-view-darwin-amd64 .
#env CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o dist/oni-view-darwin-arm64 .
env CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o dist/oni-view-windows-amd64.exe .
env CGO_ENABLED=0 GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o dist/oni-view.wasm .

# Optimize the WebAssembly binary if wasm-opt is available.
if command -v wasm-opt >/dev/null; then
  wasm-opt -Oz \
    --enable-reference-types \
    --enable-bulk-memory \
    --enable-mutable-globals \
    --enable-sign-ext \
    --enable-nontrapping-float-to-int \
    dist/oni-view.wasm -o dist/oni-view.wasm
fi

# Compress the WASM to reduce size.
gzip -9 -c dist/oni-view.wasm > dist/oni-view.wasm.gz
rm dist/oni-view.wasm

# Copy the JS runtime and HTML loader for WASM builds.
cp $(go env GOROOT)/lib/wasm/wasm_exec.js dist/
cp index.html dist/

