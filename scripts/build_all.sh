#!/bin/bash
set -euo pipefail
mkdir -p dist
env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o dist/oni-view-linux-amd64 .
#env CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o dist/oni-view-darwin-amd64 .
#env CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o dist/oni-view-darwin-arm64 .
env CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o dist/oni-view-windows-amd64.exe .
env CGO_ENABLED=0 GOOS=js GOARCH=wasm go build -o dist/oni-view.wasm .

