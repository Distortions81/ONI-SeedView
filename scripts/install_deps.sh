#!/usr/bin/env bash
# install_deps.sh
# Installs tools needed to build the oni-view application.

set -euo pipefail

sudo apt-get update
sudo apt-get install -y git curl build-essential pkg-config \
    libasound2-dev libx11-dev libxft-dev x11proto-dev \
    libxext-dev libxrandr-dev libxinerama-dev \
    libxcursor-dev libxi-dev libxrender-dev libxxf86vm-dev \
    libgl1-mesa-dev libgl1 xorg-dev xvfb webp binaryen brotli \
    protobuf-compiler

# Install Go 1.24.3 if not already available.
if ! go version 2>/dev/null | grep -q "go1.24.3"; then
    curl -LO https://go.dev/dl/go1.24.3.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.24.3.linux-amd64.tar.gz
    rm go1.24.3.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
    export PATH=$PATH:/usr/local/go/bin
fi

# Ensure any PATH updates take effect in the current session.
if [ -f ~/.profile ]; then
    # shellcheck source=/dev/null
    source ~/.profile
fi

# Fetch Go module dependencies.
go mod download

# Install protobuf tooling and regenerate generated code.
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
export PATH="$PATH:$(go env GOPATH)/bin"
protoc --go_out=. --go_opt=paths=source_relative data/pb/seed.proto

echo "Dependencies installed. Run 'go build' to compile the app."

