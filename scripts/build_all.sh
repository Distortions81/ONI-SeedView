#!/bin/bash
set -euo pipefail
mkdir -p dist
GOOS=linux GOARCH=amd64 go build -o dist/oni-view-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -o dist/oni-view-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o dist/oni-view-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o dist/oni-view-windows-amd64.exe .
GOOS=js GOARCH=wasm go build -o dist/oni-view.wasm .

