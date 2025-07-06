#!/bin/bash
set -euo pipefail
mkdir -p build
#env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o build/oni-view-linux-amd64 .
#env CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o build/oni-view-darwin-amd64 .
#env CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o build/oni-view-darwin-arm64 .
#env CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o build/oni-view-windows-amd64.exe .
env CGO_ENABLED=0 GOOS=js GOARCH=wasm go build -ldflags="-s -w" -trimpath -buildvcs=false -o build/oni-view.wasm .

# Optimize the WebAssembly binary if wasm-opt is available.
#if command -v wasm-opt >/dev/null; then
#  wasm-opt -Oz --strip-dwarf --strip-producers\
#    --enable-reference-types \
#    --enable-bulk-memory \
#    --enable-mutable-globals \
#    --enable-sign-ext \
#    --enable-nontrapping-float-to-int \
#    build/oni-view.wasm -o build/oni-view.wasm
#fi

# Compress the WASM to reduce size using Brotli.
brotli -f -q 11 build/oni-view.wasm -o build/oni-view.wasm.br
rm build/oni-view.wasm

# Copy the JS runtime and HTML loader for WASM builds.
cp -f $(go env GOROOT)/lib/wasm/wasm_exec.js build/
cp -f html/index.html build/
cp -f html/view.html build/
